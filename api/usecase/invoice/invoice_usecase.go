package invoice

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	bordir "github.com/be-sistem-informasi-konveksi/api/repository/bordir/mysql/gorm"
	invoice "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm"
	produk "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm"
	sablon "github.com/be-sistem-informasi-konveksi/api/repository/sablon/mysql/gorm"
	user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/invoice"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"golang.org/x/sync/errgroup"
)

type InvoiceUsecase interface {
	GetAll(ctx context.Context, reqFilter req.GetAll) ([]entity.Invoice, error)
	Create(ctx context.Context, reqInvoice req.Create) error
}

type invoiceUsecase struct {
	repo       invoice.InvoiceRepo
	repoUser   user.UserRepo
	repoBordir bordir.BordirRepo
	repoSablon sablon.SablonRepo
	repoProduk produk.ProdukRepo
	ulid       pkg.UlidPkg
}

func NewInvoiceUsecase(
	repo invoice.InvoiceRepo,
	repoUser user.UserRepo,
	repoBordir bordir.BordirRepo,
	repoSablon sablon.SablonRepo,
	repoProduk produk.ProdukRepo,
	ulid pkg.UlidPkg,
) InvoiceUsecase {
	return &invoiceUsecase{repo, repoUser, repoBordir, repoSablon, repoProduk, ulid}
}

// var invPool = sync.Pool{
// 	New: func() any {
// 		return &entity.Invoice{}
// 	},
// }

func convertRefNum2Uint(refNum string) (uint64, error) {
	var refNumUint uint64
	refNum = strings.TrimLeft(refNum, "0")

	// Convert the substring to an integer
	if refNum == "" {
		refNumUint = 0 // Treat empty or all-zero strings as zero
	} else if i, err := strconv.ParseUint(refNum, 10, 0); err == nil {
		refNumUint = i
	} else {
		helper.LogsError(err)
		return refNumUint, err
	}
	return refNumUint, nil
}

func createRefNumber(repo invoice.InvoiceRepo, ctx context.Context) (string, error) {
	data, err := repo.GetLastInvoiceCrrYear(ctx)
	if err != nil && err.Error() == "record not found" {
		return fmt.Sprintf("%04d", 1), nil
	}
	if err != nil {
		return "", err
	}

	refNumUint, err := convertRefNum2Uint(data.NomorReferensi)
	if err != nil {
		return "", err
	}

	nextID := refNumUint + 1
	return fmt.Sprintf("%04d", nextID), nil
}

func handleErr(err error, errMsg string) error {
	if err != nil {
		helper.LogsError(err)
		if err.Error() == "record not found" {
			return errors.New(errMsg)
		}
		return err
	}
	return nil
}

func (u *invoiceUsecase) Create(ctx context.Context, reqInvoice req.Create) error {
	tanggalDeadline, err := time.Parse(time.RFC3339, reqInvoice.TanggalDeadline)
	if err != nil {
		helper.LogsError(err)
		return err
	}
	tanggalKirim, err := time.Parse(time.RFC3339, reqInvoice.TanggalKirim)
	if err != nil {
		helper.LogsError(err)
		return err
	}

	lengthDetail := len(reqInvoice.DetailInvoice)
	detailInvoices := make([]entity.DetailInvoice, lengthDetail)

	var total_qty float64
	var paramRepo invoice.CreateParam
	invoice_id := u.ulid.MakeUlid().String()
	wg := &sync.WaitGroup{}
	g := errgroup.Group{}
	g.SetLimit(10)

	g.Go(func() error {
		_, err := u.repoUser.GetById(ctx, reqInvoice.UserID)
		return handleErr(err, message.UserNotFound)
	})

	for i, v := range reqInvoice.DetailInvoice {
		total_qty += v.Qty
		wg.Add(1)
		func(v req.ReqDetailInvoice, i int) {
			g.Go(func() error {
				var total_harga float64
				dataBordir, err := u.repoBordir.GetById(ctx, v.BordirID)
				if err := handleErr(err, fmt.Sprintf("%d %s", i, message.BordirNotFound)); err != nil {
					return err
				}
				total_harga += dataBordir.Harga
				dataSablon, err := u.repoSablon.GetById(ctx, v.SablonID)
				if err := handleErr(err, fmt.Sprintf("%d %s", i, message.SablonNotFound)); err != nil {
					return err
				}
				total_harga += dataSablon.Harga
				dataProduk, err := u.repoProduk.GetById(ctx, v.ProdukID)
				var harga_produk float64
				for _, detailHarga := range dataProduk.HargaDetails {
					if detailHarga.ID == v.HargaProdukID {
						harga_produk = detailHarga.Harga
						break
					}
				}
				if err := handleErr(err, fmt.Sprintf("%d %s", i, message.ProdukNotFound)); err != nil {
					return err
				}
				if harga_produk == 0 {
					return fmt.Errorf("%d %s", i, message.HargaDetailProdukNotFound)
				}
				total_harga += harga_produk
				detailInvoices[i] = entity.DetailInvoice{
					Base: entity.Base{
						ID: u.ulid.MakeUlid().String(),
					},
					ProdukID:     v.ProdukID,
					InvoiceID:    invoice_id,
					SablonID:     v.SablonID,
					BordirID:     v.BordirID,
					Qty:          int(v.Qty),
					GambarDesign: v.GambarDesign,
					Total:        total_harga * v.Qty,
				}
				return nil
			})
		}(v, i)
	}

	if err := g.Wait(); err != nil {
		return err
	}

	var total_harga_invoice float64

	for _, v := range detailInvoices {
		total_harga_invoice += v.Total
	}
	paramRepo.DetailInvoice = detailInvoices

	ref, err := createRefNumber(u.repo, ctx)
	if err != nil {
		return err
	}

	// setup entity invoice
	total_harga_invoice *= total_qty
	status_bayar := "BELUM_LUNAS"
	sisa_tagihan := reqInvoice.Bayar - total_harga_invoice
	if sisa_tagihan >= 0 {
		status_bayar = "LUNAS"
	}

	paramRepo.Invoice = &entity.Invoice{
		Base: entity.Base{
			ID: invoice_id,
		},
		TotalQty:        int(total_qty),
		TotalHarga:      total_harga_invoice,
		StatusProduksi:  reqInvoice.StatusProduksi,
		BuktiPembayaran: reqInvoice.BuktiPembayaran,
		UserID:          reqInvoice.UserID,
		Keterangan:      reqInvoice.Keterangan,
		TanggalDeadline: &tanggalDeadline,
		TanggalKirim:    &tanggalKirim,
		NomorReferensi:  ref,
		Kepada:          reqInvoice.Kepada,
		NoTelp:          reqInvoice.NoTelp,
		Alamat:          reqInvoice.Alamat,
		StatusBayar:     status_bayar,
		SisaTagihan:     math.Abs(sisa_tagihan),
	}

	if err := u.repo.Create(ctx, paramRepo); err != nil {
		return err
	}
	return nil
}

func (u *invoiceUsecase) GetAll(ctx context.Context, reqFilter req.GetAll) ([]entity.Invoice, error) {
	var tglDeadline, tglKirim time.Time
	var err error
	if reqFilter.TanggalDeadline != "" {
		tglDeadline, err = time.Parse(time.RFC3339, reqFilter.TanggalDeadline)
		if err != nil {
			return nil, err
		}
	}
	if reqFilter.TanggalKirim != "" {
		tglKirim, err = time.Parse(time.RFC3339, reqFilter.TanggalKirim)
		if err != nil {
			return nil, err
		}
	}

	filter := invoice.SearchFilter{
		StatusProduksi:  reqFilter.StatusProduksi,
		TanggalDeadline: tglDeadline,
		TanggalKirim:    tglKirim,
		Kepada:          reqFilter.Kepada,
		Limit:           reqFilter.Limit,
		Next:            reqFilter.Next,
	}

	if reqFilter.SortBy != "" {
		filter.Order += strings.ToLower(reqFilter.SortBy)
		if reqFilter.OrderBy != "" {
			filter.Order += " " + reqFilter.OrderBy
		} else {
			filter.Order += " ASC"
		}
	}

	if reqFilter.Limit <= 0 {
		filter.Limit = 10
	}

	invoices, err := u.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

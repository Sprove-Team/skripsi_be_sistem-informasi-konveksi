package invoice

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	akuntansi "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
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
)

type InvoiceUsecase interface {
	GetAll(ctx context.Context, reqFilter req.GetAll) ([]entity.Invoice, error)
	Create(ctx context.Context, reqInvoice req.Create) error
}

type invoiceUsecase struct {
	repo       invoice.InvoiceRepo
	repoAkun   akuntansi.AkunRepo
	repoUser   user.UserRepo
	repoBordir bordir.BordirRepo
	repoSablon sablon.SablonRepo
	repoProduk produk.ProdukRepo
	ulid       pkg.UlidPkg
}

func NewInvoiceUsecase(
	repo invoice.InvoiceRepo,
	repoAkun akuntansi.AkunRepo,
	repoUser user.UserRepo,
	repoBordir bordir.BordirRepo,
	repoSablon sablon.SablonRepo,
	repoProduk produk.ProdukRepo,
	ulid pkg.UlidPkg,
) InvoiceUsecase {
	return &invoiceUsecase{repo, repoAkun, repoUser, repoBordir, repoSablon, repoProduk, ulid}
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

	var paramRepo invoice.CreateParam
	invoiceID := u.ulid.MakeUlid().String()

	if _, err := u.repoUser.GetById(ctx, reqInvoice.UserID); err != nil {
		return handleErr(err, message.UserNotFound)
	}

	for i, v := range reqInvoice.DetailInvoice {
		detailInvoices[i] = entity.DetailInvoice{
			Base: entity.Base{
				ID: u.ulid.MakeUlid().String(),
			},
			ProdukID:     v.ProdukID,
			InvoiceID:    invoiceID,
			BordirID:     v.BordirID,
			SablonID:     v.SablonID,
			GambarDesign: v.GambarDesign,
			Qty:          v.Qty,
			Total:        v.Total,
		}
	}

	paramRepo.DetailInvoice = detailInvoices

	ref, err := createRefNumber(u.repo, ctx)
	if err != nil {
		return err
	}

	// setup entity invoice
	statusBayar := "BELUM_LUNAS"
	sisaTagihan := reqInvoice.Bayar - reqInvoice.TotalHarga
	sisaTagihanABS := math.Abs(sisaTagihan)
	if sisaTagihan >= 0 {
		statusBayar = "LUNAS"
	}

	// case in akuntansi
	// create kontak
	var kontakID string
	if reqInvoice.KontakID != "" {
		kontakID = reqInvoice.KontakID
	} else {
		kontakID = u.ulid.MakeUlid().String()
		paramRepo.Kontak = &entity.Kontak{
			Base: entity.Base{
				ID: kontakID,
			},
			Nama:       reqInvoice.NewKontak.Nama,
			NoTelp:     reqInvoice.NewKontak.NoTelp,
			Alamat:     reqInvoice.NewKontak.Alamat,
			Keterangan: reqInvoice.Keterangan,
			Email:      reqInvoice.NewKontak.Email,
		}
	}

	// create transaksi
	var transaksi_id = u.ulid.MakeUlid().String()

	// currentDate := time.Now()
	paramRepo.Transaksi = &entity.Transaksi{
		Base: entity.Base{
			ID: transaksi_id,
		},
		Tanggal:         time.Now(),
		KontakID:        kontakID,
		Total:           reqInvoice.TotalHarga,
		Keterangan:      reqInvoice.Keterangan,
		BuktiPembayaran: reqInvoice.BuktiPembayaran,
	}

	// ayat jurnals
	akunPembayaran, err := u.repoAkun.GetById(ctx, reqInvoice.AkunID)
	if err != nil {
		return errors.New(message.AkunNotFound)
	} else {
		// check akun is ASET LANCAR
		if akunPembayaran.Kode[0:2] != "11" {
			return errors.New(message.AkunNotFound)
		}
	}
	akuns, err := u.repoAkun.GetAkunByNames(ctx, []string{"piutang usaha", "pendapatan jasa"})
	if err != nil {
		return err
	}
	akuns = append(akuns, akunPembayaran)

	ayatJurnals := make([]entity.AyatJurnal, len(akuns))

	for i, akun := range akuns {
		ayatJurnal := entity.AyatJurnal{
			Base: entity.Base{
				ID: u.ulid.MakeUlid().String(),
			},
			AkunID: akun.ID,
		}

		if akun.ID == reqInvoice.AkunID {
			if statusBayar == "LUNAS" {
				ayatJurnal.Debit, ayatJurnal.Saldo = reqInvoice.Bayar, reqInvoice.Bayar
			} else {
				ayatJurnal.Debit, ayatJurnal.Saldo = reqInvoice.Bayar, reqInvoice.Bayar
			}
		}

		switch akun.Nama {
		case "pendapatan jasa":
			ayatJurnal.Kredit, ayatJurnal.Saldo = reqInvoice.TotalHarga, reqInvoice.TotalHarga
		case "piutang usaha":
			if statusBayar == "BELUM_LUNAS" {
				ayatJurnal.Debit, ayatJurnal.Saldo = sisaTagihanABS, sisaTagihanABS
			}
		}

		if ayatJurnal.Saldo != 0 {
			ayatJurnals[i] = ayatJurnal
		}
	}

	paramRepo.Transaksi.AyatJurnals = ayatJurnals

	// end ayat jurnals

	paramRepo.Invoice = &entity.Invoice{
		Base: entity.Base{
			ID: invoiceID,
		},
		TotalQty:        reqInvoice.TotalQty,
		TotalHarga:      reqInvoice.TotalHarga,
		StatusProduksi:  reqInvoice.StatusProduksi,
		BuktiPembayaran: reqInvoice.BuktiPembayaran,
		UserID:          reqInvoice.UserID,
		Keterangan:      reqInvoice.Keterangan,
		TanggalDeadline: &tanggalDeadline,
		TanggalKirim:    &tanggalKirim,
		NomorReferensi:  ref,
		StatusBayar:     statusBayar,
		SisaTagihan:     sisaTagihanABS,
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

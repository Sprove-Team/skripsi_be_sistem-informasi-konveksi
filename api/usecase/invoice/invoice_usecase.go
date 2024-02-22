package invoice

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	akunRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	kontakRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kontak"
	bordirRepo "github.com/be-sistem-informasi-konveksi/api/repository/bordir/mysql/gorm"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm"
	produkRepo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm"
	sablonRepo "github.com/be-sistem-informasi-konveksi/api/repository/sablon/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/common/message"
	reqHP "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	req "github.com/be-sistem-informasi-konveksi/common/request/invoice"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"golang.org/x/sync/errgroup"
)

type (
	ParamCreateDataInvoice struct {
		Ctx context.Context
		Req req.Create
	}
	ParamUpdateDataInvoice struct {
		Ctx context.Context
		Req req.Update
	}
	ParamCommitDB struct {
		Ctx     context.Context
		Invoice *entity.Invoice
	}
	ParamGetAll struct {
		Ctx context.Context
		Req req.GetAll
	}
	ParamGetById struct {
		Ctx context.Context
		ID  string
	}
	ParamCreateDataDetails struct {
		Ctx context.Context
		Req []req.Create
	}
	ParamCheckDataDetails struct {
		Ctx       context.Context
		ProdukIds []string
		BordirIds []string
		SablonIds []string
	}
)

type InvoiceUsecase interface {
	CreateDataInvoice(param ParamCreateDataInvoice) (*entity.Invoice, *reqHP.Create, error)
	CreateCommitDB(param ParamCommitDB) error
	UpdateDataInvoice(param ParamUpdateDataInvoice) (*entity.Invoice, error)
	SaveCommitDB(param ParamCommitDB) error
	CheckDataDetails(param ParamCheckDataDetails) error
	GetAll(param ParamGetAll) ([]entity.Invoice, error)
	GetById(param ParamGetById) (*entity.Invoice, error)
}

type invoiceUsecase struct {
	repo       repo.InvoiceRepo
	repoAkun   akunRepo.AkunRepo
	repoKontak kontakRepo.KontakRepo
	repoProduk produkRepo.ProdukRepo
	repoBordir bordirRepo.BordirRepo
	repoSablon sablonRepo.SablonRepo
	ulid       pkg.UlidPkg
}

func NewInvoiceUsecase(
	repo repo.InvoiceRepo,
	repoAkun akunRepo.AkunRepo,
	repoKontak kontakRepo.KontakRepo,
	repoProduk produkRepo.ProdukRepo,
	repoBordir bordirRepo.BordirRepo,
	repoSablon sablonRepo.SablonRepo,
	ulid pkg.UlidPkg,
) InvoiceUsecase {
	return &invoiceUsecase{repo, repoAkun, repoKontak, repoProduk,
		repoBordir, repoSablon, ulid}
}

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

func createRefNumber(repo repo.InvoiceRepo, ctx context.Context) (string, error) {
	data, err := repo.GetLastInvoiceCrrYear(ctx)
	if err != nil {
		if err.Error() == "record not found" {
			return fmt.Sprintf("%04d", 1), nil
		}
		return "", err
	}

	refNumUint, err := convertRefNum2Uint(data.NomorReferensi)
	if err != nil {
		return "", err
	}

	nextID := refNumUint + 1
	return fmt.Sprintf("%04d", nextID), nil
}

func (u *invoiceUsecase) CheckDataDetails(param ParamCheckDataDetails) error {

	g := new(errgroup.Group)

	g.Go(func() error {
		d, err := u.repoProduk.GetByIds(param.Ctx, param.ProdukIds)
		if err != nil {
			return err
		}
		if len(d) != len(param.ProdukIds) {
			return errors.New(message.ProdukNotFound)
		}
		return nil
	})

	g.Go(func() error {
		d, err := u.repoBordir.GetByIds(param.Ctx, param.BordirIds)
		if err != nil {
			return err
		}
		if len(d) != len(param.BordirIds) {
			return errors.New(message.BordirNotFound)
		}
		return nil
	})

	g.Go(func() error {
		d, err := u.repoSablon.GetByIds(param.Ctx, param.SablonIds)
		if err != nil {
			return err
		}
		if len(d) != len(param.SablonIds) {
			return errors.New(message.SablonNotFound)
		}
		return nil
	})

	if err := g.Wait(); err != nil {

		return err
	}
	return nil
}

func (u *invoiceUsecase) CreateDataInvoice(param ParamCreateDataInvoice) (*entity.Invoice, *reqHP.Create, error) {
	fmt.Println("akun_id -> ", param.Req.Bayar.AkunID)
	tanggalDeadline, err := time.Parse(time.RFC3339, param.Req.TanggalDeadline)
	if err != nil {
		helper.LogsError(err)
		return nil, nil, err
	}
	tanggalKirim, err := time.Parse(time.RFC3339, param.Req.TanggalKirim)
	if err != nil {
		helper.LogsError(err)
		return nil, nil, err
	}

	lengthDetail := len(param.Req.DetailInvoice)

	var dataInvoice *entity.Invoice

	invoiceID := u.ulid.MakeUlid().String()

	detailInvoices := make([]entity.DetailInvoice, lengthDetail)
	var produkIds = make([]string, lengthDetail)
	var bordirIds = make([]string, lengthDetail)
	var sablonIds = make([]string, lengthDetail)

	for i, v := range param.Req.DetailInvoice {
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
		produkIds[i] = v.ProdukID
		bordirIds[i] = v.BordirID
		sablonIds[i] = v.SablonID
	}

	g := new(errgroup.Group)
	g.SetLimit(10)
	g.Go(func() error {
		return u.CheckDataDetails(ParamCheckDataDetails{
			Ctx:       param.Ctx,
			ProdukIds: produkIds,
			BordirIds: bordirIds,
			SablonIds: sablonIds,
		})
	})
	var ref string
	g.Go(func() error {
		var err error
		ref, err = createRefNumber(u.repo, param.Ctx)
		if err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		_, err := u.repoAkun.GetById(param.Ctx, param.Req.Bayar.AkunID)
		if err != nil {
			if err.Error() == "record not found" {
				return errors.New(message.AkunNotFound)
			}
			helper.LogsError(err)
			return err
		}
		return nil
	})

	var akuns []entity.Akun
	g.Go(func() error {
		var err error
		akuns, err = u.repoAkun.GetAkunByNames(param.Ctx, []string{"piutang usaha", "pendapatan jasa"})
		if err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		helper.LogsError(err)
		return nil, nil, err
	}

	tanggalTr := time.Now().Format(time.RFC3339)

	sisa := param.Req.TotalHarga - param.Req.Bayar.Total

	if sisa < 0 {
		return nil, nil, errors.New(message.BayarMustLessThanTotalHargaInvoice)
	}
	reqHpUC := reqHP.Create{
		KontakID:   param.Req.KontakID,
		InvoiceID:  invoiceID,
		Jenis:      "PIUTANG",
		Keterangan: param.Req.Keterangan,
		Transaksi: reqHP.ReqTransaksi{
			Tanggal:         tanggalTr,
			BuktiPembayaran: param.Req.Bayar.BuktiPembayaran,
			AyatJurnal:      make([]reqHP.ReqAyatJurnal, 0, 2),
		},
	}

	for _, akun := range akuns {
		ayatJurnal := reqHP.ReqAyatJurnal{
			AkunID: akun.ID,
		}
		switch akun.Nama {
		case "pendapatan jasa":
			ayatJurnal.Kredit = param.Req.TotalHarga
		case "piutang usaha":
			ayatJurnal.Debit = param.Req.TotalHarga
		}
		reqHpUC.Transaksi.AyatJurnal = append(reqHpUC.Transaksi.AyatJurnal, ayatJurnal)
	}

	dataBayar := entity.DataBayarInvoice{
		Base: entity.Base{
			ID: u.ulid.MakeUlid().String(),
		},
		AkunID:          param.Req.Bayar.AkunID,
		Keterangan:      param.Req.Bayar.Keterangan,
		BuktiPembayaran: param.Req.Bayar.BuktiPembayaran,
		Total:           param.Req.Bayar.Total,
	}

	dataInvoice = &entity.Invoice{
		Base: entity.Base{
			ID: invoiceID,
		},
		KontakID:         param.Req.KontakID,
		TotalQty:         param.Req.TotalQty,
		TotalHarga:       param.Req.TotalHarga,
		StatusProduksi:   param.Req.StatusProduksi,
		UserID:           param.Req.UserID,
		Keterangan:       param.Req.Keterangan,
		TanggalDeadline:  &tanggalDeadline,
		TanggalKirim:     &tanggalKirim,
		NomorReferensi:   ref,
		DetailInvoices:   detailInvoices,
		DataBayarInvoice: []entity.DataBayarInvoice{dataBayar},
	}

	return dataInvoice, &reqHpUC, nil
}

func (u *invoiceUsecase) CreateCommitDB(param ParamCommitDB) error {
	if param.Invoice.KontakID != "" {
		_, err := u.repoKontak.GetById(kontakRepo.ParamGetById{
			Ctx: param.Ctx,
			ID:  param.Invoice.KontakID,
		})
		if err != nil {
			if err.Error() == "record not found" {
				return errors.New(message.KontakNotFound)
			}
			return err
		}
	}

	return u.repo.Create(repo.ParamCreate(param))
}

func (u *invoiceUsecase) UpdateDataInvoice(param ParamUpdateDataInvoice) (*entity.Invoice, error) {
	oldData, err := u.repo.GetByIdFullAssoc(repo.ParamGetById{
		Ctx: param.Ctx,
		ID:  param.Req.ID,
	})

	if err != nil {
		return nil, err
	}

	oldData.Keterangan = param.Req.Keterangan
	oldData.KontakID = param.Req.KontakID
	oldData.StatusProduksi = param.Req.StatusProduksi
	oldData.TotalHarga = param.Req.TotalHarga
	oldData.TotalQty = param.Req.TotalQty

	if param.Req.TanggalKirim != "" {
		tanggalKirim, err := time.Parse(time.RFC3339, param.Req.TanggalKirim)
		if err != nil {
			return nil, err
		}
		oldData.TanggalKirim = &tanggalKirim
	}

	if param.Req.TanggalDeadline != "" {
		tanggalDeadline, err := time.Parse(time.RFC3339, param.Req.TanggalDeadline)
		if err != nil {
			return nil, err
		}
		oldData.TanggalDeadline = &tanggalDeadline
	}

	var oldTotalBayar = oldData.TotalHarga - oldData.HutangPiutang.Sisa
	if param.Req.TotalHarga != 0 {
		if param.Req.TotalHarga < oldTotalBayar {
			return nil, errors.New(message.TotalBayarMustGeOrEqToTotalByr)
		}
		oldData.HutangPiutang.Sisa = param.Req.TotalHarga - oldTotalBayar
		oldData.HutangPiutang.Total = param.Req.TotalHarga

		if oldData.HutangPiutang.Sisa >= 0 {
			oldData.HutangPiutang.Status = "LUNAS"
		}
	}

	// update transaksi hutang piutang
	oldData.HutangPiutang.Transaksi.Keterangan = param.Req.Keterangan
	oldData.HutangPiutang.Transaksi.Total = param.Req.TotalHarga
	oldData.HutangPiutang.Transaksi.KontakID = param.Req.KontakID

	return oldData, nil
}

func (u *invoiceUsecase) SaveCommitDB(param ParamCommitDB) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		_, err := u.repoKontak.GetById(kontakRepo.ParamGetById{
			Ctx: param.Ctx,
			ID:  param.Invoice.KontakID,
		})
		if err != nil {
			if err.Error() == "record not found" {
				return errors.New(message.KontakNotFound)
			}
			return err
		}
		return nil
	})

	g.Go(func() error {
		_, err := u.repo.GetById(repo.ParamGetById{
			Ctx: param.Ctx,
			ID:  param.Invoice.ID,
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return u.repo.Save(repo.ParamSave(param))
}

func (u *invoiceUsecase) GetAll(param ParamGetAll) ([]entity.Invoice, error) {
	var tglDeadline, tglKirim time.Time
	var err error
	if param.Req.TanggalDeadline != "" {
		tglDeadline, err = time.Parse(time.RFC3339, param.Req.TanggalDeadline)
		if err != nil {
			return nil, err
		}
	}
	if param.Req.TanggalKirim != "" {
		tglKirim, err = time.Parse(time.RFC3339, param.Req.TanggalKirim)
		if err != nil {
			return nil, err
		}
	}

	paramRepo := repo.ParamGetAll{
		Ctx:             param.Ctx,
		StatusProduksi:  param.Req.StatusProduksi,
		TanggalDeadline: tglDeadline,
		TanggalKirim:    tglKirim,
		KontakID:        param.Req.KontakID,
		Limit:           param.Req.Limit,
		Next:            param.Req.Next,
		Order:           param.Req.OrderBy,
	}

	if param.Req.SortBy != "" {
		paramRepo.Order += strings.ToLower(param.Req.SortBy)
		if param.Req.OrderBy != "" {
			paramRepo.Order += " " + param.Req.OrderBy
		} else {
			paramRepo.Order += " ASC"
		}
	}

	if param.Req.Limit <= 0 {
		paramRepo.Limit = 10
	}

	invoices, err := u.repo.GetAll(paramRepo)
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (u *invoiceUsecase) GetById(param ParamGetById) (*entity.Invoice, error) {
	return u.repo.GetByIdFullAssoc(repo.ParamGetById(param))
}

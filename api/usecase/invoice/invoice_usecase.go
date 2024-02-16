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
	repo "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm"
	userRepo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/common/message"
	reqHP "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	req "github.com/be-sistem-informasi-konveksi/common/request/invoice"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type (
	ParamCreateDataInvoice struct {
		Ctx context.Context
		Req req.Create
	}
	ParamCreateCommitDB struct {
		Ctx     context.Context
		Invoice *entity.Invoice
	}
	ParamGetAll struct {
		Ctx context.Context
		Req req.GetAll
	}
)

type InvoiceUsecase interface {
	CreateDataInvoice(param ParamCreateDataInvoice) (*entity.Invoice, *reqHP.Create, error)
	CreateCommitDB(param ParamCreateCommitDB) error
	GetAll(param ParamGetAll) ([]entity.Invoice, error)
}

type invoiceUsecase struct {
	repo       repo.InvoiceRepo
	repoAkun   akunRepo.AkunRepo
	repoUser   userRepo.UserRepo
	repoKontak kontakRepo.KontakRepo
	ulid       pkg.UlidPkg
}

func NewInvoiceUsecase(
	repo repo.InvoiceRepo,
	repoAkun akunRepo.AkunRepo,
	repoUser userRepo.UserRepo,
	repoKontak kontakRepo.KontakRepo,
	ulid pkg.UlidPkg,
) InvoiceUsecase {
	return &invoiceUsecase{repo, repoAkun, repoUser, repoKontak, ulid}
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

func (u *invoiceUsecase) CreateDataInvoice(param ParamCreateDataInvoice) (*entity.Invoice, *reqHP.Create, error) {
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
	}

	ref, err := createRefNumber(u.repo, param.Ctx)
	if err != nil {
		return nil, nil, err
	}

	// ayat jurnals
	akunPembayaran, err := u.repoAkun.GetById(param.Ctx, param.Req.Bayar.AkunBayarID)
	if err != nil {
		return nil, nil, errors.New(message.AkunNotFound)
	}

	akuns, err := u.repoAkun.GetAkunByNames(param.Ctx, []string{"piutang usaha", "pendapatan jasa"})
	if err != nil {
		return nil, nil, err
	}

	tanggalTr := time.Now().Format(time.RFC3339)

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
		BayarAwal: make([]reqHP.ReqBayar, 0, 1),
	}

	akuns = append(akuns, akunPembayaran)

	for _, akun := range akuns {

		if akun.ID == param.Req.Bayar.AkunBayarID {
			reqHpUC.BayarAwal = append(reqHpUC.BayarAwal, reqHP.ReqBayar{
				Tanggal:         tanggalTr,
				BuktiPembayaran: param.Req.Bayar.BuktiPembayaran,
				Keterangan:      param.Req.Bayar.Keterangan,
				AkunBayarID:     param.Req.Bayar.AkunBayarID,
				Total:           param.Req.Bayar.Total,
			})
		} else {
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

	}

	// end ayat jurnals
	dataInvoice = &entity.Invoice{
		Base: entity.Base{
			ID: invoiceID,
		},
		KontakID:        param.Req.KontakID,
		TotalQty:        param.Req.TotalQty,
		TotalHarga:      param.Req.TotalHarga,
		StatusProduksi:  param.Req.StatusProduksi,
		UserID:          param.Req.UserID,
		Keterangan:      param.Req.Keterangan,
		TanggalDeadline: &tanggalDeadline,
		TanggalKirim:    &tanggalKirim,
		NomorReferensi:  ref,
		DetailInvoices:  detailInvoices,
	}

	return dataInvoice, &reqHpUC, nil
}

func (u *invoiceUsecase) CreateCommitDB(param ParamCreateCommitDB) error {
	_, err := u.repoUser.GetById(userRepo.ParamGetById{
		Ctx: param.Ctx,
		ID:  param.Invoice.UserID,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New(message.UserNotFound)
		}
		return err
	}

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
	// param.Kontak

	return u.repo.Create(repo.CreateParam(param))
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

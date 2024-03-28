package uc_invoice

import (
	"context"
	"errors"
	"fmt"
	"sort"
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
		Ctx    context.Context
		Claims *pkg.Claims
		Req    req.Create
	}
	ParamUpdateDataInvoice struct {
		Ctx    context.Context
		Claims *pkg.Claims
		Req    req.Update
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
	ParamDelete struct {
		Ctx context.Context
		ID  string
	}
)

type InvoiceUsecase interface {
	CreateDataInvoice(param ParamCreateDataInvoice) (*entity.Invoice, *reqHP.Create, error)
	CreateCommitDB(param ParamCommitDB) error
	UpdateDataInvoice(param ParamUpdateDataInvoice) (*entity.Invoice, error)
	SaveCommitDB(param ParamCommitDB) error
	CheckDataDetails(param ParamCheckDataDetails) error
	Delete(param ParamDelete) error
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
		count := helper.CountUniqueElements(param.ProdukIds)

		d, err := u.repoProduk.GetByIds(param.Ctx, param.ProdukIds)
		if err != nil {
			return err
		}

		if len(d) != count {
			return errors.New(message.ProdukNotFound)
		}
		return nil
	})
	if len(param.SablonIds) > 0 {
		g.Go(func() error {
			count := helper.CountUniqueElements(param.BordirIds)
			d, err := u.repoBordir.GetByIds(param.Ctx, param.BordirIds)
			if err != nil {
				return err
			}
			if len(d) != count {
				return errors.New(message.BordirNotFound)
			}
			return nil
		})
	}
	if len(param.SablonIds) > 0 {
		g.Go(func() error {
			count := helper.CountUniqueElements(param.SablonIds)
			d, err := u.repoSablon.GetByIds(param.Ctx, param.SablonIds)
			if err != nil {
				return err
			}
			if len(d) != count {
				return errors.New(message.SablonNotFound)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
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

	if param.Req.NewKontak.Nama != "" {
		param.Req.KontakID = u.ulid.MakeUlid().String()
		err := u.repoKontak.Create(kontakRepo.ParamCreate{
			Ctx: param.Ctx,
			Kontak: &entity.Kontak{
				Base: entity.Base{
					ID: param.Req.KontakID,
				},
				Nama:       param.Req.NewKontak.Nama,
				NoTelp:     param.Req.NewKontak.NoTelp,
				Alamat:     param.Req.NewKontak.Alamat,
				Email:      param.Req.NewKontak.Email,
				Keterangan: "pelanggan pengguna jasa",
			},
		})
		if err != nil {
			return nil, nil, err
		}
	}

	lengthDetail := len(param.Req.DetailInvoice)

	var dataInvoice *entity.Invoice

	invoiceID := u.ulid.MakeUlid().String()

	detailInvoice := make([]entity.DetailInvoice, lengthDetail)
	var produkIds = make([]string, lengthDetail)
	var bordirIds = make([]string, lengthDetail)
	var sablonIds = make([]string, lengthDetail)

	var totalHarga float64
	var totalQty int
	for i, v := range param.Req.DetailInvoice {
		detailInvoice[i] = entity.DetailInvoice{
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
		totalHarga += v.Total
		totalQty += v.Qty
	}

	g := new(errgroup.Group)
	g.SetLimit(10)
	g.Go(func() error {
		param := ParamCheckDataDetails{
			Ctx:       param.Ctx,
			ProdukIds: produkIds,
		}
		if len(bordirIds) > 0 {
			param.BordirIds = bordirIds
		}
		if len(sablonIds) > 0 {
			param.BordirIds = sablonIds
		}

		return u.CheckDataDetails(param)
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

	sisa := totalHarga - param.Req.Bayar.Total

	if sisa < 0 {
		return nil, nil, errors.New(message.BayarMustLessThanTotalHargaInvoice)
	}
	reqHpUC := reqHP.Create{
		KontakID:   param.Req.KontakID,
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
			ayatJurnal.Kredit = totalHarga
		case "piutang usaha":
			ayatJurnal.Debit = totalHarga
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
		TotalQty:         totalQty,
		TotalHarga:       totalHarga,
		UserID:           param.Claims.ID,
		Keterangan:       param.Req.Keterangan,
		TanggalDeadline:  &tanggalDeadline,
		TanggalKirim:     &tanggalKirim,
		NomorReferensi:   ref,
		DetailInvoice:    detailInvoice,
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

	oldData.Base = entity.Base{
		ID: oldData.ID,
	}
	switch param.Claims.Role {
	case entity.RolesById[4]: // manager produksi
		return &entity.Invoice{
			Base:           oldData.Base,
			UserID:         param.Claims.ID,
			KontakID:       oldData.KontakID,
			StatusProduksi: param.Req.StatusProduksi,
		}, nil
	case entity.RolesById[1]: // direktur
		if param.Req.StatusProduksi != "" {
			oldData.StatusProduksi = param.Req.StatusProduksi
		}
	default:
		if param.Req.StatusProduksi != "" {
			return nil, errors.New(strings.ToLower(param.Claims.Role) + message.UserNotAllowedToModifiedStatusProdusi)
		}
	}

	oldData.UserID = param.Claims.ID
	if param.Req.Keterangan != "" {
		oldData.Keterangan = param.Req.Keterangan
	}

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

	lengthDetail := len(param.Req.DetailInvoice)

	if lengthDetail > 0 {

		sort.Slice(oldData.DetailInvoice, func(i, j int) bool {
			return oldData.DetailInvoice[i].ID < oldData.DetailInvoice[j].ID
		})

		for _, detail := range param.Req.DetailInvoice {
			idx := sort.Search(len(oldData.DetailInvoice), func(i int) bool {
				return oldData.DetailInvoice[i].ID >= detail.ID
			})

			if idx < len(oldData.DetailInvoice) && oldData.DetailInvoice[idx].ID == detail.ID {

				oldData.TotalHarga = (oldData.TotalHarga - oldData.DetailInvoice[idx].Total) + detail.Total
				oldData.TotalQty = (oldData.TotalQty - oldData.DetailInvoice[idx].Qty) + detail.Qty

				oldData.DetailInvoice[idx].Qty = detail.Qty
				oldData.DetailInvoice[idx].Total = detail.Total

			} else {
				return nil, errors.New(message.DetailInvoiceNotFound)
			}
		}

		// update oldData total harga & qty
		var oldTotalBayar = oldData.HutangPiutang.Total - oldData.HutangPiutang.Sisa

		newTotalHarga := oldData.TotalHarga

		if newTotalHarga < oldTotalBayar {
			return nil, errors.New(message.TotalBayarMustGeOrEqToTotalByr)
		}

		oldData.HutangPiutang.Sisa = newTotalHarga - oldTotalBayar
		oldData.HutangPiutang.Total = newTotalHarga

		if oldData.HutangPiutang.Sisa <= 0 {
			oldData.HutangPiutang.Status = "LUNAS"
		}

		// update tr
		oldData.HutangPiutang.Transaksi.Total = newTotalHarga
		for i, ay := range oldData.HutangPiutang.Transaksi.AyatJurnals {
			// update debit kredit
			if ay.Debit != 0 {
				ay.Debit = newTotalHarga
			} else {
				ay.Kredit = newTotalHarga
			}
			// update saldo
			if ay.Saldo > 0 {
				ay.Saldo = newTotalHarga
			} else {
				ay.Saldo = -newTotalHarga
			}
			oldData.HutangPiutang.Transaksi.AyatJurnals[i] = ay
		}
	}

	// update transaksi hutang piutang
	oldData.HutangPiutang.Transaksi.Keterangan = param.Req.Keterangan
	return oldData, nil
}

func (u *invoiceUsecase) SaveCommitDB(param ParamCommitDB) error {

	if param.Invoice != nil && param.Invoice.KontakID != "" {
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

	return u.repo.UpdateFullAssoc(repo.ParamUpdateFullAssoc(param))
}

func (u *invoiceUsecase) Delete(param ParamDelete) error {
	data, err := u.repo.GetByIdFullAssoc(repo.ParamGetById(param))
	if err != nil {
		return err
	}

	if err := u.repo.Delete(repo.ParamDelete{Ctx: param.Ctx, Invoice: data}); err != nil {
		return err
	}
	return nil
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
		Ctx:            param.Ctx,
		StatusProduksi: param.Req.StatusProduksi,
		KontakID:       param.Req.KontakID,
		Limit:          param.Req.Limit,
		Next:           param.Req.Next,
	}

	if !tglDeadline.IsZero() {
		paramRepo.TanggalDeadline = tglDeadline
	}
	if !tglKirim.IsZero() {
		paramRepo.TanggalKirim = tglKirim
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
	return u.repo.GetById(repo.ParamGetById(param))
}

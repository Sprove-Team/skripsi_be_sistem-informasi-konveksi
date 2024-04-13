package uc_akuntansi_transaksi

import (
	"context"
	"errors"
	"sync"
	"time"

	repoAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	repoDataBayarHutangPiutang "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/data_bayar_hutang_piutang"
	repoHutangPiutang "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/hutang_piutang"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/transaksi"
	repo_invoice_data_bayar "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm/data_bayar"
	pkgAkuntansiLogic "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/transaksi"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type TransaksiUsecase interface {
	Create(ctx context.Context, reqTransaksi req.Create) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, reqTransaksi req.Update) error
	GetAll(ctx context.Context, reqTransaksi req.GetAll) ([]entity.Transaksi, error)
	GetHistory(ctx context.Context, reqTransaksi req.GetHistory) ([]entity.Transaksi, error)
	GetById(ctx context.Context, id string) (entity.Transaksi, error)
}

type transaksiUsecase struct {
	repo                       repo.TransaksiRepo
	repoAkun                   repoAkun.AkunRepo
	repoHutangPiutang          repoHutangPiutang.HutangPiutangRepo
	repoDataBayarHutangPiutang repoDataBayarHutangPiutang.DataBayarHutangPiutangRepo
	repoDataBayarInvoice       repo_invoice_data_bayar.DataBayarInvoiceRepo
	ulid                       pkg.UlidPkg
}

func NewTransaksiUsecase(
	repo repo.TransaksiRepo,
	repoAkun repoAkun.AkunRepo,
	repoHutangPiutang repoHutangPiutang.HutangPiutangRepo,
	repoDataBayarHutangPiutang repoDataBayarHutangPiutang.DataBayarHutangPiutangRepo,
	repoDataBayarInvoice repo_invoice_data_bayar.DataBayarInvoiceRepo,
	ulid pkg.UlidPkg) TransaksiUsecase {
	return &transaksiUsecase{repo, repoAkun, repoHutangPiutang, repoDataBayarHutangPiutang, repoDataBayarInvoice, ulid}
}

func (u *transaksiUsecase) Delete(ctx context.Context, id string) error {
	hp, err := u.repoHutangPiutang.GetByTrId(repoHutangPiutang.ParamGetByTrId{
		Ctx: ctx,
		ID:  id,
	})
	if err != nil {
		if err.Error() != "record not found" {
			return err
		}
	}

	dataByarInvoice, err := u.repoDataBayarInvoice.GetByInvoiceID(repo_invoice_data_bayar.ParamGetByInvoiceID{
		Ctx:       ctx,
		Status:    "BELUM_TERKONFIRMASI",
		InvoiceID: hp.InvoiceID,
	})
	if err != nil {
		if err.Error() != "record not found" {
			return err
		}
	}

	if len(dataByarInvoice) > 0 {
		return errors.New(message.CantDeleteTrIfDataByrBlmTerkonfirmasiExist)
	}

	return u.repo.Delete(ctx, id)
}

func (u *transaksiUsecase) Update(ctx context.Context, reqTransaksi req.Update) error {

	hp, err := u.repoHutangPiutang.GetByTrId(repoHutangPiutang.ParamGetByTrId{
		Ctx: ctx,
		ID:  reqTransaksi.ID,
	})

	if err != nil {
		if err.Error() != "record not found" {
			return err
		}
	} else {
		if hp.ID != "" && len(reqTransaksi.AyatJurnal) != 2 {
			return errors.New(message.AkunHutangPiutangNotEq2)
		}
	}

	// validasi transaksi jika termasuk hutang piutang
	var hpFromByr entity.DataBayarHutangPiutang
	if hp.ID == "" {
		hpFromByr, err = u.repoDataBayarHutangPiutang.GetByTrId(ctx, reqTransaksi.ID)
		if err != nil && err.Error() != "record not found" {
			return err
		}
	}

	lengthReqAyatJurnals := len(reqTransaksi.AyatJurnal)

	var wg sync.WaitGroup
	wg.Add(2)

	var reqAy []entity.AyatJurnal
	var oldTr entity.Transaksi
	var oldAy []entity.AyatJurnal

	go func() {
		defer wg.Done()
		var reqAyData = make([]entity.AyatJurnal, lengthReqAyatJurnals)
		for i, ay := range reqTransaksi.AyatJurnal {
			reqAyData[i] = entity.AyatJurnal{
				AkunID: ay.AkunID,
				Debit:  ay.Debit,
				Kredit: ay.Kredit,
			}
		}
		reqAy = reqAyData
	}()

	// Fetch oldTr
	go func() {
		defer wg.Done()
		oldTr, err = u.repo.GetById(ctx, reqTransaksi.ID)
		if err != nil {
			return
		}
		oldAyData := make([]entity.AyatJurnal, len(oldTr.AyatJurnals))
		for i, ay := range oldTr.AyatJurnals {
			oldAyData[i] = entity.AyatJurnal{
				AkunID: ay.AkunID,
				Debit:  ay.Debit,
				Kredit: ay.Kredit,
			}
		}
		oldAy = oldAyData
	}()
	wg.Wait()

	// Handle errors if any
	if err != nil {
		return err
	}

	// setup param tr update
	repoParam := repo.UpdateParam{
		UpdateTr: &entity.Transaksi{
			BaseSoftDelete: entity.BaseSoftDelete{
				ID: reqTransaksi.ID,
			},
			Keterangan:      reqTransaksi.Keterangan,
			BuktiPembayaran: reqTransaksi.BuktiPembayaran,
		},
	}

	// update tr tanggal
	if reqTransaksi.Tanggal != "" {
		tanggalTr, err := time.Parse(time.RFC3339, reqTransaksi.Tanggal)
		if err != nil {
			return err
		}
		repoParam.UpdateTr.Tanggal = tanggalTr
	}

	// change the saldo ayat jurnal if the ayat jurnals is difference
	isSameAy := pkgAkuntansiLogic.IsSameReqAyJurnals(oldAy, reqAy)
	if !isSameAy {

		// u.repo.GetById()
		akunIDs := make([]string, lengthReqAyatJurnals)
		for i, ay := range reqTransaksi.AyatJurnal {
			akunIDs[i] = ay.AkunID
		}

		akuns, err := u.repoAkun.GetByIds(ctx, akunIDs)
		if err != nil {
			return err
		}

		if len(akuns) != helper.CountUniqueElements(akunIDs) {
			return errors.New(message.AkunNotFound)
		}

		akunsMap := map[string]entity.Akun{}

		if hp.ID != "" || hpFromByr.ID != "" {
			// validasi ayat jurnal jika termasuk hutang piutang
			if err := pkgAkuntansiLogic.IsAkunHutangPiutangExist(akuns); err != nil {
				return err
			}
		}
		for _, v := range akuns {
			akunsMap[v.ID] = v
		}

		for _, v := range oldTr.AyatJurnals {
			// setup data akun, if there a same id with the old one
			if _, ok := akunsMap[v.AkunID]; !ok {
				akunsMap[v.AkunID] = entity.Akun{
					Base: entity.Base{
						ID: v.Akun.ID,
					},
					Nama:           v.Akun.Nama,
					Kode:           v.Akun.Kode,
					KelompokAkunID: v.Akun.KelompokAkunID,
				}
			}
		}

		// validate kredit, debit and get totalTransaksi
		var totalTransaksi float64

		if err := pkgAkuntansiLogic.IsKreditEqualToDebit(reqAy, &totalTransaksi); err != nil {
			return err
		}
		if err := pkgAkuntansiLogic.IsDuplicateAkun(reqAy); err != nil {
			return err
		}

		// validate total bayar with total tr hp
		if hp.ID != "" {
			totalByr := hp.Total - hp.Sisa
			if totalTransaksi < totalByr {
				return errors.New(message.TotalHPMustGeOrEqToTotalByr)
			}
		}
		// validate total bayar dengan sisa tagihan
		if hpFromByr.ID != "" {
			currentSisa := hpFromByr.HutangPiutang.Sisa + hpFromByr.Total
			if totalTransaksi > currentSisa {
				return errors.New(message.BayarMustLessThanSisaTagihan)
			}
		}

		// set new total for update tr
		repoParam.UpdateTr.Total = totalTransaksi

		// add newAyatJurnals
		newAyatJurnals := make([]*entity.AyatJurnal, len(reqTransaksi.AyatJurnal))

		for i, ay := range reqTransaksi.AyatJurnal {
			akun := akunsMap[ay.AkunID]
			if err := pkgAkuntansiLogic.IsValidAkunHp(hp, &akun, &ay); err != nil {
				return err
			}
			if err := pkgAkuntansiLogic.IsValidAkunByrHP(&hpFromByr, &akun, &ay); err != nil {
				return err
			}

			// logic calculate saldo ayatJurnal
			var saldo = pkgAkuntansiLogic.UpdateSaldo(ay.Kredit, ay.Debit, akun.SaldoNormal)
			ayatJurnal := entity.AyatJurnal{
				Base: entity.Base{
					ID: u.ulid.MakeUlid().String(),
				},
				TransaksiID: reqTransaksi.ID,
				AkunID:      ay.AkunID,
				Kredit:      ay.Kredit,
				Debit:       ay.Debit,
				Saldo:       saldo,
			}

			newAyatJurnals[i] = &ayatJurnal
			akunsMap[ay.AkunID] = akun
		}

		repoParam.NewAyatJurnals = newAyatJurnals

		if hp.Total != totalTransaksi {
			// set data hutang piutang
			if hp.ID != "" {
				repoParam.UpdateHutangPiutang = hp
				repoParam.UpdateHutangPiutang.Total = repoParam.UpdateTr.Total
				// total baru - total dibayar didapat dari old total - old sisa
				repoParam.UpdateHutangPiutang.Sisa = repoParam.UpdateTr.Total - (oldTr.Total - hp.Sisa)

			}

			if hpFromByr.ID != "" {
				// update data hp
				repoParam.UpdateHutangPiutang = &hpFromByr.HutangPiutang
				// logic update sisa hp, dengan logika (sisa_old + total_byr old) - total_byr skrng
				repoParam.UpdateHutangPiutang.Sisa = (hpFromByr.HutangPiutang.Sisa + hpFromByr.Total) - repoParam.UpdateTr.Total
				// update data bayar
				repoParam.UpdateDataBayarHutangPiutang = &hpFromByr
				repoParam.UpdateDataBayarHutangPiutang.Total = repoParam.UpdateTr.Total
			}

			if repoParam.UpdateHutangPiutang != nil {
				if repoParam.UpdateHutangPiutang.Sisa <= 0 && repoParam.UpdateHutangPiutang.Status != "LUNAS" {
					repoParam.UpdateHutangPiutang.Status = "LUNAS"
				} else if repoParam.UpdateHutangPiutang.Status != "BELUM_LUNAS" {
					repoParam.UpdateHutangPiutang.Status = "BELUM_LUNAS"
				}
			}
		}

	}

	if err := u.repo.Update(ctx, repoParam); err != nil {
		return err
	}

	return nil
}

func (u *transaksiUsecase) Create(ctx context.Context, reqTransaksi req.Create) error {

	reqAy := make([]entity.AyatJurnal, len(reqTransaksi.AyatJurnal))

	for i, ay := range reqTransaksi.AyatJurnal {
		reqAy[i] = entity.AyatJurnal{
			AkunID: ay.AkunID,
			Debit:  ay.Debit,
			Kredit: ay.Kredit,
		}
	}

	// get total transaksi
	var totalTransaksi float64
	if err := pkgAkuntansiLogic.IsKreditEqualToDebit(reqAy, &totalTransaksi); err != nil {
		return err
	}
	if err := pkgAkuntansiLogic.IsDuplicateAkun(reqAy); err != nil {
		return err
	}

	// get akun that haved by ayatjurnal
	akunIds := make([]string, len(reqTransaksi.AyatJurnal))

	for i, v := range reqTransaksi.AyatJurnal {
		akunIds[i] = v.AkunID
	}

	akuns, err := u.repoAkun.GetByIds(ctx, akunIds)
	if err != nil {
		if len(akuns) != len(akunIds) {
			return errors.New(message.AkunNotFound)
		}
		return err
	}

	// map to make essier to get by id without call db inside looping
	mapAkuns := make(map[string]entity.Akun, len(akuns))

	for _, akun := range akuns {
		mapAkuns[akun.ID] = akun
	}

	dataAyatJurnals := make([]*entity.AyatJurnal, len(reqTransaksi.AyatJurnal))
	transaksiID := u.ulid.MakeUlid().String()

	for i, v := range reqTransaksi.AyatJurnal {
		akun, ok := mapAkuns[v.AkunID]
		if !ok {
			return errors.New(message.AkunNotFound)
		}

		// logic calculate saldo ayatJurnal
		var saldo = pkgAkuntansiLogic.UpdateSaldo(v.Kredit, v.Debit, akun.SaldoNormal)

		ayatJurnal := entity.AyatJurnal{
			Base: entity.Base{
				ID: u.ulid.MakeUlid().String(),
			},
			TransaksiID: transaksiID,
			AkunID:      v.AkunID,
			Kredit:      v.Kredit,
			Debit:       v.Debit,
			Saldo:       saldo,
		}
		dataAyatJurnals[i] = &ayatJurnal

	}

	parsedTime, err := time.Parse(time.RFC3339, reqTransaksi.Tanggal)
	if err != nil {
		helper.LogsError(err)
		return err
	}

	dataTransaksi := entity.Transaksi{
		BaseSoftDelete: entity.BaseSoftDelete{
			ID: transaksiID,
		},
		Tanggal:         parsedTime,
		Keterangan:      reqTransaksi.Keterangan,
		KontakID:        reqTransaksi.KontakID,
		BuktiPembayaran: reqTransaksi.BuktiPembayaran,
		Total:           totalTransaksi,
	}

	err = u.repo.Create(ctx, repo.CreateParam{
		AyatJurnals: dataAyatJurnals,
		Transaksi:   &dataTransaksi,
	})

	if err != nil {
		helper.LogsError(err)
		return err
	}

	return nil
}

func (u *transaksiUsecase) GetById(ctx context.Context, id string) (entity.Transaksi, error) {
	data, err := u.repo.GetById(ctx, id)
	if err != nil {
		return entity.Transaksi{}, err
	}
	return data, nil
}

func (u *transaksiUsecase) GetAll(ctx context.Context, reqTransaksi req.GetAll) ([]entity.Transaksi, error) {
	timeZone, _ := time.LoadLocation(reqTransaksi.TimeZone)
	endDate, err := time.ParseInLocation(time.DateOnly, reqTransaksi.EndDate, timeZone)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	startDate, err := time.ParseInLocation(time.DateOnly, reqTransaksi.StartDate, timeZone)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	
	searchFilter := repo.SearchTransaksi{
		EndDate:   endDate,
		StartDate: startDate,
		TimeZone:  reqTransaksi.TimeZone,
	}

	dataTransaksi, err := u.repo.GetAll(ctx, searchFilter)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}

	return dataTransaksi, nil
}

func (u *transaksiUsecase) GetHistory(ctx context.Context, reqTransaksi req.GetHistory) ([]entity.Transaksi, error) {
	timeZone, _ := time.LoadLocation(reqTransaksi.TimeZone)
	endDate, err := time.ParseInLocation(time.DateOnly, reqTransaksi.EndDate, timeZone)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	startDate, err := time.ParseInLocation(time.DateOnly, reqTransaksi.StartDate, timeZone)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	searchFilter := repo.SearchTransaksi{
		EndDate:   endDate,
		StartDate: startDate,
		TimeZone:  reqTransaksi.TimeZone,
	}

	dataHistory, err := u.repo.GetHistory(ctx, searchFilter)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}

	return dataHistory, nil
}

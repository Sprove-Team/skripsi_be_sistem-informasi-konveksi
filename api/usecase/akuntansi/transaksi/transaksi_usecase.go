package akuntansi

import (
	"context"
	"errors"
	"time"

	repoAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	repoDataBayarHutangPiutang "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/data_bayar_hutang_piutang"
	repoHutangPiutang "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/hutang_piutang"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/transaksi"
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
	ulid                       pkg.UlidPkg
}

func NewTransaksiUsecase(
	repo repo.TransaksiRepo,
	repoAkun repoAkun.AkunRepo,
	repoHutangPiutang repoHutangPiutang.HutangPiutangRepo,
	repoDataBayarHutangPiutang repoDataBayarHutangPiutang.DataBayarHutangPiutangRepo,
	ulid pkg.UlidPkg) TransaksiUsecase {
	return &transaksiUsecase{repo, repoAkun, repoHutangPiutang, repoDataBayarHutangPiutang, ulid}
}

func (u *transaksiUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *transaksiUsecase) Update(ctx context.Context, reqTransaksi req.Update) error {
	// g := &errgroup.Group{}
	// g.SetLimit(10)

	hp, err := u.repoHutangPiutang.GetByTrId(ctx, reqTransaksi.ID)
	if err != nil {
		if err.Error() != "record not found" {
			return err
		}
	}

	var hpFromByr entity.DataBayarHutangPiutang
	// validasi transaksi jika termasuk hutang piutang
	if hp.ID != "" {
		if len(reqTransaksi.AyatJurnal) != 2 {
			return errors.New(message.AkunHutangPiutangNotEq2)
		}
		// validasi transaksi jika termasuk bayar hutang piutang
	} else {
		hpFromByr, err = u.repoDataBayarHutangPiutang.GetByTrId(ctx, reqTransaksi.ID)
		if err != nil {
			if err.Error() != "record not found" {
				return err
			}
		}

	}

	lengthReqAyatJurnals := len(reqTransaksi.AyatJurnal)
	reqAy := make([]entity.AyatJurnal, lengthReqAyatJurnals)

	for i, ay := range reqTransaksi.AyatJurnal {
		reqAy[i] = entity.AyatJurnal{
			AkunID: ay.AkunID,
			Debit:  ay.Debit,
			Kredit: ay.Kredit,
		}
	}

	oldTr, err := u.repo.GetById(ctx, reqTransaksi.ID)
	if err != nil {
		return err
	}

	lengthOldAyatJurnals := len(oldTr.AyatJurnals)

	// setup param tr update
	repoParam := repo.UpdateParam{
		UpdateTr: &entity.Transaksi{
			Base: entity.Base{
				ID: reqTransaksi.ID,
			},
			Keterangan:      reqTransaksi.Keterangan,
			BuktiPembayaran: reqTransaksi.BuktiPembayaran,
		},
	}

	oldAy := make([]entity.AyatJurnal, lengthOldAyatJurnals)
	for i, ay := range oldTr.AyatJurnals {
		oldAy[i] = entity.AyatJurnal{
			AkunID: ay.AkunID,
			Debit:  ay.Debit,
			Kredit: ay.Kredit,
		}
	}

	// change the saldo ayat jurnal if the ayat jurnals is difference
	if !pkgAkuntansiLogic.IsSameReqAyJurnals(oldAy, reqAy) {

		// u.repo.GetById()
		akunIDs := make([]string, lengthReqAyatJurnals)
		for i, ay := range reqTransaksi.AyatJurnal {
			akunIDs[i] = ay.AkunID
		}

		akuns, err := u.repoAkun.GetByIds(ctx, akunIDs)
		if err != nil {
			return err
		}

		akunsMap := map[string]entity.Akun{}

		for _, v := range akuns {
			// validasi ayat jurnal jika termasuk hutang piutang
			if hp.ID != "" || hpFromByr.ID != "" {
				if !pkgAkuntansiLogic.IsValidAkunHutangPiutang(v.KelompokAkun.Nama) {
					return errors.New(message.InvalidAkunHutangPiutang)
				}
			}

			akunsMap[v.ID] = v
		}

		for _, v := range oldTr.AyatJurnals {
			// setup data akun, if there a same id with the old one
			akun, ok := akunsMap[v.AkunID]
			if !ok {
				akun = entity.Akun{
					Base: entity.Base{
						ID: v.Akun.ID,
					},
					Nama:           v.Akun.Nama,
					Kode:           v.Akun.Kode,
					KelompokAkunID: v.Akun.KelompokAkunID,
				}
			}
			akunsMap[v.AkunID] = akun
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

		repoParam.UpdateTr.Total = totalTransaksi

		// add newAyatJurnals
		newAyatJurnals := make([]*entity.AyatJurnal, len(reqTransaksi.AyatJurnal))

		for i, ay := range reqTransaksi.AyatJurnal {
			akun := akunsMap[ay.AkunID]

			if err := isValidAkunHp(&hp, &akun, &ay); err != nil {
				return err
			}
			if err := isValidAkunByrHP(&hpFromByr, &akun, &ay); err != nil {
				return err
			}

			// logic calculate saldo ayatJurnal
			var saldo float64
			pkgAkuntansiLogic.UpdateSaldo(&saldo, ay.Kredit, ay.Debit, akun.SaldoNormal)
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

		if repoParam.UpdateTr.Total != oldTr.Total {
			// set data hutang piutang
			if hp.ID != "" {
				repoParam.UpdateHutangPiutang = &hp
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
		var saldo float64
		pkgAkuntansiLogic.UpdateSaldo(&saldo, v.Kredit, v.Debit, akun.SaldoNormal)

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
		Base: entity.Base{
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
	endDate, err := time.Parse(time.DateOnly, reqTransaksi.EndDate)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	startDate, err := time.Parse(time.DateOnly, reqTransaksi.StartDate)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}

	searchFilter := repo.SearchTransaksi{
		EndDate:   endDate,
		StartDate: startDate,
	}

	dataTransaksi, err := u.repo.GetAll(ctx, searchFilter)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}

	return dataTransaksi, nil
}

func (u *transaksiUsecase) GetHistory(ctx context.Context, reqTransaksi req.GetHistory) ([]entity.Transaksi, error) {
	endDate, err := time.Parse(time.DateOnly, reqTransaksi.EndDate)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	startDate, err := time.Parse(time.DateOnly, reqTransaksi.StartDate)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}

	searchFilter := repo.SearchTransaksi{
		EndDate:   endDate,
		StartDate: startDate,
	}

	dataHistory, err := u.repo.GetHistory(ctx, searchFilter)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}

	return dataHistory, nil
}

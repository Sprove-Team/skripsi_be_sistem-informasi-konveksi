package akuntansi

import (
	"context"
	"errors"
	"sync"
	"time"

	repoAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/transaksi"
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
	repo     repo.TransaksiRepo
	repoAkun repoAkun.AkunRepo
	ulid     pkg.UlidPkg
}

func NewTransaksiUsecase(repo repo.TransaksiRepo, repoAkun repoAkun.AkunRepo, ulid pkg.UlidPkg) TransaksiUsecase {
	return &transaksiUsecase{repo, repoAkun, ulid}
}

func (u *transaksiUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *transaksiUsecase) Update(ctx context.Context, reqTransaksi req.Update) error {
	oldTr, err := u.repo.GetById(ctx, reqTransaksi.ID)
	if err != nil {
		return err
	}

	lengthReqAyatJurnals := len(reqTransaksi.AyatJurnal)
	lengthOldAyatJurnals := len(oldTr.AyatJurnals)

	updateTrParam := repo.UpdateParam{}

	newTr := entity.Transaksi{
		Base: entity.Base{
			ID: reqTransaksi.ID,
		},
		Keterangan:      reqTransaksi.Keterangan,
		BuktiPembayaran: reqTransaksi.BuktiPembayaran,
	}

	oldAy := make([]entity.AyatJurnal, lengthOldAyatJurnals)
	reqAy := make([]entity.AyatJurnal, lengthReqAyatJurnals)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i, ay := range oldTr.AyatJurnals {
			oldAy[i] = entity.AyatJurnal{
				AkunID: ay.AkunID,
				Debit:  ay.Debit,
				Kredit: ay.Kredit,
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i, ay := range reqTransaksi.AyatJurnal {
			reqAy[i] = entity.AyatJurnal{
				AkunID: ay.AkunID,
				Debit:  ay.Debit,
				Kredit: ay.Kredit,
			}
		}
	}()
	wg.Wait()

	// change the saldo ayat jurnal if the ayat jurnals is difference
	if !helper.IsSameReqAyJurnals(oldAy, reqAy) {

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

		if err := helper.IsKreditEqualToDebit(reqAy, &totalTransaksi); err != nil {
			return err
		}
		if err := helper.IsDuplicateAkun(reqAy); err != nil {
			return err
		}

		newTr.Total = totalTransaksi

		// add newAyatJurnals
		newAyatJurnals := make([]*entity.AyatJurnal, len(reqTransaksi.AyatJurnal))

		for i, ay := range reqTransaksi.AyatJurnal {
			akun := akunsMap[ay.AkunID]
			// logic calculate saldo ayatJurnal
			var saldo float64
			helper.UpdateSaldo(&saldo, ay.Kredit, ay.Debit, akun.SaldoNormal)
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

		updateTrParam.NewAyatJurnals = newAyatJurnals
	}
	updateTrParam.UpdateTr = &newTr

	if err := u.repo.Update(ctx, updateTrParam); err != nil {
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
	if err := helper.IsKreditEqualToDebit(reqAy, &totalTransaksi); err != nil {
		return err
	}
	if err := helper.IsDuplicateAkun(reqAy); err != nil {
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
		helper.UpdateSaldo(&saldo, v.Kredit, v.Debit, akun.SaldoNormal)

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

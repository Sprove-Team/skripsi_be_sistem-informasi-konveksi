package akuntansi

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	repoAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/transaksi"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/transaksi"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"golang.org/x/sync/errgroup"
)

type TransaksiUsecase interface {
	Create(ctx context.Context, reqTransaksi req.Create) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, reqTransaksi req.Update) error
	GetAll(ctx context.Context, reqTransaksi req.GetAll) ([]entity.Transaksi, error)
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
	detailAkuns, err := u.repoAkun.GetAkunDetailsByTransactionID(ctx, id)
	if err != nil {
		return err
	}

	saldoAkunValues := make([]string, len(detailAkuns))

	for i, v := range detailAkuns {
		newUp := entity.Akun{
			ID:             v.ID,
			GolonganAkunID: v.GolID,
			Saldo:          v.Saldo,
			Nama:           v.Nama,
			Kode:           v.Kode,
		}

		newUp.Saldo -= v.TotalSaldoTr

		saldoAkunValues[i] = fmt.Sprintf("('%s', %.2f, '%s', '%s', '%s')", newUp.ID, newUp.Saldo, newUp.GolonganAkunID, newUp.Nama, newUp.Kode)
	}

	return u.repo.Delete(ctx, repo.DeleteParam{
		ID:              id,
		SaldoAkunValues: saldoAkunValues,
	})
}

func (u *transaksiUsecase) Update(ctx context.Context, reqTransaksi req.Update) error {
	oldTr, err := u.repo.GetById(ctx, reqTransaksi.ID)
	if err != nil {
		return err
	}

	lengthReqAyatJurnals := len(reqTransaksi.AyatJurnals)
	lengthOldAyatJurnals := len(oldTr.AyatJurnals)

	// isReqAyatJurnalsEmpty := lengthReqAyatJurnals == 0

	updateTrParam := repo.UpdateParam{}

	newTr := entity.Transaksi{
		ID:              reqTransaksi.ID,
		Keterangan:      reqTransaksi.Keterangan,
		BuktiPembayaran: reqTransaksi.BuktiPembayaran,
	}

	oldAy := make([]req.ReqAyatJurnal, lengthOldAyatJurnals)

	// isLengthAreSame := lengthReqAyatJurnals == lengthOldAyatJurnals
	// countSameAyJurnals := 0

	for i, ay := range oldTr.AyatJurnals {
		oldAy[i] = req.ReqAyatJurnal{
			AkunID: ay.AkunID,
			Debit:  ay.Debit,
			Kredit: ay.Kredit,
		}
	}

	if !isSameReqAyJurnals(oldAy, reqTransaksi.AyatJurnals) {

		akunIDs := make([]string, lengthReqAyatJurnals)
		for i, ay := range reqTransaksi.AyatJurnals {
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

		// data oldData map
		// oldSaldoAkunMap := map[string]entity.Akun{}
		//
		// for _, v := range oldTr.AyatJurnals {
		// 	oldSaldoAkunMap[v.AkunID] = v.Akun
		// }

		for _, v := range oldTr.AyatJurnals {
			// upOldSaldoAkunMap := oldSaldoAkunMap[v.AkunID]
			// upOldSaldoAkunMap.Saldo -= v.Saldo
			// oldSaldoAkunMap[v.AkunID] = upOldSaldoAkunMap
			// update if there a same id with the old one
			akun, ok := akunsMap[v.AkunID]
			if !ok {
				akun = entity.Akun{
					ID:             v.Akun.ID,
					Saldo:          v.Akun.Saldo,
					Nama:           v.Akun.Nama,
					Kode:           v.Akun.Kode,
					GolonganAkunID: v.Akun.GolonganAkunID,
				}
			}

			akun.Saldo -= v.Saldo
			akunsMap[v.AkunID] = akun
		}

		wg := &sync.WaitGroup{}
		m := &sync.Mutex{}

		// validate kredit, debit and get totalTransaksi
		var totalTransaksi float64

		if err := isKreditEqualToDebit(reqTransaksi.AyatJurnals, &totalTransaksi, wg, m); err != nil {
			return err
		}
		if err := isDuplicateAkun(reqTransaksi.AyatJurnals); err != nil {
			return err
		}

		newTr.Total = totalTransaksi

		// add newAyatJurnals
		newAyatJurnals := make([]*entity.AyatJurnal, len(reqTransaksi.AyatJurnals))

		for i, ay := range reqTransaksi.AyatJurnals {
			wg.Add(1)
			go func(i int, ay req.ReqAyatJurnal) {
				defer wg.Done()
				akun := akunsMap[ay.AkunID]
				// logic calculate saldo ayatJurnal
				var saldo float64
				if ay.Kredit != 0 {
					updateSaldo(&akun.Saldo, &saldo, ay.Kredit, akun.SaldoNormal == "DEBIT")
				}
				if ay.Debit != 0 {
					updateSaldo(&akun.Saldo, &saldo, ay.Debit, akun.SaldoNormal == "KREDIT")
				}

				ayatJurnal := entity.AyatJurnal{
					ID:          u.ulid.MakeUlid().String(),
					TransaksiID: reqTransaksi.ID,
					AkunID:      ay.AkunID,
					Kredit:      ay.Kredit,
					Debit:       ay.Debit,
					Saldo:       saldo,
				}
				m.Lock()
				newAyatJurnals[i] = &ayatJurnal
				akunsMap[ay.AkunID] = akun
				m.Unlock()
			}(i, ay)
		}
		wg.Wait()

		newSaldoAkun := make([]string, len(akunsMap))

		i := 0
		for _, v := range akunsMap {
			newSaldoAkun[i] = fmt.Sprintf("('%s', %.2f, '%s', '%s', '%s')", v.ID, v.Saldo, v.GolonganAkunID, v.Nama, v.Kode)
			i++
		}

		updateTrParam.NewAyatJurnals = newAyatJurnals
		updateTrParam.UpdateTr = &newTr
		updateTrParam.NewSaldoAkunValues = newSaldoAkun
	}

	if err := u.repo.Update(ctx, updateTrParam); err != nil {
		return err
	}

	return nil
}

func (u *transaksiUsecase) Create(ctx context.Context, reqTransaksi req.Create) error {
	// validate kredit and debit
	wg := &sync.WaitGroup{}
	m := &sync.Mutex{}

	var totalTransaksi float64
	if err := isKreditEqualToDebit(reqTransaksi.AyatJurnals, &totalTransaksi, wg, m); err != nil {
		return err
	}
	if err := isDuplicateAkun(reqTransaksi.AyatJurnals); err != nil {
		return err
	}
	// akun, err := u.repoAkun.GetById(ctx, )

	akunIds := make([]string, len(reqTransaksi.AyatJurnals))

	for i, v := range reqTransaksi.AyatJurnals {
		wg.Add(1)
		go func(i int, v req.ReqAyatJurnal) {
			defer wg.Done()
			akunIds[i] = v.AkunID
		}(i, v)
	}
	wg.Wait()

	akuns, err := u.repoAkun.GetByIds(ctx, akunIds)
	if err != nil {
		return err
	}

	// map to make essier to get by id without call db inside looping
	mapAkuns := make(map[string]entity.Akun, len(akuns))

	for _, akun := range akuns {
		mapAkuns[akun.ID] = akun
	}

	dataAyatJurnals := make([]*entity.AyatJurnal, len(reqTransaksi.AyatJurnals))
	transaksiID := u.ulid.MakeUlid().String()

	updateAkuns := []*entity.Akun{}
	g := errgroup.Group{}

	for i, v := range reqTransaksi.AyatJurnals {
		i := i
		v := v
		g.Go(func() error {
			akun, ok := mapAkuns[v.AkunID]
			if !ok {
				return errors.New(message.AkunIdNotFound)
			}
			// logic calculate saldo ayatJurnal
			var saldo float64
			if v.Kredit != 0 {
				updateSaldo(&akun.Saldo, &saldo, v.Kredit, akun.SaldoNormal == "DEBIT")
			}
			if v.Debit != 0 {
				updateSaldo(&akun.Saldo, &saldo, v.Debit, akun.SaldoNormal == "KREDIT")
			}
			// mapAkuns[v.AkunID] = repoAkun
			updateAkuns = append(updateAkuns, &akun)
			ayatJurnal := entity.AyatJurnal{
				ID:          u.ulid.MakeUlid().String(),
				TransaksiID: transaksiID,
				AkunID:      v.AkunID,
				Kredit:      v.Kredit,
				Debit:       v.Debit,
				Saldo:       saldo,
			}
			m.Lock()
			dataAyatJurnals[i] = &ayatJurnal
			m.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339, reqTransaksi.Tanggal)
	// log.Println(parsedTime)
	if err != nil {
		helper.LogsError(err)
		return err
	}

	dataTransaksi := entity.Transaksi{
		ID:              transaksiID,
		Tanggal:         parsedTime,
		Keterangan:      reqTransaksi.Keterangan,
		BuktiPembayaran: reqTransaksi.BuktiPembayaran,
		Total:           totalTransaksi,
	}

	err = u.repo.Create(ctx, repo.CreateParam{
		AyatJurnals: dataAyatJurnals,
		Transaksi:   &dataTransaksi,
		UpdateAkuns: updateAkuns,
	})

	if err != nil {
		helper.LogsError(err)
		return err
	}

	return nil
}

func (u *transaksiUsecase) GetById(ctx context.Context, id string) (entity.Transaksi, error) {
	detailTr, err := u.repo.GetById(ctx, id)
	if err != nil {
		return entity.Transaksi{}, err
	}
	return detailTr, nil
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

	// datasRes := make([]res.DataGetAllTransaksi, len(dataTransaksi))

	// wg := &sync.WaitGroup{}
	// for i, tr := range dataTransaksi {
	// 	dataAyatJurnals := make([]res.DataAyatJurnalTR, len(tr.AyatJurnals))
	// 	dataTransaksi := res.DataGetAllTransaksi{}
	// 	wg.Add(1)
	// 	go func(ays []entity.AyatJurnal) {
	// 		defer wg.Done()
	// 		for j, ay := range ays {
	// 			dataAyatJurnals[j] = res.DataAyatJurnalTR{
	// 				AkunID:       ay.AkunID,
	// 				AyatJurnalID: ay.ID,
	// 				Kredit:       ay.Kredit,
	// 				Debit:        ay.Debit,
	// 				Saldo:        ay.Saldo,
	// 			}
	// 		}
	// 	}(tr.AyatJurnals)
	// 	wg.Wait()
	// 	dataTransaksi.ID = tr.ID
	// 	dataTransaksi.Tanggal = tr.Tanggal.Format(time.DateTime)
	// 	dataTransaksi.Keterangan = tr.Keterangan
	// 	dataTransaksi.BuktiPembayaran = tr.BuktiPembayaran
	// 	dataTransaksi.TotalDebit = tr.Total
	// 	dataTransaksi.TotalKredit = tr.Total
	// 	dataTransaksi.AyatJurnals = dataAyatJurnals
	// 	datasRes[i] = dataTransaksi
	// }

	return dataTransaksi, nil
}

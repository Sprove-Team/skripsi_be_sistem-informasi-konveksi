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
	"golang.org/x/sync/errgroup"
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
	detailAkuns, err := u.repoAkun.GetAkunDetailsByTransactionID(ctx, id)
	if err != nil {
		return err
	}

	newAkuns := make([]entity.Akun, len(detailAkuns))

	for i, v := range detailAkuns {
		newAkuns[i] = entity.Akun{
			Base: entity.Base{
				ID: v.ID,
			},
			KelompokAkunID: v.KelID,
			Saldo:          v.Saldo - v.TotalSaldoTr,
			Nama:           v.Nama,
			Kode:           v.Kode,
		}
	}

	return u.repo.Delete(ctx, repo.DeleteParam{
		ID:          id,
		UpdateAkuns: newAkuns,
	})
}

func (u *transaksiUsecase) Update(ctx context.Context, reqTransaksi req.Update) error {
	oldTr, err := u.repo.GetById(ctx, reqTransaksi.ID)
	if err != nil {
		return err
	}

	lengthReqAyatJurnals := len(reqTransaksi.AyatJurnal)
	lengthOldAyatJurnals := len(oldTr.AyatJurnals)

	// isReqAyatJurnalsEmpty := lengthReqAyatJurnals == 0

	updateTrParam := repo.UpdateParam{}

	newTr := entity.Transaksi{
		Base: entity.Base{
			ID: reqTransaksi.ID,
		},
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

	// change the saldo akun if the ayat jurnals is difference
	if !isSameReqAyJurnals(oldAy, reqTransaksi.AyatJurnal) {

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
			// update if there a same id with the old one
			akun, ok := akunsMap[v.AkunID]
			if !ok {
				akun = entity.Akun{
					Base: entity.Base{
						ID: v.Akun.ID,
					},
					Saldo:          v.Akun.Saldo,
					Nama:           v.Akun.Nama,
					Kode:           v.Akun.Kode,
					KelompokAkunID: v.Akun.KelompokAkunID,
				}
			}

			akun.Saldo -= v.Saldo
			akunsMap[v.AkunID] = akun
		}

		// wg := &sync.WaitGroup{}
		// m := &sync.Mutex{}

		// validate kredit, debit and get totalTransaksi
		var totalTransaksi float64

		if err := isKreditEqualToDebit(reqTransaksi.AyatJurnal, &totalTransaksi); err != nil {
			return err
		}
		if err := isDuplicateAkun(reqTransaksi.AyatJurnal); err != nil {
			return err
		}

		newTr.Total = totalTransaksi

		// add newAyatJurnals
		newAyatJurnals := make([]*entity.AyatJurnal, len(reqTransaksi.AyatJurnal))

		// maxConcurent := make(chan struct{}, 10)
		for i, ay := range reqTransaksi.AyatJurnal {
			akun := akunsMap[ay.AkunID]
			// logic calculate saldo ayatJurnal
			var saldo float64
			// saldo += ay.Debit
			if ay.Kredit != 0 {
				updateSaldo(&akun.Saldo, &saldo, ay.Kredit, akun.SaldoNormal == "DEBIT")
			}
			if ay.Debit != 0 {
				updateSaldo(&akun.Saldo, &saldo, ay.Debit, akun.SaldoNormal == "KREDIT")
			}

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
		// wg.Wait()

		// create new values for saldo akun
		// newSaldoAkun := make([]string, len(akunsMap))

		newAkuns := make([]entity.Akun, len(akunsMap))
		i := 0
		for _, v := range akunsMap {
			newAkuns[i] = v
			i++
			// newSaldoAkun[i] = fmt.Sprintf("('%s', %.2f, '%s', '%s', '%s')", v.ID, v.Saldo, v.KelompokAkunID, v.Nama, v.Kode)
		}

		updateTrParam.NewAyatJurnals = newAyatJurnals
		updateTrParam.UpdateAkuns = newAkuns
		// updateTrParam.NewSaldoAkunValues = newSaldoAkun
	}
	updateTrParam.UpdateTr = &newTr
	// fmt.Println(updateTrParam.NewAyatJurnals)

	if err := u.repo.Update(ctx, updateTrParam); err != nil {
		return err
	}

	return nil
}

func (u *transaksiUsecase) Create(ctx context.Context, reqTransaksi req.Create) error {
	// validate kredit and debit
	// wg := &sync.WaitGroup{}

	var totalTransaksi float64
	if err := isKreditEqualToDebit(reqTransaksi.AyatJurnal, &totalTransaksi); err != nil {
		return err
	}
	if err := isDuplicateAkun(reqTransaksi.AyatJurnal); err != nil {
		return err
	}
	// akun, err := u.repoAkun.GetById(ctx, )

	akunIds := make([]string, len(reqTransaksi.AyatJurnal))

	for i, v := range reqTransaksi.AyatJurnal {
		akunIds[i] = v.AkunID
	}

	akuns, err := u.repoAkun.GetByIds(ctx, akunIds)
	if err != nil {
		return err
	}

	// map to make essier to get by id without call db inside looping
	mapAkuns := make(map[string]entity.Akun, len(akuns))

	for _, akun := range akuns {
		mapAkuns[akun.ID] = akun
	}

	dataAyatJurnals := make([]*entity.AyatJurnal, len(reqTransaksi.AyatJurnal))
	transaksiID := u.ulid.MakeUlid().String()

	updateAkuns := []*entity.Akun{}
	m := &sync.Mutex{}
	g := errgroup.Group{}

	for i, v := range reqTransaksi.AyatJurnal {
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
				Base: entity.Base{
					ID: u.ulid.MakeUlid().String(),
				},
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

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
}

type transaksiUsecase struct {
	repo     repo.TransaksiRepo
	repoAkun repoAkun.AkunRepo
	ulid     pkg.UlidPkg
}

func NewTransaksiUsecase(repo repo.TransaksiRepo, repoAkun repoAkun.AkunRepo, ulid pkg.UlidPkg) TransaksiUsecase {
	return &transaksiUsecase{repo, repoAkun, ulid}
}

func updateSaldo(saldoAkhir *float64, saldo *float64, amount float64, isDebit bool) {
	if isDebit {
		*saldoAkhir -= amount
		*saldo -= amount
	} else {
		*saldoAkhir += amount
		*saldo += amount
	}
}

func isKreditEqualToDebit(ayatJurnals []req.ReqAyatJurnal, wg *sync.WaitGroup, m *sync.Mutex) error {
	var totalKredit, totalDebit float64
	for _, v := range ayatJurnals {
		wg.Add(1)
		go func(v req.ReqAyatJurnal) {
			defer wg.Done()
			if v.Debit != 0 {
				m.Lock()
				totalDebit += v.Debit
				m.Unlock()
			}
			if v.Kredit != 0 {
				m.Lock()
				totalKredit += v.Kredit
				m.Unlock()
			}
		}(v)
	}
	wg.Wait()

	if totalDebit != totalKredit {
		return errors.New(message.CreditDebitNotSame)
	}
	return nil
}

func isDuplicateAkun(ayatJurnals []req.ReqAyatJurnal) error {
	akunCount := make(map[string]int)
	var debit, kredit float64

	for _, v := range ayatJurnals {
		akunCount[v.AkunID]++

		if akunCount[v.AkunID] > 1 {
			return errors.New(message.AkunCannotBeSame)
		}

		debit += v.Debit
		kredit += v.Kredit

		if debit == kredit {
			// Reset counters and map for the next pair
			debit = 0
			kredit = 0
			akunCount = make(map[string]int)
		}
	}

	return nil
}

func (u *transaksiUsecase) Create(ctx context.Context, reqTransaksi req.Create) error {
	parsedTime, err := time.Parse(time.DateOnly, reqTransaksi.Tanggal)
	if err != nil {
		helper.LogsError(err)
		return err
	}

	// validate kredit and debit
	wg := &sync.WaitGroup{}
	m := &sync.Mutex{}

	if err := isKreditEqualToDebit(reqTransaksi.AyatJurnals, wg, m); err != nil {
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
		akun := akun
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Lock()
			mapAkuns[akun.ID] = akun
			m.Unlock()
		}()
	}
	wg.Wait()

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
			var saldo float64
			if v.Kredit != 0 {
				updateSaldo(&akun.SaldoAkhir, &saldo, v.Kredit, akun.SaldoNormal == "DEBIT")
			}
			if v.Debit != 0 {
				updateSaldo(&akun.SaldoAkhir, &saldo, v.Debit, akun.SaldoNormal == "KREDIT")
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

	dataTransaksi := entity.Transaksi{
		ID:              transaksiID,
		Tanggal:         parsedTime,
		Keterangan:      reqTransaksi.Keterangan,
		BuktiPembayaran: reqTransaksi.BuktiPembayaran,
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

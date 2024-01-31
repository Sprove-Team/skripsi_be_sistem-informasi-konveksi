package akuntansi

import (
	"context"
	"errors"
	"math"
	"strings"
	"sync"
	"time"

	repoAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/hutang_piutang"
	repoKontak "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kontak"
	pkgAkuntansiLogic "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	res "github.com/be-sistem-informasi-konveksi/common/response/akuntansi/hutang_piutang"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"golang.org/x/sync/errgroup"
)

type HutangPiutangUsecase interface {
	Create(ctx context.Context, reqHutangPiutang req.Create) error
	GetAll(ctx context.Context, reqHutangPiutang req.GetAll) ([]res.GetAll, error)
	// CreateBayarHutangPiutang(ctx context.Context, reqBayarHutangPiutang req.BayarHutangPiutang) error
}

type hutangPiutangUsecase struct {
	repo       repo.HutangPiutangRepo
	repoAkun   repoAkun.AkunRepo
	repoKontak repoKontak.KontakRepo
	ulid       pkg.UlidPkg
}

func NewHutangPiutangUsecase(repo repo.HutangPiutangRepo, repoAkun repoAkun.AkunRepo, repoKontak repoKontak.KontakRepo, ulid pkg.UlidPkg) HutangPiutangUsecase {
	return &hutangPiutangUsecase{repo, repoAkun, repoKontak, ulid}
}

func (u *hutangPiutangUsecase) Create(ctx context.Context, reqHutangPiutang req.Create) error {

	g := &errgroup.Group{}
	g.SetLimit(10)
	// check kontak is founded
	var kontak entity.Kontak
	g.Go(func() error {
		data, err := u.repoKontak.GetById(ctx, reqHutangPiutang.KontakID)
		if err != nil {
			if err.Error() == "record not found" {
				return errors.New(message.KontakNotFound)
			}
			return err
		}
		kontak = data
		return nil
	})

	hutangPiutangID := u.ulid.MakeUlid().String()
	repoParam := &entity.HutangPiutang{
		Base: entity.Base{
			ID: hutangPiutangID,
		},
		Jenis: reqHutangPiutang.Jenis,
	}

	ayHP1 := entity.AyatJurnal{
		Base: entity.Base{
			ID: u.ulid.MakeUlid().String(),
		},
		AkunID: reqHutangPiutang.Transaksi.AyatJurnal[0].AkunID,
		Debit:  reqHutangPiutang.Transaksi.AyatJurnal[0].Debit,
		Kredit: reqHutangPiutang.Transaksi.AyatJurnal[0].Kredit,
	}
	ayHP2 := entity.AyatJurnal{
		Base: entity.Base{
			ID: u.ulid.MakeUlid().String(),
		},
		AkunID: reqHutangPiutang.Transaksi.AyatJurnal[1].AkunID,
		Debit:  reqHutangPiutang.Transaksi.AyatJurnal[1].Debit,
		Kredit: reqHutangPiutang.Transaksi.AyatJurnal[1].Kredit,
	}

	lengthAyByr := len(reqHutangPiutang.BayarAwal)
	akunIds := make([]string, 2+lengthAyByr)
	akunIds[0] = ayHP1.AkunID
	akunIds[1] = ayHP2.AkunID

	for i, akunByr := range reqHutangPiutang.BayarAwal {
		akunIds[i+2] = akunByr.AkunBayarID
	}

	var akuns []entity.Akun
	g.Go(func() error {
		datas, err := u.repoAkun.GetByIds(ctx, akunIds)
		if err != nil {
			if len(datas) != 3 {
				return errors.New(message.AkunNotFound)
			}
			return err
		}
		akuns = datas
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	akunMap := map[string]entity.Akun{}

	// validate akun + add ayat jurnal saldo untk create tr hutang piutang
	ayTagihan := entity.AyatJurnal{} // ay hutang/piutang yg akan dibayar bila bayar tidak kosong
	akunTagihan := entity.Akun{}
	for _, akun := range akuns {
		if !pkgAkuntansiLogic.IsValidAkunHutangPiutang(akun.KelompokAkun.Nama) {
			return errors.New(message.InvalidAkunHutangPiutang)
		}

		// pasti ada 1 saldo debit jika jenis piutang, jika tidak error, begitu sebaliknya untuk hutang
		switch akun.ID {
		case ayHP1.AkunID:
			pkgAkuntansiLogic.UpdateSaldo(&ayHP1.Saldo, ayHP1.Kredit, ayHP1.Debit, akun.SaldoNormal)
			if reqHutangPiutang.Jenis == "PIUTANG" {
				if akun.SaldoNormal == "DEBIT" {
					ayTagihan = ayHP2
					akunTagihan = akun
				}
			} else {
				if akun.SaldoNormal == "KREDIT" {
					ayTagihan = ayHP2
					akunTagihan = akun
				}
			}

		case ayHP2.AkunID:
			pkgAkuntansiLogic.UpdateSaldo(&ayHP2.Saldo, ayHP2.Kredit, ayHP2.Debit, akun.SaldoNormal)
			if reqHutangPiutang.Jenis == "PIUTANG" {
				if akun.SaldoNormal == "DEBIT" {
					ayTagihan = ayHP2
					akunTagihan = akun
				}
			} else {
				if akun.SaldoNormal == "KREDIT" {
					ayTagihan = ayHP2
					akunTagihan = akun
				}
			}
		}

		akunMap[akun.ID] = akun
	}

	// cek kelompok akun sesuai dengan jenis hutang/piutang
	if akunTagihan.KelompokAkun == nil {
		return errors.New(message.AkunNotMatchWithJenisHP)
	}
	if !strings.EqualFold(akunTagihan.KelompokAkun.Nama, reqHutangPiutang.Jenis) {
		return errors.New(message.AkunNotMatchWithJenisHP)
	}

	// validate credit, debit hp is equal
	if ayHP1.Saldo < 0 || ayHP2.Saldo < 0 {
		return errors.New(message.IncorrectPlacementOfCreditAndDebit)
	}
	if ayHP1.Saldo-ayHP2.Saldo != 0 {
		return errors.New(message.CreditDebitNotSame)
	}

	// parse tanggal
	tanggalTrHp, err := time.Parse(time.RFC3339, reqHutangPiutang.Transaksi.Tanggal)
	if err != nil {
		helper.LogsError(err)
		return err
	}

	// create transaksi HP, based on ay1 and ay2 variable
	transaksiHpID := u.ulid.MakeUlid().String()
	transaksiHP := entity.Transaksi{
		Base: entity.Base{
			ID: transaksiHpID,
		},
		Keterangan:      reqHutangPiutang.Keterangan,
		BuktiPembayaran: reqHutangPiutang.Transaksi.BuktiPembayaran,
		Total:           math.Abs(ayHP1.Saldo),
		Tanggal:         tanggalTrHp,
		KontakID:        reqHutangPiutang.KontakID,
		AyatJurnals:     []entity.AyatJurnal{ayHP1, ayHP2},
	}

	// create transaksi bayar hp
	dataByrHutangPiutang := make([]entity.DataBayarHutangPiutang, lengthAyByr)
	var totalAllTrBayar float64
	if lengthAyByr > 0 {
		for i, trByr := range reqHutangPiutang.BayarAwal {
			totalAllTrBayar += trByr.Total
			akun := akunMap[trByr.AkunBayarID]
			if akun.KelompokAkun.Nama != "kas & bank" {
				return errors.New(message.InvalidAkunBayar)
			}
			byrHP, err := autoCreateDataBayarAwal(trByr, ayTagihan, reqHutangPiutang.KontakID, kontak.Nama+" melakukan pembayaran dengan akun "+akun.Nama, u.ulid)
			if err != nil {
				return err
			}
			dataByrHutangPiutang[i] = *byrHP
		}
	}

	if totalAllTrBayar > ayTagihan.Debit {
		return errors.New(message.BayarMustLessThanSisaTagihan)
	}

	repoParam.Total = transaksiHP.Total
	repoParam.Sisa = transaksiHP.Total - totalAllTrBayar
	if repoParam.Sisa <= 0 {
		repoParam.Status = "LUNAS"
	}
	repoParam.DataBayarHutangPiutang = dataByrHutangPiutang
	repoParam.Transaksi = transaksiHP

	// fmt.Println(len(repoParam.DataBayarHutangPiutang))

	if err := u.repo.Create(ctx, repoParam); err != nil {
		return err
	}
	return nil
}

func (u *hutangPiutangUsecase) GetAll(ctx context.Context, reqHutangPiutang req.GetAll) ([]res.GetAll, error) {
	repoSearch := repo.SearchParam{
		KontakID: reqHutangPiutang.KontakID,
		Status:   []string{reqHutangPiutang.Status},
		Jenis:    []string{reqHutangPiutang.Jenis},
	}

	if reqHutangPiutang.Jenis == "ALL" || reqHutangPiutang.Jenis == "" {
		repoSearch.Jenis = []string{"PIUTANG", "HUTANG"}
	}
	if reqHutangPiutang.Status == "ALL" || reqHutangPiutang.Status == "" {
		repoSearch.Status = []string{"BELUM_LUNAS", "LUNAS"}
	}

	datas, err := u.repo.GetAll(ctx, repoSearch)
	if err != nil {
		return nil, err
	}

	groupDataByKontak := make(map[string]res.GetAll)

	m := new(sync.RWMutex)
	wg := new(sync.WaitGroup)
	maxGoroutines := 12
	guard := make(chan struct{}, maxGoroutines)
	for _, dataHpMixUp := range datas {
		wg.Add(1)
		guard <- struct{}{} // would block if guard channel is already filled
		go func(data entity.HutangPiutang) {
			defer func() {
				<-guard
				wg.Done()
			}()
			m.Lock()
			defer m.Unlock()
			dat, ok := groupDataByKontak[data.Transaksi.Kontak.Nama]
			if !ok {
				dat = res.GetAll{
					Nama: data.Transaksi.Kontak.Nama,
				}
			}

			dataHp := res.ResDataHutangPiutang{
				ID:          data.ID,
				InvoiceSlug: data.InvoiceSlug,
				Jenis:       data.Jenis,
				TransaksiID: data.TransaksiID,
				Status:      data.Status,
				Total:       data.Total,
				Sisa:        data.Sisa,
			}

			if data.Jenis == "PIUTANG" {
				dat.TotalPiutang += data.Total
				dat.SisaPiutang += data.Sisa
			} else {
				dat.TotalHutang += data.Total
				dat.SisaHutang += data.Sisa
			}

			dat.HutangPiutang = append(dat.HutangPiutang, dataHp)
			groupDataByKontak[data.Transaksi.Kontak.Nama] = dat

		}(dataHpMixUp)
	}
	wg.Wait()

	resData := make([]res.GetAll, len(groupDataByKontak))
	i := 0
	for _, v := range groupDataByKontak {
		resData[i] = v
		i++
	}
	return resData, nil
}

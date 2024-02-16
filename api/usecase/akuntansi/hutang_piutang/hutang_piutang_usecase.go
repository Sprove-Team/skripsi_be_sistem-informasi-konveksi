package akuntansi

import (
	"context"
	"errors"
	"math"
	"strings"
	"sync"
	"time"

	repoAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	repoBayarHP "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/data_bayar_hutang_piutang"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/hutang_piutang"
	repoKontak "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kontak"
	pkgAkuntansiLogic "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	res "github.com/be-sistem-informasi-konveksi/common/response/akuntansi/hutang_piutang"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type (
	ParamCreateDataHp struct {
		Ctx context.Context
		Req req.Create
	}
)

type HutangPiutangUsecase interface {
	CreateDataHP(param ParamCreateDataHp) (*entity.HutangPiutang, error)
	CreateCommitDB(ctx context.Context, hp *entity.HutangPiutang) error
	CreateBayar(ctx context.Context, reqHutangPiutang req.CreateBayar) error
	GetAll(ctx context.Context, reqHutangPiutang req.GetAll) ([]res.GetAll, error)
}

type hutangPiutangUsecase struct {
	repo        repo.HutangPiutangRepo
	repoBayarHP repoBayarHP.DataBayarHutangPiutangRepo
	repoAkun    repoAkun.AkunRepo
	repoKontak  repoKontak.KontakRepo
	ulid        pkg.UlidPkg
}

func NewHutangPiutangUsecase(repo repo.HutangPiutangRepo, repoBayarHP repoBayarHP.DataBayarHutangPiutangRepo, repoAkun repoAkun.AkunRepo, repoKontak repoKontak.KontakRepo, ulid pkg.UlidPkg) HutangPiutangUsecase {
	return &hutangPiutangUsecase{repo, repoBayarHP, repoAkun, repoKontak, ulid}
}

func (u *hutangPiutangUsecase) CreateDataHP(param ParamCreateDataHp) (*entity.HutangPiutang, error) {
	hutangPiutangID := u.ulid.MakeUlid().String()
	dataHP := &entity.HutangPiutang{
		Base: entity.Base{
			ID: hutangPiutangID,
		},
		Jenis: param.Req.Jenis,
	}

	ayHP1 := entity.AyatJurnal{
		Base: entity.Base{
			ID: u.ulid.MakeUlid().String(),
		},
		AkunID: param.Req.Transaksi.AyatJurnal[0].AkunID,
		Debit:  param.Req.Transaksi.AyatJurnal[0].Debit,
		Kredit: param.Req.Transaksi.AyatJurnal[0].Kredit,
	}
	ayHP2 := entity.AyatJurnal{
		Base: entity.Base{
			ID: u.ulid.MakeUlid().String(),
		},
		AkunID: param.Req.Transaksi.AyatJurnal[1].AkunID,
		Debit:  param.Req.Transaksi.AyatJurnal[1].Debit,
		Kredit: param.Req.Transaksi.AyatJurnal[1].Kredit,
	}

	lengthAyByr := len(param.Req.BayarAwal)
	akunIds := make([]string, 2+lengthAyByr)
	akunIds[0] = ayHP1.AkunID
	akunIds[1] = ayHP2.AkunID

	for i, akunByr := range param.Req.BayarAwal {
		akunIds[i+2] = akunByr.AkunBayarID
	}

	akuns, err := u.repoAkun.GetByIds(param.Ctx, akunIds)
	if err != nil {
		if len(akuns) != 3 {
			return nil, errors.New(message.AkunNotFound)
		}
		return nil, err
	}

	akunMap := map[string]entity.Akun{}

	// validate akun + add ayat jurnal saldo untk create tr hutang piutang
	ayTagihan := entity.AyatJurnal{} // ay hutang/piutang yg akan dibayar bila bayar tidak kosong
	akunTagihan := entity.Akun{}
	for _, akun := range akuns {
		if err := pkgAkuntansiLogic.IsValidAkunHutangPiutang(akun.KelompokAkun.Nama); err != nil {
			return nil, err
		}

		// pasti ada 1 saldo debit jika jenis piutang, jika tidak error, begitu sebaliknya untuk hutang
		switch akun.ID {
		case ayHP1.AkunID:
			pkgAkuntansiLogic.UpdateSaldo(&ayHP1.Saldo, ayHP1.Kredit, ayHP1.Debit, akun.SaldoNormal)
			if param.Req.Jenis == "PIUTANG" {
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
			if param.Req.Jenis == "PIUTANG" {
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
		return nil, errors.New(message.AkunNotMatchWithJenisHP)
	}
	if !strings.EqualFold(akunTagihan.KelompokAkun.Nama, param.Req.Jenis) {
		return nil, errors.New(message.AkunNotMatchWithJenisHP)
	}

	// validate credit, debit hp is equal
	if ayHP1.Saldo < 0 || ayHP2.Saldo < 0 {
		return nil, errors.New(message.IncorrectPlacementOfCreditAndDebit)
	}
	if ayHP1.Saldo-ayHP2.Saldo != 0 {
		return nil, errors.New(message.CreditDebitNotSame)
	}

	// parse tanggal
	tanggalTrHp, err := time.Parse(time.RFC3339, param.Req.Transaksi.Tanggal)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}

	// create transaksi HP, based on ay1 and ay2 variable
	transaksiHpID := u.ulid.MakeUlid().String()
	transaksiHP := entity.Transaksi{
		Base: entity.Base{
			ID: transaksiHpID,
		},
		Keterangan:      param.Req.Keterangan,
		BuktiPembayaran: param.Req.Transaksi.BuktiPembayaran,
		Total:           math.Abs(ayHP1.Saldo),
		Tanggal:         tanggalTrHp,
		KontakID:        param.Req.KontakID,
		AyatJurnals:     []entity.AyatJurnal{ayHP1, ayHP2},
	}

	// create transaksi bayar awal hp
	dataByrHutangPiutang := make([]entity.DataBayarHutangPiutang, lengthAyByr)
	var totalAllTrBayar float64
	if lengthAyByr > 0 {
		for i, trByr := range param.Req.BayarAwal {
			totalAllTrBayar += trByr.Total
			akun := akunMap[trByr.AkunBayarID]
			if akun.KelompokAkun.Nama != "kas & bank" {
				return nil, errors.New(message.InvalidAkunBayar)
			}
			byrHP, err := pkgAkuntansiLogic.CreateDataBayarHP(trByr, ayTagihan, param.Req.KontakID, trByr.Keterangan, u.ulid)
			if err != nil {
				return nil, err
			}
			dataByrHutangPiutang[i] = *byrHP
		}
	}

	if totalAllTrBayar > ayTagihan.Debit {
		return nil, errors.New(message.BayarMustLessThanSisaTagihan)
	}

	dataHP.Total = transaksiHP.Total
	dataHP.Sisa = transaksiHP.Total - totalAllTrBayar
	if dataHP.Sisa <= 0 {
		dataHP.Status = "LUNAS"
	}
	dataHP.DataBayarHutangPiutang = dataByrHutangPiutang
	dataHP.Transaksi = transaksiHP

	return dataHP, nil
}

func (u *hutangPiutangUsecase) CreateCommitDB(ctx context.Context, hp *entity.HutangPiutang) error {
	_, err := u.repoKontak.GetById(repoKontak.ParamGetById{
		Ctx: ctx,
		ID:  hp.Transaksi.KontakID,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New(message.KontakNotFound)
		}
		return err
	}
	return u.repo.Create(repo.ParamCreate{
		Ctx:           ctx,
		HutangPiutang: hp,
	})
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

	datas, err := u.repo.GetAll(repo.ParamGetAll{
		Ctx:    ctx,
		Search: repoSearch,
	})
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
				InvoiceID:   data.InvoiceID,
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

func (u *hutangPiutangUsecase) CreateBayar(ctx context.Context, reqHutangPiutang req.CreateBayar) error {

	hpChan := make(chan entity.HutangPiutang, 1)
	errChan := make(chan error, 1)
	go func() {
		hp, err := u.repo.GetHPForBayar(repo.ParamGetHPForBayar{
			Ctx: ctx,
			ID:  reqHutangPiutang.HutangPiutangID,
		})
		if err != nil {
			if err.Error() == "record not found" {
				errChan <- errors.New(message.HutangPiutangNotFound)
				return
			}
			errChan <- err
			return
		}
		if reqHutangPiutang.Total > hp.Sisa {
			errChan <- errors.New(message.BayarMustLessThanSisaTagihan)
			return
		}
		hpChan <- *hp
	}()
	go func() {
		akun, err := u.repoAkun.GetById(ctx, reqHutangPiutang.AkunBayarID)
		if err != nil {
			if err.Error() == "record not found" {
				errChan <- errors.New(message.AkunNotFound)
				return
			}
			errChan <- err
			return
		}
		if akun.KelompokAkun.Nama != "kas & bank" {
			errChan <- errors.New(message.InvalidAkunBayar)
			return
		}
	}()

	var ayTagihan entity.AyatJurnal
	var byrHP *entity.DataBayarHutangPiutang
	select {
	case hp := <-hpChan:
		for _, ay := range hp.Transaksi.AyatJurnals {
			// Check if the account ID matches and transaction type is 'PIUTANG' or 'HUTANG'
			kodeKlompokAkun := ay.Akun.Kode[0:2]
			if (kodeKlompokAkun == "12" && hp.Jenis == "PIUTANG") ||
				(kodeKlompokAkun == "27" && hp.Jenis == "HUTANG") {
				ayTagihan.AkunID = ay.AkunID
				ayTagihan.Saldo = -reqHutangPiutang.Total

				// Set Kredit or Debit based on transaction type
				if hp.Jenis == "PIUTANG" {
					ayTagihan.Kredit = reqHutangPiutang.Total
				} else if hp.Jenis == "HUTANG" {
					ayTagihan.Debit = reqHutangPiutang.Total
				}
				break
			}
		}
		var err error

		byrHP, err = pkgAkuntansiLogic.CreateDataBayarHP(reqHutangPiutang.ReqBayar, ayTagihan, hp.Transaksi.KontakID, reqHutangPiutang.Keterangan, u.ulid)

		if err != nil {
			return err
		}

		byrHP.HutangPiutang = hp
		byrHP.HutangPiutang.Sisa = hp.Sisa - reqHutangPiutang.Total
		if byrHP.HutangPiutang.Sisa <= 0 {
			byrHP.HutangPiutang.Status = "LUNAS"
		}

	case err := <-errChan:
		return err
	}

	if err := u.repoBayarHP.Create(ctx, byrHP); err != nil {
		return err
	}

	return nil
}

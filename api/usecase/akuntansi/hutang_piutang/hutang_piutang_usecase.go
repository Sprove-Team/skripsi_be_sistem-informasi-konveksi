package uc_akuntansi_hp

import (
	"context"
	"errors"
	"fmt"
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
	CreateDataBayar(ctx context.Context, reqHutangPiutang req.CreateBayar) (*entity.DataBayarHutangPiutang, error)
	CreateBayarCommitDB(ctx context.Context, byrHP *entity.DataBayarHutangPiutang) error
	GetAll(ctx context.Context, reqHutangPiutang req.GetAll) ([]res.GetAll, error)
	GetHPByInvoiceID(ctx context.Context, id string) (*entity.HutangPiutang, error)
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

	akunIds := []string{ayHP1.AkunID, ayHP2.AkunID}

	akuns, err := u.repoAkun.GetByIds(param.Ctx, akunIds)
	if err != nil {
		return nil, err
	}

	if len(akuns) != len(akunIds) {
		return nil, errors.New(message.AkunNotFound)
	}

	// validate akun + add ayat jurnal saldo untk create tr hutang piutang
	akunMap := map[string]entity.Akun{}
	akunYgAkanDibayar := entity.Akun{} // akun yang akan dibayar nantinya
	for _, akun := range akuns {
		if err := pkgAkuntansiLogic.IsValidAkunHutangPiutang(akun.KelompokAkun.Nama); err != nil {
			return nil, err
		}

		switch akun.ID {
		case ayHP1.AkunID:
			pkgAkuntansiLogic.UpdateSaldo(&ayHP1.Saldo, ayHP1.Kredit, ayHP1.Debit, akun.SaldoNormal)

		case ayHP2.AkunID:
			pkgAkuntansiLogic.UpdateSaldo(&ayHP2.Saldo, ayHP2.Kredit, ayHP2.Debit, akun.SaldoNormal)
		}

		//! mencari akun pembayaran berdasarkan jenis dari requestnya
		//? pasti ada 1 saldo debit jika jenis piutang, jika tidak error, begitu sebaliknya untuk hutang
		if (param.Req.Jenis == "PIUTANG" && akun.SaldoNormal == "DEBIT") ||
			(param.Req.Jenis == "HUTANG" && akun.SaldoNormal == "KREDIT") {
			akunYgAkanDibayar = akun
		}

		akunMap[akun.ID] = akun
	}

	fmt.Println("akun pembayaran -> ", akunYgAkanDibayar.Nama)
	// cek kelompok akun yang akaun dibayar sesuai dengan jenis hutang/piutang
	// validasi ini ada untuk menghindari pengurangan atau penambahan ayat jurnal pada akun hp yang salah
	if akunYgAkanDibayar.ID == "" || !strings.EqualFold(akunYgAkanDibayar.KelompokAkun.Nama, param.Req.Jenis) {
		return nil, errors.New(message.IncorrectEntryAkunHP)
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
	transaksiHP := entity.Transaksi{
		Base: entity.Base{
			ID: u.ulid.MakeUlid().String(),
		},
		Keterangan:      param.Req.Keterangan,
		BuktiPembayaran: param.Req.Transaksi.BuktiPembayaran,
		Total:           math.Abs(ayHP1.Saldo),
		Tanggal:         tanggalTrHp,
		KontakID:        param.Req.KontakID,
		AyatJurnals:     []entity.AyatJurnal{ayHP1, ayHP2},
	}

	dataHP.Total = transaksiHP.Total
	dataHP.Sisa = transaksiHP.Total
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
				Tanggal:     data.Transaksi.Tanggal.Format(time.RFC3339),
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

func (u *hutangPiutangUsecase) GetHPByInvoiceID(ctx context.Context, id string) (*entity.HutangPiutang, error) {
	dataHp, err := u.repo.GetByInvoiceId(repo.ParamGetByInvoiceId{
		Ctx: ctx,
		ID:  id,
	})
	if err != nil {
		return nil, err
	}
	return dataHp, nil
}

func (u *hutangPiutangUsecase) CreateDataBayar(ctx context.Context, reqHutangPiutang req.CreateBayar) (*entity.DataBayarHutangPiutang, error) {

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
			return nil, err
		}

		byrHP.HutangPiutang = hp
		byrHP.HutangPiutang.Sisa = hp.Sisa - reqHutangPiutang.Total
		if byrHP.HutangPiutang.Sisa <= 0 {
			byrHP.HutangPiutang.Status = "LUNAS"
		}

	case err := <-errChan:
		return nil, err
	}
	return byrHP, nil

}

func (u *hutangPiutangUsecase) CreateBayarCommitDB(ctx context.Context, byrHP *entity.DataBayarHutangPiutang) error {
	if err := u.repoBayarHP.Create(ctx, byrHP); err != nil {
		return err
	}
	return nil
}

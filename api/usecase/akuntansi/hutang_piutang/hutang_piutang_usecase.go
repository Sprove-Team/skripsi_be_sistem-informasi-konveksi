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
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"golang.org/x/sync/errgroup"
)

type HutangPiutangUsecase interface {
	Create(ctx context.Context, reqHutangPiutang req.Create) error
	// GetAll(ctx context.Context, reqHutangPiutang req.GetAll) ([]entity.HutangPiutang, error)
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

func isValidAkunHutangPiutang(nama string) bool {
	validNames := []string{"piutang", "hutang", "kas & bank"}
	if strings.HasPrefix(nama, "pendapatan") {
		return true
	}
	for _, validName := range validNames {
		if nama == validName {
			return true
		}
	}
	return false
}

func (u *hutangPiutangUsecase) Create(ctx context.Context, reqHutangPiutang req.Create) error {

	g := &errgroup.Group{}
	m := &sync.Mutex{}
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

	hutangPiuatangID := u.ulid.MakeUlid().String()
	repoParam := &entity.HutangPiutang{
		Base: entity.Base{
			ID: hutangPiuatangID,
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
	ayTagihan := entity.AyatJurnal{} // ay utang/piutang yg akan dibayar bila bayar tidak kosong
	for _, akun := range akuns {
		func(akun entity.Akun) {
			g.Go(func() error {
				if !isValidAkunHutangPiutang(akun.KelompokAkun.Nama) {
					return errors.New(message.InvalidAkunHutangPiutang)
				}
				switch akun.ID {
				case ayHP1.AkunID:
					helper.UpdateSaldo(&ayHP1.Saldo, ayHP1.Kredit, ayHP1.Debit, akun.SaldoNormal)
					if akun.SaldoNormal == "DEBIT" {
						ayTagihan = ayHP1
					}
				case ayHP2.AkunID:
					helper.UpdateSaldo(&ayHP2.Saldo, ayHP2.Kredit, ayHP2.Debit, akun.SaldoNormal)
					if akun.SaldoNormal == "DEBIT" {
						ayTagihan = ayHP2
					}
				}
				m.Lock()
				akunMap[akun.ID] = akun
				m.Unlock()
				return nil
			})
		}(akun)
	}
	if err := g.Wait(); err != nil {
		return err
	}

	// validate credit, debit hp is equal
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
			tanggal, err := time.Parse(time.RFC3339, trByr.Tanggal)
			if err != nil {
				helper.LogsError(err)
				return err
			}
			trByr.Total = math.Abs(trByr.Total)
			// cek jika total bayarnya lebih besar dari sisa tagihan
			if trByr.Total > ayTagihan.Debit {
				return errors.New(message.BayarMustLessThanSisaTagihan)
			}

			totalAllTrBayar += trByr.Total

			byrHP := entity.DataBayarHutangPiutang{
				Base: entity.Base{
					ID: u.ulid.MakeUlid().String(),
				},
				Transaksi: entity.Transaksi{
					Base: entity.Base{
						ID: u.ulid.MakeUlid().String(),
					},
					Keterangan:      kontak.Nama + " melakukan pembayaran dengan akun " + akunMap[trByr.AkunBayarID].Nama,
					BuktiPembayaran: trByr.BuktiPembayaran,
					Total:           trByr.Total,
					Tanggal:         tanggal,
					KontakID:        reqHutangPiutang.KontakID,
					AyatJurnals: []entity.AyatJurnal{
						{
							Base: entity.Base{
								ID: u.ulid.MakeUlid().String(),
							},
							AkunID: trByr.AkunBayarID,
							Debit:  trByr.Total,
							Saldo:  trByr.Total,
						},
						// ay kredit, bayar mengurangi akun utang/piutangnya
						{
							Base: entity.Base{
								ID: u.ulid.MakeUlid().String(),
							},
							AkunID: ayTagihan.AkunID,
							Kredit: trByr.Total,
							Saldo:  -trByr.Total,
						},
					},
				},
			}

			dataByrHutangPiutang[i] = byrHP
		}
	}

	if totalAllTrBayar > ayTagihan.Debit {
		return errors.New(message.BayarMustLessThanSisaTagihan)
	}

	// default status = "BELUM_LUNAS"
	sisaHutangPiutang := ayTagihan.Debit - totalAllTrBayar
	if sisaHutangPiutang == 0 {
		repoParam.Status = "LUNAS"
	}

	repoParam.DataBayarHutangPiutang = dataByrHutangPiutang
	repoParam.Transaksi = transaksiHP

	if err := u.repo.Create(ctx, repoParam); err != nil {
		return err
	}
	return nil
}

// func (u *hutangPiutangUsecase) GetAll(ctx context.Context, reqHutangPiutang req.GetAll) ([]entity.HutangPiutang, error) {
// 	// if reqHutangPiutang.Jenis ==
// }

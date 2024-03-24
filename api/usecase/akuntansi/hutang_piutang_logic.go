package uc_akuntansi

import (
	"errors"
	"math"
	"time"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	reqTr "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/transaksi"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

func IsAkunHutangPiutangExist(akuns []entity.Akun) error {
	for _, akun := range akuns {
		if akun.KelompokAkun.Nama == "piutang" || akun.KelompokAkun.Nama == "hutang" {
			return nil
		}
	}
	return errors.New(message.AkunHPDoesNotExist)
}

func IsValidAkunHp(hp *entity.HutangPiutang, akun *entity.Akun, ay *reqTr.ReqAyatJurnal) error {
	// validasi akun sesuai dengan jenis hutang piutang, tr hp
	if hp.ID == "" {
		return nil
	}

	if hp.Jenis == "PIUTANG" && ay.Debit != 0 {
		if akun.KelompokAkun.Nama != "piutang" {
			return errors.New(message.AkunNotMatchWithJenisHPTr)
		}
	}
	if hp.Jenis == "HUTANG" && ay.Kredit != 0 {
		if akun.KelompokAkun.Nama != "hutang" {
			return errors.New(message.AkunNotMatchWithJenisHPTr)
		}
	}

	return nil
}

func IsValidAkunByrHP(dataByrHP *entity.DataBayarHutangPiutang, akun *entity.Akun, ay *reqTr.ReqAyatJurnal) error {

	// untuk tr bayr hp
	if dataByrHP.HutangPiutang.ID != "" && ay.Debit != 0 {
		if dataByrHP.HutangPiutang.Jenis == "PIUTANG" {
			if ay.Kredit != 0 && akun.KelompokAkun.Nama != "piutang" {
				return errors.New(message.AkunNotMatchWithJenisHPTr)
			}
			if ay.Debit != 0 && akun.KelompokAkun.Nama != "kas & bank" {
				return errors.New(message.InvalidAkunBayar)
			}
		}
		if dataByrHP.HutangPiutang.Jenis == "HUTANG" {
			if ay.Debit != 0 && akun.KelompokAkun.Nama != "hutang" {
				return errors.New(message.AkunNotMatchWithJenisHPTr)
			}
			if ay.Kredit != 0 && akun.KelompokAkun.Nama != "kas & bank" {
				return errors.New(message.InvalidAkunBayar)
			}
		}
	}

	return nil
}

func CreateDataBayarHP(trByr req.ReqBayar, ayTagihan entity.AyatJurnal, kontakId, keterangan string, ulid pkg.UlidPkg) (*entity.DataBayarHutangPiutang, error) {
	tanggal, err := time.Parse(time.RFC3339, trByr.Tanggal)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	trByr.Total = math.Abs(trByr.Total)
	ayatJurnals := []entity.AyatJurnal{
		{
			Base: entity.Base{
				ID: ulid.MakeUlid().String(),
			},
			AkunID: trByr.AkunBayarID,
			Debit:  trByr.Total,
			Saldo:  trByr.Total,
		},
		// ay kredit, bayar mengurangi akun utang/piutangnya
		{
			Base: entity.Base{
				ID: ulid.MakeUlid().String(),
			},
			AkunID: ayTagihan.AkunID,
			Debit:  ayTagihan.Debit,
			Kredit: ayTagihan.Kredit,
			Saldo:  ayTagihan.Saldo,
		},
	}

	// cek jika total bayarnya lebih besar dari sisa tagihan
	if trByr.Total > math.Abs(ayTagihan.Saldo) {
		return nil, errors.New(message.BayarMustLessThanSisaTagihan)

	}

	byrHP := entity.DataBayarHutangPiutang{
		Base: entity.Base{
			ID: ulid.MakeUlid().String(),
		},
		Total: trByr.Total,
		Transaksi: entity.Transaksi{
			Base: entity.Base{
				ID: ulid.MakeUlid().String(),
			},
			Keterangan:      keterangan,
			BuktiPembayaran: trByr.BuktiPembayaran,
			Total:           trByr.Total,
			Tanggal:         tanggal,
			KontakID:        kontakId,
			AyatJurnals:     ayatJurnals,
		},
	}

	return &byrHP, nil
}

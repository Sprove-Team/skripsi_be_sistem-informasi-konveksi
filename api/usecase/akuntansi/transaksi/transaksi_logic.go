package akuntansi

import (
	"errors"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/transaksi"
	"github.com/be-sistem-informasi-konveksi/entity"
)

func isValidAkunHp(hp *entity.HutangPiutang, akun *entity.Akun, ay *req.ReqAyatJurnal) error {
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

func isValidAkunByrHP(hpFromByr *entity.DataBayarHutangPiutang, akun *entity.Akun, ay *req.ReqAyatJurnal) error {

	// untuk tr bayr hp
	if hpFromByr.HutangPiutang.ID != "" && ay.Debit != 0 {
		if hpFromByr.HutangPiutang.Jenis == "PIUTANG" {
			if ay.Kredit != 0 && akun.KelompokAkun.Nama != "piutang" {
				return errors.New(message.AkunNotMatchWithJenisHPTr)
			}
			if ay.Debit != 0 && akun.KelompokAkun.Nama != "kas & bank" {
				return errors.New(message.InvalidAkunBayar)
			}
		}
		if hpFromByr.HutangPiutang.Jenis == "HUTANG" {
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

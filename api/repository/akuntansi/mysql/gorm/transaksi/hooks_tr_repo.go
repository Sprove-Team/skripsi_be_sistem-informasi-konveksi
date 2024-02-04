package akuntansi

import (
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

func (r *transaksiRepo) OnDeleteTrDataBayarHP(tx *gorm.DB, tr *entity.Transaksi) error {

	var hp entity.HutangPiutang
	if err := tx.First(&hp, "id = ?", tr.DataBayarHutangPiutang.HutangPiutangID).Error; err != nil {
		return err
	}

	hp.Sisa += tr.Total
	if hp.Sisa > 0 && hp.Status != "BELUM_LUNAS" {
		hp.Status = "BELUM_LUNAS"
	}
	if err := tx.Select("sisa", "status").Updates(&hp).Error; err != nil {
		return err
	}

	return nil
}

func (r *transaksiRepo) OnDeleteTrHP(tx *gorm.DB, tr *entity.Transaksi) error {

	var dataByrHp []entity.DataBayarHutangPiutang
	err := tx.Find(&dataByrHp, "hutang_piutang_id = ?", tr.HutangPiutang.ID).Error

	if err != nil {
		helper.LogsError(err)
		return err
	}

	lengthByr := len(dataByrHp)
	if lengthByr > 0 {
		var idsTrByr = make([]string, lengthByr)
		for i, v := range dataByrHp {
			idsTrByr[i] = v.TransaksiID
		}
		var errs = make(chan error, 3)
		go func() {
			if err := tx.Delete(&entity.Transaksi{}, "id IN (?)", idsTrByr).Error; err != nil {
				errs <- err
			}
			errs <- nil
		}()
		go func() {
			if err := tx.Delete(&entity.AyatJurnal{}, "transaksi_id IN (?)", idsTrByr).Error; err != nil {
				errs <- err
			}
			errs <- nil
		}()

		go func() {
			if err := tx.Delete(&dataByrHp).Error; err != nil {
				errs <- err
			}
			errs <- nil
		}()

		for i := 0; i < 3; i++ {
			if err := <-errs; err != nil {
				return err
			}
		}

	}

	return nil
}

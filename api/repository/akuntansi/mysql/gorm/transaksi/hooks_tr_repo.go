package repo_akuntansi_transaksi

import (
	"github.com/be-sistem-informasi-konveksi/entity"
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

	lengthReq := 2
	errs := make(chan error, lengthReq)
	subQuery := tx.Model(&entity.DataBayarHutangPiutang{}).Where("hutang_piutang_id = ?", tr.HutangPiutang.ID).Select("transaksi_id")
	go func() {
		if err := tx.Delete(&entity.Transaksi{}, "id IN (?)", subQuery).Error; err != nil {
			errs <- err
			return
		}
		errs <- nil
	}()
	go func() {
		if err := tx.Delete(&entity.AyatJurnal{}, "transaksi_id IN (?)", subQuery).Error; err != nil {
			errs <- err
			return
		}
		errs <- nil
	}()

	for i := 0; i < lengthReq; i++ {
		if err := <-errs; err != nil {
			return err
		}
	}

	if err := tx.Delete(&entity.DataBayarHutangPiutang{}, "hutang_piutang_id = ?", tr.HutangPiutang.ID).Error; err != nil {
		return err
	}

	return nil
}

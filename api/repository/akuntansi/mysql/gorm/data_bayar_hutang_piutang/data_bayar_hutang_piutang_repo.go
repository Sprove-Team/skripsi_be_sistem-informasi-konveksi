package repo_akuntansi_data_bayar_hp

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type DataBayarHutangPiutangRepo interface {
	Create(ctx context.Context, dataBayar *entity.DataBayarHutangPiutang) error
	GetByTrId(ctx context.Context, id string) (entity.DataBayarHutangPiutang, error)
}

type dataBayarHutangPiutangRepo struct {
	DB *gorm.DB
}

func NewDataBayarHutangPiutangRepo(DB *gorm.DB) DataBayarHutangPiutangRepo {
	return &dataBayarHutangPiutangRepo{DB}
}

func (r *dataBayarHutangPiutangRepo) GetByTrId(ctx context.Context, id string) (entity.DataBayarHutangPiutang, error) {
	data := entity.DataBayarHutangPiutang{}
	err := r.DB.WithContext(ctx).Model(&entity.DataBayarHutangPiutang{}).Preload("Transaksi").Preload("HutangPiutang").First(&data, "transaksi_id = ?", id).Error

	if err != nil {
		helper.LogsError(err)
		return data, err
	}

	return data, nil
}
func (r *dataBayarHutangPiutangRepo) Create(ctx context.Context, dataBayar *entity.DataBayarHutangPiutang) error {
	tx := r.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(dataBayar).Error; err != nil {
			return err
		}
		if err := tx.Select("sisa", "status").Updates(&dataBayar.HutangPiutang).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

package akuntansi

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type DataBayarHutangPiutangRepo interface {
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

package akuntansi

import (
	"context"
	"time"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type CreateParam struct {
	Transaksi   *entity.Transaksi
	AyatJurnals []*entity.AyatJurnal
	UpdateAkuns []*entity.Akun
}

type SearchTransaksi struct {
	Nama      string
	StartDate time.Time
	EndDate   time.Time
	Next      string
	Limit     int
	// Offset           int
}

type TransaksiRepo interface {
	Create(ctx context.Context, param CreateParam) error
	GetAll(ctx context.Context, param SearchTransaksi) ([]entity.Transaksi, error)
	// Add more methods as needed for your repository operations
}

type transaksiRepo struct {
	DB *gorm.DB
}

func NewTransaksiRepo(DB *gorm.DB) TransaksiRepo {
	return &transaksiRepo{DB}
}

func (r *transaksiRepo) Create(ctx context.Context, param CreateParam) error {
	tx := r.DB.WithContext(ctx)

	err := tx.Transaction(func(tx *gorm.DB) error {
		tx = tx.Session(&gorm.Session{
			NewDB: true,
		})
		if err := tx.Create(param.Transaksi).Error; err != nil {
			helper.LogsError(err)
			return err
		}

		if err := tx.Model(&entity.AyatJurnal{}).Create(param.AyatJurnals).Error; err != nil {
			helper.LogsError(err)
			return err
		}
		for _, akun := range param.UpdateAkuns {
			if err := tx.Model(&entity.Akun{}).Where("id = ?", akun.ID).Update("saldo_akhir", akun.SaldoAkhir).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *transaksiRepo) GetAll(ctx context.Context, param SearchTransaksi) ([]entity.Transaksi, error) {
	datas := []entity.Transaksi{}
	tx := r.DB.WithContext(ctx).Model(&datas).Order("transaksi.id ASC")

	if param.Next != "" {
		tx = tx.Where("transaksi.id > ?", param.Next)
	}

	tx = tx.Where("created_at >= ? AND created_at < ?", param.StartDate, param.EndDate)

	tx = tx.InnerJoins("AyatJurnals").Preload("AyatJurnals")

	err := tx.Limit(param.Limit).Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return datas, err
	}
	return datas, nil
}

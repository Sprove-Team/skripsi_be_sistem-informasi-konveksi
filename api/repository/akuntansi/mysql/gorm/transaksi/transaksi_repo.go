package akuntansi

import (
	"context"
	"fmt"
	"strings"
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
	StartDate time.Time
	EndDate   time.Time
}

type DeleteTransaksi struct {
	ID              string
	SaldoAkunValues []string
}

type TransaksiRepo interface {
	Create(ctx context.Context, param CreateParam) error
	GetAll(ctx context.Context, param SearchTransaksi) ([]entity.Transaksi, error)
	Delete(ctx context.Context, param DeleteTransaksi) error
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
			if err := tx.Model(&entity.Akun{}).Where("id = ?", akun.ID).Update("saldo", akun.Saldo).Error; err != nil {
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

func (r *transaksiRepo) Delete(ctx context.Context, param DeleteTransaksi) error {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&entity.Transaksi{}, "id = ?", param.ID).Error; err != nil {
			return err
		}
		newValueAkun := strings.Join(param.SaldoAkunValues, ",")
		query := fmt.Sprintf("INSERT INTO akun (id, saldo, golongan_akun_id, nama, kode) VALUES %s ON DUPLICATE KEY UPDATE saldo = VALUES(saldo)", newValueAkun)
		if err := tx.Exec(query).Error; err != nil {
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

func (r *transaksiRepo) GetAll(ctx context.Context, param SearchTransaksi) ([]entity.Transaksi, error) {
	datas := []entity.Transaksi{}
	tx := r.DB.WithContext(ctx).Model(&datas).Order("transaksi.id ASC")

	tx = tx.Where("DATE(tanggal) >= ? AND DATE(tanggal) <= ?", param.StartDate, param.EndDate)

	err := tx.Preload("AyatJurnals").Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return datas, err
	}
	return datas, nil
}

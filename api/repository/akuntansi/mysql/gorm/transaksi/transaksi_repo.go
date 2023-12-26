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

type DeleteParam struct {
	ID              string
	SaldoAkunValues []string
}

type UpdateParam struct {
	UpdateTr           *entity.Transaksi
	NewAyatJurnals     []*entity.AyatJurnal
	NewSaldoAkunValues []string
}

type SearchTransaksi struct {
	StartDate time.Time
	EndDate   time.Time
}

type TransaksiRepo interface {
	Create(ctx context.Context, param CreateParam) error
	Update(ctx context.Context, param UpdateParam) error
	GetAll(ctx context.Context, param SearchTransaksi) ([]entity.Transaksi, error)
	Delete(ctx context.Context, param DeleteParam) error
	GetById(ctx context.Context, id string) (entity.Transaksi, error)
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

func (r *transaksiRepo) Update(ctx context.Context, param UpdateParam) error {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// if err := tx.Session(&gorm.Session{
		// 	FullSaveAssociations: true,
		// }).Updates(param.UpdateTr).Error; err != nil {
		// 	return err
		// }

		if err := tx.Where("id =? ", param.UpdateTr.ID).Updates(param.UpdateTr).Error; err != nil {
			helper.LogsError(err)
			return err
		}
		if param.NewAyatJurnals != nil {
			if err := tx.Unscoped().Where("transaksi_id IN (?)", param.UpdateTr.ID).Delete(&entity.AyatJurnal{}).Error; err != nil {
				helper.LogsError(err)
				return err
			}
			if err := tx.Model(&entity.AyatJurnal{}).Create(param.NewAyatJurnals).Error; err != nil {
				helper.LogsError(err)
				return err
			}
		}

		// update multiple lines with duplicate IDs in the 'akun' table.
		// the SQL query uses INSERT INTO ... ON DUPLICATE KEY UPDATE to efficiently handle duplicates.

		newValueAkun := strings.Join(param.NewSaldoAkunValues, ",")
		query2 := fmt.Sprintf("INSERT INTO akun (id, saldo, golongan_akun_id, nama, kode) VALUES %s ON DUPLICATE KEY UPDATE saldo = VALUES(saldo)", newValueAkun)
		if err := tx.Exec(query2).Error; err != nil {
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

func (r *transaksiRepo) Delete(ctx context.Context, param DeleteParam) error {
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

	tx := r.DB.WithContext(ctx).Model(&datas).Order("tanggal ASC").Omit("created_at", "deleted_at", "updated_at")

	tx = tx.Where("DATE(tanggal) >= ? AND DATE(tanggal) <= ?", param.StartDate, param.EndDate)

	err := tx.Preload("AyatJurnals").
		Preload("AyatJurnals.Akun", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "deleted_at", "updated_at").
				Select("nama", "id", "saldo_normal", "saldo", "kode", "golongan_akun_id")
		}).
		Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return datas, err
	}
	return datas, nil
}

func (r *transaksiRepo) GetById(ctx context.Context, id string) (entity.Transaksi, error) {
	data := entity.Transaksi{}
	err := r.DB.WithContext(ctx).Model(&entity.Transaksi{}).Omit("created_at", "deleted_at", "updated_at").Where("id = ?", id).Preload("AyatJurnals", func(db *gorm.DB) *gorm.DB {
		return db.Omit("created_at", "deleted_at", "updated_at").Preload("Akun", func(db2 *gorm.DB) *gorm.DB {
			return db2.Select("nama", "id", "saldo_normal", "saldo", "kode", "golongan_akun_id")
		})
	}).First(&data).Error
	if err != nil {
		helper.LogsError(err)
		return entity.Transaksi{}, err
	}

	return data, nil
}

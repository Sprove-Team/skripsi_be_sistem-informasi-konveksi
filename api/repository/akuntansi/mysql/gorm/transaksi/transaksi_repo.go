package akuntansi

import (
	"context"
	"time"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CreateParam struct {
	Transaksi   *entity.Transaksi
	AyatJurnals []*entity.AyatJurnal
}

type UpdateParam struct {
	UpdateTr                     *entity.Transaksi
	NewAyatJurnals               []*entity.AyatJurnal
	UpdateHutangPiutang          *entity.HutangPiutang
	UpdateDataBayarHutangPiutang *entity.DataBayarHutangPiutang
}

type SearchTransaksi struct {
	StartDate time.Time
	EndDate   time.Time
}

type TransaksiRepo interface {
	Create(ctx context.Context, param CreateParam) error
	Update(ctx context.Context, param UpdateParam) error
	GetAll(ctx context.Context, param SearchTransaksi) ([]entity.Transaksi, error)
	GetHistory(ctx context.Context, param SearchTransaksi) ([]entity.Transaksi, error)
	Delete(ctx context.Context, id string) error
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
		// update transaksi
		if err := tx.Where("id =? ", param.UpdateTr.ID).Updates(param.UpdateTr).Error; err != nil {
			helper.LogsError(err)
			return err
		}

		if param.UpdateHutangPiutang != nil {
			if err := tx.Select("status", "sisa", "total").Updates(param.UpdateHutangPiutang).Error; err != nil {
				helper.LogsError(err)
				return err
			}
		}

		if param.UpdateDataBayarHutangPiutang != nil {
			if err := tx.Updates(param.UpdateDataBayarHutangPiutang).Error; err != nil {
				helper.LogsError(err)
				return err
			}
		}
		// update ayat jurnals
		if len(param.NewAyatJurnals) > 0 {
			if err := tx.Unscoped().Where("transaksi_id IN (?)", param.UpdateTr.ID).Delete(&entity.AyatJurnal{}).Error; err != nil {
				helper.LogsError(err)
				return err
			}

			if err := tx.Model(&entity.AyatJurnal{}).Create(param.NewAyatJurnals).Error; err != nil {
				helper.LogsError(err)
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

func (r *transaksiRepo) Delete(ctx context.Context, id string) error {

	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var transaksi entity.Transaksi
		// delete transaksi by id
		if err := tx.Preload("HutangPiutang").Preload("DataBayarHutangPiutang").First(&transaksi, "id = ?", id).Error; err != nil {
			helper.LogsError(err)
			return err
		}

		if transaksi.DataBayarHutangPiutang != nil {
			if err := r.OnDeleteTrDataBayarHP(tx, &transaksi); err != nil {
				return err
			}
		}

		if err := tx.Select(clause.Associations).Delete(&transaksi).Error; err != nil {
			helper.LogsError(err)
			return err
		}

		if transaksi.HutangPiutang != nil {
			if err := r.OnDeleteTrHP(tx, &transaksi); err != nil {
				return err
			}

		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *transaksiRepo) GetAll(ctx context.Context, param SearchTransaksi) ([]entity.Transaksi, error) {
	datas := []entity.Transaksi{}

	tx := r.DB.WithContext(ctx).Model(&datas).Order("tanggal DESC").Omit("created_at", "deleted_at", "updated_at", "bukti_pembayaran")

	err := tx.Where("DATE(tanggal) >= ? AND DATE(tanggal) <= ?", param.StartDate, param.EndDate).Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return datas, err
	}
	return datas, nil
}

func (r *transaksiRepo) GetHistory(ctx context.Context, param SearchTransaksi) ([]entity.Transaksi, error) {
	datas := []entity.Transaksi{}

	tx := r.DB.WithContext(ctx).Model(&datas).Unscoped().Order("tanggal DESC").Omit("created_at", "deleted_at", "updated_at", "bukti_pembayaran")

	err := tx.Where("DATE(tanggal) >= ? AND DATE(tanggal) <= ?", param.StartDate, param.EndDate).Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return datas, err
	}
	return datas, nil
}

func (r *transaksiRepo) GetById(ctx context.Context, id string) (entity.Transaksi, error) {
	data := entity.Transaksi{}
	err := r.DB.WithContext(ctx).Model(&entity.Transaksi{}).Omit("created_at", "deleted_at", "updated_at").
		Where("id = ?", id).
		Preload("AyatJurnals", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "deleted_at", "updated_at").Preload("Akun", func(db2 *gorm.DB) *gorm.DB {
				return db2.Select("nama", "id", "saldo_normal", "kode", "kelompok_akun_id")
			})
		}).First(&data).Error
	if err != nil {
		helper.LogsError(err)
		return entity.Transaksi{}, err
	}

	return data, nil
}

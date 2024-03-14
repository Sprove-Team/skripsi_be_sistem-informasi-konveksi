package repo_tugas

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	ParamCreate struct {
		Ctx   context.Context
		Tugas *entity.Tugas
	}
	ParamUpdate struct {
		Ctx      context.Context
		NewTugas *entity.Tugas
	}
	ParamGetById struct {
		Ctx context.Context
		ID  string
	}
	ParamDelete struct {
		Ctx context.Context
		ID  string
	}
	ParamGetAll struct {
		Ctx   context.Context
		Tahun uint
		Bulan uint
		Limit int
		Next  string
	}
	ParamGetByInvoiceId struct {
		Ctx       context.Context
		InvoiceID string
	}
)

type TugasRepo interface {
	Create(param ParamCreate) error
	Update(param ParamUpdate) error
	Delete(param ParamDelete) error
	GetById(param ParamGetById) (*entity.Tugas, error)
	GetAll(param ParamGetAll) ([]entity.Tugas, error)
	GetByInvoiceId(param ParamGetByInvoiceId) ([]entity.Tugas, error)
}

type tugasRepo struct {
	DB *gorm.DB
}

func NewTugasRepo(DB *gorm.DB) TugasRepo {
	return &tugasRepo{DB}
}

func (r *tugasRepo) Create(param ParamCreate) error {
	err := r.DB.WithContext(param.Ctx).Create(param.Tugas).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *tugasRepo) Update(param ParamUpdate) error {
	err := r.DB.WithContext(param.Ctx).Transaction(func(tx *gorm.DB) error {
		if len(param.NewTugas.Users) > 0 {
			if err := tx.Model(param.NewTugas).Association("Users").Replace(param.NewTugas.Users); err != nil {
				return err
			}
		}
		if err := tx.Updates(param.NewTugas).Error; err != nil {
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

func (r *tugasRepo) GetById(param ParamGetById) (*entity.Tugas, error) {
	data := new(entity.Tugas)
	err := r.DB.WithContext(param.Ctx).Preload("Invoice").
		Preload("JenisSpv").
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nama", "username", "role", "jenis_spv_id")
		}).
		Preload("Users.JenisSpv").
		Preload("SubTugas").
		First(data, "id = ?", param.ID).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return data, nil
}

func (r *tugasRepo) Delete(param ParamDelete) error {
	err := r.DB.WithContext(param.Ctx).Transaction(func(tx *gorm.DB) error {
		data := new(entity.Tugas)
		if err := tx.First(data, "id = ?", param.ID).Error; err != nil {
			return err
		}
		if err := tx.Select(clause.Associations).Delete(data).Error; err != nil {
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

func (r *tugasRepo) GetAll(param ParamGetAll) ([]entity.Tugas, error) {
	datas := []entity.Tugas{}
	tx := r.DB.WithContext(param.Ctx).Model(&datas)

	if param.Bulan != 0 {
		tx = tx.Where("MONTH(tanggal_deadline) = ?", param.Bulan)
	}
	if param.Tahun != 0 {
		tx = tx.Where("YEAR(tanggal_deadline) = ?", param.Tahun)
	}

	err := tx.
		Preload("Invoice").
		Preload("JenisSpv").
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nama", "username", "role", "jenis_spv_id")
		}).
		Preload("Users.JenisSpv").
		Find(&datas).Error

	if err != nil {
		helper.LogsError(err)
		return nil, err
	}

	return datas, nil
}

func (r *tugasRepo) GetByInvoiceId(param ParamGetByInvoiceId) ([]entity.Tugas, error) {
	datas := []entity.Tugas{}
	err := r.DB.WithContext(param.Ctx).
		Model(&datas).
		Where("invoice_id = ?", param.InvoiceID).
		Preload("JenisSpv").
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nama", "username", "role", "jenis_spv_id")
		}).
		Preload("Users.JenisSpv").
		Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, nil
}

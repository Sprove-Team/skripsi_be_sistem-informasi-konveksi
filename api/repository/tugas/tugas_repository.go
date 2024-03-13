package repo_tugas

import (
	"context"
	"time"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type SearchFilter struct {
	Tahun *time.Time
	Bulan *time.Time
}

type (
	ParamCreate struct {
		Ctx   context.Context
		Tugas *entity.Tugas
	}
	ParamGet struct {
		Ctx          context.Context
		SearchFilter SearchFilter
	}
	ParamGetByInvoiceId struct {
		Ctx       context.Context
		InvoiceID string
	}
)

type TugasRepo interface {
	Create(param ParamCreate) error
	// Get(param ParamGet) ([]entity.Tugas, error)
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

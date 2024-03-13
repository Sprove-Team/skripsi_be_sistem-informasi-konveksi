package repo_tugas

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type (
	ParamCreate struct {
		Ctx   context.Context
		Tugas *entity.Tugas
	}
)

type TugasRepo interface {
	Create(param ParamCreate) error
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

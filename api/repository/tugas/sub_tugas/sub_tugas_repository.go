package repo_sub_tugas

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type (
	ParamGetById struct {
		Ctx context.Context
		ID  string
	}
	ParamCreateByTugasId struct {
		Ctx      context.Context
		SubTugas *entity.SubTugas
	}
	ParamUpdate struct {
		Ctx         context.Context
		NewSubTugas *entity.SubTugas
	}
	ParamDelete struct {
		Ctx context.Context
		ID  string
	}
)

type SubTugasRepo interface {
	GetById(param ParamGetById) (*entity.SubTugas, error)
	CreateByTugasId(param ParamCreateByTugasId) error
	Update(param ParamUpdate) error
	Delete(param ParamDelete) error
}

type subTugasRepo struct {
	DB *gorm.DB
}

func NewRepoSubTugasRepo(DB *gorm.DB) SubTugasRepo {
	return &subTugasRepo{DB}
}

func (r *subTugasRepo) GetById(param ParamGetById) (*entity.SubTugas, error) {
	data := new(entity.SubTugas)
	err := r.DB.WithContext(param.Ctx).First(data, "id = ?", param.ID).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *subTugasRepo) CreateByTugasId(param ParamCreateByTugasId) error {
	err := r.DB.WithContext(param.Ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&entity.Tugas{}, "id = ?", param.SubTugas.TugasID).Error; err != nil {
			return err
		}
		if err := tx.Create(param.SubTugas).Error; err != nil {
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

func (r *subTugasRepo) Update(param ParamUpdate) error {
	err := r.DB.WithContext(param.Ctx).Omit("id").Updates(param.NewSubTugas).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *subTugasRepo) Delete(param ParamDelete) error {
	err := r.DB.WithContext(param.Ctx).Delete(new(entity.SubTugas), "id = ?", param.ID).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

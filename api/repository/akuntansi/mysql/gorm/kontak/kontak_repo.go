package akuntansi

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type (
	ParamCreate struct {
		Ctx    context.Context
		Kontak *entity.Kontak
	}
	ParamDelete struct {
		Ctx context.Context
		ID  string
	}
	ParamUpdate struct {
		Ctx    context.Context
		Kontak *entity.Kontak
	}
	ParamGetById struct {
		Ctx context.Context
		ID  string
	}
	ParamGetAll struct {
		Ctx    context.Context
		Limit  int
		Nama   string
		Email  string
		NoTelp string
		Next   string
	}
)

type KontakRepo interface {
	Create(param ParamCreate) error
	Delete(param ParamDelete) error
	Update(param ParamUpdate) error
	GetById(param ParamGetById) (entity.Kontak, error)
	GetAll(param ParamGetAll) ([]entity.Kontak, error)
}

type kontakRepo struct {
	DB *gorm.DB
}

func NewKontakRepo(DB *gorm.DB) KontakRepo {
	return &kontakRepo{DB}
}

func (r *kontakRepo) Create(param ParamCreate) error {
	if err := r.DB.WithContext(param.Ctx).Create(param.Kontak).Error; err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *kontakRepo) Delete(param ParamDelete) error {
	if err := r.DB.WithContext(param.Ctx).Delete(&entity.Kontak{}, "id = ?", param.ID).Error; err != nil {
		return err
	}
	return nil
}

func (r *kontakRepo) Update(param ParamUpdate) error {
	if err := r.DB.WithContext(param.Ctx).Model(&entity.Kontak{}).Where("id = ?", param.Kontak.ID).Updates(param.Kontak).Error; err != nil {
		return err
	}
	return nil
}

func (r *kontakRepo) GetById(param ParamGetById) (entity.Kontak, error) {
	data := entity.Kontak{}
	if err := r.DB.WithContext(param.Ctx).Where("id = ?", param.ID).First(&data).Error; err != nil {
		return entity.Kontak{}, err
	}
	return data, nil
}

func (r *kontakRepo) GetAll(param ParamGetAll) ([]entity.Kontak, error) {
	datas := []entity.Kontak{}

	tx := r.DB.WithContext(param.Ctx).Model(&datas).Order("id ASC")

	conditions := map[string]interface{}{
		"id > ?":         param.Next,
		"nama LIKE ?":    "%" + param.Nama + "%",
		"email LIKE ?":   param.Email + "%",
		"no_telp LIKE ?": param.NoTelp + "%",
	}

	for condition, value := range conditions {
		if value != "" {
			tx = tx.Where(condition, value)
		}
	}

	if err := tx.Limit(param.Limit).Find(&datas).Error; err != nil {
		return nil, err
	}

	return datas, nil

}

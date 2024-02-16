package user

import (
	"context"

	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
)

type JenisSpvRepo interface {
	Create(ctx context.Context, jenisSpv *entity.JenisSpv) error
	GetById(ctx context.Context, id string) (*entity.JenisSpv, error)
	GetByIds(ctx context.Context, ids []string) ([]entity.JenisSpv, error)
	GetAll(ctx context.Context) ([]entity.JenisSpv, error)
	Update(ctx context.Context, jenisSpv *entity.JenisSpv) error
	Delete(ctx context.Context, id string) error
}

type jenisSpvRepo struct {
	DB *gorm.DB
}

func NewJenisSpvRepo(DB *gorm.DB) JenisSpvRepo {
	return &jenisSpvRepo{DB}
}

func (r *jenisSpvRepo) Create(ctx context.Context, jenisSpv *entity.JenisSpv) error {
	return r.DB.WithContext(ctx).Create(jenisSpv).Error
}

func (r *jenisSpvRepo) Update(ctx context.Context, jenisSpv *entity.JenisSpv) error {
	return r.DB.WithContext(ctx).Omit("id").Updates(jenisSpv).Error
}

func (r *jenisSpvRepo) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&entity.JenisSpv{}, "id = ?", id).Error
}

func (r *jenisSpvRepo) GetById(ctx context.Context, id string) (*entity.JenisSpv, error) {
	data := entity.JenisSpv{}
	err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error
	return &data, err
}

func (r *jenisSpvRepo) GetByIds(ctx context.Context, ids []string) ([]entity.JenisSpv, error) {
	datas := make([]entity.JenisSpv, 0, len(ids))

	err := r.DB.WithContext(ctx).Model(&entity.JenisSpv{}).Where("id IN (?)", ids).Find(&datas).Error
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (r *jenisSpvRepo) GetAll(ctx context.Context) ([]entity.JenisSpv, error) {
	datas := []entity.JenisSpv{}
	err := r.DB.WithContext(ctx).Find(&datas).Error
	return datas, err
}

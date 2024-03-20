package repo_user_jenis_spv

import (
	"context"

	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
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
	if err := r.DB.WithContext(ctx).Omit("id").Updates(jenisSpv).Error; err != nil {
		if err != gorm.ErrDuplicatedKey {
			helper.LogsError(err)
		}
		return err
	}
	return nil
}

func (r *jenisSpvRepo) Delete(ctx context.Context, id string) error {
	if err := r.DB.WithContext(ctx).Delete(&entity.JenisSpv{}, "id = ?", id).Error; err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *jenisSpvRepo) GetById(ctx context.Context, id string) (*entity.JenisSpv, error) {
	data := entity.JenisSpv{}
	if err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			helper.LogsError(err)
		}
		return nil, err
	}
	return &data, nil
}

func (r *jenisSpvRepo) GetByIds(ctx context.Context, ids []string) ([]entity.JenisSpv, error) {
	datas := make([]entity.JenisSpv, 0, len(ids))

	err := r.DB.WithContext(ctx).Model(&entity.JenisSpv{}).Where("id IN (?)", ids).Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, nil
}

func (r *jenisSpvRepo) GetAll(ctx context.Context) ([]entity.JenisSpv, error) {
	datas := []entity.JenisSpv{}
	if err := r.DB.WithContext(ctx).Find(&datas).Error; err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, nil
}

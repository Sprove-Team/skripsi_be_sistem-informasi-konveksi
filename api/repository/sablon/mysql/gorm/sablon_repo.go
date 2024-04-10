package repo_sablon

import (
	"context"

	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
)

type SablonRepo interface {
	Create(ctx context.Context, sablon *entity.Sablon) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, sablon *entity.Sablon) error
	GetById(ctx context.Context, id string) (entity.Sablon, error)
	GetByIds(ctx context.Context, ids []string) ([]entity.Sablon, error)
	GetAll(ctx context.Context, param SearchSablon) ([]entity.Sablon, error)
}

type sablonRepo struct {
	DB *gorm.DB
}

func NewSablonRepo(DB *gorm.DB) SablonRepo {
	return &sablonRepo{DB}
}

func (r *sablonRepo) Create(ctx context.Context, sablon *entity.Sablon) error {
	err := r.DB.WithContext(ctx).Create(sablon).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *sablonRepo) Delete(ctx context.Context, id string) error {
	err := r.DB.WithContext(ctx).Delete(&entity.Sablon{}, "id = ?", id).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *sablonRepo) Update(ctx context.Context, sablon *entity.Sablon) error {
	err := r.DB.WithContext(ctx).Updates(sablon).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *sablonRepo) GetById(ctx context.Context, id string) (entity.Sablon, error) {
	data := entity.Sablon{}
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	if err != nil {
		helper.LogsError(err)
		return data, err
	}
	return data, nil
}

func (r *sablonRepo) GetByIds(ctx context.Context, ids []string) ([]entity.Sablon, error) {
	datas := make([]entity.Sablon, 0, len(ids))
	err := r.DB.WithContext(ctx).Where("id IN (?)", ids).Find(&datas).Error
	if err != nil {
		return nil, err
	}
	return datas, nil
}

type SearchSablon struct {
	Nama  string
	Limit int
	Next  string
}

func (r *sablonRepo) GetAll(ctx context.Context, param SearchSablon) ([]entity.Sablon, error) {
	datas := []entity.Sablon{}
	// var totalData int64

	tx := r.DB.WithContext(ctx).Model(&entity.Sablon{}).Order("id ASC")
	if param.Next != "" {
		tx = tx.Where("id > ?", param.Next)
	}
	tx = tx.Where("nama LIKE ?", "%"+param.Nama+"%")

	err := tx.Limit(param.Limit).Find(&datas).Error

	if err != nil {
		helper.LogsError(err)
		return datas, err
	}
	return datas, nil
}

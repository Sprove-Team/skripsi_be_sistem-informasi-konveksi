package bordir

import (
	"context"

	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
)

type BordirRepo interface {
	GetAll(ctx context.Context, param SearchBordir) ([]entity.Bordir, error)
	GetById(ctx context.Context, id string) (entity.Bordir, error)
	GetByIds(ctx context.Context, ids []string) ([]entity.Bordir, error)
	Create(ctx context.Context, bordir *entity.Bordir) error
	Update(ctx context.Context, bordir *entity.Bordir) error
	Delete(ctx context.Context, id string) error
}

type bordirRepo struct {
	DB *gorm.DB
}

func NewBordirRepo(DB *gorm.DB) BordirRepo {
	return &bordirRepo{DB}
}

func (r *bordirRepo) Create(ctx context.Context, bordir *entity.Bordir) error {
	return r.DB.WithContext(ctx).Create(bordir).Error
}

func (r *bordirRepo) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&entity.Bordir{}, "id = ?", id).Error
}

func (r *bordirRepo) Update(ctx context.Context, bordir *entity.Bordir) error {
	return r.DB.WithContext(ctx).Omit("id").Updates(bordir).Error
}

func (r *bordirRepo) GetById(ctx context.Context, id string) (entity.Bordir, error) {
	data := entity.Bordir{}
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	if err != nil {
		helper.LogsError(err)
		return entity.Bordir{}, err
	}
	return data, err
}

func (r *bordirRepo) GetByIds(ctx context.Context, ids []string) ([]entity.Bordir, error) {
	datas := make([]entity.Bordir, 0, len(ids))
	err := r.DB.WithContext(ctx).Where("id IN (?)", ids).Find(&datas).Error
	if err != nil {
		return nil, err
	}
	return datas, nil
}

type SearchBordir struct {
	Nama  string
	Next  string
	Limit int
}

func (r *bordirRepo) GetAll(ctx context.Context, param SearchBordir) ([]entity.Bordir, error) {
	datas := []entity.Bordir{}

	tx := r.DB.WithContext(ctx).Model(&entity.Bordir{}).Order("id ASC")
	if param.Next != "" {
		tx = tx.Where("id > ?", param.Next)
	}
	tx = tx.Where("nama LIKE ?", "%"+param.Nama+"%")
	err := tx.Limit(param.Limit).Find(&datas).Error
	return datas, err
}

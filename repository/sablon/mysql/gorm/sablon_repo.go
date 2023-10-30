package sablon

import (
	"context"

	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
)

type SablonRepo interface {
	GetAll(ctx context.Context, param SearchParams) ([]entity.Sablon, int64, error)
	GetById(ctx context.Context, id string) (entity.Sablon, error)
	Create(ctx context.Context, sablon *entity.Sablon) error
	Update(ctx context.Context, sablon *entity.Sablon) error
	Delete(ctx context.Context, id string) error
}

type sablonRepo struct {
	DB *gorm.DB
}

func NewSablonRepo(DB *gorm.DB) SablonRepo {
	return &sablonRepo{DB}
}

func (r *sablonRepo) Create(ctx context.Context, sablon *entity.Sablon) error {
	return r.DB.WithContext(ctx).Create(sablon).Error
}

func (r *sablonRepo) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&entity.Sablon{}, "id = ?", id).Error
}

func (r *sablonRepo) Update(ctx context.Context, sablon *entity.Sablon) error {
	return r.DB.WithContext(ctx).Updates(sablon).Error
}

func (r *sablonRepo) GetById(ctx context.Context, id string) (entity.Sablon, error) {
	data := entity.Sablon{}
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

type SearchParams struct {
	Nama   string
	Limit  int
	Offset int
}

func (r *sablonRepo) GetAll(ctx context.Context, param SearchParams) ([]entity.Sablon, int64, error) {
	datas := []entity.Sablon{}
	var totalData int64

	tx := r.DB.WithContext(ctx).Model(&entity.Sablon{})
	tx = tx.Where("nama LIKE ?", "%"+param.Nama+"%")
	err := tx.Count(&totalData).Limit(param.Limit).Offset(param.Offset).Find(&datas).Error
	return datas, totalData, err
}

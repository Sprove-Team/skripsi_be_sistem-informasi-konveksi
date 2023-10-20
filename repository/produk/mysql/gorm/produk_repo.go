package produk

import (
	"context"

	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
)

type ProdukRepo interface {
	GetById(ctx context.Context, id string) (entity.Produk, error)
	Create(produk *entity.Produk) error
	Update(produk *entity.Produk) error
	Delete(id string) error
}

type produkRepo struct {
	DB *gorm.DB
}

func NewProdukRepo(DB *gorm.DB) ProdukRepo {
	return &produkRepo{DB}
}

func (r *produkRepo) Create(produk *entity.Produk) error {
	return r.DB.Create(&produk).Error
}

func (r *produkRepo) Delete(id string) error {
	return r.DB.Delete(&entity.Produk{}, "id = ?", id).Error
}

func (r *produkRepo) Update(produk *entity.Produk) error {
	return r.DB.Updates(produk).Error
}

func (r *produkRepo) GetById(ctx context.Context, id string) (entity.Produk, error) {
	produkD := entity.Produk{}
	err := r.DB.WithContext(ctx).Model(&produkD).Where("id = ?", id).Preload("HargaDetails").First(&produkD).Error
	return produkD, err
}

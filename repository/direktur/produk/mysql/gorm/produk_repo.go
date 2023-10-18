package direktur

import (
	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
)

type ProdukRepo interface {
	Create(produk *entity.Produk) error
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

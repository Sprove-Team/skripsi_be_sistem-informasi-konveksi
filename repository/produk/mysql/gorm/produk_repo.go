package produk 

import (
	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
)

type ProdukRepo interface {
	GetById(id string) (entity.Produk, error)
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

func (r *produkRepo) GetById(id string) (entity.Produk, error) {
	produkD := entity.Produk{}
	err := r.DB.Preload("HargaDetails").First(&produkD, "id = ?", id).Error
	return produkD, err
}

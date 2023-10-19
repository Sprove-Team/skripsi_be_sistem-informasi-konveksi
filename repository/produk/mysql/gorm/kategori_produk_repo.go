package produk 

import (
	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
)

type KategoriProdukRepo interface {
	Create(kategori *entity.KategoriProduk) error
	Delete(id uint64) error
  GetById(id uint64) (entity.KategoriProduk, error)
}

type kategoriRepo struct {
	DB *gorm.DB
}

func NewKategoriProdukRepo(DB *gorm.DB) KategoriProdukRepo {
	return &kategoriRepo{DB}
}

func (r *kategoriRepo) Create(kategori *entity.KategoriProduk) error {
	return r.DB.Create(kategori).Error
}

func (r *kategoriRepo) GetById(id uint64) (entity.KategoriProduk, error) {
  var data entity.KategoriProduk
  err := r.DB.First(&data, "id = ?", id).Error
  return data, err
}

func (r *kategoriRepo) Delete(id uint64) error {
	return r.DB.Delete(&entity.KategoriProduk{}, "id = ?", id).Error
}

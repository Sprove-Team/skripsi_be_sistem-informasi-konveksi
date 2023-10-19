package direktur

import (
	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
)

type HargaDetailProdukRepo interface {
	Create(hargaDetailProduk *entity.HargaDetailProduk) error
}

type hargaDetailProdukRepo struct {
	DB *gorm.DB
}

func NewHargaDetailProdukRepo(DB *gorm.DB) HargaDetailProdukRepo {
	return &hargaDetailProdukRepo{DB}
}

func (r *hargaDetailProdukRepo) Create(hargaDetailProduk *entity.HargaDetailProduk) error {
	return r.DB.Create(hargaDetailProduk).Error
}

// func (r *hargaDetailProdukRepo) Create(hargaDetailProduk *entity.HargaDetailProduk) error {
// 	return r.DB.Create(hargaDetailProduk).Error
// }

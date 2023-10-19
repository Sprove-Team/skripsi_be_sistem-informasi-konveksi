package produk

import (
	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
)

type HargaDetailProdukRepo interface {
	GetByQty(qty uint) ([]entity.HargaDetailProduk, error)
	Create(hargaDetailProduk *entity.HargaDetailProduk) error
}

type hargaDetailProdukRepo struct {
	DB *gorm.DB
}

func NewHargaDetailProdukRepo(DB *gorm.DB) HargaDetailProdukRepo {
	return &hargaDetailProdukRepo{DB}
}

func (r *hargaDetailProdukRepo) Create(hargaDetailProduk *entity.HargaDetailProduk) error {
	return r.DB.Save(hargaDetailProduk).Error
}

func (r *hargaDetailProdukRepo) GetByQty(qty uint) ([]entity.HargaDetailProduk, error) {
	data := []entity.HargaDetailProduk{}
	err := r.DB.Find(&data, "qty = ?", qty).Error
	return data, err
}

// func (r *hargaDetailProdukRepo) Create(hargaDetailProduk *entity.HargaDetailProduk) error {
// 	return r.DB.Create(hargaDetailProduk).Error
// }

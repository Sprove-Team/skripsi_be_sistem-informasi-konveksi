package produk

import (
	"context"

	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
)

type HargaDetailProdukRepo interface {
	GetByQtyProdukId(ctx context.Context, qty uint, produkId string) (entity.HargaDetailProduk, error)
	GetById(ctx context.Context, id string) (entity.HargaDetailProduk, error)
	Delete(ctx context.Context, id string) error
	DeleteByProdukId(ctx context.Context, produk_id string) error
	UpdateById(ctx context.Context, hargaDetailProduk *entity.HargaDetailProduk) error
    GetByProdukId(ctx context.Context, id string) ([]entity.HargaDetailProduk, error)
	Create(ctx context.Context, hargaDetailProduk *entity.HargaDetailProduk) error
}

type hargaDetailProdukRepo struct {
	DB *gorm.DB
}

func NewHargaDetailProdukRepo(DB *gorm.DB) HargaDetailProdukRepo {
	return &hargaDetailProdukRepo{DB}
}

func (r *hargaDetailProdukRepo) Create(ctx context.Context, hargaDetailProduk *entity.HargaDetailProduk) error {
	return r.DB.WithContext(ctx).Create(hargaDetailProduk).Error
}

func (r *hargaDetailProdukRepo) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&entity.HargaDetailProduk{}, "id = ?", id).Error
}

func (r *hargaDetailProdukRepo) DeleteByProdukId(ctx context.Context, produk_id string) error {
	return r.DB.WithContext(ctx).Delete(&entity.HargaDetailProduk{}, "produk_id = ?", produk_id).Error
}

func (r *hargaDetailProdukRepo) UpdateById(ctx context.Context, hargaDetailProduk *entity.HargaDetailProduk) error {
	return r.DB.WithContext(ctx).Omit("id").Updates(hargaDetailProduk).Error
}

func (r *hargaDetailProdukRepo) GetByQtyProdukId(ctx context.Context, qty uint, produkId string) (entity.HargaDetailProduk, error) {
	data := entity.HargaDetailProduk{}
	err := r.DB.WithContext(ctx).Find(&data, "qty = ? AND produk_id = ?", qty, produkId).Error
	return data, err
}

func (r *hargaDetailProdukRepo) GetByProdukId(ctx context.Context, produk_id string) ([]entity.HargaDetailProduk, error) {
  datas := []entity.HargaDetailProduk{}
  err := r.DB.WithContext(ctx).Find(&datas,"produk_id = ?", produk_id).Error
  return datas, err
}

func (r *hargaDetailProdukRepo) GetById(ctx context.Context, id string) (entity.HargaDetailProduk, error) {
	var data entity.HargaDetailProduk
	err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error
	return data, err
}
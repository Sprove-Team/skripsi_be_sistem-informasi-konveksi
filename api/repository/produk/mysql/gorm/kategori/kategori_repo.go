package repo_produk_kategori

import (
	"context"

	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
)

type KategoriProdukRepo interface {
	Create(ctx context.Context, kategori *entity.KategoriProduk) error
	Update(ctx context.Context, kategori *entity.KategoriProduk) error
	Delete(ctx context.Context, id string) error
	GetById(ctx context.Context, id string) (entity.KategoriProduk, error)
	GetAll(ctx context.Context, param SearchKategoriProduk) ([]entity.KategoriProduk, error)
}

type kategoriRepo struct {
	DB *gorm.DB
}

func NewKategoriProdukRepo(DB *gorm.DB) KategoriProdukRepo {
	return &kategoriRepo{DB}
}

func (r *kategoriRepo) Create(ctx context.Context, kategori *entity.KategoriProduk) error {
	return r.DB.WithContext(ctx).Create(kategori).Error
}

func (r *kategoriRepo) Update(ctx context.Context, kategori *entity.KategoriProduk) error {
	return r.DB.WithContext(ctx).Omit("id").Updates(kategori).Error
}

func (r *kategoriRepo) GetById(ctx context.Context, id string) (entity.KategoriProduk, error) {
	var data entity.KategoriProduk
	err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error
	return data, err
}

func (r *kategoriRepo) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&entity.KategoriProduk{}, "id = ?", id).Error
}

type SearchKategoriProduk struct {
	Nama  string
	Next  string
	Limit int
}

func (r *kategoriRepo) GetAll(ctx context.Context, param SearchKategoriProduk) ([]entity.KategoriProduk, error) {
	datas := []entity.KategoriProduk{}
	tx := r.DB.WithContext(ctx).Model(&entity.KategoriProduk{}).Order("id ASC").Where("nama LIKE ?", "%"+param.Nama+"%")

	if param.Next != "" {
		tx = tx.Where("id > ?", param.Next)
	}

	err := tx.Limit(param.Limit).Find(&datas).Error

	return datas, err
}

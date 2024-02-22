package produk

import (
	"context"

	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
)

type ProdukRepo interface {
	Create(ctx context.Context, produk *entity.Produk) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, produk *entity.Produk) error
	GetById(ctx context.Context, id string) (entity.Produk, error)
	GetByIds(ctx context.Context, ids []string) ([]entity.Produk, error)
	GetAll(ctx context.Context, param SearchProduk) ([]entity.Produk, error)
}

type produkRepo struct {
	DB *gorm.DB
}

func NewProdukRepo(DB *gorm.DB) ProdukRepo {
	return &produkRepo{DB}
}

func (r *produkRepo) Create(ctx context.Context, produk *entity.Produk) error {
	return r.DB.WithContext(ctx).Create(produk).Error
}

func (r *produkRepo) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&entity.Produk{}, "id = ?", id).Error
}

func (r *produkRepo) Update(ctx context.Context, produk *entity.Produk) error {
	return r.DB.WithContext(ctx).Omit("id").Updates(produk).Error
}

func (r *produkRepo) GetById(ctx context.Context, id string) (entity.Produk, error) {
	data := entity.Produk{}
	err := r.DB.WithContext(ctx).Where("id = ?", id).Preload("HargaDetails", func(db *gorm.DB) *gorm.DB {
		return db.Order("qty ASC")
	}).Preload("KategoriProduk").First(&data).Error
	return data, err
}

func (r *produkRepo) GetByIds(ctx context.Context, ids []string) ([]entity.Produk, error) {
	datas := make([]entity.Produk, 0, len(ids))
	err := r.DB.WithContext(ctx).Where("id IN (?)", ids).Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, nil
}

type SearchProduk struct {
	Nama             string
	KategoriProdukId string
	HasHargaDetail   int
	Next             string
	Limit            int
}

func (r *produkRepo) GetAll(ctx context.Context, param SearchProduk) ([]entity.Produk, error) {
	datas := make([]entity.Produk, param.Limit)

	tx := r.DB.WithContext(ctx).Model(&entity.Produk{}).Order("produk.id ASC")

	if param.Next != "" {
		tx = tx.Where("produk.id > ?", param.Next)
	}

	tx = tx.Where("produk.nama LIKE ?", "%"+param.Nama+"%")

	if param.KategoriProdukId != "" {
		tx = tx.Where("kategori_produk_id = ?", param.KategoriProdukId)
	}

	switch param.HasHargaDetail {
	case 0:
		tx = tx.Joins("LEFT JOIN harga_detail_produk hd on hd.produk_id = produk.id")
	case 1:
		tx = tx.Joins("JOIN harga_detail_produk hd on hd.produk_id = produk.id")
	}

	tx = tx.Group("produk.id")

	err := tx.Limit(param.Limit).Preload("HargaDetails").Preload("KategoriProduk").Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, err
}

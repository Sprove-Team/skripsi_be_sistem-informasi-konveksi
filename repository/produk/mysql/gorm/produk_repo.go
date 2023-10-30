package produk

import (
	"context"

	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
)

type ProdukRepo interface {
	GetAll(ctx context.Context, param SearchProduk) ([]entity.Produk, int64, error)
	GetById(ctx context.Context, id string) (entity.Produk, error)
	Create(ctx context.Context, produk *entity.Produk) error
	Update(ctx context.Context, produk *entity.Produk) error
	Delete(ctx context.Context, id string) error
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
	err := r.DB.WithContext(ctx).Where("id = ?", id).Preload("HargaDetails").First(&data).Error
	return data, err
}

type SearchProduk struct {
	Nama             string
	KategoriProdukId string
	HasHargaDetail  bool 
	Limit            int
	Offset           int
}

func (r *produkRepo) GetAll(ctx context.Context, param SearchProduk) ([]entity.Produk, int64, error) {
	datas := []entity.Produk{}
	var totalData int64

	tx := r.DB.WithContext(ctx).Model(&entity.Produk{})
	tx = tx.Where("nama LIKE ?", "%"+param.Nama+"%")

	if param.KategoriProdukId != "" {
		tx = tx.Where("kategori_produk_id = ?", param.KategoriProdukId)
	}
	var err error
	if param.HasHargaDetail{
		err = tx.InnerJoins("HargaDetails").Preload("HargaDetails").Count(&totalData).Find(&datas).Error
	} else {
		err = tx.Preload("HargaDetails").Joins("LEFT JOIN harga_detail_produk hd on hd.produk_id = produk.id").Where("hd.id IS NULL").Count(&totalData).Find(&datas).Error
	}

	return datas, totalData, err
}

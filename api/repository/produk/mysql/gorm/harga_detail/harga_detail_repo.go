package produk

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
)

type HargaDetailProdukRepo interface {
	GetByQtyProdukId(ctx context.Context, qty uint, produkId string) (entity.HargaDetailProduk, error)
	GetByInQtyProdukId(ctx context.Context, qty []uint, produkId string) ([]entity.HargaDetailProduk, error)
	GetById(ctx context.Context, id string) (entity.HargaDetailProduk, error)
	Delete(ctx context.Context, id string) error
	DeleteByProdukId(ctx context.Context, produk_id string) error
	UpdateById(ctx context.Context, hargaDetailProduk *entity.HargaDetailProduk) error
	UpdateByProdukId(ctx context.Context, produkId string, hargaDetails []entity.HargaDetailProduk) error
	GetByProdukId(ctx context.Context, id string) ([]entity.HargaDetailProduk, error)
	Create(ctx context.Context, hargaDetailProduk []*entity.HargaDetailProduk) error
}

type hargaDetailProdukRepo struct {
	DB *gorm.DB
}

func NewHargaDetailProdukRepo(DB *gorm.DB) HargaDetailProdukRepo {
	return &hargaDetailProdukRepo{DB}
}

func (r *hargaDetailProdukRepo) Create(ctx context.Context, hargaDetailProduk []*entity.HargaDetailProduk) error {

	clauses := clause.OnConflict{
		Columns:   []clause.Column{{Name: "produk_id"}, {Name: "qty"}},
		DoUpdates: clause.AssignmentColumns([]string{"harga"}),
	}
	err := r.DB.WithContext(ctx).Clauses(clauses).Create(hargaDetailProduk).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *hargaDetailProdukRepo) Delete(ctx context.Context, id string) error {
	err := r.DB.WithContext(ctx).Delete(&entity.HargaDetailProduk{}, "id = ?", id).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *hargaDetailProdukRepo) DeleteByProdukId(ctx context.Context, produk_id string) error {
	err := r.DB.WithContext(ctx).Model(&entity.HargaDetailProduk{}).Where("produk_id = ?", produk_id).Delete(&entity.HargaDetailProduk{}).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *hargaDetailProdukRepo) UpdateById(ctx context.Context, hargaDetailProduk *entity.HargaDetailProduk) error {
	err := r.DB.WithContext(ctx).Omit("id").Updates(hargaDetailProduk).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *hargaDetailProdukRepo) UpdateByProdukId(ctx context.Context, produkId string, hargaDetails []entity.HargaDetailProduk) error {
	clauses := clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"harga", "qty"}),
	}
	err := r.DB.WithContext(ctx).Where("produk_id = ?", produkId).Clauses(clauses).Create(&hargaDetails).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *hargaDetailProdukRepo) GetByQtyProdukId(ctx context.Context, qty uint, produkId string) (entity.HargaDetailProduk, error) {
	data := entity.HargaDetailProduk{}
	err := r.DB.WithContext(ctx).Find(&data, "qty = ? AND produk_id = ?", qty, produkId).Error
	if err != nil {
		helper.LogsError(err)
		return data, err
	}
	return data, nil
}

func (r *hargaDetailProdukRepo) GetByInQtyProdukId(ctx context.Context, qty []uint, produkId string) ([]entity.HargaDetailProduk, error) {
	datas := []entity.HargaDetailProduk{}

	err := r.DB.WithContext(ctx).Where("produk_id = ? AND qty IN (?)", produkId, qty).Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, err
}

func (r *hargaDetailProdukRepo) GetByProdukId(ctx context.Context, produk_id string) ([]entity.HargaDetailProduk, error) {
	datas := []entity.HargaDetailProduk{}
	err := r.DB.WithContext(ctx).Find(&datas, "produk_id = ?", produk_id).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, nil
}

func (r *hargaDetailProdukRepo) GetById(ctx context.Context, id string) (entity.HargaDetailProduk, error) {
	var data entity.HargaDetailProduk
	err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error
	if err != nil {
		helper.LogsError(err)
		return data, err
	}
	return data, nil
}

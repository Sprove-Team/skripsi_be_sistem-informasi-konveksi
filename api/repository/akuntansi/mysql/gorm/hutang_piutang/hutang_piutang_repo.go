package akuntansi

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

// type CreateParam struct {
// 	Transaksi     *entity.Transaksi
// 	HutangPiutang *entity.HutangPiutang
// }

type SearchParam struct {
	Jenis  []string
	Status []string
}

type HutangPiutangRepo interface {
	Create(ctx context.Context, hutangPiutang *entity.HutangPiutang) error
	GetAll(ctx context.Context, search SearchParam) ([]entity.Transaksi, error)
	// CreateBayarHutangPiutang(ctx context.Context, )
	// Update(ctx context.Context, param UpdateParam) error
	// GetHistory(ctx context.Context, param SearchTransaksi) ([]entity.Transaksi, error)
	// Delete(ctx context.Context, param DeleteParam) error
	// GetById(ctx context.Context, id string) (entity.Transaksi, error)
	// Add more methods as needed for your repository operations
}

type hutangPiutangRepo struct {
	DB *gorm.DB
}

func NewHutangPiutangRepo(DB *gorm.DB) HutangPiutangRepo {
	return &hutangPiutangRepo{DB}
}

func (r *hutangPiutangRepo) Create(ctx context.Context, hutangPiutang *entity.HutangPiutang) error {
	// fmt.Println()
	if err := r.DB.WithContext(ctx).Create(hutangPiutang).Error; err != nil {
		return err
	}
	return nil
}

func (r *hutangPiutangRepo) GetAll(ctx context.Context, search SearchParam) ([]entity.Transaksi, error) {
	datas := []entity.Transaksi{}
	tx := r.DB.WithContext(ctx).Model(&datas).Order("id ASC")

	if search.Jenis != nil {
		tx = tx.Where("jenis IN (?)", search.Jenis)
	}
	if search.Status != nil {
		tx = tx.Where("status IN (?)", search.Status)
	}

	if err := tx.Find(&datas).Error; err != nil {
		helper.LogsError(err)
		return nil, err
	}

	return datas, nil

}

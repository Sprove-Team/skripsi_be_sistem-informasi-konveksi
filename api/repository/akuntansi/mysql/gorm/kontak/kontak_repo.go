package akuntansi

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/entity"
	"gorm.io/gorm"
)

// type CreateParam struct {
// 	Transaksi     *entity.Transaksi
// 	HutangPiutang *entity.HutangPiutang
// }

type KontakRepo interface {
	// Create(ctx context.Context, hutangPiutang *entity.HutangPiutang) error
	// GetAll(ctx context.Context, search SearchParam) ([]entity.Transaksi, error)
	// CreateBayarHutangPiutang(ctx context.Context, )
	// Update(ctx context.Context, param UpdateParam) error
	// GetHistory(ctx context.Context, param SearchTransaksi) ([]entity.Transaksi, error)
	// Delete(ctx context.Context, param DeleteParam) error
	GetById(ctx context.Context, id string) (entity.Kontak, error)
	// Add more methods as needed for your repository operations
}

type kontakRepo struct {
	DB *gorm.DB
}

func NewKontakRepo(DB *gorm.DB) KontakRepo {
	return &kontakRepo{DB}
}

// func (r *kontakRepo) Create(ctx context.Context, hutangPiutang *entity.HutangPiutang) error {
// 	if err := r.DB.WithContext(ctx).Create(hutangPiutang).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *kontakRepo) GetById(ctx context.Context, id string) (entity.Kontak, error) {
	data := entity.Kontak{}
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&data).Error; err != nil {
		return entity.Kontak{}, err
	}
	return data, nil

}

package akuntansi

import (
	"context"
	"log"

	"github.com/be-sistem-informasi-konveksi/entity"
	"gorm.io/gorm"
)

type GolonganAkunRepo interface {
	Create(ctx context.Context, golonganAkun *entity.GolonganAkun) error
	GetById(ctx context.Context, id string) (entity.GolonganAkun, error)
	// Add more methods as needed for your repository operations
}

type golonganAkunRepo struct {
	DB *gorm.DB
}

func NewGolonganAkunRepo(DB *gorm.DB) GolonganAkunRepo {
	return &golonganAkunRepo{DB}
}

func (r *golonganAkunRepo) Create(ctx context.Context, golonganAkun *entity.GolonganAkun) error {
	return r.DB.WithContext(ctx).Create(golonganAkun).Error
}

func (r *golonganAkunRepo) GetById(ctx context.Context, id string) (entity.GolonganAkun, error) {
	data := entity.GolonganAkun{}
	err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error
	if err != nil {
		log.Println("error in golongan_akun_repo GetById -> ", err.Error())
		return data, err
	}
	return data, nil
}

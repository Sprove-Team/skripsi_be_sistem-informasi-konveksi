package akuntansi

import (
	"context"
	"log"

	"github.com/be-sistem-informasi-konveksi/entity"
	"gorm.io/gorm"
)

type KelompokAkunRepo interface {
	Create(ctx context.Context, akun *entity.KelompokAkun) error
	GetById(ctx context.Context, id string) (entity.KelompokAkun, error)
}

type kelompokAkunRepo struct {
	DB *gorm.DB
}

func NewKelompokAkunRepo(DB *gorm.DB) KelompokAkunRepo {
	return &kelompokAkunRepo{DB}
}

func (r *kelompokAkunRepo) Create(ctx context.Context, kelompokAkun *entity.KelompokAkun) error {
	return r.DB.WithContext(ctx).Create(kelompokAkun).Error
}

func (r *kelompokAkunRepo) GetById(ctx context.Context, id string) (entity.KelompokAkun, error) {
	data := entity.KelompokAkun{}
	err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error
	if err != nil {
		log.Println("error in kelompok_akun_repo GetById -> ", err.Error())
		return data, err
	}
	return data, nil
}

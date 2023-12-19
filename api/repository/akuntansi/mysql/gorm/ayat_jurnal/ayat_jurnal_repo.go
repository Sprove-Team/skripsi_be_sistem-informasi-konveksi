package akuntansi

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/entity"
	"gorm.io/gorm"
)

type AyatJurnalRepo interface {
	Create(ctx context.Context, ayatJurnal *entity.AyatJurnal) error
	CreateBatch(ctx context.Context, ayatJurnals []*entity.AyatJurnal) error
}

type ayatJurnalRepo struct {
	DB *gorm.DB
}

func NewAyatJurnal(DB *gorm.DB) AyatJurnalRepo {
	return &ayatJurnalRepo{DB}
}

func (r *ayatJurnalRepo) Create(ctx context.Context, ayatJurnal *entity.AyatJurnal) error {
	return r.DB.WithContext(ctx).Create(ayatJurnal).Error
}

func (r *ayatJurnalRepo) CreateBatch(ctx context.Context, ayatJurnals []*entity.AyatJurnal) error {
	return r.DB.WithContext(ctx).Create(ayatJurnals).Error
}

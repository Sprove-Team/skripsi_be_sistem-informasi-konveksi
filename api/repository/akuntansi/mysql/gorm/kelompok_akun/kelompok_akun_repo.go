package repo_akuntansi_kelompok_akun

import (
	"context"
	"errors"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type KelompokAkunRepo interface {
	Create(ctx context.Context, kelompokAkun *entity.KelompokAkun) error
	Update(ctx context.Context, kelompokAkun *entity.KelompokAkun) error
	Delete(ctx context.Context, id string) error
	GetById(ctx context.Context, id string) (entity.KelompokAkun, error)
	GetAll(ctx context.Context, searchKelompokAkun SearchKelompokAkun) ([]entity.KelompokAkun, error)
	// GetPreloadedAssByJenisAkun(ctx context.Context, jenisAkun string) ([]entity.KelompokAkun, error)
}

type kelompokAkunRepo struct {
	DB *gorm.DB
}

func NewKelompokAkunRepo(DB *gorm.DB) KelompokAkunRepo {
	return &kelompokAkunRepo{DB}
}

func (r *kelompokAkunRepo) Create(ctx context.Context, kelompokAkun *entity.KelompokAkun) error {
	err := r.DB.WithContext(ctx).Create(kelompokAkun).Error
	if err != nil {
		if err != gorm.ErrDuplicatedKey {
			helper.LogsError(err)
		}
		return errors.New(err.Error())
	}
	return err
}

func (r *kelompokAkunRepo) Update(ctx context.Context, kelompokAkun *entity.KelompokAkun) error {
	err := r.DB.WithContext(ctx).Omit("id").Updates(kelompokAkun).Error
	if err != nil {
		if err != gorm.ErrDuplicatedKey {
			helper.LogsError(err)
		}
		return err
	}
	return nil
}

func (r *kelompokAkunRepo) Delete(ctx context.Context, id string) error {
	err := r.DB.WithContext(ctx).Where("id = ?", id).Delete(&entity.KelompokAkun{}).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

type SearchKelompokAkun struct {
	Nama         string
	KategoriAkun string
	Kode         string
	Next         string
	Limit        int
}

func (r *kelompokAkunRepo) GetAll(ctx context.Context, searchKelompokAkun SearchKelompokAkun) ([]entity.KelompokAkun, error) {
	datas := []entity.KelompokAkun{}

	tx := r.DB.WithContext(ctx).Model(&entity.KelompokAkun{}).Order("id ASC").Omit("deleted_at", "updated_at")

	conditions := map[string]interface{}{
		"id > ?":            searchKelompokAkun.Next,
		"kategori_akun = ?": searchKelompokAkun.KategoriAkun,
		"nama LIKE ?":       "%" + searchKelompokAkun.Nama + "%",
		"kode LIKE ?":       searchKelompokAkun.Kode + "%",
	}

	for condition, value := range conditions {
		if value != "" {
			tx = tx.Where(condition, value)
		}
	}

	err := tx.Limit(searchKelompokAkun.Limit).Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return datas, err
	}
	return datas, nil
}

func (r *kelompokAkunRepo) GetById(ctx context.Context, id string) (entity.KelompokAkun, error) {
	data := entity.KelompokAkun{}
	err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			helper.LogsError(err)
		}
		return data, err
	}
	return data, nil
}

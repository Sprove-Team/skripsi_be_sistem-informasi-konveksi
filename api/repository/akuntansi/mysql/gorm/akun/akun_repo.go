package akuntansi

import (
	"context"
	"log"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type AkunRepo interface {
	Create(ctx context.Context, akun *entity.Akun) error
	Update(ctx context.Context, akun *entity.Akun) error
	GetAll(ctx context.Context, search SearchAkun) ([]entity.Akun, error)
	GetById(ctx context.Context, id string) (entity.Akun, error)
	GetByIds(ctx context.Context, ids []string) ([]entity.Akun, error)
}

type akunRepo struct {
	DB *gorm.DB
}

func NewAkunRepo(DB *gorm.DB) AkunRepo {
	return &akunRepo{DB}
}

func (r *akunRepo) GetById(ctx context.Context, id string) (entity.Akun, error) {
	data := entity.Akun{}
	err := r.DB.WithContext(ctx).First(&data, "id = ?", id).Error
	if err != nil {
		log.Println("akun_repo GetById -> ", err.Error())
		return data, err
	}
	return data, nil
}

func (r *akunRepo) GetByIds(ctx context.Context, ids []string) ([]entity.Akun, error) {
	datas := []entity.Akun{}

	err := r.DB.WithContext(ctx).Where("id IN ?", ids).Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return datas, err
	}
	return datas, nil
}

func (r *akunRepo) Create(ctx context.Context, akun *entity.Akun) error {
	err := r.DB.WithContext(ctx).Create(akun).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *akunRepo) Update(ctx context.Context, akun *entity.Akun) error {
	err := r.DB.WithContext(ctx).Omit("id, created_at, updated_at, deleted_at").Updates(akun).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

type SearchAkun struct {
	Nama  string
	Kode  string
	Next  string
	Limit int
}

func (r *akunRepo) GetAll(ctx context.Context, searchAkun SearchAkun) ([]entity.Akun, error) {
	datas := []entity.Akun{}

	tx := r.DB.WithContext(ctx).Model(&entity.Akun{}).Order("id ASC")

	conditions := map[string]interface{}{
		"id > ?":      searchAkun.Next,
		"nama LIKE ?": "%" + searchAkun.Nama + "%",
		"kode = ?":    searchAkun.Kode,
	}

	for condition, value := range conditions {
		if value != "" {
			tx = tx.Where(condition, value)
		}
	}

	err := tx.Limit(searchAkun.Limit).Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return datas, err
	}
	return datas, nil
}

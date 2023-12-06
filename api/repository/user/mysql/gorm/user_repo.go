package user

import (
	"context"

	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
)

type UserRepo interface {
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, param SearchUser) ([]entity.User, error)
	GetByJenisSpvId(ctx context.Context, jenisSpvId string) (entity.User, error)
	GetByUsername(ctx context.Context, username string) (entity.User, error)
	GetById(ctx context.Context, id string) (entity.User, error)
}

type userRepo struct {
	DB *gorm.DB
}

func NewUserRepo(DB *gorm.DB) UserRepo {
	return &userRepo{DB}
}

func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	err := r.DB.WithContext(ctx).Create(user).Error
	return err
}

type SearchUser struct {
	Nama       string
	Role       string
	Username   string
	NoTelp     string
	Alamat     string
	JenisSpvId string
	Limit      int
	Next       string
	// Offset     int
}

func (r *userRepo) GetAll(ctx context.Context, param SearchUser) ([]entity.User, error) {
	datas := []entity.User{}

	tx := r.DB.Model(&entity.User{}).Order("id ASC")

	if param.Role != "" {
		tx = tx.Where("role = ?", param.Role)
	}

	if param.JenisSpvId != "" {
		tx = tx.Where("jenis_spv_id = ?", param.JenisSpvId)
	}
	if param.Nama != "" {
		tx = tx.Where("nama LIKE ?", "%"+param.Nama+"%")
	}

	if param.Username != "" {
		tx = tx.Where("username LIKE ?", "%"+param.Username+"%")
	}

	if param.NoTelp != "" {
		tx = tx.Where("no_telp LIKE ?", "%"+param.NoTelp+"%")
	}

	if param.Alamat != "" {
		tx = tx.Where("alamat LIKE ?", "%"+param.Alamat+"%")
	}

	if param.Next != "" {
		tx = tx.Where("id > ?", param.Next)
	}

	// err := tx.Count(&totalData).Limit(param.Limit).Offset(param.Offset).Find(&datas, "role != ?", "DIREKTUR").Error

	err := tx.Limit(param.Limit).Find(&datas, "role != ?", "DIREKTUR").Error
	return datas, err
}

func (r *userRepo) GetById(ctx context.Context, id string) (entity.User, error) {
	data := entity.User{}
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&data).Error
	return data, err
}

func (r *userRepo) Update(ctx context.Context, user *entity.User) error {
	tx := r.DB.Session(&gorm.Session{
		Context: ctx,
	})
	err := tx.Transaction(func(tx *gorm.DB) error {
		if user.Role != "" && user.Role != "SUPERVISOR" {
			err := tx.Model(&entity.User{}).Where("id = ?", user.ID).Update("jenis_spv_id", nil).Error
			if err != nil {
				return err
			}
			user.JenisSpvID = ""
		}
		err := tx.Omit("id").Updates(user).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error
}

func (r *userRepo) GetByJenisSpvId(ctx context.Context, jenisSpvId string) (entity.User, error) {
	data := entity.User{}

	err := r.DB.WithContext(ctx).Where("jenis_spv_id = ?", jenisSpvId).Error

	return data, err
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	data := entity.User{}

	err := r.DB.WithContext(ctx).Where("username = ?").First(&data).Error

	return data, err
}

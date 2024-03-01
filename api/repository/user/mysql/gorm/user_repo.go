package repo_user

import (
	"context"

	"gorm.io/gorm"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
)

type (
	ParamCreate struct {
		Ctx  context.Context
		User *entity.User
	}

	// ParamUpdate struct represents the parameters for Update method.
	ParamUpdate struct {
		Ctx  context.Context
		User *entity.User
	}

	// ParamDelete struct represents the parameters for Delete method.
	ParamDelete struct {
		Ctx context.Context
		ID  string
	}

	// ParamGetAll struct represents the parameters for GetAll method.
	ParamGetAll struct {
		Ctx    context.Context
		Search SearchParam
	}

	// ParamGetByJenisSpvId struct represents the parameters for GetByJenisSpvId method.
	ParamGetByJenisSpvId struct {
		Ctx        context.Context
		JenisSpvID string
	}

	// ParamGetByUsername struct represents the parameters for GetByUsername method.
	ParamGetByUsername struct {
		Ctx      context.Context
		Username string
	}

	// ParamGetById struct represents the parameters for GetById method.
	ParamGetById struct {
		Ctx context.Context
		ID  string
	}
)

type SearchParam struct {
	Nama       string
	Role       string
	Username   string
	NoTelp     string
	Alamat     string
	JenisSpvId string
	Limit      int
	Next       string
}
type UserRepo interface {
	Create(param ParamCreate) error
	Update(param ParamUpdate) error
	Delete(param ParamDelete) error
	GetAll(param ParamGetAll) ([]entity.User, error)
	GetByJenisSpvId(param ParamGetByJenisSpvId) (*entity.User, error)
	GetByUsername(param ParamGetByUsername) (*entity.User, error)
	GetById(param ParamGetById) (*entity.User, error)
}

type userRepo struct {
	DB *gorm.DB
}

func NewUserRepo(DB *gorm.DB) UserRepo {
	return &userRepo{DB}
}

func (r *userRepo) Create(param ParamCreate) error {
	err := r.DB.WithContext(param.Ctx).Create(param.User).Error
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return err
}

func (r *userRepo) GetAll(param ParamGetAll) ([]entity.User, error) {
	datas := []entity.User{}

	tx := r.DB.Model(&entity.User{}).Order("id ASC")

	if param.Search.Role != "" {
		tx = tx.Where("role = ?", param.Search.Role)
	}

	if param.Search.JenisSpvId != "" {
		tx = tx.Where("jenis_spv_id = ?", param.Search.JenisSpvId)
	}
	if param.Search.Nama != "" {
		tx = tx.Where("nama LIKE ?", "%"+param.Search.Nama+"%")
	}

	if param.Search.Username != "" {
		tx = tx.Where("username LIKE ?", "%"+param.Search.Username+"%")
	}

	if param.Search.NoTelp != "" {
		tx = tx.Where("no_telp LIKE ?", "%"+param.Search.NoTelp+"%")
	}

	if param.Search.Alamat != "" {
		tx = tx.Where("alamat LIKE ?", "%"+param.Search.Alamat+"%")
	}

	if param.Search.Next != "" {
		tx = tx.Where("id > ?", param.Search.Next)
	}

	err := tx.Limit(param.Search.Limit).Preload("JenisSpv").Find(&datas, "role != ?", "DIREKTUR").Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, err
}

func (r *userRepo) GetById(param ParamGetById) (*entity.User, error) {
	data := entity.User{}
	err := r.DB.WithContext(param.Ctx).Where("id = ?", param.ID).First(&data).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return &data, err
}

func (r *userRepo) Update(param ParamUpdate) error {
	tx := r.DB.Session(&gorm.Session{
		Context: param.Ctx,
	})
	err := tx.Transaction(func(tx *gorm.DB) error {
		if param.User.Role != "" && param.User.Role != "SUPERVISOR" {
			err := tx.Model(&entity.User{}).Where("id = ?", param.User.ID).Update("jenis_spv_id", nil).Error
			if err != nil {
				return err
			}
			param.User.JenisSpvID = ""
		}
		err := tx.Omit("id").Updates(param.User).Error
		if err != nil {
			helper.LogsError(err)
			return err
		}
		return nil
	})
	return err
}

func (r *userRepo) Delete(param ParamDelete) error {
	return r.DB.WithContext(param.Ctx).Delete(&entity.User{}, "id = ?", param.ID).Error
}

func (r *userRepo) GetByJenisSpvId(param ParamGetByJenisSpvId) (*entity.User, error) {
	data := entity.User{}

	err := r.DB.WithContext(param.Ctx).Where("jenis_spv_id = ?", param.JenisSpvID).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return &data, err
}

func (r *userRepo) GetByUsername(param ParamGetByUsername) (*entity.User, error) {
	data := entity.User{}
	err := r.DB.WithContext(param.Ctx).Where("username = ?", param.Username).First(&data).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return &data, err
}

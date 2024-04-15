package repo_user

import (
	"context"

	"gorm.io/gorm"

	res_user "github.com/be-sistem-informasi-konveksi/common/response/user"
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
		Ctx      context.Context
		WithJoin bool
		ID       string
	}
	ParamGetUserSpvByIds struct {
		Ctx        context.Context
		JenisSpvID string
		IDs        []string
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
	GetAll(param ParamGetAll) ([]res_user.DataGetUserRes, error)
	GetByJenisSpvId(param ParamGetByJenisSpvId) (*entity.User, error)
	GetByUsername(param ParamGetByUsername) (*entity.User, error)
	GetById(param ParamGetById) (*res_user.DataGetUserRes, error)
	GetUserSpvByIds(param ParamGetUserSpvByIds) ([]entity.User, error)
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
		if err != gorm.ErrDuplicatedKey {
			helper.LogsError(err)
		}
		return err
	}
	return err
}

func (r *userRepo) GetAll(param ParamGetAll) ([]res_user.DataGetUserRes, error) {
	datas := []res_user.DataGetUserRes{}

	tx := r.DB.WithContext(param.Ctx).Model(&entity.User{}).Order("id ASC")

	if param.Search.Role != "" {
		tx = tx.Where("user.role = ?", param.Search.Role)
	}

	if param.Search.JenisSpvId != "" {
		tx = tx.Where("user.jenis_spv_id = ?", param.Search.JenisSpvId)
	}

	if param.Search.Nama != "" {
		tx = tx.Where("user.nama LIKE ?", "%"+param.Search.Nama+"%")
	}

	if param.Search.Username != "" {
		tx = tx.Where("user.username LIKE ?", "%"+param.Search.Username+"%")
	}

	if param.Search.NoTelp != "" {
		tx = tx.Where("user.no_telp LIKE ?", param.Search.NoTelp+"%")
	}

	if param.Search.Alamat != "" {
		tx = tx.Where("user.alamat LIKE ?", "%"+param.Search.Alamat+"%")
	}

	tx = tx.Limit(param.Search.Limit).
		Preload("JenisSpv").
		Joins("LEFT JOIN user_tugas ON user.id = user_tugas.user_id").
		Joins(`LEFT JOIN  (SELECT t.id FROM tugas t
			   	LEFT JOIN sub_tugas st ON t.id = st.tugas_id
				group by t.id
				HAVING COUNT(*) <> SUM(CASE WHEN st.status = ? THEN 1 ELSE 0 END))
				filtered_tugas ON user_tugas.tugas_id = filtered_tugas.id
		`, "SELESAI").
		Select("user.*, COALESCE(SUM(CASE WHEN filtered_tugas.id IS NOT NULL THEN 1 ELSE 0 END), 0) AS total_tugas")

	if param.Search.Next != "" {
		tx = tx.Where("user.id > ?", param.Search.Next)
	}

	err := tx.Group("user.id").Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, err
}

func (r *userRepo) GetById(param ParamGetById) (*res_user.DataGetUserRes, error) {
    data := res_user.DataGetUserRes{}
    // Start building the query with the base user model
    tx := r.DB.WithContext(param.Ctx).Model(&entity.User{}).Where("user.id = ?", param.ID)

    if param.WithJoin {
        // The Joins and subquery need correct SQL syntax and string placement
        tx = tx.
			Preload("JenisSpv").
            Joins("LEFT JOIN user_tugas ON user.id = user_tugas.user_id").
            Joins(`LEFT JOIN (
                SELECT t.id 
                FROM tugas t
                LEFT JOIN sub_tugas st ON t.id = st.tugas_id
                GROUP BY t.id
                HAVING COUNT(*) <> SUM(CASE WHEN st.status = ? THEN 1 ELSE 0 END)
            ) AS filtered_tugas ON user_tugas.tugas_id = filtered_tugas.id`, "SELESAI").
            Select("user.*, COALESCE(SUM(CASE WHEN filtered_tugas.id IS NOT NULL THEN 1 ELSE 0 END), 0) AS total_tugas").
            Group("user.id")
    }

	result := tx.Scan(&data)
    // Execute the query and handle the result
    if err := result.Error; err != nil {
        helper.LogsError(err)
        return nil, err
    }

	if result.RowsAffected == 0{
		return nil, gorm.ErrRecordNotFound
	}

    return &data, nil
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
			if err != gorm.ErrDuplicatedKey {
				helper.LogsError(err)
			}
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

func (r *userRepo) GetUserSpvByIds(param ParamGetUserSpvByIds) ([]entity.User, error) {
	datas := make([]entity.User, 0, len(param.IDs))

	tx := r.DB.WithContext(param.Ctx).Model(&entity.User{}).Where("id IN (?) AND role = 'SUPERVISOR'", param.IDs)
	if param.JenisSpvID != "" {
		tx = tx.Where("jenis_spv_id = ?", param.JenisSpvID)
	}
	err := tx.Find(&datas).Error
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, nil
}

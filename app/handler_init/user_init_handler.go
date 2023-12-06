package handler_init

import (
	"gorm.io/gorm"

	userHandler "github.com/be-sistem-informasi-konveksi/api/handler/user"
	jenisSpvHandler "github.com/be-sistem-informasi-konveksi/api/handler/user/jenis_spv"
	userRepo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	jenisSpvRepo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm/jenis_spv"
	userUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/user"
	jenisSpvUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/user/jenis_spv"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type UserHandlerInit interface {
	UserHandler() userHandler.UserHandler
	JenisSpvHandler() jenisSpvHandler.JenisSpvHandler
}
type userHandlerInit struct {
	DB        *gorm.DB
	validator pkg.Validator
	// uuidGen   pkg.UuidGenerator
	ulid pkg.UlidPkg
	paginate  helper.Paginate
	encryptor helper.Encryptor
}

func NewUserHandlerInit(DB *gorm.DB, validator pkg.Validator, ulid pkg.UlidPkg, paginate helper.Paginate, encryptor helper.Encryptor) UserHandlerInit {
	return &userHandlerInit{DB, validator, ulid, paginate, encryptor}
}

func (d *userHandlerInit) UserHandler() userHandler.UserHandler {
	r := userRepo.NewUserRepo(d.DB)
	rJenisSpv := jenisSpvRepo.NewJenisSpvRepo(d.DB)
	uc := userUsecase.NewUserUsecase(r, rJenisSpv, d.ulid, d.paginate, d.encryptor)
	h := userHandler.NewUserHandler(uc, d.validator)
	return h
}

func (d *userHandlerInit) JenisSpvHandler() jenisSpvHandler.JenisSpvHandler {
	r := jenisSpvRepo.NewJenisSpvRepo(d.DB)
	uc := jenisSpvUsecase.NewJenisSpvUsecase(r, d.ulid, d.paginate)
	h := jenisSpvHandler.NewJenisSpvHandler(uc, d.validator)
	return h
}

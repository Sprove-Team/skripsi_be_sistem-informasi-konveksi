package handler_init

import (
	"gorm.io/gorm"

	handler "github.com/be-sistem-informasi-konveksi/api/handler/user"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/user"
	"github.com/be-sistem-informasi-konveksi/helper"
)

type UserHandlerInit interface {
	UserHandler() handler.UserHandler
	// JenisSpv() handler.JenisSpvHandler
}
type userHandlerInit struct {
	DB        *gorm.DB
	validator helper.Validator
	uuidGen   helper.UuidGenerator
	paginate  helper.Paginate
	encryptor helper.Encryptor
}

func NewUserHandlerInit(DB *gorm.DB, validator helper.Validator, uuidGen helper.UuidGenerator, paginate helper.Paginate, encryptor helper.Encryptor) UserHandlerInit {
	return &userHandlerInit{DB, validator, uuidGen, paginate, encryptor}
}

func (d *userHandlerInit) UserHandler() handler.UserHandler {
	r := repo.NewUserRepo(d.DB)
	rJenisSpv := repo.NewJenisSpvRepo(d.DB)
	uc := usecase.NewUserUsecase(r, rJenisSpv, d.uuidGen, d.paginate, d.encryptor)
	h := handler.NewUserHandler(uc, d.validator)
	return h
}

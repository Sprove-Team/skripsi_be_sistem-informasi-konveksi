package handler_init

import (
	"gorm.io/gorm"

	handler "github.com/be-sistem-informasi-konveksi/handler/user"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/user/mysql/gorm"
	usecase "github.com/be-sistem-informasi-konveksi/usecase/user"
)

type UserHandlerInit interface {
	UserHandler() handler.UserHandler
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
	uc := usecase.NewUserUsecase(r, d.uuidGen, d.paginate ,d.encryptor)
	h := handler.NewUserHandler(uc, d.validator)
	return h
}

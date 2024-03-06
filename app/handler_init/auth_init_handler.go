package handler_init

import (
	handler_auth "github.com/be-sistem-informasi-konveksi/api/handler/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	uc_auth "github.com/be-sistem-informasi-konveksi/api/usecase/auth"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"gorm.io/gorm"
)

type ucAuthInit struct {
	uc uc_auth.AuthUsecase
}

type AuthHandlerInit interface {
	Auth() handler_auth.AuthHandler
}

type authHandlerInit struct {
	DB        *gorm.DB
	jwt       pkg.JwtC
	encryptor helper.Encryptor
	validator pkg.Validator
	uc        ucAuthInit
}

func NewAuthHandlerInit(DB *gorm.DB, jwt pkg.JwtC, validator pkg.Validator, encryptor helper.Encryptor) AuthHandlerInit {
	userRepo := repo_user.NewUserRepo(DB)
	uc := ucAuthInit{
		uc: uc_auth.NewAuthUsecase(userRepo, encryptor, jwt),
	}
	return &authHandlerInit{DB, jwt, encryptor, validator, uc}
}

func (d *authHandlerInit) Auth() handler_auth.AuthHandler {
	h := handler_auth.NewAuthHandler(d.uc.uc, d.validator)
	return h
}

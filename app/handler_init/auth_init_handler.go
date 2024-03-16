package handler_init

import (
	handler_auth "github.com/be-sistem-informasi-konveksi/api/handler/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	uc_auth "github.com/be-sistem-informasi-konveksi/api/usecase/auth"
	uc_user "github.com/be-sistem-informasi-konveksi/api/usecase/user"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"gorm.io/gorm"
)

type ucAuthInit struct {
	uc     uc_auth.AuthUsecase
	ucUser uc_user.UserUsecase
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

func NewAuthHandlerInit(DB *gorm.DB, jwt pkg.JwtC, validator pkg.Validator, ulid pkg.UlidPkg, encryptor helper.Encryptor) AuthHandlerInit {
	userRepo := repo_user.NewUserRepo(DB)
	ucUser := uc_user.NewUserUsecase(userRepo, ulid, encryptor)
	uc := ucAuthInit{
		uc:     uc_auth.NewAuthUsecase(userRepo, encryptor, jwt),
		ucUser: ucUser,
	}
	return &authHandlerInit{DB, jwt, encryptor, validator, uc}
}

func (d *authHandlerInit) Auth() handler_auth.AuthHandler {
	h := handler_auth.NewAuthHandler(d.uc.uc, d.uc.ucUser, d.validator)
	return h
}

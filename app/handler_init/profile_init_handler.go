package handler_init

import (
	"gorm.io/gorm"

	handler_profile "github.com/be-sistem-informasi-konveksi/api/handler/profile"
	userRepo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	uc_profile "github.com/be-sistem-informasi-konveksi/api/usecase/profile"
	uc_user "github.com/be-sistem-informasi-konveksi/api/usecase/user"
	userUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/user"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type ProfileHandlerInit interface {
	ProfileHandler() handler_profile.ProfileHandler
}

type ucProfileInit struct {
	uc     uc_profile.ProfileUsecase
	ucUser uc_user.UserUsecase
}

type profileHandlerInit struct {
	DB        *gorm.DB
	validator pkg.Validator
	ulid      pkg.UlidPkg
	uc        ucProfileInit
}

func NewProfileHandlerInit(DB *gorm.DB, validator pkg.Validator, ulid pkg.UlidPkg, encryptor helper.Encryptor) ProfileHandlerInit {
	userRepo := userRepo.NewUserRepo(DB)
	uc := uc_profile.NewProfileUsecase(userRepo, encryptor)
	ucUser := userUsecase.NewUserUsecase(userRepo, ulid, encryptor)
	ucInit := ucProfileInit{
		uc:     uc,
		ucUser: ucUser,
	}
	return &profileHandlerInit{DB, validator, ulid, ucInit}
}

func (d *profileHandlerInit) ProfileHandler() handler_profile.ProfileHandler {
	h := handler_profile.NewProfileHandler(d.uc.uc, d.uc.ucUser, d.validator)
	return h
}

package uc_auth

import (
	userRepo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	req "github.com/be-sistem-informasi-konveksi/common/request/auth"
)

type AuthUsecase interface {
	Login(reqLogin req.Login) error
}

type authUsecase struct {
	userRepo userRepo.UserRepo
}

// func NewAuthUsecase(userRepo userRepo.UserRepo) AuthUsecase {
// 	return &authUsecase{userRepo}
// }

// func (u *authUsecase) Login(reqLogin req.Login) error {
// }

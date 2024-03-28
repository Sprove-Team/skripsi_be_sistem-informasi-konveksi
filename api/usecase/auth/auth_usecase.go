package uc_auth

import (
	"context"
	"errors"
	"strings"
	"time"

	userRepo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/auth"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"golang.org/x/sync/errgroup"
)

type (
	ParamLogin struct {
		Ctx context.Context
		Req req.Login
	}
	ParamRefreshToken struct {
		Ctx context.Context
		Req req.GetNewToken
	}
)

type AuthUsecase interface {
	Login(param ParamLogin) (token *string, refreshToken *string, err error)
	RefreshToken(param ParamRefreshToken) (newToken *string, err error)
}

type authUsecase struct {
	userRepo  userRepo.UserRepo
	encryptor helper.Encryptor
	jwt       pkg.JwtC
}

func NewAuthUsecase(userRepo userRepo.UserRepo, encryptor helper.Encryptor, jwt pkg.JwtC) AuthUsecase {
	return &authUsecase{userRepo, encryptor, jwt}
}

func (u *authUsecase) Login(param ParamLogin) (token *string, refreshToken *string, err error) {
	userData, err := u.userRepo.GetByUsername(userRepo.ParamGetByUsername{
		Ctx:      param.Ctx,
		Username: param.Req.Username,
	})

	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil, errors.New(message.InvalidUsernameOrPassword)
		}
		return nil, nil, err
	}

	if !u.encryptor.CheckPasswordHash(param.Req.Password, userData.Password) {
		return nil, nil, errors.New(message.InvalidUsernameOrPassword)
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		claims := new(pkg.Claims)
		claims.ID = userData.ID
		claims.Nama = userData.Nama
		claims.Username = userData.Username
		claims.Role = userData.Role
		claims.Subject = "access_token"
		tokenData, err := u.jwt.CreateToken(false, claims, time.Now().Add(time.Second*30))
		if err != nil {
			return err
		}
		token = &tokenData
		return nil
	})

	g.Go(func() error {
		claims := new(pkg.Claims)
		claims.ID = userData.ID
		claims.Subject = "refresh_token"
		refTokenData, err := u.jwt.CreateToken(true, claims, time.Now().Add(time.Minute))
		if err != nil {
			return err
		}
		refreshToken = &refTokenData
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}
	return
}

func (u *authUsecase) RefreshToken(param ParamRefreshToken) (newToken *string, err error) {
	parse, err := u.jwt.ParseToken(true, param.Req.RefreshToken, &pkg.Claims{})
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, errors.New(message.RefreshTokenExpired)
		}
		return nil, errors.New(message.InvalidRefreshToken)
	}

	if !parse.Valid {
		return nil, errors.New(message.InvalidRefreshToken)
	}

	claims := parse.Claims.(*pkg.Claims)

	userData, err := u.userRepo.GetById(userRepo.ParamGetById{
		Ctx: param.Ctx,
		ID:  claims.ID,
	})

	if err != nil {
		if err.Error() == "record not found" {
			return nil, errors.New(message.InvalidRefreshToken)
		}
		return nil, err
	}

	claims.Subject = "access_token"
	claims.ID = userData.ID
	claims.Nama = userData.Nama
	claims.Username = userData.Username
	claims.Role = userData.Role

	newTokenData, err := u.jwt.CreateToken(false, claims, time.Now().Add(time.Hour*8))

	if err != nil {
		return nil, err
	}
	return &newTokenData, nil
}

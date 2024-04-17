package uc_profile

import (
	"context"
	"errors"

	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_profile "github.com/be-sistem-informasi-konveksi/common/request/profile"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"golang.org/x/sync/errgroup"
)

type (
	ParamUpdate struct {
		Ctx    context.Context
		Claims *pkg.Claims
		Req    req_profile.Update
	}
)

type ProfileUsecase interface {
	Update(param ParamUpdate) error
}

type profileUsecase struct {
	userRepo  repo_user.UserRepo
	encryptor helper.Encryptor
}

func NewProfileUsecase(userRepo repo_user.UserRepo, encryptor helper.Encryptor) ProfileUsecase {
	return &profileUsecase{userRepo, encryptor}
}

func (u *profileUsecase) Update(param ParamUpdate) error {
	g := new(errgroup.Group)
	g.SetLimit(10)
	var userData *entity.User
	g.Go(func() error {
		var err error
		userData, err = u.userRepo.GetById(repo_user.ParamGetById{
			Ctx: param.Ctx,
			ID:  param.Claims.ID,
		})
		if err != nil {
			return err
		}
		if param.Req.PasswordLama != "" {
			if !u.encryptor.CheckPasswordHash(param.Req.PasswordLama, userData.Password) {
				return errors.New(message.NotFitOldPassword)
			}
		}
		return nil
	})

	var newPass string
	if param.Req.PasswordBaru != "" {
		g.Go(func() error {
			var err error
			newPass, err = u.encryptor.HashPassword(param.Req.PasswordBaru)
			if err != nil {
				helper.LogsError(err)
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	err := u.userRepo.Update(repo_user.ParamUpdate{
		Ctx: param.Ctx,
		User: &entity.User{
			Base: entity.Base{
				ID: param.Claims.ID,
			},
			Nama:     param.Req.Nama,
			Password: newPass,
			Username: param.Req.Username,
			NoTelp:   param.Req.NoTelp,
			Alamat:   param.Req.Alamat,
		},
	})

	if err != nil {
		if err.Error() == "duplicated key not allowed" {
			return errors.New(message.UsernameConflict)
		}
		return err
	}

	return nil
}

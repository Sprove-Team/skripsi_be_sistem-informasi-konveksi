package uc_user

import (
	"context"
	"errors"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/user"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type (
	ParamCreateUserData struct {
		Ctx context.Context
		Req req.Create
	}
	ParamCreateCommitDB struct {
		Ctx  context.Context
		User *entity.User
	}

	ParamUpdate struct {
		Ctx context.Context
		Req req.Update
	}
	ParamDelete struct {
		Ctx context.Context
		ID  string
	}
	ParamGetAll struct {
		Ctx context.Context
		Req req.GetAll
	}

	ParamGetById struct {
		Ctx context.Context
		ID  string
	}

	ParamUpdateCommitDB struct {
		Ctx  context.Context
		User *entity.User
	}
)

type UserUsecase interface {
	CreateUserData(param ParamCreateUserData) (*entity.User, error)
	CreateCommitDB(param ParamCreateCommitDB) error
	UpdateUserData(param ParamUpdate) (*entity.User, error)
	UpdateCommitDB(param ParamUpdateCommitDB) error
	GetAll(param ParamGetAll) ([]entity.User, error)
	GetById(param ParamGetById) (*entity.User, error)
	Delete(param ParamDelete) error
}

type userUsecase struct {
	repo      repo.UserRepo
	ulid      pkg.UlidPkg
	encryptor helper.Encryptor
}

func NewUserUsecase(
	repo repo.UserRepo,
	ulid pkg.UlidPkg,
	encryptor helper.Encryptor,
) UserUsecase {
	return &userUsecase{repo, ulid, encryptor}
}

func (u *userUsecase) CreateUserData(param ParamCreateUserData) (*entity.User, error) {
	id := u.ulid.MakeUlid().String()
	pass, err := u.encryptor.HashPassword(param.Req.Password)
	if err != nil {
		return nil, err
	}
	user := entity.User{
		Base: entity.Base{
			ID: id,
		},
		Nama:     param.Req.Nama,
		Role:     param.Req.Role,
		NoTelp:   param.Req.NoTelp,
		Alamat:   param.Req.Alamat,
		Username: param.Req.Username,
		Password: pass,
	}
	if param.Req.JenisSpvID != "" {
		user.JenisSpvID = param.Req.JenisSpvID
	}
	return &user, nil
}

func (u *userUsecase) CreateCommitDB(param ParamCreateCommitDB) error {
	return u.repo.Create(repo.ParamCreate{Ctx: param.Ctx, User: param.User})
}

func (u *userUsecase) UpdateUserData(param ParamUpdate) (*entity.User, error) {

	user := entity.User{
		Base: entity.Base{
			ID: param.Req.ID,
		},
		Nama:     param.Req.Nama,
		Role:     param.Req.Role,
		NoTelp:   param.Req.NoTelp,
		Alamat:   param.Req.Alamat,
		Username: param.Req.Username,
	}

	if param.Req.Password != "" {
		pass, err := u.encryptor.HashPassword(param.Req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = pass
	}

	if param.Req.JenisSpvID != "" {
		user.JenisSpvID = param.Req.JenisSpvID
	}

	return &user, nil
}

func (u *userUsecase) UpdateCommitDB(param ParamUpdateCommitDB) error {
	_, err := u.repo.GetById(repo.ParamGetById{
		Ctx: param.Ctx,
		ID:  param.User.ID,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New(message.UserNotFound)
		}
		return err
	}

	err = u.repo.Update(repo.ParamUpdate{
		Ctx:  param.Ctx,
		User: param.User,
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) Delete(param ParamDelete) error {
	_, err := u.repo.GetById(repo.ParamGetById(param))
	if err != nil {
		return err
	}
	err = u.repo.Delete(repo.ParamDelete(param))
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) GetAll(param ParamGetAll) ([]entity.User, error) {
	if param.Req.Limit <= 0 {
		param.Req.Limit = 10
	}
	datas, err := u.repo.GetAll(repo.ParamGetAll{
		Ctx: param.Ctx,
		Search: repo.SearchParam{
			Nama:       param.Req.Search.Nama,
			Alamat:     param.Req.Search.Alamat,
			Username:   param.Req.Search.Username,
			NoTelp:     param.Req.Search.NoTelp,
			Role:       param.Req.Search.Role,
			JenisSpvId: param.Req.Search.JenisSpvID,
			Limit:      param.Req.Limit,
			Next:       param.Req.Next,
		},
	})
	if err != nil {
		return nil, err
	}

	return datas, err
}

func (u *userUsecase) GetById(param ParamGetById) (*entity.User, error) {
	return u.repo.GetById(repo.ParamGetById(param))
}

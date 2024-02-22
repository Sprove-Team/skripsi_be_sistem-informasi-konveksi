package akuntansi

import (
	"context"
	"errors"
	"reflect"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kontak"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/kontak"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type (
	ParamCreateDataKontak struct {
		Ctx context.Context
		Req req.Create
	}
	ParamCreateCommitDB struct {
		Ctx    context.Context
		Kontak *entity.Kontak
	}
	ParamUpdate struct {
		Ctx context.Context
		Req req.Update
	}
	ParamGetAll struct {
		Ctx context.Context
		Req req.GetAll
	}
	ParamGetById struct {
		Ctx context.Context
		ID  string
	}
	ParamDelete struct {
		Ctx context.Context
		ID  string
	}
)

type KontakUsecase interface {
	CreateDataKontak(param ParamCreateDataKontak) (*entity.Kontak, error)
	CreateCommitDB(param ParamCreateCommitDB) error
	Update(param ParamUpdate) error
	GetAll(param ParamGetAll) ([]entity.Kontak, error)
	GetById(param ParamGetById) (*entity.Kontak, error)
	Delete(param ParamDelete) error
}

type kontakUsecase struct {
	repo repo.KontakRepo
	ulid pkg.UlidPkg
}

func NewKontakUsecase(repo repo.KontakRepo, ulid pkg.UlidPkg) KontakUsecase {
	return &kontakUsecase{repo, ulid}
}

func (u *kontakUsecase) CreateDataKontak(param ParamCreateDataKontak) (*entity.Kontak, error) {
	return &entity.Kontak{
		Base: entity.Base{
			ID: u.ulid.MakeUlid().String(),
		},
		Nama:       param.Req.Nama,
		NoTelp:     param.Req.NoTelp,
		Alamat:     param.Req.Alamat,
		Keterangan: param.Req.Keterangan,
		Email:      param.Req.Email,
	}, nil
}

func (u *kontakUsecase) CreateCommitDB(param ParamCreateCommitDB) error {
	return u.repo.Create(repo.ParamCreate{
		Ctx:    param.Ctx,
		Kontak: param.Kontak,
	})
}

func (u *kontakUsecase) Update(param ParamUpdate) error {
	oldData, err := u.repo.GetById(repo.ParamGetById{Ctx: param.Ctx, ID: param.Req.ID})

	if err != nil {
		if err.Error() == "record not found" {
			return errors.New(message.KontakNotFound)
		}
		return err
	}

	paramRepo := repo.ParamUpdate{
		Ctx: param.Ctx,
		Kontak: &entity.Kontak{
			Base: entity.Base{
				ID: param.Req.ID,
			},
			Nama:       param.Req.Nama,
			NoTelp:     param.Req.NoTelp,
			Alamat:     param.Req.Alamat,
			Keterangan: param.Req.Keterangan,
			Email:      param.Req.Email,
		},
	}

	if reflect.DeepEqual(oldData, *paramRepo.Kontak) {
		return errors.New(message.Conflict)
	}
	return u.repo.Update(paramRepo)
}

func (u *kontakUsecase) GetAll(param ParamGetAll) ([]entity.Kontak, error) {
	if param.Req.Limit <= 0 {
		param.Req.Limit = 10
	}

	datas, err := u.repo.GetAll(repo.ParamGetAll{
		Ctx:    param.Ctx,
		Limit:  param.Req.Limit,
		Nama:   param.Req.Nama,
		Email:  param.Req.Email,
		NoTelp: param.Req.NoTelp,
		Next:   param.Req.Next,
	})
	if err != nil {
		return nil, err
	}

	return datas, nil
}

func (u *kontakUsecase) GetById(param ParamGetById) (*entity.Kontak, error) {
	data, err := u.repo.GetById(repo.ParamGetById{
		Ctx: param.Ctx,
		ID:  param.ID,
	})

	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (u *kontakUsecase) Delete(param ParamDelete) error {
	return u.repo.Delete(repo.ParamDelete(param))
}

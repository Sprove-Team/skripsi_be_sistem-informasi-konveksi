package uc_sub_tugas

import (
	"context"
	"errors"

	repo_sub_tugas "github.com/be-sistem-informasi-konveksi/api/repository/tugas/sub_tugas"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_global "github.com/be-sistem-informasi-konveksi/common/request/global"
	req_sub_tugas "github.com/be-sistem-informasi-konveksi/common/request/tugas/sub_tugas"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type (
	ParamCreateByTugasId struct {
		Ctx context.Context
		Req req_sub_tugas.CreateByTugasId
	}
	ParamUpdate struct {
		Ctx    context.Context
		Claims *pkg.Claims
		Req    req_sub_tugas.Update
	}
	ParamDelete struct {
		Ctx context.Context
		Req req_global.ParamByID
	}
)

type SubTugasUsecase interface {
	CreateByTugasId(param ParamCreateByTugasId) error
	Update(param ParamUpdate) error
	Delete(param ParamDelete) error
}

type subTugasUsecase struct {
	repo repo_sub_tugas.SubTugasRepo
	ulid pkg.UlidPkg
}

func NewTugasUsecase(
	repo repo_sub_tugas.SubTugasRepo,
	ulid pkg.UlidPkg) SubTugasUsecase {
	return &subTugasUsecase{repo, ulid}
}

func (u *subTugasUsecase) CreateByTugasId(param ParamCreateByTugasId) error {
	err := u.repo.CreateByTugasId(repo_sub_tugas.ParamCreateByTugasId{
		Ctx: param.Ctx,
		SubTugas: &entity.SubTugas{
			Base: entity.Base{
				ID: u.ulid.MakeUlid().String(),
			},
			Nama:      param.Req.Nama,
			TugasID:   param.Req.TugasID,
			Status:    param.Req.Status,
			Deskripsi: param.Req.Deskripsi,
		},
	})
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New(message.TugasNotFound)
		}
		return err
	}
	return nil
}

func (u *subTugasUsecase) Update(param ParamUpdate) error {
	if _, err := u.repo.GetById(repo_sub_tugas.ParamGetById{
		Ctx: param.Ctx,
		ID:  param.Req.ID,
	}); err != nil {
		return err
	}

	newSubTugas := &entity.SubTugas{
		Base: entity.Base{
			ID: param.Req.ID,
		},
	}

	// supervisor can only update status
	if param.Claims.Role == entity.RolesById[5] {
		newSubTugas.Status = param.Req.Status
	} else {
		newSubTugas.Nama = param.Req.Nama
		newSubTugas.Deskripsi = param.Req.Deskripsi
		newSubTugas.Status = param.Req.Status
	}

	err := u.repo.Update(repo_sub_tugas.ParamUpdate{
		Ctx:         param.Ctx,
		NewSubTugas: newSubTugas,
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *subTugasUsecase) Delete(param ParamDelete) error {
	if _, err := u.repo.GetById(repo_sub_tugas.ParamGetById{
		Ctx: param.Ctx,
		ID:  param.Req.ID,
	}); err != nil {
		return err
	}

	err := u.repo.Delete(repo_sub_tugas.ParamDelete{
		Ctx: param.Ctx,
		ID:  param.Req.ID,
	})

	if err != nil {
		return err
	}

	return nil
}

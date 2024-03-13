package uc_tugas

import (
	"context"
	"errors"
	"time"

	repo_tugas "github.com/be-sistem-informasi-konveksi/api/repository/tugas"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	repo_user_jenis_spv "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm/jenis_spv"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_tugas "github.com/be-sistem-informasi-konveksi/common/request/tugas"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type (
	ParamCreate struct {
		Ctx context.Context
		Req req_tugas.Create
	}
)

type TugasUsecase interface {
	Create(param ParamCreate) error
}

type tugasUsecase struct {
	repo         repo_tugas.TugasRepo
	userRepo     repo_user.UserRepo
	jenisSpvRepo repo_user_jenis_spv.JenisSpvRepo
	ulid         pkg.UlidPkg
}

func NewTugasUsecase(repo repo_tugas.TugasRepo,
	userRepo repo_user.UserRepo,
	jenisSpvRepo repo_user_jenis_spv.JenisSpvRepo,
	ulid pkg.UlidPkg) TugasUsecase {
	return &tugasUsecase{repo, userRepo, jenisSpvRepo, ulid}
}

func countUniqueElements(slice []string) int {
	uniqueMap := make(map[string]bool)

	// Iterate over the slice and add each element to the map
	for _, element := range slice {
		uniqueMap[element] = true
	}

	// Return the number of unique elements (length of the map)
	return len(uniqueMap)
}

func (u *tugasUsecase) Create(param ParamCreate) error {

	users, err := u.userRepo.GetUserSpvByIds(repo_user.ParamGetByIds{
		Ctx: param.Ctx,
		IDs: param.Req.UserID,
	})

	if err != nil {
		return err
	}

	count := countUniqueElements(param.Req.UserID)

	if len(param.Req.UserID) != count {
		return errors.New(message.UserNotFound)
	}

	tanggal_deadline, _ := time.Parse(time.RFC3339, param.Req.TanggalDeadline)

	err = u.repo.Create(repo_tugas.ParamCreate{
		Ctx: param.Ctx,
		Tugas: &entity.Tugas{
			Base: entity.Base{
				ID: u.ulid.MakeUlid().String(),
			},
			InvoiceID:       param.Req.InvoiceID,
			JenisSpvID:      param.Req.JenisSpvID,
			TanggalDeadline: &tanggal_deadline,
			Users:           users,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

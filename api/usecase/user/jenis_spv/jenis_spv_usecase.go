package uc_user_jenis_spv

import (
	"context"
	"log"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm/jenis_spv"
	req "github.com/be-sistem-informasi-konveksi/common/request/user/jenis_spv"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type JenisSpvUsecase interface {
	Create(ctx context.Context, reqJenisSpv req.Create) error
	Update(ctx context.Context, reqJenisSpv req.Update) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]entity.JenisSpv, error)
}

type jenisSpvUsecase struct {
	repo repo.JenisSpvRepo
	ulid pkg.UlidPkg
}

func NewJenisSpvUsecase(repo repo.JenisSpvRepo, ulid pkg.UlidPkg) JenisSpvUsecase {
	return &jenisSpvUsecase{repo, ulid}
}

func (u *jenisSpvUsecase) Create(ctx context.Context, reqJenisSpv req.Create) error {
	id := u.ulid.MakeUlid().String()
	jenis_spv := entity.JenisSpv{
		Base: entity.Base{
			ID: id,
		},
		Nama: reqJenisSpv.Nama,
	}
	return u.repo.Create(ctx, &jenis_spv)
}

func (u *jenisSpvUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.repo.Delete(ctx, id)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (u *jenisSpvUsecase) GetAll(ctx context.Context) ([]entity.JenisSpv, error) {
	return u.repo.GetAll(ctx)
}

func (u *jenisSpvUsecase) Update(ctx context.Context, reqJenisSpv req.Update) error {
	jenisSpv := entity.JenisSpv{
		Base: entity.Base{
			ID: reqJenisSpv.ID,
		},
		Nama: reqJenisSpv.Nama,
	}
	return u.repo.Update(ctx, &jenisSpv)
}

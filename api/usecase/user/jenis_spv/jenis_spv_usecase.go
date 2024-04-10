package uc_user_jenis_spv

import (
	"context"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm/jenis_spv"
	req "github.com/be-sistem-informasi-konveksi/common/request/user/jenis_spv"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"golang.org/x/sync/errgroup"
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
		return err
	}
	return nil
}

func (u *jenisSpvUsecase) GetAll(ctx context.Context) ([]entity.JenisSpv, error) {
	return u.repo.GetAll(ctx)
}

func (u *jenisSpvUsecase) Update(ctx context.Context, reqJenisSpv req.Update) error {
	g := new(errgroup.Group)
	g.SetLimit(3)
	g.Go(func() error {
		_, err := u.repo.GetById(ctx, reqJenisSpv.ID)
		if err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		jenisSpv := entity.JenisSpv{
			Base: entity.Base{
				ID: reqJenisSpv.ID,
			},
			Nama: reqJenisSpv.Nama,
		}
		err := u.repo.Update(ctx, &jenisSpv)
		if err != nil {
			return err
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

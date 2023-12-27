package akuntansi

import (
	"context"
	"errors"
	"strings"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	repoGolonganAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/golongan_akun"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/akun"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type AkunUsecase interface {
	Create(ctx context.Context, reqAkun req.Create) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, reqAkun req.Update) error
	GetAll(ctx context.Context, reqAkun req.GetAll) ([]entity.Akun, error)
}

type akunUsecase struct {
	repo             repo.AkunRepo
	repoGolonganAkun repoGolonganAkun.GolonganAkunRepo
	ulid             pkg.UlidPkg
}

func NewAkunUsecase(repo repo.AkunRepo, ulid pkg.UlidPkg, repoGolonganAkun repoGolonganAkun.GolonganAkunRepo) AkunUsecase {
	return &akunUsecase{repo, repoGolonganAkun, ulid}
}

func (u *akunUsecase) Create(ctx context.Context, reqAkun req.Create) error {
	golonganAkun, err := u.repoGolonganAkun.GetById(ctx, reqAkun.GolonganAkunID)
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New(message.GolonganAkunIdNotFound)
		}
		helper.LogsError(err)
		return err
	}

	kode := golonganAkun.Kode + reqAkun.Kode

	data := entity.Akun{
		ID:             u.ulid.MakeUlid().String(),
		Nama:           reqAkun.Nama,
		Kode:           kode,
		GolonganAkunID: reqAkun.GolonganAkunID,
	}

	return u.repo.Create(ctx, &data)
}

func (u *akunUsecase) Update(ctx context.Context, reqAkun req.Update) error {
	return u.repo.Update(ctx, &entity.Akun{
		ID:             reqAkun.ID,
		GolonganAkunID: reqAkun.GolonganAkunID,
		Nama:           reqAkun.Nama,
		SaldoNormal:    reqAkun.SaldoNormal,
	})
}

func (u *akunUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *akunUsecase) GetAll(ctx context.Context, reqAkun req.GetAll) ([]entity.Akun, error) {
	if reqAkun.Limit <= 0 {
		reqAkun.Limit = 10
	}
	datas, err := u.repo.GetAll(ctx, repo.SearchAkun{
		Nama:  strings.ToLower(reqAkun.Nama),
		Kode:  reqAkun.Kode,
		Limit: reqAkun.Limit,
		Next:  reqAkun.Next,
	})
	if err != nil {
		return nil, err
	}
	return datas, err
}

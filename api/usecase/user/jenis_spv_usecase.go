package user

import (
	"context"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	req "github.com/be-sistem-informasi-konveksi/common/request/user"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
)

type JenisSpvUsecase interface {
	Create(ctx context.Context, reqJenisSpv req.CreateJenisSpv) error
	GetAll(ctx context.Context) ([]entity.JenisSpv, error)
}

type jenisSpvUsecase struct {
	repo     repo.JenisSpvRepo
	uuidGen  helper.UuidGenerator
	paginate helper.Paginate
}

func (u *jenisSpvUsecase) Create(ctx context.Context, reqJenisSpv req.CreateJenisSpv) error {
	id, _ := u.uuidGen.GenerateUUID()
	jenis_spv := entity.JenisSpv{
		ID:   id,
		Nama: reqJenisSpv.Nama,
	}
	return u.repo.Create(ctx, &jenis_spv)
}

func (u *jenisSpvUsecase) GetAll(ctx context.Context) ([]entity.JenisSpv, error) {
	return u.repo.GetAll(ctx)
}

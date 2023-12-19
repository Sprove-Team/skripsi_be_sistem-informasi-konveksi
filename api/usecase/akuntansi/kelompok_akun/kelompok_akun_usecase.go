package akuntansi

import (
	"context"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kelompok_akun"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/kelompok_akun"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type KelompokAkunUsecase interface {
	Create(ctx context.Context, reqKelompokAKun req.Create) error
}

type kelompokAkunUsecase struct {
	repo repo.KelompokAkunRepo
	ulid pkg.UlidPkg
}

func NewKelompokAkunUsecase(repo repo.KelompokAkunRepo, ulid pkg.UlidPkg) KelompokAkunUsecase {
	return &kelompokAkunUsecase{repo, ulid}
}

func (u *kelompokAkunUsecase) Create(ctx context.Context, reqKelompokAKun req.Create) error {
	data := entity.KelompokAkun{
		ID:   u.ulid.MakeUlid().String(),
		Nama: reqKelompokAKun.Nama,
		Kode: reqKelompokAKun.Kode,
	}

	return u.repo.Create(ctx, &data)
}

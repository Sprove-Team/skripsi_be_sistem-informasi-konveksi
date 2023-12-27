package akuntansi

import (
	"context"
	"errors"
	"strings"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/golongan_akun"
	repoKelompokAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kelompok_akun"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/golongan_akun"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type GolonganAkunUsecase interface {
	Create(ctx context.Context, reqGolonganAkun req.Create) error
}

type golonganAkunUsecase struct {
	repo             repo.GolonganAkunRepo
	repoKelompokAkun repoKelompokAkun.KelompokAkunRepo
	ulid             pkg.UlidPkg
}

func NewGolonganAkunUsecase(repo repo.GolonganAkunRepo, repoKelompokAkun repoKelompokAkun.KelompokAkunRepo, ulid pkg.UlidPkg) GolonganAkunUsecase {
	return &golonganAkunUsecase{repo, repoKelompokAkun, ulid}
}

func (u *golonganAkunUsecase) Create(ctx context.Context, reqGolonganAkun req.Create) error {
	// Combine KelompokAkunID from the request with Kode from the request

	kelompokAkun, err := u.repoKelompokAkun.GetById(ctx, reqGolonganAkun.KelompokAkunID)
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New(message.KelompokAkunIdNotFound)
		}
		helper.LogsError(err)
		return err
	}

	kode := kelompokAkun.Kode + reqGolonganAkun.Kode

	data := entity.GolonganAkun{
		ID:             u.ulid.MakeUlid().String(),
		Nama:           strings.ToLower(reqGolonganAkun.Nama),
		Kode:           kode,
		KelompokAkunID: reqGolonganAkun.KelompokAkunID,
	}

	return u.repo.Create(ctx, &data)
}

package akuntansi

import (
	"context"
	"log"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	repoGolonganAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/golongan_akun"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/akun"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type AkunUsecase interface {
	Create(ctx context.Context, reqAkun req.Create) error
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
		log.Println("akun_usecase Create -> ", err)
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

func (u *akunUsecase) GetAll(ctx context.Context, reqAkun req.GetAll) ([]entity.Akun, error) {
	if reqAkun.Limit <= 0 {
		reqAkun.Limit = 10
	}
	datas, err := u.repo.GetAll(ctx, repo.SearchAkun{
		Nama:  reqAkun.Nama,
		Kode:  reqAkun.Kode,
		Limit: reqAkun.Limit,
		Next:  reqAkun.Next,
	})
	if err != nil {
		log.Println("akun_usecase GetAll -> ", err)
		return nil, err
	}
	return datas, err
}

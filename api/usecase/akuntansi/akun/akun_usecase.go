package uc_akuntansi_akun

import (
	"context"
	"errors"
	"strings"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	repoKelompokAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kelompok_akun"
	dataDefault "github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/akun"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type AkunUsecase interface {
	Create(ctx context.Context, reqAkun req.Create) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, reqAkun req.Update) error
	GetAll(ctx context.Context, reqAkun req.GetAll) ([]entity.Akun, error)
	GetById(ctx context.Context, id string) (entity.Akun, error)
}

type akunUsecase struct {
	repo             repo.AkunRepo
	repoKelompokAkun repoKelompokAkun.KelompokAkunRepo
	ulid             pkg.UlidPkg
}

func NewAkunUsecase(repo repo.AkunRepo, ulid pkg.UlidPkg, repoKelompokAkun repoKelompokAkun.KelompokAkunRepo) AkunUsecase {
	return &akunUsecase{repo, repoKelompokAkun, ulid}
}

func (u *akunUsecase) Create(ctx context.Context, reqAkun req.Create) error {
	kelompokAkun, err := u.repoKelompokAkun.GetById(ctx, reqAkun.KelompokAkunID)
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New(message.KelompokAkunNotFound)
		}
		return err
	}

	kode := kelompokAkun.Kode + reqAkun.Kode

	data := entity.Akun{
		Base: entity.Base{
			ID: u.ulid.MakeUlid().String(),
		},
		Nama:           reqAkun.Nama,
		Kode:           kode,
		KelompokAkunID: reqAkun.KelompokAkunID,
		Deskripsi:      reqAkun.Deskripsi,
	}

	return u.repo.Create(ctx, &data)
}

func (u *akunUsecase) GetById(ctx context.Context, id string) (entity.Akun, error) {
	return u.repo.GetById(ctx, id)
}

func (u *akunUsecase) Update(ctx context.Context, reqAkun req.Update) error {
	defaultData := dataDefault.DefaultKodeAkunNKelompokAkun
	if _, ok := defaultData[reqAkun.ID]; ok {
		return errors.New(message.CantModifiedDefaultData)
	}
	akun, err := u.repo.GetById(ctx, reqAkun.ID)
	if err != nil {
		return err
	}

	newKode := ""

	// kelompok akun code
	if reqAkun.KelompokAkunID != "" {
		klmpAkun, err := u.repoKelompokAkun.GetById(ctx, reqAkun.KelompokAkunID)
		if err != nil {
			if err.Error() == "record not found" {
				return errors.New(message.KelompokAkunNotFound)
			}
		}
		newKode += klmpAkun.Kode
	} else {
		newKode += akun.KelompokAkun.Kode
	}

	// akun code
	if reqAkun.Kode != "" {
		newKode += reqAkun.Kode
	} else {
		// remove kelompok akun kode in akun kode
		newKode += strings.Replace(akun.Kode, newKode, "", len(newKode))
	}

	// make empty if `newKode` is the same as the old code
	if newKode == akun.Kode {
		newKode = ""
	}

	return u.repo.Update(ctx, &entity.Akun{
		Base: entity.Base{
			ID: reqAkun.ID,
		},
		KelompokAkunID: reqAkun.KelompokAkunID,
		Nama:           reqAkun.Nama,
		Kode:           newKode,
		SaldoNormal:    reqAkun.SaldoNormal,
	})
}

func (u *akunUsecase) Delete(ctx context.Context, id string) error {
	if _, err := u.repo.GetById(ctx, id); err != nil {
		return err
	}
	defaultData := dataDefault.DefaultKodeAkunNKelompokAkun
	if _, ok := defaultData[id]; ok {
		return errors.New(message.CantModifiedDefaultData)
	}
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

package uc_akuntansi_kelompok_akun

import (
	"context"
	"errors"
	"strings"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kelompok_akun"
	dataDefault "github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/kelompok_akun"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type KelompokAkunUsecase interface {
	Create(ctx context.Context, reqKelompokAKun req.Create) error
	Update(ctx context.Context, reqKelompokAKun req.Update) error
	GetAll(ctx context.Context, reqKelompokAkun req.GetAll) ([]entity.KelompokAkun, error)
	GetById(ctx context.Context, id string) (entity.KelompokAkun, error)
	Delete(ctx context.Context, id string) error
}

type kelompokAkunUsecase struct {
	repo repo.KelompokAkunRepo
	ulid pkg.UlidPkg
}

func NewKelompokAkunUsecase(repo repo.KelompokAkunRepo, ulid pkg.UlidPkg) KelompokAkunUsecase {
	return &kelompokAkunUsecase{repo, ulid}
}

func (u *kelompokAkunUsecase) Create(ctx context.Context, reqKelompokAkun req.Create) error {
	kodeKategori, ok := entity.KategoriAkun[reqKelompokAkun.KategoriAkun]
	if !ok {
		helper.LogsError(errors.New("kategori akun not found"))
		return errors.New(message.KategoriAkunNotFound)
	}
	data := entity.KelompokAkun{
		Base: entity.Base{
			ID: u.ulid.MakeUlid().String(),
		},
		Nama:         strings.ToLower(reqKelompokAkun.Nama),
		Kode:         kodeKategori + reqKelompokAkun.Kode,
		KategoriAkun: reqKelompokAkun.KategoriAkun,
	}

	err := u.repo.Create(ctx, &data)
	if err != nil {
		return err
	}
	return nil
}

func (u *kelompokAkunUsecase) Update(ctx context.Context, reqKelompokAkun req.Update) error {
	defaultData := dataDefault.DefaultKodeAkunNKelompokAkun
	if _, ok := defaultData[reqKelompokAkun.ID]; ok {
		return errors.New(message.CantDeleteDefaultData)
	}
	kelompokAkun, err := u.repo.GetById(ctx, reqKelompokAkun.ID)
	if err != nil {
		return err
	}

	newKode := ""
	// kategori akun code
	if reqKelompokAkun.KategoriAkun != "" {
		kodeKategori, ok := entity.KategoriAkun[reqKelompokAkun.KategoriAkun]
		if !ok {
			return errors.New(message.KategoriAkunNotFound)
		}
		newKode += kodeKategori
	} else {
		newKode += entity.KategoriAkun[kelompokAkun.KategoriAkun]
	}

	// kelompok akun code
	if reqKelompokAkun.Kode != "" {
		newKode += reqKelompokAkun.Kode
	} else {
		// remove kategori akun kode in kelompok akun kode
		newKode += strings.Replace(kelompokAkun.Kode, newKode, "", len(newKode))
	}

	// make empty if `newKode` is the same as the old code
	if newKode == kelompokAkun.Kode {
		newKode = ""
	}

	dataUp := entity.KelompokAkun{
		Base: entity.Base{
			ID: reqKelompokAkun.ID,
		},
		Nama:         strings.ToLower(reqKelompokAkun.Nama),
		KategoriAkun: reqKelompokAkun.KategoriAkun,
		Kode:         newKode,
	}
	return u.repo.Update(ctx, &dataUp)
}

func (u *kelompokAkunUsecase) Delete(ctx context.Context, id string) error {
	defaultData := dataDefault.DefaultKodeAkunNKelompokAkun
	if _, ok := defaultData[id]; ok {
		return errors.New(message.CantDeleteDefaultData)
	}
	return u.repo.Delete(ctx, id)
}

func (u *kelompokAkunUsecase) GetAll(ctx context.Context, reqKelompokAkun req.GetAll) ([]entity.KelompokAkun, error) {
	if reqKelompokAkun.Limit <= 0 {
		reqKelompokAkun.Limit = 10
	}
	datas, err := u.repo.GetAll(ctx, repo.SearchKelompokAkun{
		Nama:         strings.ToLower(reqKelompokAkun.Nama),
		KategoriAkun: reqKelompokAkun.KategoriAkun,
		Kode:         reqKelompokAkun.Kode,
		Limit:        reqKelompokAkun.Limit,
		Next:         reqKelompokAkun.Next,
	})
	if err != nil {
		return nil, err
	}
	return datas, err
}

func (u *kelompokAkunUsecase) GetById(ctx context.Context, id string) (entity.KelompokAkun, error) {
	return u.repo.GetById(ctx, id)
}

package uc_produk_kategori

import (
	"context"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm/kategori"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk/kategori"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type KategoriProdukUsecase interface {
	Create(ctx context.Context, reqKategoriProduk req.Create) error
	Update(ctx context.Context, reqKategoriProduk req.Update) error
	Delete(ctx context.Context, id string) error
	GetById(ctx context.Context, id string) (entity.KategoriProduk, error)
	GetAll(ctx context.Context, reqKategoriProduk req.GetAll) ([]entity.KategoriProduk, error)
}

type kategoriProdukUsecase struct {
	repo repo.KategoriProdukRepo
	ulid pkg.UlidPkg
}

func NewKategoriProdukUsecase(repo repo.KategoriProdukRepo, ulid pkg.UlidPkg) KategoriProdukUsecase {
	return &kategoriProdukUsecase{repo, ulid}
}

func (u *kategoriProdukUsecase) Create(ctx context.Context, reqKategoriProduk req.Create) error {
	id := u.ulid.MakeUlid().String()
	data := entity.KategoriProduk{
		Base: entity.Base{
			ID: id,
		},
		Nama: reqKategoriProduk.Nama,
	}
	return u.repo.Create(ctx, &data)
}

func (u *kategoriProdukUsecase) Update(ctx context.Context, reqKategoriProduk req.Update) error {
	_, err := u.repo.GetById(ctx, reqKategoriProduk.ID)
	if err != nil {
		return err
	}
	data := entity.KategoriProduk{
		Base: entity.Base{
			ID: reqKategoriProduk.ID,
		},
		Nama: reqKategoriProduk.Nama,
	}

	return u.repo.Update(ctx, &data)
}

func (u *kategoriProdukUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.repo.Delete(ctx, id)
	return err
}

func (u *kategoriProdukUsecase) GetById(ctx context.Context, id string) (entity.KategoriProduk, error) {
	data, err := u.repo.GetById(ctx, id)
	return data, err
}

func (u *kategoriProdukUsecase) GetAll(ctx context.Context, reqKategoriProduk req.GetAll) ([]entity.KategoriProduk, error) {
	// currentPage, offset, limit := u.paginate.GetPaginateData(reqKategoriProduk.Page, reqKategoriProduk.Limit)
	if reqKategoriProduk.Limit <= 0 {
		reqKategoriProduk.Limit = 10
	}

	datas, err := u.repo.GetAll(ctx, repo.SearchKategoriProduk{
		Nama:  reqKategoriProduk.Nama,
		Limit: reqKategoriProduk.Limit,
		Next:  reqKategoriProduk.Next,
	})
	if err != nil {
		return nil, err
	}

	return datas, err
}

package produk

import (
	"context"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm/kategori"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk/kategori"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type KategoriProdukUsecase interface {
	Create(ctx context.Context, reqKategoriProduk req.Create) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, reqKategoriProduk req.Update) error
	GetAll(ctx context.Context, reqKategoriProduk req.GetAll) ([]entity.KategoriProduk, int, int, error)
	GetById(ctx context.Context, id string) (entity.KategoriProduk, error)
}

type kategoriProdukUsecase struct {
	repo     repo.KategoriProdukRepo
	uuidGen  pkg.UuidGenerator
	paginate helper.Paginate
}

func NewKategoriProdukUsecase(repo repo.KategoriProdukRepo, uuidGen pkg.UuidGenerator, paginate helper.Paginate) KategoriProdukUsecase {
	return &kategoriProdukUsecase{repo, uuidGen, paginate}
}

func (u *kategoriProdukUsecase) Create(ctx context.Context, reqKategoriProduk req.Create) error {
	id, _ := u.uuidGen.GenerateUUID()
	data := entity.KategoriProduk{
		ID:   id,
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
		ID:   reqKategoriProduk.ID,
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

func (u *kategoriProdukUsecase) GetAll(ctx context.Context, reqKategoriProduk req.GetAll) ([]entity.KategoriProduk, int, int, error) {
	currentPage, offset, limit := u.paginate.GetPaginateData(reqKategoriProduk.Page, reqKategoriProduk.Limit)

	datas, totalData, err := u.repo.GetAll(ctx, repo.SearchKategoriProduk{
		Nama:   reqKategoriProduk.Search.Nama,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, currentPage, 0, err
	}
	totalPage := u.paginate.GetTotalPages(int(totalData), limit)

	return datas, currentPage, totalPage, err
}

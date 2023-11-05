package produk

import (
	"errors"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/produk/mysql/gorm"
)

type ProdukUsecase interface {
	GetById(ctx context.Context, id string) (entity.Produk, error)
	Create(ctx context.Context, reqProduk req.CreateProduk) error
	Update(ctx context.Context, reqProduk req.UpdateProduk) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, reqProduk req.GetAllProduk) ([]entity.Produk, int, int, error)
}

type produkUsecase struct {
	repo         repo.ProdukRepo
	kategoriRepo repo.KategoriProdukRepo
	uuidGen      helper.UuidGenerator
	paginate     helper.Paginate
}

func NewProdukUsecase(
	repo repo.ProdukRepo,
	kategoriRepo repo.KategoriProdukRepo,
	uuidGen helper.UuidGenerator,
	paginate helper.Paginate,
) ProdukUsecase {
	return &produkUsecase{repo, kategoriRepo, uuidGen, paginate}
}

func (u *produkUsecase) Create(ctx context.Context, produk req.CreateProduk) error {
	_, err := u.kategoriRepo.GetById(ctx, produk.KategoriID)
	if err != nil {
		return err
	}
	id, _ := u.uuidGen.GenerateUUID()
	data := entity.Produk{
		ID:               id,
		Nama:             produk.Nama,
		KategoriProdukID: produk.KategoriID,
	}

	err = u.repo.Create(ctx, &data)
	if err != nil {
		return err
	}

	return err
}

func (u *produkUsecase) GetAll(ctx context.Context, reqProduk req.GetAllProduk) ([]entity.Produk, int, int, error) {
	currentPage, offset, limit := u.paginate.GetPaginateData(reqProduk.Page, reqProduk.Limit)
	hargaDetailFilter := reqProduk.Search.HargaDetail != "EMPTY"
	datas, totalData, err := u.repo.GetAll(ctx, repo.SearchProduk{
		Nama:             reqProduk.Search.Nama,
		KategoriProdukId: reqProduk.Search.KategoriID,
		HasHargaDetail:   hargaDetailFilter,
		Limit:            limit,
		Offset:           offset,
	})
	if err != nil {
		return nil, currentPage, 0, err
	}
	totalPage := u.paginate.GetTotalPages(int(totalData), limit)
	return datas, currentPage, totalPage, err
}

func (u *produkUsecase) GetById(ctx context.Context, id string) (entity.Produk, error) {
	data, err := u.repo.GetById(ctx, id)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (u *produkUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.repo.Delete(ctx, id)
	return err
}

func (u *produkUsecase) Update(ctx context.Context, reqProduk req.UpdateProduk) error {
	g := errgroup.Group{}
	// default data not found
	g.Go(func() error {
		_, err := u.repo.GetById(ctx, reqProduk.ID)
		if err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		_, err := u.kategoriRepo.GetById(ctx, reqProduk.KategoriID)
		if err != nil {
			if err.Error() == "record not found" {
				return errors.New(message.KategoriNotFound)
			}
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	data := entity.Produk{
		ID:               reqProduk.ID,
		Nama:             reqProduk.Nama,
		KategoriProdukID: reqProduk.KategoriID,
	}

	return u.repo.Update(ctx, &data)
}

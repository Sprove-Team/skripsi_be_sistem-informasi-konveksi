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
	Create(ctx context.Context, produk req.CreateProduk) error
	Update(ctx context.Context, produk req.UpdateProduk) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, get req.GetAll) ([]entity.Produk, int, int, error)
}

type produkUsecase struct {
	repo            repo.ProdukRepo
	kategoriRepo    repo.KategoriProdukRepo
	uuidGen         helper.UuidGenerator
	paginate        helper.Paginate
}

func NewProdukUsecase(repo repo.ProdukRepo, kategoriRepo repo.KategoriProdukRepo, uuidGen helper.UuidGenerator, paginate helper.Paginate) ProdukUsecase {
	return &produkUsecase{repo, kategoriRepo, uuidGen, paginate}
}

func (u *produkUsecase) Create(ctx context.Context, produk req.CreateProduk) error {
	_, err := u.kategoriRepo.GetById(ctx, produk.KategoriID)
	if err != nil {
		return err
	}
	id, _ := u.uuidGen.GenerateUUID()
	produkR := entity.Produk{
		ID:               id,
		Nama:             produk.Nama,
		KategoriProdukID: produk.KategoriID,
	}
	return u.repo.Create(ctx, &produkR)
}

func (u *produkUsecase) GetAll(ctx context.Context, get req.GetAll) ([]entity.Produk, int, int, error) {
	currentPage, offset, limit := u.paginate.GetPaginateData(get.Page, get.Limit)

	produkRs, totalData, err := u.repo.GetAll(ctx, repo.SearchParams{
		Nama:             get.Search.Nama,
		KategoriProdukId: get.Search.KategoriID,
		Limit:            limit,
		Offset:           offset,
	})
	
	if err != nil {
		return produkRs, currentPage, 0, err
	}
	totalPage := u.paginate.GetTotalPages(int(totalData), limit)

	return produkRs, currentPage, totalPage, err
}

func (u *produkUsecase) GetById(ctx context.Context, id string) (entity.Produk, error) {
	produkR, err := u.repo.GetById(ctx, id)
	if err != nil {
		return produkR, err
	}
	return produkR, nil
}

func (u *produkUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.repo.Delete(ctx, id)
	return err
}

func (u *produkUsecase) Update(ctx context.Context, produk req.UpdateProduk) error {
	g := errgroup.Group{}
	// default data not found
	g.Go(func() error {
		_, err := u.repo.GetById(ctx, produk.ID)
		if err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		_, err := u.kategoriRepo.GetById(ctx, produk.KategoriID)
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

	produkR := entity.Produk{
		ID:               produk.ID,
		Nama:             produk.Nama,
		KategoriProdukID: produk.KategoriID,
	}

	return u.repo.Update(ctx, &produkR)
}

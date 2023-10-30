package produk

import (
	"errors"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	"github.com/be-sistem-informasi-konveksi/common/response/produk"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/produk/mysql/gorm"
)

type ProdukUsecase interface {
	GetById(ctx context.Context, id string) (entity.Produk, error)
	Create(ctx context.Context, produk req.CreateProduk) error
	Update(ctx context.Context, produk req.UpdateProduk) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, get req.GetAllProduk) ([]produk.ProdukRes, int, int, error)
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
	produkR := entity.Produk{
		ID:               id,
		Nama:             produk.Nama,
		KategoriProdukID: produk.KategoriID,
	}

	err = u.repo.Create(ctx, &produkR)
	if err != nil {
		return err
	}

	return err
}

func (u *produkUsecase) GetAll(ctx context.Context, get req.GetAllProduk) ([]produk.ProdukRes, int, int, error) {
	currentPage, offset, limit := u.paginate.GetPaginateData(get.Page, get.Limit)
	hargaDetailFilter := get.Search.HargaDetail != "EMPTY"
	log.Println(get.Search.HargaDetail)
	produkRs, totalData, err := u.repo.GetAll(ctx, repo.SearchProduk{
		Nama:             get.Search.Nama,
		KategoriProdukId: get.Search.KategoriID,
		HasHargaDetail:   hargaDetailFilter,
		Limit:            limit,
		Offset:           offset,
	})
	if err != nil {
		return nil, currentPage, 0, err
	}
	totalPage := u.paginate.GetTotalPages(int(totalData), limit)

	g := errgroup.Group{}
	datas := make([]produk.ProdukRes, len(produkRs))
	for i, produkDat := range produkRs {
		i := i
		produkDat := produkDat
		g.Go(func() error {
			kategori, err := u.kategoriRepo.GetById(ctx, produkDat.KategoriProdukID)
			if err != nil {
				return err
			}
			datas[i] = produk.ProdukRes{
				ID:          produkDat.ID,
				Nama:        produkDat.Nama,
				HargaDetail: produkDat.HargaDetails,
				Kategori:    kategori,
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, currentPage, 0, err
	}

	log.Println(datas)

	return datas, currentPage, totalPage, err
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

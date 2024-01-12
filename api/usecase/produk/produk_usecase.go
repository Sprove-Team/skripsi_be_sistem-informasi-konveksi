package produk

import (
	"errors"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm"
	kategoriRepo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm/kategori"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type ProdukUsecase interface {
	GetById(ctx context.Context, id string) (entity.Produk, error)
	Create(ctx context.Context, reqProduk req.Create) error
	Update(ctx context.Context, reqProduk req.Update) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, reqProduk req.GetAll) ([]entity.Produk, error)
}

type produkUsecase struct {
	repo         repo.ProdukRepo
	kategoriRepo kategoriRepo.KategoriProdukRepo
	ulid         pkg.UlidPkg
	paginate     helper.Paginate
}

func NewProdukUsecase(
	repo repo.ProdukRepo,
	kategoriRepo kategoriRepo.KategoriProdukRepo,
	ulid pkg.UlidPkg,
	paginate helper.Paginate,
) ProdukUsecase {
	return &produkUsecase{repo, kategoriRepo, ulid, paginate}
}

func (u *produkUsecase) Create(ctx context.Context, produk req.Create) error {
	_, err := u.kategoriRepo.GetById(ctx, produk.KategoriID)
	if err != nil {
		return err
	}
	id := u.ulid.MakeUlid().String()
	data := entity.Produk{
		Base: entity.Base{
			ID: id,
		},
		Nama:             produk.Nama,
		KategoriProdukID: produk.KategoriID,
	}

	err = u.repo.Create(ctx, &data)
	if err != nil {
		return err
	}

	return err
}

func (u *produkUsecase) GetAll(ctx context.Context, reqProduk req.GetAll) ([]entity.Produk, error) {
	hargaDetailFilter := reqProduk.HargaDetail != "EMPTY"
	if reqProduk.Limit <= 0 {
		reqProduk.Limit = 10
	}
	datas, err := u.repo.GetAll(ctx, repo.SearchProduk{
		Nama:             reqProduk.Nama,
		KategoriProdukId: reqProduk.KategoriID,
		HasHargaDetail:   hargaDetailFilter,
		Next:             reqProduk.Next,
		Limit:            reqProduk.Limit,
	})
	if err != nil {
		return nil, err
	}
	return datas, err
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

func (u *produkUsecase) Update(ctx context.Context, reqProduk req.Update) error {
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
		Base: entity.Base{
			ID: reqProduk.ID,
		},
		Nama:             reqProduk.Nama,
		KategoriProdukID: reqProduk.KategoriID,
	}

	return u.repo.Update(ctx, &data)
}

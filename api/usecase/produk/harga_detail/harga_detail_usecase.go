package uc_produk_harga_detail

import (
	"context"
	"errors"

	produkRepo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm/harga_detail"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk/harga_detail"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"golang.org/x/sync/errgroup"
)

type HargaDetailProdukUsecase interface {
	Create(ctx context.Context, reqHargaDetailProduk req.Create) error
	Update(ctx context.Context, reqHargaDetailProduk req.Update) error
	Delete(ctx context.Context, id string) error
	GetByProdukId(ctx context.Context, reqHargaDetailProduk req.GetByProdukId) ([]entity.HargaDetailProduk, error)
}

type hargaDetailProdukUsecase struct {
	repo    repo.HargaDetailProdukRepo
	produkR produkRepo.ProdukRepo
	ulid    pkg.UlidPkg
}

func NewHargaDetailProdukUsecase(repo repo.HargaDetailProdukRepo, produkR produkRepo.ProdukRepo, ulid pkg.UlidPkg) HargaDetailProdukUsecase {
	return &hargaDetailProdukUsecase{repo, produkR, ulid}
}

func (u *hargaDetailProdukUsecase) Create(ctx context.Context, reqHargaDetailProduk req.Create) error {
	_, err := u.produkR.GetById(ctx, reqHargaDetailProduk.ProdukId)
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New(message.ProdukNotFound)
		}
		return err
	}

	if err := u.repo.Create(ctx, &entity.HargaDetailProduk{
		Base: entity.Base{
			ID: u.ulid.MakeUlid().String(),
		},
		ProdukID: reqHargaDetailProduk.ProdukId,
		QTY:      reqHargaDetailProduk.QTY,
		Harga:    reqHargaDetailProduk.Harga,
	}); err != nil {
		return err
	}

	return nil
}

func (u *hargaDetailProdukUsecase) Update(ctx context.Context, reqHargaDetailProduk req.Update) error {

	if _, err := u.repo.GetById(ctx, reqHargaDetailProduk.ID); err != nil {
		return err
	}

	err := u.repo.Update(ctx, &entity.HargaDetailProduk{
		Base: entity.Base{
			ID: reqHargaDetailProduk.ID,
		},
		QTY:   reqHargaDetailProduk.QTY,
		Harga: reqHargaDetailProduk.Harga,
	})

	if err != nil {
		return err
	}
	return nil
}

func (u *hargaDetailProdukUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *hargaDetailProdukUsecase) GetByProdukId(ctx context.Context, reqHargaDetailProduk req.GetByProdukId) ([]entity.HargaDetailProduk, error) {
	g := new(errgroup.Group)
	g.SetLimit(3)
	g.Go(func() error {
		err := u.produkR.CheckProduk(ctx, reqHargaDetailProduk.ProdukId)
		if err != nil {
			if err.Error() == "record not found" {
				return errors.New(message.ProdukNotFound)
			}
			return err
		}
		return nil
	})
	var datas []entity.HargaDetailProduk
	g.Go(func() error {
		var err error
		datas, err = u.repo.GetByProdukId(ctx, reqHargaDetailProduk.ProdukId)
		if err != nil {
			return err
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return datas, nil
}

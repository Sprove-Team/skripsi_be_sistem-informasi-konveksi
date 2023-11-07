package produk

import (
	"context"
	"errors"
	"log"

	"golang.org/x/sync/errgroup"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
)

type HargaDetailProdukUsecase interface {
	Create(ctx context.Context, hargaDetailProduk req.CreateHargaDetailProduk) error
	UpdateById(ctx context.Context, hargaDetailProduk req.UpdateHargaDetailProdukById) error
	Delete(ctx context.Context, id string) error
	DeleteByProdukId(ctx context.Context, produkId string) error
	GetByProdukId(ctx context.Context, get req.GetByProdukId) ([]entity.HargaDetailProduk, error)
}

type hargaDetailProdukUsecase struct {
	repo    repo.HargaDetailProdukRepo
	produkR repo.ProdukRepo
	uuidGen helper.UuidGenerator
}

func NewHargaDetailProdukUsecase(repo repo.HargaDetailProdukRepo, produkR repo.ProdukRepo, uuidGen helper.UuidGenerator) HargaDetailProdukUsecase {
	return &hargaDetailProdukUsecase{repo, produkR, uuidGen}
}

func (u *hargaDetailProdukUsecase) Create(ctx context.Context, hargaDetailProduk req.CreateHargaDetailProduk) error {
	datas, _ := u.repo.GetByProdukId(ctx, hargaDetailProduk.ProdukId)
	if len(datas) <= 0 {
		return errors.New(message.ProdukNotFound)
	}
	g := errgroup.Group{}

	if len(hargaDetailProduk.HargaDetail) > 0 {
		for i := 0; i < len(hargaDetailProduk.HargaDetail); i++ {
			i := i
			g.Go(func() error {
				dat, err := u.repo.GetByQtyProdukId(ctx, hargaDetailProduk.HargaDetail[i].QTY, hargaDetailProduk.ProdukId)
				if err != nil {
					return err
				}
				emptyDat := entity.HargaDetailProduk{}
				if dat != emptyDat {
					return errors.New("duplicated key not allowed")
				}
				id, _ := u.uuidGen.GenerateUUID()
				data := entity.HargaDetailProduk{
					ID:       id,
					ProdukID: hargaDetailProduk.ProdukId,
					QTY:      hargaDetailProduk.HargaDetail[i].QTY,
					Harga:    hargaDetailProduk.HargaDetail[i].Harga,
				}
				err = u.repo.Create(ctx, &data)
				if err != nil {
					return err
				}
				return nil
			})
		}
	}
	if err := g.Wait(); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (u *hargaDetailProdukUsecase) UpdateById(ctx context.Context, hargaDetailProduk req.UpdateHargaDetailProdukById) error {
	data := entity.HargaDetailProduk{
		ID:    hargaDetailProduk.ID,
		QTY:   hargaDetailProduk.QTY,
		Harga: float64(hargaDetailProduk.Harga),
		// ProdukID: hargaDetailProduk.ProdukId,
	}
	return u.repo.UpdateById(ctx, &data)
}

func (u *hargaDetailProdukUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.repo.Delete(ctx, id)
	return err
}

func (u *hargaDetailProdukUsecase) DeleteByProdukId(ctx context.Context, produkId string) error {
	_, err := u.repo.GetByProdukId(ctx, produkId)
	if err != nil {
		return err
	}
	err = u.repo.DeleteByProdukId(ctx, produkId)
	return err
}

func (u *hargaDetailProdukUsecase) GetByProdukId(ctx context.Context, get req.GetByProdukId) ([]entity.HargaDetailProduk, error) {
	datas, err := u.repo.GetByProdukId(ctx, get.ProdukId)
	return datas, err
}

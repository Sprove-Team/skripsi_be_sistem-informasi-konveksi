package produk

import (
	"context"
	"errors"
	"log"

	"golang.org/x/sync/errgroup"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/produk/mysql/gorm"
)

type HargaDetailProdukUsecase interface {
	Create(ctx context.Context, hargaDetailProduk req.CreateHargaDetailProduk) error
	UpdateById(ctx context.Context, hargaDetailProduk req.UpdateHargaDetailProdukById) error
	Delete(ctx context.Context, id string) error
	DeleteByProdukId(ctx context.Context, produkId string) error
	GetAll(ctx context.Context) ([]entity.HargaDetailProduk, error)
	GetByProdukId(ctx context.Context, produk_id string) ([]entity.HargaDetailProduk, error)
}

type hargaDetailProdukUsecase struct {
	repo repo.HargaDetailProdukRepo
	produkR repo.ProdukRepo
	uuidGen         helper.UuidGenerator
}

func NewHargaDetailProdukUsecase(repo repo.HargaDetailProdukRepo, produkR repo.ProdukRepo, uuidGen helper.UuidGenerator) HargaDetailProdukUsecase {
	return &hargaDetailProdukUsecase{repo, produkR, uuidGen}
}

func (u *hargaDetailProdukUsecase) Create(ctx context.Context, hargaDetailProduk req.CreateHargaDetailProduk) error {
  _, err := u.repo.GetByProdukId(ctx, hargaDetailProduk.ProdukId)

  if err != nil {
    if err.Error() == "record not found" {
      return errors.New(message.ProdukNotFound)
    }
    return err
  }

	g := errgroup.Group{}

	if len(hargaDetailProduk.HargaProduks) > 0 {
		for i := 0; i < len(hargaDetailProduk.HargaProduks); i++ {
			i := i
			g.Go(func() error {
				dat, err := u.repo.GetByQtyProdukId(ctx, hargaDetailProduk.HargaProduks[i].QTY, hargaDetailProduk.ProdukId)
				if err != nil {
					return err
				}
				emptyDat := entity.HargaDetailProduk{}
				if dat != emptyDat {
					return errors.New("duplicated key not allowed")
				}
				id, _ := u.uuidGen.GenerateUUID()
				hargaDetailProdukR := entity.HargaDetailProduk{
					ID: id,
					ProdukID: hargaDetailProduk.ProdukId,
					QTY:      hargaDetailProduk.HargaProduks[i].QTY,
					Harga:    hargaDetailProduk.HargaProduks[i].Harga,
				}
				return u.repo.Create(ctx, &hargaDetailProdukR)
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
	hargaDetailProdukR := entity.HargaDetailProduk{
		ID:       hargaDetailProduk.ID,
		QTY:      hargaDetailProduk.QTY,
		Harga:    float64(hargaDetailProduk.Harga),
		ProdukID: hargaDetailProduk.ProdukId,
	}
	return u.repo.UpdateById(ctx, &hargaDetailProdukR)
}

func (u *hargaDetailProdukUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.repo.Delete(ctx, id)
	return err
}
func (u *hargaDetailProdukUsecase) DeleteByProdukId(ctx context.Context, produkId string) error{
	_, err := u.repo.GetByProdukId(ctx, produkId)
	if err != nil {
		return err
	}
	err = u.repo.DeleteByProdukId(ctx, produkId)
	return err
}

func (u *hargaDetailProdukUsecase) GetAll(ctx context.Context) ([]entity.HargaDetailProduk, error) {
	hargaDetailProduks, err := u.repo.GetAll(ctx)
	return hargaDetailProduks, err
}

func (u *hargaDetailProdukUsecase) GetByProdukId(ctx context.Context, produk_id string) ([]entity.HargaDetailProduk, error) {
	hargaDetailProduks, err := u.repo.GetByProdukId(ctx, produk_id)
	return hargaDetailProduks, err
}
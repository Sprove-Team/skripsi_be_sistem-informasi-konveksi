package produk

import (
	"sync"

	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/produk/mysql/gorm"
	"golang.org/x/net/context"
)

type ProdukUsecase interface {
	GetById(ctx context.Context, id string) (entity.Produk, error)
	Create(ctx context.Context,produk req.CreateProduk) error
}

type produkUsecase struct {
	repo         repo.ProdukRepo
	kategoriRepo repo.KategoriProdukRepo
	uuidGen      helper.UuidGenerator
}

func NewProdukUsecase(repo repo.ProdukRepo, kategoriRepo repo.KategoriProdukRepo, uuidGen helper.UuidGenerator) ProdukUsecase {
	return &produkUsecase{repo, kategoriRepo, uuidGen}
}

func (u *produkUsecase) Create(ctx context.Context,produk req.CreateProduk) error {
	_, err := u.kategoriRepo.GetById(uint64(produk.IDKategori))
	if err != nil {
		return err
	}
	id, _ := u.uuidGen.GenerateUUID()
	produkR := entity.Produk{
		ID:               id,
		Nama:             produk.Nama,
		KategoriProdukID: produk.IDKategori,
	}
	return u.repo.Create(&produkR)
}

func (u *produkUsecase) GetById(ctx context.Context, id string) (entity.Produk, error) {
	produkR, err := u.repo.GetById(ctx, id)
  
  if err != nil {
    return produkR, err
  }

	wg := new(sync.WaitGroup)

	for i, d := range produkR.HargaDetails {
		wg.Add(1)
		go func(i int, d entity.HargaDetailProduk) {
			defer wg.Done()
      produkR.HargaDetails[i] = entity.HargaDetailProduk{
        QTY: d.QTY,
        ID: d.ID,
        Harga: d.Harga,
      }	
		}(i, d)
	}
	wg.Wait()

	return produkR, nil
}

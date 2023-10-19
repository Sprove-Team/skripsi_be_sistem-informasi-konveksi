package produk

import (
	"sync"

	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	res "github.com/be-sistem-informasi-konveksi/common/response/produk"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/produk/mysql/gorm"
)

type ProdukUsecase interface {
	GetById(id string) (res.DataProdukRes, error)
	Create(produk req.CreateProduk) error
}

type produkUsecase struct {
	repo         repo.ProdukRepo
	kategoriRepo repo.KategoriProdukRepo
	uuidGen      helper.UuidGenerator
}

func NewProdukUsecase(repo repo.ProdukRepo, kategoriRepo repo.KategoriProdukRepo, uuidGen helper.UuidGenerator) ProdukUsecase {
	return &produkUsecase{repo, kategoriRepo, uuidGen}
}

func (u *produkUsecase) Create(produk req.CreateProduk) error {
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

func (u *produkUsecase) GetById(id string) (res.DataProdukRes, error) {
	produkR, err := u.repo.GetById(id)

	produkRes := res.DataProdukRes{}

	if err != nil {
		return produkRes, err
	}

	produkRes.ID = produkR.ID
	produkRes.Nama = produkR.Nama
	produkRes.KategoriId = produkR.KategoriProdukID

	hargaProdukDetails := make([]res.HargaDetailsRes, len(produkR.HargaDetails))

	wg := new(sync.WaitGroup)

	for i, d := range produkR.HargaDetails {
		wg.Add(1)
		go func(i int, d entity.HargaDetailProduk) {
			defer wg.Done()
			hargaProdukDetails[i] = res.HargaDetailsRes{
				ID:    d.ID,
				QTY:   d.QTY,
				Harga: d.Harga,
			}
		}(i, d)
	}
	wg.Wait()

	produkRes.HargaDetail = hargaProdukDetails

	return produkRes, nil
}

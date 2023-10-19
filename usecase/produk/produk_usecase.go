package direktur

import (
	req "github.com/be-sistem-informasi-konveksi/common/request/direktur"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/produk/mysql/gorm"
)

type ProdukUsecase interface {
	Create(produk req.CreateProduk) error
}

type produkUsecase struct {
	repo    repo.ProdukRepo
  kategoriRepo repo.KategoriProdukRepo
	uuidGen helper.UuidGenerator
}

func NewProdukUsecase(repo repo.ProdukRepo,  kategoriRepo repo.KategoriProdukRepo ,uuidGen helper.UuidGenerator) ProdukUsecase {
	return &produkUsecase{repo, kategoriRepo ,uuidGen}
}

func (u *produkUsecase) Create(produk req.CreateProduk) error {
  _,err := u.kategoriRepo.GetById(uint64(produk.IDKategori))
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

// func (u *produkUsecase) GetAll(produk req)

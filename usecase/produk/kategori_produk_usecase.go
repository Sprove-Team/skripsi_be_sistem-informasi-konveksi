package produk 

import (
	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	"github.com/be-sistem-informasi-konveksi/entity"
	repo "github.com/be-sistem-informasi-konveksi/repository/produk/mysql/gorm"
)

type KategoriProdukUsecase interface {
	Create(kategoriProduk req.CreateKategoriProduk) error
	Delete(id uint64) error
}

type kategoriProdukUsecase struct {
	repo repo.KategoriProdukRepo
}

func NewKategoriProdukUsecase(repo repo.KategoriProdukRepo) KategoriProdukUsecase {
	return &kategoriProdukUsecase{repo}
}

func (u *kategoriProdukUsecase) Create(kategoriProduk req.CreateKategoriProduk) error {
	kategoriProdukR := entity.KategoriProduk{
		Nama: kategoriProduk.Nama,
	}
	return u.repo.Create(&kategoriProdukR)
}

func (u *kategoriProdukUsecase) Delete(id uint64) error {
	_, err := u.repo.GetById(id)
	if err != nil {
		return err
	}
	err = u.repo.Delete(id)
	return err
}

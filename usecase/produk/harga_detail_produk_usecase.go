package direktur

import (
	req "github.com/be-sistem-informasi-konveksi/common/request/direktur"
	"github.com/be-sistem-informasi-konveksi/entity"
	repo "github.com/be-sistem-informasi-konveksi/repository/produk/mysql/gorm"
)

type HargaDetailProdukUsecase interface {
	Create(hargaDetailProduk req.CreateHargaDetailProduk) error
}

type hargaDetailProdukUsecase struct {
	repo repo.HargaDetailProdukRepo
}

func NewHargaDetailProdukUsecase(repo repo.HargaDetailProdukRepo) HargaDetailProdukUsecase {
	return &hargaDetailProdukUsecase{repo}
}

func (u *hargaDetailProdukUsecase) Create(hargaDetailProduk req.CreateHargaDetailProduk) error {
	hargaDetailProdukR := entity.HargaDetailProduk{
		ProdukID: hargaDetailProduk.ProdukId,
		QTY:      hargaDetailProduk.QTY,
		Harga:    hargaDetailProduk.Harga,
	}
	return u.repo.Create(&hargaDetailProdukR)
}


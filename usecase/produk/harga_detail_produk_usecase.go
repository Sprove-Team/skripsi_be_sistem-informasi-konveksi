package produk

import (

	"golang.org/x/sync/errgroup"

	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
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
	g := errgroup.Group{}
  
	if len(hargaDetailProduk.ProdukId) > 0 {
		for i := 0; i < len(hargaDetailProduk.HargaProduks); i++ {
			i := i
			g.Go(func() error {
				hargaDetailProdukR := entity.HargaDetailProduk{
					ProdukID: hargaDetailProduk.ProdukId,
					QTY:      hargaDetailProduk.HargaProduks[i].QTY,
					Harga:    hargaDetailProduk.HargaProduks[i].Harga,
				}
				return u.repo.Create(&hargaDetailProdukR)
			})
		}
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

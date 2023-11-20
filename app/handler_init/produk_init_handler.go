package handler_init

import (
	"gorm.io/gorm"

	produkHandler "github.com/be-sistem-informasi-konveksi/api/handler/produk"
	hargaDetailHandler "github.com/be-sistem-informasi-konveksi/api/handler/produk/harga_detail"
	kategoriHandler "github.com/be-sistem-informasi-konveksi/api/handler/produk/kategori"

	produkRepo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm"
	hargaDetailRepo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm/harga_detail"
	kategoriRepo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm/kategori"

	produkUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/produk"
	hargaDetailUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/produk/harga_detail"
	kategoriUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/produk/kategori"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type ProdukHandlerInit interface {
	ProdukHandler() produkHandler.ProdukHandler
	KategoriProdukHandler() kategoriHandler.KategoriProdukHandler
	HargaDetailProdukHandler() hargaDetailHandler.HargaDetailProdukHandler
}
type produkHandlerInit struct {
	DB        *gorm.DB
	validator pkg.Validator
	uuidGen   pkg.UuidGenerator
	paginate  helper.Paginate
}

func NewProdukHandlerInit(DB *gorm.DB, validator pkg.Validator, uuidGen pkg.UuidGenerator, paginate helper.Paginate) ProdukHandlerInit {
	return &produkHandlerInit{DB, validator, uuidGen, paginate}
}

func (d *produkHandlerInit) ProdukHandler() produkHandler.ProdukHandler {
	r := produkRepo.NewProdukRepo(d.DB)
	kategoriR := kategoriRepo.NewKategoriProdukRepo(d.DB)
	uc := produkUsecase.NewProdukUsecase(r, kategoriR, d.uuidGen, d.paginate)
	h := produkHandler.NewProdukHandler(uc, d.validator)
	return h
}

func (d *produkHandlerInit) KategoriProdukHandler() kategoriHandler.KategoriProdukHandler {
	r := kategoriRepo.NewKategoriProdukRepo(d.DB)
	uc := kategoriUsecase.NewKategoriProdukUsecase(r, d.uuidGen, d.paginate)
	h := kategoriHandler.NewKategoriProdukHandler(uc, d.validator)
	return h
}

func (d *produkHandlerInit) HargaDetailProdukHandler() hargaDetailHandler.HargaDetailProdukHandler {
	r := hargaDetailRepo.NewHargaDetailProdukRepo(d.DB)
	produkR := produkRepo.NewProdukRepo(d.DB)
	uc := hargaDetailUsecase.NewHargaDetailProdukUsecase(r, produkR, d.uuidGen)
	h := hargaDetailHandler.NewHargaDetailProdukHandler(uc, d.validator)
	return h
}

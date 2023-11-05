package handler_init

import (
	"gorm.io/gorm"

	handler "github.com/be-sistem-informasi-konveksi/handler/produk"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/produk/mysql/gorm"
	usecase "github.com/be-sistem-informasi-konveksi/usecase/produk"
)

type ProdukHandlerInit interface {
	ProdukHandler() handler.ProdukHandler
	KategoriProdukHandler() handler.KategoriProdukHandler
	HargaDetailProdukHandler() handler.HargaDetailProdukHandler
}
type produkHandlerInit struct {
	DB        *gorm.DB
	validator helper.Validator
	uuidGen   helper.UuidGenerator
	paginate  helper.Paginate
}

func NewProdukHandlerInit(DB *gorm.DB, validator helper.Validator, uuidGen helper.UuidGenerator, paginate helper.Paginate) ProdukHandlerInit {
	return &produkHandlerInit{DB, validator, uuidGen, paginate}
}

func (d *produkHandlerInit) ProdukHandler() handler.ProdukHandler {
	r := repo.NewProdukRepo(d.DB)
	kategoriR := repo.NewKategoriProdukRepo(d.DB)
	uc := usecase.NewProdukUsecase(r, kategoriR, d.uuidGen, d.paginate)
	h := handler.NewProdukHandler(uc, d.validator)
	return h
}

func (d *produkHandlerInit) KategoriProdukHandler() handler.KategoriProdukHandler {
	r := repo.NewKategoriProdukRepo(d.DB)
	uc := usecase.NewKategoriProdukUsecase(r, d.uuidGen, d.paginate)
	h := handler.NewKategoriProdukHandler(uc, d.validator)
	return h
}

func (d *produkHandlerInit) HargaDetailProdukHandler() handler.HargaDetailProdukHandler {
	r := repo.NewHargaDetailProdukRepo(d.DB)
	produkR := repo.NewProdukRepo(d.DB)
	uc := usecase.NewHargaDetailProdukUsecase(r, produkR, d.uuidGen)
	h := handler.NewHargaDetailProdukHandler(uc, d.validator)
	return h
}

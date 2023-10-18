package direktur

import (
	"gorm.io/gorm"

	handler "github.com/be-sistem-informasi-konveksi/handler/direktur/produk"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/direktur/produk/mysql/gorm"
	usecase "github.com/be-sistem-informasi-konveksi/usecase/direktur/produk"
)

type DirekturHandlerInit interface {
	ProdukHandler() handler.ProdukHandler
	KategoriProdukHandler() handler.KategoriProdukHandler
	HargaDetailProdukHandler() handler.HargaDetailProdukHandler
}
type direkturHandlerInit struct {
	DB        *gorm.DB
	validator helper.Validator
	uuidGen   helper.UuidGenerator
}

func NewDirekturHandlerInit(DB *gorm.DB, validator helper.Validator, uuidGen helper.UuidGenerator) DirekturHandlerInit {
	return &direkturHandlerInit{DB, validator, uuidGen}
}

func (d *direkturHandlerInit) ProdukHandler() handler.ProdukHandler {
	r := repo.NewProdukRepo(d.DB)
	uc := usecase.NewProdukUsecase(r, d.uuidGen)
	h := handler.NewProdukHandler(uc, d.validator)
	return h
}

func (d *direkturHandlerInit) KategoriProdukHandler() handler.KategoriProdukHandler {
	r := repo.NewKategoriProdukRepo(d.DB)
	uc := usecase.NewKategoriProdukUsecase(r)
	h := handler.NewKategoriProdukHandler(uc, d.validator)
	return h
}

func (d *direkturHandlerInit) HargaDetailProdukHandler() handler.HargaDetailProdukHandler {
	r := repo.NewHargaDetailProdukRepo(d.DB)
	uc := usecase.NewHargaDetailProdukUsecase(r)
	h := handler.NewHargaDetailProdukHandler(uc, d.validator)
	return h
}

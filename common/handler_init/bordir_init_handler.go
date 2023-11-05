package handler_init

import (
	"gorm.io/gorm"

	handler "github.com/be-sistem-informasi-konveksi/handler/bordir"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/bordir/mysql/gorm"
	usecase "github.com/be-sistem-informasi-konveksi/usecase/bordir"
)

type BordirHandlerInit interface {
	BordirHandler() handler.BordirHandler
}
type bordirHandlerInit struct {
	DB        *gorm.DB
	validator helper.Validator
	uuidGen   helper.UuidGenerator
	paginate  helper.Paginate
}

func NewBordirHandlerInit(DB *gorm.DB, validator helper.Validator, uuidGen helper.UuidGenerator, paginate helper.Paginate) BordirHandlerInit {
	return &bordirHandlerInit{DB, validator, uuidGen, paginate}
}

func (d *bordirHandlerInit) BordirHandler() handler.BordirHandler {
	r := repo.NewProdukRepo(d.DB)
	uc := usecase.NewBordirUsecase(r, d.uuidGen, d.paginate)
	h := handler.NewBordirHandler(uc, d.validator)
	return h
}

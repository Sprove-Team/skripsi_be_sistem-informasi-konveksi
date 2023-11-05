package handler_init

import (
	"gorm.io/gorm"

	handler "github.com/be-sistem-informasi-konveksi/handler/sablon"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/sablon/mysql/gorm"
	usecase "github.com/be-sistem-informasi-konveksi/usecase/sablon"
)

type SablonHandlerInit interface {
	SablonHandler() handler.SablonHandler
}
type sablonHandlerInit struct {
	DB        *gorm.DB
	validator helper.Validator
	uuidGen   helper.UuidGenerator
	paginate  helper.Paginate
}

func NewSablonHandlerInit(DB *gorm.DB, validator helper.Validator, uuidGen helper.UuidGenerator, paginate helper.Paginate) SablonHandlerInit {
	return &sablonHandlerInit{DB, validator, uuidGen, paginate}
}

func (d *sablonHandlerInit) SablonHandler() handler.SablonHandler {
	r := repo.NewSablonRepo(d.DB)
	uc := usecase.NewSablonUsecase(r, d.uuidGen, d.paginate)
	h := handler.NewSablonHandler(uc, d.validator)
	return h
}

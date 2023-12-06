package handler_init

import (
	handler "github.com/be-sistem-informasi-konveksi/api/handler/bordir"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/bordir/mysql/gorm"
	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/bordir"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"gorm.io/gorm"
)

type BordirHandlerInit interface {
	BordirHandler() handler.BordirHandler
}
type bordirHandlerInit struct {
	DB        *gorm.DB
	validator pkg.Validator
	// uuidGen   pkg.UuidGenerator
	ulid pkg.UlidPkg
	paginate  helper.Paginate
}

func NewBordirHandlerInit(DB *gorm.DB, validator pkg.Validator, ulid pkg.UlidPkg, paginate helper.Paginate) BordirHandlerInit {
	return &bordirHandlerInit{DB, validator, ulid, paginate}
}

func (d *bordirHandlerInit) BordirHandler() handler.BordirHandler {
	r := repo.NewProdukRepo(d.DB)
	uc := usecase.NewBordirUsecase(r, d.ulid, d.paginate)
	h := handler.NewBordirHandler(uc, d.validator)
	return h
}

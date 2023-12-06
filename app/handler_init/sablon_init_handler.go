package handler_init

import (
	"gorm.io/gorm"

	handler "github.com/be-sistem-informasi-konveksi/api/handler/sablon"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/sablon/mysql/gorm"
	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/sablon"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type SablonHandlerInit interface {
	SablonHandler() handler.SablonHandler
}
type sablonHandlerInit struct {
	DB        *gorm.DB
	validator pkg.Validator
	// uuidGen   pkg.UuidGenerator
	ulid pkg.UlidPkg
	paginate  helper.Paginate
}

func NewSablonHandlerInit(DB *gorm.DB, validator pkg.Validator, ulid pkg.UlidPkg, paginate helper.Paginate) SablonHandlerInit {
	return &sablonHandlerInit{DB, validator, ulid, paginate}
}

func (d *sablonHandlerInit) SablonHandler() handler.SablonHandler {
	r := repo.NewSablonRepo(d.DB)
	uc := usecase.NewSablonUsecase(r, d.ulid, d.paginate)
	h := handler.NewSablonHandler(uc, d.validator)
	return h
}

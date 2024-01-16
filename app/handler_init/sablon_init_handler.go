package handler_init

import (
	"gorm.io/gorm"

	handler "github.com/be-sistem-informasi-konveksi/api/handler/sablon"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/sablon/mysql/gorm"
	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/sablon"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type SablonHandlerInit interface {
	SablonHandler() handler.SablonHandler
}
type sablonHandlerInit struct {
	DB        *gorm.DB
	validator pkg.Validator
	ulid pkg.UlidPkg
}

func NewSablonHandlerInit(DB *gorm.DB, validator pkg.Validator, ulid pkg.UlidPkg) SablonHandlerInit {
	return &sablonHandlerInit{DB, validator, ulid}
}

func (d *sablonHandlerInit) SablonHandler() handler.SablonHandler {
	r := repo.NewSablonRepo(d.DB)
	uc := usecase.NewSablonUsecase(r, d.ulid)
	h := handler.NewSablonHandler(uc, d.validator)
	return h
}

package handler_init

import (
	"gorm.io/gorm"

	handler_tugas "github.com/be-sistem-informasi-konveksi/api/handler/tugas"
	repo_invoice "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm"
	repo_tugas "github.com/be-sistem-informasi-konveksi/api/repository/tugas"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	uc_tugas "github.com/be-sistem-informasi-konveksi/api/usecase/tugas"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type TugasHandlerInit interface {
	TugasHandler() handler_tugas.TugasHandler
}

type ucTugasInit struct {
	uc uc_tugas.TugasUsecase
}

type tugasHandlerInit struct {
	DB        *gorm.DB
	validator pkg.Validator
	ulid      pkg.UlidPkg
	uc        ucTugasInit
}

func NewTugasHandlerInit(DB *gorm.DB, validator pkg.Validator, ulid pkg.UlidPkg) TugasHandlerInit {
	repo := repo_tugas.NewTugasRepo(DB)
	userRepo := repo_user.NewUserRepo(DB)
	invoiceRepo := repo_invoice.NewInvoiceRepo(DB)
	uc_init := ucTugasInit{
		uc: uc_tugas.NewTugasUsecase(repo, userRepo, invoiceRepo, ulid),
	}
	return &tugasHandlerInit{DB, validator, ulid, uc_init}
}

func (d *tugasHandlerInit) TugasHandler() handler_tugas.TugasHandler {
	h := handler_tugas.NewTugasHandler(d.uc.uc, d.validator)
	return h
}

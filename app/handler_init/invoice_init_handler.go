package handler_init

import (
	handler "github.com/be-sistem-informasi-konveksi/api/handler/invoice"
	repoAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	repoBordir "github.com/be-sistem-informasi-konveksi/api/repository/bordir/mysql/gorm"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm"
	repoProduk "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm"
	repoSablon "github.com/be-sistem-informasi-konveksi/api/repository/sablon/mysql/gorm"
	repoUser "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/invoice"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"gorm.io/gorm"
)

type InvoiceHandlerInit interface {
	InvoiceHandler() handler.InvoiceHandler
}

type invoiceHandlerInit struct {
	DB        *gorm.DB
	validator pkg.Validator
	ulid      pkg.UlidPkg
}

func NewInvoiceHandlerInit(DB *gorm.DB, validator pkg.Validator, ulid pkg.UlidPkg) InvoiceHandlerInit {
	return &invoiceHandlerInit{DB, validator, ulid}
}

func (d *invoiceHandlerInit) InvoiceHandler() handler.InvoiceHandler {
	rp := repoProduk.NewProdukRepo(d.DB)
	ra := repoAkun.NewAkunRepo(d.DB)
	rb := repoBordir.NewBordirRepo(d.DB)
	rs := repoSablon.NewSablonRepo(d.DB)
	ru := repoUser.NewUserRepo(d.DB)
	r := repo.NewInvoiceRepo(d.DB)
	uc := usecase.NewInvoiceUsecase(r, ra, ru, rb, rs, rp, d.ulid)
	h := handler.NewInvoiceHandler(uc, d.validator)
	return h
}

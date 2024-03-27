package handler_init

import (
	handler "github.com/be-sistem-informasi-konveksi/api/handler/invoice"
	handler_invoice_data_bayar "github.com/be-sistem-informasi-konveksi/api/handler/invoice/data_bayar"
	repoAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	repoBayarHP "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/data_bayar_hutang_piutang"
	repoHP "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/hutang_piutang"
	repoKelompokAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kelompok_akun"
	repoKontak "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kontak"
	repoBordir "github.com/be-sistem-informasi-konveksi/api/repository/bordir/mysql/gorm"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm"
	repo_invoice_data_bayar "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm/data_bayar"
	repoProduk "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm"
	repoSablon "github.com/be-sistem-informasi-konveksi/api/repository/sablon/mysql/gorm"
	repoUser "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	ucAkun "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/akun"
	ucHP "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/hutang_piutang"
	ucKontak "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/kontak"
	uc "github.com/be-sistem-informasi-konveksi/api/usecase/invoice"
	uc_invoice_data_bayar "github.com/be-sistem-informasi-konveksi/api/usecase/invoice/data_bayar"
	ucUser "github.com/be-sistem-informasi-konveksi/api/usecase/user"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"gorm.io/gorm"
)

type ucInvoiceInit struct {
	uc                 uc.InvoiceUsecase
	ucAkun             ucAkun.AkunUsecase
	ucUser             ucUser.UserUsecase
	ucKontak           ucKontak.KontakUsecase
	ucHP               ucHP.HutangPiutangUsecase
	ucDataBayarInvoice uc_invoice_data_bayar.DataBayarInvoice
}

type InvoiceHandlerInit interface {
	InvoiceHandler() handler.InvoiceHandler
	DataBayarInvoiceHandler() handler_invoice_data_bayar.DataBayarInvoiceHandler
}

type invoiceHandlerInit struct {
	DB        *gorm.DB
	validator pkg.Validator
	ulid      pkg.UlidPkg
	uc        ucInvoiceInit
}

func NewInvoiceHandlerInit(DB *gorm.DB, validator pkg.Validator, ulid pkg.UlidPkg, encryptor helper.Encryptor) InvoiceHandlerInit {

	repo := repo.NewInvoiceRepo(DB)
	akunRepo := repoAkun.NewAkunRepo(DB)
	kelompokAkunRepo := repoKelompokAkun.NewKelompokAkunRepo(DB)
	userRepo := repoUser.NewUserRepo(DB)
	kontakRepo := repoKontak.NewKontakRepo(DB)
	hpRepo := repoHP.NewHutangPiutangRepo(DB)
	bayarRepo := repoBayarHP.NewDataBayarHutangPiutangRepo(DB)
	produkRepo := repoProduk.NewProdukRepo(DB)
	bordirRepo := repoBordir.NewBordirRepo(DB)
	sablonRepo := repoSablon.NewSablonRepo(DB)
	dataBayarInvoiceRepo := repo_invoice_data_bayar.NewDataBayarInvoiceRepo(DB)

	uc := ucInvoiceInit{
		uc:                 uc.NewInvoiceUsecase(repo, akunRepo, kontakRepo, produkRepo, bordirRepo, sablonRepo, ulid),
		ucAkun:             ucAkun.NewAkunUsecase(akunRepo, ulid, kelompokAkunRepo),
		ucUser:             ucUser.NewUserUsecase(userRepo, ulid, encryptor),
		ucKontak:           ucKontak.NewKontakUsecase(kontakRepo, ulid),
		ucHP:               ucHP.NewHutangPiutangUsecase(hpRepo, bayarRepo, akunRepo, kontakRepo, ulid),
		ucDataBayarInvoice: uc_invoice_data_bayar.NewDataInvoice(dataBayarInvoiceRepo, repo, hpRepo, ulid),
	}

	return &invoiceHandlerInit{DB, validator, ulid, uc}
}

func (d *invoiceHandlerInit) InvoiceHandler() handler.InvoiceHandler {
	h := handler.NewInvoiceHandler(d.uc.uc, d.uc.ucUser, d.uc.ucKontak, d.uc.ucHP, d.validator)
	return h
}

func (d *invoiceHandlerInit) DataBayarInvoiceHandler() handler_invoice_data_bayar.DataBayarInvoiceHandler {
	h := handler_invoice_data_bayar.NewDataBayarInvoiceHandler(d.uc.ucDataBayarInvoice, d.uc.ucHP, d.validator)
	return h
}

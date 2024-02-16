package handler_init

import (
	handler "github.com/be-sistem-informasi-konveksi/api/handler/invoice"
	repoAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	repoBayarHP "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/data_bayar_hutang_piutang"
	repoHP "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/hutang_piutang"
	repoKelompokAkun "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kelompok_akun"
	repoKontak "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kontak"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/invoice/mysql/gorm"
	repoUser "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	ucAkun "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/akun"
	ucHP "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/hutang_piutang"
	ucKontak "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/kontak"
	uc "github.com/be-sistem-informasi-konveksi/api/usecase/invoice"
	ucUser "github.com/be-sistem-informasi-konveksi/api/usecase/user"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"gorm.io/gorm"
)

type ucInvoiceInit struct {
	uc       uc.InvoiceUsecase
	ucAkun   ucAkun.AkunUsecase
	ucUser   ucUser.UserUsecase
	ucKontak ucKontak.KontakUsecase
	ucHP     ucHP.HutangPiutangUsecase
}

type InvoiceHandlerInit interface {
	InvoiceHandler() handler.InvoiceHandler
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

	uc := ucInvoiceInit{
		uc:       uc.NewInvoiceUsecase(repo, akunRepo, userRepo, kontakRepo, ulid),
		ucAkun:   ucAkun.NewAkunUsecase(akunRepo, ulid, kelompokAkunRepo),
		ucUser:   ucUser.NewUserUsecase(userRepo, ulid, encryptor),
		ucKontak: ucKontak.NewKontakUsecase(kontakRepo, ulid),
		ucHP:     ucHP.NewHutangPiutangUsecase(hpRepo, bayarRepo, akunRepo, kontakRepo, ulid),
	}

	return &invoiceHandlerInit{DB, validator, ulid, uc}
}

func (d *invoiceHandlerInit) InvoiceHandler() handler.InvoiceHandler {
	h := handler.NewInvoiceHandler(d.uc.uc, d.uc.ucUser, d.uc.ucKontak, d.uc.ucHP, d.validator)
	return h
}

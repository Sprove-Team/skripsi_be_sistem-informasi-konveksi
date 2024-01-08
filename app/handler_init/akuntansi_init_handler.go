package handler_init

import (
	akuntansiHandler "github.com/be-sistem-informasi-konveksi/api/handler/akuntansi"
	akuntansiRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm"
	akuntansiUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi"

	transaksiHandler "github.com/be-sistem-informasi-konveksi/api/handler/akuntansi/transaksi"
	transaksiRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/transaksi"
	transaksiUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/transaksi"

	kelompokAkunHandler "github.com/be-sistem-informasi-konveksi/api/handler/akuntansi/kelompok_akun"
	kelompokAkunRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kelompok_akun"
	kelompokAkunUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/kelompok_akun"

	akunHandler "github.com/be-sistem-informasi-konveksi/api/handler/akuntansi/akun"
	akunRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	akunUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/akun"

	"github.com/be-sistem-informasi-konveksi/pkg"
	"gorm.io/gorm"
)

type AkuntansiHandlerInit interface {
	AkunHandler() akunHandler.AkunHandler
	// GolonganAkunHandler() golonganAkunHandler.GolonganAkunHandler
	KelompokAkunHandler() kelompokAkunHandler.KelompokAkunHandler
	Transaksi() transaksiHandler.TransaksiHandler
	Akuntansi() akuntansiHandler.AkuntansiHandler
}

type akuntansiHandlerInit struct {
	DB        *gorm.DB
	validator pkg.Validator
	ulid      pkg.UlidPkg
	ac        pkg.AccountingPkg
}

func NewAkuntansiHandlerInit(DB *gorm.DB, validator pkg.Validator, ulid pkg.UlidPkg, ac pkg.AccountingPkg) AkuntansiHandlerInit {
	return &akuntansiHandlerInit{DB, validator, ulid, ac}
}

func (d *akuntansiHandlerInit) AkunHandler() akunHandler.AkunHandler {
	r := akunRepo.NewAkunRepo(d.DB)
	// rg := golonganAkunRepo.NewGolonganAkunRepo(d.DB)
	rk := kelompokAkunRepo.NewKelompokAkunRepo(d.DB)

	uc := akunUsecase.NewAkunUsecase(r, d.ulid, rk)
	h := akunHandler.NewAkunHandler(uc, d.validator)
	return h
}

func (d *akuntansiHandlerInit) KelompokAkunHandler() kelompokAkunHandler.KelompokAkunHandler {
	r := kelompokAkunRepo.NewKelompokAkunRepo(d.DB)
	uc := kelompokAkunUsecase.NewKelompokAkunUsecase(r, d.ulid)
	h := kelompokAkunHandler.NewKelompokAkunHandler(uc, d.validator)
	return h
}

func (d *akuntansiHandlerInit) Transaksi() transaksiHandler.TransaksiHandler {
	r := transaksiRepo.NewTransaksiRepo(d.DB)
	rk := akunRepo.NewAkunRepo(d.DB)
	uc := transaksiUsecase.NewTransaksiUsecase(r, rk, d.ulid)
	h := transaksiHandler.NewTransaksiHandler(uc, d.validator)
	return h
}

func (d *akuntansiHandlerInit) Akuntansi() akuntansiHandler.AkuntansiHandler {
	ra := akuntansiRepo.NewAkuntansiRepo(d.DB)
	uc := akuntansiUsecase.NewAkuntansiUsecase(ra)
	h := akuntansiHandler.NewAkuntansiHandler(uc, d.validator)
	return h
}

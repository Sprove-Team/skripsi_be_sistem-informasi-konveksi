package handler_init

import (
	transaksiHandler "github.com/be-sistem-informasi-konveksi/api/handler/akuntansi/transaksi"
	transaksiRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/transaksi"
	transaksiUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/transaksi"

	golonganAkunHandler "github.com/be-sistem-informasi-konveksi/api/handler/akuntansi/golongan_akun"
	golonganAkunRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/golongan_akun"
	golonganAkunUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/golongan_akun"

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
	GolonganAkunHandler() golonganAkunHandler.GolonganAkunHandler
	KelompokAkunHandler() kelompokAkunHandler.KelompokAkunHandler
	Transaksi() transaksiHandler.TransaksiHandler
}

type akuntansiHandlerInit struct {
	DB        *gorm.DB
	validator pkg.Validator
	ulid      pkg.UlidPkg
}

func NewAkuntansiHandlerInit(DB *gorm.DB, validator pkg.Validator, ulid pkg.UlidPkg) AkuntansiHandlerInit {
	return &akuntansiHandlerInit{DB, validator, ulid}
}

func (d *akuntansiHandlerInit) AkunHandler() akunHandler.AkunHandler {
	r := akunRepo.NewAkunRepo(d.DB)
	rg := golonganAkunRepo.NewGolonganAkunRepo(d.DB)

	uc := akunUsecase.NewAkunUsecase(r, d.ulid, rg)
	h := akunHandler.NewAkunHandler(uc, d.validator)
	return h
}

func (d *akuntansiHandlerInit) GolonganAkunHandler() golonganAkunHandler.GolonganAkunHandler {
	r := golonganAkunRepo.NewGolonganAkunRepo(d.DB)
	rk := kelompokAkunRepo.NewKelompokAkunRepo(d.DB)
	uc := golonganAkunUsecase.NewGolonganAkunUsecase(r, rk, d.ulid)
	h := golonganAkunHandler.NewGolonganAkunHandler(uc, d.validator)
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

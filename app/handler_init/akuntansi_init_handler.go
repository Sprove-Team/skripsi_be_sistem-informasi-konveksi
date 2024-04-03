package handler_init

import (
	akuntansiHandler "github.com/be-sistem-informasi-konveksi/api/handler/akuntansi"
	akuntansiRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm"
	akuntansiUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi"

	transaksiHandler "github.com/be-sistem-informasi-konveksi/api/handler/akuntansi/transaksi"
	transaksiRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/transaksi"
	transaksiUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/transaksi"

	hutangPiutangHandler "github.com/be-sistem-informasi-konveksi/api/handler/akuntansi/hutang_piutang"
	hutangPiutangRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/hutang_piutang"
	hutangPiutangUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/hutang_piutang"

	kelompokAkunHandler "github.com/be-sistem-informasi-konveksi/api/handler/akuntansi/kelompok_akun"
	kelompokAkunRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kelompok_akun"
	kelompokAkunUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/kelompok_akun"

	akunHandler "github.com/be-sistem-informasi-konveksi/api/handler/akuntansi/akun"
	akunRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/akun"
	akunUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/akun"

	kontakHandler "github.com/be-sistem-informasi-konveksi/api/handler/akuntansi/kontak"
	dataBayarHutangPiutangRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/data_bayar_hutang_piutang"
	kontakRepo "github.com/be-sistem-informasi-konveksi/api/repository/akuntansi/mysql/gorm/kontak"
	kontakUsecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/kontak"

	"github.com/be-sistem-informasi-konveksi/pkg"
	"gorm.io/gorm"
)

type AkuntansiHandlerInit interface {
	Akun() akunHandler.AkunHandler
	HutangPiutang() hutangPiutangHandler.HutangPiutangHandler
	KelompokAkun() kelompokAkunHandler.KelompokAkunHandler
	Transaksi() transaksiHandler.TransaksiHandler
	Akuntansi() akuntansiHandler.AkuntansiHandler
	Kontak() kontakHandler.KontakHandler
}

type repoAkuntanInit struct {
	akun                   akunRepo.AkunRepo
	kelompokAkun           kelompokAkunRepo.KelompokAkunRepo
	transaksi              transaksiRepo.TransaksiRepo
	hutangPiutang          hutangPiutangRepo.HutangPiutangRepo
	kontak                 kontakRepo.KontakRepo
	dataBayarHutangPiutang dataBayarHutangPiutangRepo.DataBayarHutangPiutangRepo
	akuntansi              akuntansiRepo.AkuntansiRepo
}

type akuntansiHandlerInit struct {
	DB        *gorm.DB
	validator pkg.Validator
	ulid      pkg.UlidPkg
	excelize  pkg.ExcelizePkg
	repo      repoAkuntanInit
}

func NewAkuntansiHandlerInit(DB *gorm.DB, validator pkg.Validator, ulid pkg.UlidPkg, excelize pkg.ExcelizePkg) AkuntansiHandlerInit {
	r := repoAkuntanInit{
		akun:                   akunRepo.NewAkunRepo(DB),
		kelompokAkun:           kelompokAkunRepo.NewKelompokAkunRepo(DB),
		transaksi:              transaksiRepo.NewTransaksiRepo(DB),
		hutangPiutang:          hutangPiutangRepo.NewHutangPiutangRepo(DB),
		dataBayarHutangPiutang: dataBayarHutangPiutangRepo.NewDataBayarHutangPiutangRepo(DB),
		akuntansi:              akuntansiRepo.NewAkuntansiRepo(DB),
		kontak:                 kontakRepo.NewKontakRepo(DB),
	}
	return &akuntansiHandlerInit{DB, validator, ulid, excelize, r}
}

func (d *akuntansiHandlerInit) Akun() akunHandler.AkunHandler {
	uc := akunUsecase.NewAkunUsecase(d.repo.akun, d.ulid, d.repo.kelompokAkun)
	h := akunHandler.NewAkunHandler(uc, d.validator)
	return h
}

func (d *akuntansiHandlerInit) KelompokAkun() kelompokAkunHandler.KelompokAkunHandler {
	uc := kelompokAkunUsecase.NewKelompokAkunUsecase(d.repo.kelompokAkun, d.ulid)
	h := kelompokAkunHandler.NewKelompokAkunHandler(uc, d.validator)
	return h
}

func (d *akuntansiHandlerInit) Transaksi() transaksiHandler.TransaksiHandler {
	uc := transaksiUsecase.NewTransaksiUsecase(d.repo.transaksi, d.repo.akun, d.repo.hutangPiutang, d.repo.dataBayarHutangPiutang, d.ulid)
	h := transaksiHandler.NewTransaksiHandler(uc, d.validator)
	return h
}
func (d *akuntansiHandlerInit) HutangPiutang() hutangPiutangHandler.HutangPiutangHandler {
	uc := hutangPiutangUsecase.NewHutangPiutangUsecase(d.repo.hutangPiutang, d.repo.dataBayarHutangPiutang, d.repo.akun, d.repo.kontak, d.ulid)
	h := hutangPiutangHandler.NewHutangPiutangHandler(uc, d.validator)
	return h
}

func (d *akuntansiHandlerInit) Kontak() kontakHandler.KontakHandler {
	uc := kontakUsecase.NewKontakUsecase(d.repo.kontak, d.ulid)
	h := kontakHandler.NewKontakHandler(uc, d.validator)
	return h
}

func (d *akuntansiHandlerInit) Akuntansi() akuntansiHandler.AkuntansiHandler {
	uc := akuntansiUsecase.NewAkuntansiUsecase(d.repo.akuntansi, d.excelize)
	h := akuntansiHandler.NewAkuntansiHandler(uc, d.validator)
	return h
}

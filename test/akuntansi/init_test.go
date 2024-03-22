package test_akuntansi

import (
	"os"
	"testing"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

var dbt *gorm.DB
var token string
var app = fiber.New()

func cleanUp() {
	var ids []string
	for _, kelompok := range static_data.DataKelompokAkun {
		ids = append(ids, kelompok.ID)
	}
	dbt.Unscoped().Delete(&entity.KelompokAkun{}, "id NOT IN (?)", ids)
	dbt.Unscoped().Where("1 = 1").Delete(&entity.Kontak{})
	dbt.Unscoped().Where("1 = 1").Delete(&entity.Transaksi{})
}

func TestMain(m *testing.M) {
	dbt = test.GetDB()
	akuntansiH := handler_init.NewAkuntansiHandlerInit(dbt, test.Validator, test.UlidPkg)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	akuntansiRoute := route.NewAkuntansiRoute(akuntansiH, authMid)

	token = test.GetToken(dbt, authMid)

	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	akuntansiGroup := v1.Group("/akuntansi")
	akuntansiGroup.Route("", akuntansiRoute.Akuntansi)
	akuntansiGroup.Route("/akun", akuntansiRoute.Akun)
	akuntansiGroup.Route("/kontak", akuntansiRoute.Kontak)
	akuntansiGroup.Route("/kelompok_akun", akuntansiRoute.KelompokAkun)
	akuntansiGroup.Route("/transaksi", akuntansiRoute.Transaksi)
	akuntansiGroup.Route("/hutang_piutang", akuntansiRoute.HutangPiutang)
	// Run tests
	exitVal := m.Run()
	cleanUp()
	os.Exit(exitVal)
}

func TestEndPointAkuntansi(t *testing.T) {
	// kelompok akun
	AkuntansiCreateKelompokAkun(t)
	AkuntansiUpdateKelompokAkun(t)
	AkuntansiGetAllKelompokAkun(t)
	AkuntansiGetKelompokAkun(t)

	// akun
	AkuntansiCreateAkun(t)
	AkuntansiUpdateAkun(t)
	AkuntansiGetAllAkun(t)
	AkuntansiGetAkun(t)

	// kontak
	AkuntansiCreateKontak(t)
	AkuntansiUpdateKontak(t)
	AkuntansiGetAllKontak(t)

	// transaksi
	AkuntansiCreateTransaksi(t)
	AkuntansiUpdateTransaksi(t)
}

func TestEndPointDelete(t *testing.T) {
	AkuntansiDeleteKelompokAkun(t)
	AkuntansiDeleteAkun(t)
	AkuntansiDeleteKontak(t)
}

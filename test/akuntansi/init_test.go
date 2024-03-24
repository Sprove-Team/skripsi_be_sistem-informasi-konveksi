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
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

var dbt *gorm.DB
var tokens map[string]string
var app = fiber.New()

func cleanUp() {
	var ids []string
	for _, kelompok := range static_data.DataKelompokAkun {
		ids = append(ids, kelompok.ID)
	}
	if err := dbt.Unscoped().Delete(&entity.KelompokAkun{}, "id NOT IN (?)", ids).Error; err != nil {
		helper.LogsError(err)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.HutangPiutang{}).Error; err != nil {
		helper.LogsError(err)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.Transaksi{}).Error; err != nil {
		helper.LogsError(err)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.Kontak{}).Error; err != nil {
		helper.LogsError(err)
	}
}

// func setupInvoiceData() {
// 	payload := req_invoice.Create{
// 		KontakID: idKontak,
// 	}
// 	token := tokens[entity.RolesById[1]]
// 	test.GetJsonTestRequestResponse(app, "POST", "/api/v1/invoice", payload, &token)
// }

func TestMain(m *testing.M) {
	test.GetDB()
	dbt = test.DBT
	akuntansiH := handler_init.NewAkuntansiHandlerInit(dbt, test.Validator, test.UlidPkg)
	invoiceH := handler_init.NewInvoiceHandlerInit(dbt, test.Validator, test.UlidPkg, test.Encryptor)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	akuntansiRoute := route.NewAkuntansiRoute(akuntansiH, authMid)
	invoiceRoute := route.NewInvoiceRoute(invoiceH, authMid)

	test.GetTokens(dbt, authMid)
	tokens = test.Tokens
	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	akuntansiGroup := v1.Group("/akuntansi")
	akuntansiGroup.Route("/akun", akuntansiRoute.Akun)
	akuntansiGroup.Route("/kontak", akuntansiRoute.Kontak)
	akuntansiGroup.Route("/kelompok_akun", akuntansiRoute.KelompokAkun)
	akuntansiGroup.Route("/transaksi", akuntansiRoute.Transaksi)
	akuntansiGroup.Route("/hutang_piutang", akuntansiRoute.HutangPiutang)
	akuntansiGroup.Route("", akuntansiRoute.Akuntansi)
	invoiceGroup := v1.Group("/invoice")
	invoiceGroup.Route("/", invoiceRoute.Invoice)
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

	// hutang piutang
	AkuntansiCreateHutangPiutang(t)

	// transaksi
	AkuntansiCreateTransaksi(t)
	AkuntansiUpdateTransaksi(t) // TODO: add validasi tr hp
	AkuntansiGetAllTransaksi(t)
	AkuntansiGetTransaksi(t)
	AkuntansiGetHistoryTransaksi(t)

}

func TestEndPointDelete(t *testing.T) {
	AkuntansiDeleteKelompokAkun(t)
	AkuntansiDeleteAkun(t)
	AkuntansiDeleteTransaksi(t)
	AkuntansiDeleteKontak(t)
}

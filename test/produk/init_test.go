package test_produk

import (
	"os"
	"testing"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

var dbt *gorm.DB
var tokens map[string]string
var app = fiber.New()

func TestMain(m *testing.M) {
	test.GetDB()
	dbt = test.DBT
	produkH := handler_init.NewProdukHandlerInit(dbt, test.Validator, test.UlidPkg)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	produkRoute := route.NewProdukRoute(produkH, authMid)

	test.GetTokens(dbt, authMid)
	tokens = test.Tokens
	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	produkGroup := v1.Group("/produk")
	produkGroup.Route("/harga_detail", produkRoute.HargaDetailProduk)
	produkGroup.Route("/kategori", produkRoute.KategoriProduk)
	produkGroup.Route("/", produkRoute.Produk)

	// Run tests
	exitVal := m.Run()
	dbt.Unscoped().Where("1 = 1").Delete(&entity.KategoriProduk{})
	os.Exit(exitVal)
}

func TestEndPointProduk(t *testing.T) {
	//? kategori
	ProdukCreateKategori(t)
	ProdukUpdateKategori(t)
	ProdukGetAllKategori(t)
	ProdukGetKategori(t)

	//? produk
	ProdukCreate(t)
	ProdukUpdate(t)
	ProdukGetAll(t)
	ProdukGet(t)

	//? harga detail
	ProdukCreateHargaDetail(t)
	ProdukUpdateHargaDetail(t)
	ProdukGetAllHargaDetailByProdukId(t)

	//? delete
	ProdukDeleteKategori(t)
	ProdukDelete(t)
	ProdukDeleteHargaDetail(t)
}

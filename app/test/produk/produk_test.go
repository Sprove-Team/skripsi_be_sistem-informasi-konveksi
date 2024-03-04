package test_produk

import (
	"os"
	"testing"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/config"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/app/test"
	"github.com/be-sistem-informasi-konveksi/common/message"
	produkKategori "github.com/be-sistem-informasi-konveksi/common/request/produk/kategori"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var DBT *gorm.DB
var ulidPkg = pkg.NewUlidPkg()
var produkRoute route.ProdukRoute

func TestMain(m *testing.M) {
	// Initialize environment and dependencies
	config.LoadEnv()
	dbGormConfTest := config.DBGorm{
		DB_Username: os.Getenv("DB_USERNAME"),
		DB_Password: os.Getenv("DB_PASSWORD"),
		DB_Name:     os.Getenv("DB_NAME_TEST"),
		DB_Port:     os.Getenv("DB_PORT"),
		DB_Test:     true,
	}
	DBT = dbGormConfTest.InitDBGorm(ulidPkg)
	produkH := handler_init.NewProdukHandlerInit(DBT, pkg.NewValidator(), ulidPkg)
	userRepo := repo_user.NewUserRepo(DBT)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	produkRoute = route.NewProdukRoute(produkH, authMid)

	// Run tests
	exitVal := m.Run()

	// Clean up resources if needed
	// For example, you can close the database connection here

	// Exit with the test result
	os.Exit(exitVal)
}
func TestProdukCreateKategori(t *testing.T) {

	tests := []struct {
		name         string
		payload      produkKategori.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: produkKategori.Create{
				Nama: "kaos",
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name: "err: conflict",
			payload: produkKategori.Create{
				Nama: "kaos",
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name: "err: nama wajib diisi",
			payload: produkKategori.Create{
				Nama: "",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"nama wajib diisi"},
			},
		},
	}

	app := fiber.New()
	app.Use(recover.New())
	api := app.Group("/api/v1")
	produkGroup := api.Group("/produk")
	produkGroup.Route("/kategori", produkRoute.KategoriProduk)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/produk/kategori", tt.payload)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
	DBT.Unscoped().Where("1=1").Delete(&entity.KategoriProduk{})
}

// func TestProdukCreate(t *testing.T) {
// 	var c *fiber.Ctx
// 	test := test.TestData{
// 		Name: "sukses",
// 		Payload: produk.Create{
// 			Nama:       "kaos 45 dr",
// 			KategoriID: idKategori,
// 		},
// 		ExpectedCode: 201,
// 	}
// 	c.BodyParser(&data)
// 	err := produkH.HargaDetailProdukHandler().Create(c)
// 	if err != nil {

// 	}
// }

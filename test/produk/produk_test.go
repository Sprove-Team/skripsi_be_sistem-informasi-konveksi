package test_produk

import (
	"os"
	"testing"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_produk "github.com/be-sistem-informasi-konveksi/common/request/produk"
	produkKategori "github.com/be-sistem-informasi-konveksi/common/request/produk/kategori"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var produkRoute route.ProdukRoute
var dbt *gorm.DB
var token string

func TestMain(m *testing.M) {
	dbt = test.GetDB()
	produkH := handler_init.NewProdukHandlerInit(dbt, test.Validator, test.UlidPkg)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	produkRoute = route.NewProdukRoute(produkH, authMid)

	token = test.GetToken(dbt, authMid)

	// Run tests
	exitVal := m.Run()
	dbt.Unscoped().Where("1 = 1").Delete(&entity.KategoriProduk{})
	dbt.Unscoped().Where("1 = 1").Delete(&entity.Produk{})
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
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/produk/kategori", tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

func TestProdukUpdateKategori(t *testing.T) {
	data := new(entity.KategoriProduk)
	err := dbt.Select("id").First(data).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	tests := []struct {
		name         string
		payload      produkKategori.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: produkKategori.Update{
				ID:   data.ID,
				Nama: "jaket",
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name: "err: tidak ditemukan",
			payload: produkKategori.Update{
				ID:   "01HM4B8QBH7MWAVAYP10WN6PKA",
				Nama: "jaket",
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name: "err: ulid tidak valid",
			payload: produkKategori.Update{
				ID:   data.ID + "123",
				Nama: "jaket",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/produk/kategori/"+tt.payload.ID, tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

func TestProdukGetAllKategori(t *testing.T) {
	tests := []struct {
		name         string
		queryBody    string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses limit 1",
			queryBody:    "?limit=1",
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: ulid tidak valid",
			expectedCode: 400,
			queryBody:    "?next=01HQVTTJ1S2606JGTYYZ5NDKNR123",
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"next tidak berupa ulid yang valid"},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/produk/kategori"+tt.queryBody, nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			// fmt.Println("err -> ",body)

			var res []map[string]interface{}
			switch tt.name {
			case "sukses":
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.Greater(t, len(res), 0)
				assert.NotEmpty(t, res[0])
				assert.NotEmpty(t, res[0]["id"])
				assert.NotEmpty(t, res[0]["created_at"])
				assert.Equal(t, "jaket", res[0]["nama"])
				assert.Equal(t, tt.expectedBody.Status, body.Status)
			case "sukses limit 1":
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.Len(t, res, 1)
				assert.Equal(t, tt.expectedBody.Status, body.Status)
			default:
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

func TestProdukGetKategori(t *testing.T) {
	data := new(entity.KategoriProduk)
	err := dbt.Select("id").First(data).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	tests := []struct {
		id           string
		name         string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			id:           data.ID,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: tidak ditemukan",
			id:           "01HM4B8QBH7MWAVAYP10WN6PKA",
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name:         "err: ulid tidak valid",
			id:           "01HQVTTJ1S2606JGTYYZ5NDKNR123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/produk/kategori/"+tt.id, nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)

			var res map[string]interface{}
			switch tt.name {
			case "sukses":
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				assert.NotEmpty(t, res["id"])
				assert.NotEmpty(t, res["created_at"])
				assert.Equal(t, "jaket", res["nama"])
				assert.Equal(t, tt.expectedBody.Status, body.Status)
			default:
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

func TestProdukCreate(t *testing.T) {
	kategori := new(entity.KategoriProduk)
	err := dbt.Select("id").First(kategori).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	tests := []struct {
		name         string
		payload      req_produk.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_produk.Create{
				Nama:       "apparel premium jaket 24",
				KategoriID: kategori.ID,
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name: "err: conflict",
			payload: req_produk.Create{
				Nama:       "apparel premium jaket 24",
				KategoriID: kategori.ID,
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name: "err: kategori produk not found",
			payload: req_produk.Create{
				Nama:       "apparel premium jaket 25",
				KategoriID: "01HQVTTJ1S2606JGTYYZ5NDKNZ",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kategori produk tidak ditemukan"},
			},
		},
		{
			name: "err: ulid tidak valid",
			payload: req_produk.Create{
				Nama:       "apparel premium jaket 25",
				KategoriID: "01HQVTTJ1S2606JGTYYZ5NDKNZ123",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kategori id tidak berupa ulid yang valid"},
			},
		},
		{
			name: "err: wajib diisi",
			payload: req_produk.Create{
				Nama:       "",
				KategoriID: "",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"nama wajib diisi", "kategori id wajib diisi"},
			},
		},
	}

	app := fiber.New()
	app.Use(recover.New())
	api := app.Group("/api/v1")
	produkGroup := api.Group("/produk")
	produkGroup.Route("/", produkRoute.Produk)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/produk", tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

func TestProdukUpdate(t *testing.T) {
	produk := new(entity.Produk)
	err := dbt.Select("id").First(produk).Error
	if err != nil {
		helper.LogsError(err)
		return
	}

	id2 := test.UlidPkg.MakeUlid().String()
	err = dbt.Create(&entity.KategoriProduk{
		Base: entity.Base{
			ID: id2,
		},
		Nama: "test",
	}).Error

	if err != nil {
		helper.LogsError(err)
		return
	}
	tests := []struct {
		name         string
		payload      req_produk.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_produk.Update{
				Nama:       "apparel premium jaket 25",
				ID:         produk.ID,
				KategoriID: id2,
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name: "err: tidak ditemukan",
			payload: req_produk.Update{
				Nama:       "apparel premium jaket 25",
				ID:         "01HQVTTJ1S2606JGTYYZ5NDKNZ",
				KategoriID: id2,
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name: "err: kategori produk tidak ditemukan",
			payload: req_produk.Update{
				Nama:       "apparel premium jaket 25",
				ID:         produk.ID,
				KategoriID: "01HQVTTJ1S2606JGTYYZ5NDKNZ",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kategori produk tidak ditemukan"},
			},
		},
		{
			name: "err: ulid tidak valid",
			payload: req_produk.Update{
				Nama:       "apparel premium jaket 25",
				ID:         "01HQVTTJ1S2606JGTYYZ5NDKNZ123",
				KategoriID: "01HQVTTJ1S2606JGTYYZ5NDKNZ123",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid", "kategori id tidak berupa ulid yang valid"},
			},
		},
	}

	app := fiber.New()
	app.Use(recover.New())
	api := app.Group("/api/v1")
	produkGroup := api.Group("/produk")
	produkGroup.Route("/", produkRoute.Produk)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/produk/"+tt.payload.ID, tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

// func TestProdukDelete(t *testing.T) {
// 	kategori := new(entity.KategoriProduk)
// 	err := dbt.Select("id").First(kategori).Error
// 	if err != nil {
// 		helper.LogsError(err)
// 		return
// 	}
// 	tests := []struct {
// 		name         string
// 		payload      req_produk.Create
// 		expectedBody test.Response
// 		expectedCode int
// 	}{
// 		{
// 			name: "sukses",
// 			payload: req_produk.Create{
// 				Nama:       "apparel premium jaket 24",
// 				KategoriID: kategori.ID,
// 			},
// 			expectedCode: 201,
// 			expectedBody: test.Response{
// 				Status: message.Created,
// 				Code:   201,
// 			},
// 		},
// 	}

// 	app := fiber.New()
// 	app.Use(recover.New())
// 	api := app.Group("/api/v1")
// 	produkGroup := api.Group("/produk")
// 	produkGroup.Route("/", produkRoute.Produk)

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/produk", tt.payload, &token)
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expectedCode, code)
// 			assert.Equal(t, tt.expectedBody, body)
// 		})
// 	}
// }

package test_produk

import (
	"os"
	"strings"
	"testing"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_produk "github.com/be-sistem-informasi-konveksi/common/request/produk"
	req_produk_harga_detail "github.com/be-sistem-informasi-konveksi/common/request/produk/harga_detail"
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
var app = fiber.New()

func TestMain(m *testing.M) {
	dbt = test.GetDB()
	produkH := handler_init.NewProdukHandlerInit(dbt, test.Validator, test.UlidPkg)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	produkRoute = route.NewProdukRoute(produkH, authMid)

	token = test.GetToken(dbt, authMid)

	// app kategori
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/produk/kategori", tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

var idKategori string

func TestProdukUpdateKategori(t *testing.T) {
	kategori := new(entity.KategoriProduk)
	err := dbt.Select("id").First(kategori).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	idKategori = kategori.ID
	tests := []struct {
		name         string
		payload      req_produk.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_produk.Update{
				ID:   kategori.ID,
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
			payload: req_produk.Update{
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
			payload: req_produk.Update{
				ID:   kategori.ID + "123",
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
			name:         "sukses with: filter nama",
			queryBody:    "?nama=jaket",
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/produk/kategori"+tt.queryBody, nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)

			var res []map[string]interface{}
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.Greater(t, len(res), 0)
				assert.NotEmpty(t, res[0])
				assert.NotEmpty(t, res[0]["id"])
				assert.NotEmpty(t, res[0]["created_at"])
				assert.NotEmpty(t, res[0]["nama"])
				assert.Equal(t, tt.expectedBody.Status, body.Status)
				switch tt.name {
				case "sukses with: filter nama":
					assert.Contains(t, res[0]["nama"], "jaket")
				case "sukses limit 1":
					assert.Len(t, res, 1)
				}
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}

		})
	}
}

func TestProdukGetKategori(t *testing.T) {
	tests := []struct {
		id           string
		name         string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			id:           idKategori,
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
				KategoriID: idKategori,
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
				KategoriID: idKategori,
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/produk", tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

var idProduk string
var idKategori2 string

func TestProdukUpdate(t *testing.T) {
	produk := new(entity.Produk)
	err := dbt.Select("id").First(produk).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	idProduk = produk.ID

	idKategori2 = test.UlidPkg.MakeUlid().String()
	err = dbt.Create(&entity.KategoriProduk{
		Base: entity.Base{
			ID: idKategori2,
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
				KategoriID: idKategori2,
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
				KategoriID: idKategori2,
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/produk/"+tt.payload.ID, tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

var idProdukHasHargaDetail string

func TestProdukGetAll(t *testing.T) {
	idProdukHasHargaDetail = test.UlidPkg.MakeUlid().String()
	produkCreate := &entity.Produk{
		Base: entity.Base{
			ID: idProdukHasHargaDetail,
		},
		KategoriProdukID: idKategori,
		Nama:             "produk coba",
		HargaDetails: []entity.HargaDetailProduk{
			{
				Base: entity.Base{
					ID: test.UlidPkg.MakeUlid().String(),
				},
				QTY:   10,
				Harga: 200000,
			},
		},
	}

	err := dbt.Create(produkCreate).Error
	if err != nil {
		helper.LogsError(err)
		return
	}

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
			name:         "sukses with: filter nama",
			queryBody:    "?nama=apparel",
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses with: filter harga detail `NOT_EMPTY`",
			queryBody:    "?harga_detail=NOT_EMPTY",
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses with: filter harga detail `EMPTY`",
			queryBody:    "?harga_detail=EMPTY",
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
		{
			name:         "err: value harus berupa [EMPTY,NOT_EMPTY]",
			expectedCode: 400,
			queryBody:    "?harga_detail=abcd",
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"harga detail harus berupa salah satu dari [EMPTY,NOT_EMPTY]"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/produk"+tt.queryBody, nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			var res []map[string]any
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.Greater(t, len(res), 0)
				assert.NotEmpty(t, res[0])
				assert.NotEmpty(t, res[0]["id"])
				assert.NotEmpty(t, res[0]["created_at"])
				assert.NotEmpty(t, res[0]["nama"])
				assert.NotEmpty(t, res[0]["kategori"])
				assert.NotEmpty(t, res[0]["kategori"].(map[string]any)["id"])
				assert.NotEmpty(t, res[0]["kategori"].(map[string]any)["created_at"])
				assert.NotEmpty(t, res[0]["kategori"].(map[string]any)["nama"])
				assert.Equal(t, tt.expectedBody.Status, body.Status)
				switch tt.name {
				case "sukses with: filter harga detail `NOT_EMPTY`":
					assert.NotEmpty(t, res[0]["harga_detail"])
					hargaDetail := res[0]["harga_detail"].([]any)
					assert.Greater(t, len(hargaDetail), 0)
					for _, v := range hargaDetail {
						assert.NotEmpty(t, v.(map[string]any)["id"])
						assert.NotEmpty(t, v.(map[string]any)["created_at"])
						assert.NotEmpty(t, v.(map[string]any)["qty"])
						assert.NotEmpty(t, v.(map[string]any)["harga"])
					}
				case "sukses with: filter harga detail `EMPTY`":
					assert.Empty(t, res[0]["harga_detail"])
				case "sukses with: filter nama":
					assert.Contains(t, res[0]["nama"], "apparel")
				case "sukses limit 1":
					assert.Len(t, res, 1)
				}
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}

		})
	}
}

func TestProdukGet(t *testing.T) {
	tests := []struct {
		id           string
		name         string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses without: harga detail",
			id:           idProduk,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses with: harga detail",
			id:           idProdukHasHargaDetail,
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/produk/"+tt.id, nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)

			var res map[string]any
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.Greater(t, len(res), 0)
				assert.NotEmpty(t, res)
				assert.NotEmpty(t, res["id"])
				assert.NotEmpty(t, res["created_at"])
				assert.NotEmpty(t, res["nama"])
				assert.NotEmpty(t, res["kategori"])
				assert.NotEmpty(t, res["kategori"].(map[string]any)["id"])
				assert.NotEmpty(t, res["kategori"].(map[string]any)["created_at"])
				assert.NotEmpty(t, res["kategori"].(map[string]any)["nama"])
				assert.Equal(t, tt.expectedBody.Status, body.Status)
				switch tt.name {
				case "sukses without: harga detail":
					assert.Empty(t, res["harga_detail"])
				case "sukses with: harga detail":
					assert.NotEmpty(t, res["harga_detail"])
					hargaDetail := res["harga_detail"].([]any)
					assert.Greater(t, len(hargaDetail), 0)
					for _, v := range hargaDetail {
						assert.NotEmpty(t, v.(map[string]any)["id"])
						assert.NotEmpty(t, v.(map[string]any)["created_at"])
						assert.NotEmpty(t, v.(map[string]any)["qty"])
						assert.NotEmpty(t, v.(map[string]any)["harga"])
					}
				}
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

func TestProdukCreateHargaDetail(t *testing.T) {
	tests := []struct {
		name         string
		payload      req_produk_harga_detail.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_produk_harga_detail.Create{
				ProdukId: idProduk,
				QTY:      100,
				Harga:    10000,
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name: "err: conflict",
			payload: req_produk_harga_detail.Create{
				ProdukId: idProduk,
				QTY:      100,
				Harga:    10000,
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name: "err: produk not found",
			payload: req_produk_harga_detail.Create{
				ProdukId: idKategori,
				QTY:      101,
				Harga:    999,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.ProdukNotFound},
			},
		},
		{
			name: "err: ulid tidak valid",
			payload: req_produk_harga_detail.Create{
				ProdukId: idProduk + "123",
				QTY:      102,
				Harga:    998,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"produk id tidak berupa ulid yang valid"},
			},
		},
		{
			name: "err: wajib diisi",
			payload: req_produk_harga_detail.Create{
				ProdukId: "",
				QTY:      0,
				Harga:    0,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"produk id wajib diisi", "qty wajib diisi", "harga wajib diisi"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/produk/harga_detail", tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

var idHargaDetail string

func TestProdukUpdateHargaDetail(t *testing.T) {
	hargaDetail := new(entity.HargaDetailProduk)
	err := dbt.Select("id").First(hargaDetail).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	idHargaDetail = hargaDetail.ID

	tests := []struct {
		name         string
		payload      req_produk_harga_detail.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_produk_harga_detail.Update{
				ID:    idHargaDetail,
				QTY:   200,
				Harga: 50000,
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name: "err: tidak ditemukan",
			payload: req_produk_harga_detail.Update{
				ID:    idKategori,
				QTY:   201,
				Harga: 4999,
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name: "err: ulid tidak valid",
			payload: req_produk_harga_detail.Update{
				ID:    idHargaDetail + "123",
				QTY:   202,
				Harga: 4998,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/produk/harga_detail/"+tt.payload.ID, tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

func TestProdukGetAllHargaDetailByProdukId(t *testing.T) {
	tests := []struct {
		name         string
		idProduk     string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			idProduk:     idProdukHasHargaDetail,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: produk not found",
			expectedCode: 400,
			idProduk:     idKategori,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.ProdukNotFound},
			},
		},
		{
			name:         "err: ulid tidak valid",
			expectedCode: 400,
			idProduk:     idProduk + "123",
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"produk id tidak berupa ulid yang valid"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/produk/harga_detail/"+tt.idProduk, nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			var res []map[string]any
			switch tt.name {
			case "sukses":
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.Greater(t, len(res), 0)
				for _, v := range res {
					assert.NotEmpty(t, v)
					assert.NotEmpty(t, v["id"])
					assert.NotEmpty(t, v["produk_id"])
					assert.NotEmpty(t, v["qty"])
					assert.NotEmpty(t, v["harga"])
				}
				assert.Equal(t, tt.expectedBody.Status, body.Status)
			default:
				assert.Equal(t, tt.expectedBody, body)
			}

		})
	}
}

// ! All delete case
func TestProdukDeleteKategori(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			id:           idKategori,
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
			id:           idKategori + "123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/produk/kategori/"+tt.id, nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

func TestProdukDeleteHargaDetail(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			id:           idHargaDetail,
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
			id:           idHargaDetail + "123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/produk/harga_detail/"+tt.id, nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

func TestProdukDelete(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			id:           idProduk,
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
			id:           idProduk + "123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/produk/"+tt.id, nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

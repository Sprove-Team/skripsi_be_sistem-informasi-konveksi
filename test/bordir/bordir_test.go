package test_bordir

import (
	"net/url"
	"os"
	"strings"
	"testing"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_bordir "github.com/be-sistem-informasi-konveksi/common/request/bordir"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var bordirRoute route.BordirRoute
var dbt *gorm.DB
var token string
var app = fiber.New()

func TestMain(m *testing.M) {
	dbt = test.GetDB()
	bordirH := handler_init.NewBordirHandlerInit(dbt, test.Validator, test.UlidPkg)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	bordirRoute = route.NewBordirRoute(bordirH, authMid)

	token = test.GetToken(dbt, authMid)

	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	produkGroup := v1.Group("/bordir")
	produkGroup.Route("/", bordirRoute.Bordir)
	// Run tests
	exitVal := m.Run()
	dbt.Unscoped().Where("1 = 1").Delete(&entity.Bordir{})
	os.Exit(exitVal)
}

func TestBordirCreate(t *testing.T) {
	tests := []struct {
		name         string
		payload      req_bordir.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_bordir.Create{
				Nama:  "Bordir 2 baris",
				Harga: 40000,
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name: "err: conflict",
			payload: req_bordir.Create{
				Nama:  "Bordir 2 baris",
				Harga: 40000,
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name: "err: nama wajib diisi",
			payload: req_bordir.Create{
				Nama:  "",
				Harga: 0,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"nama wajib diisi", "harga wajib diisi"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/bordir", tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

var idBordir string

func TestBordirUpdate(t *testing.T) {
	bordir := new(entity.Bordir)
	err := dbt.Select("id").First(bordir).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	idBordir = bordir.ID
	tests := []struct {
		name         string
		payload      req_bordir.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_bordir.Update{
				ID:    bordir.ID,
				Nama:  "Bordir 1 baris",
				Harga: 20000,
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name: "err: tidak ditemukan",
			payload: req_bordir.Update{
				ID:    "01HM4B8QBH7MWAVAYP10WN6PKA",
				Nama:  "Bordir 3 baris",
				Harga: 20001,
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name: "err: ulid tidak valid",
			payload: req_bordir.Update{
				ID:    bordir.ID + "123",
				Nama:  "Bordir 4 baris",
				Harga: 20004,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/bordir/"+tt.payload.ID, tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

func TestBordirGetAll(t *testing.T) {
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
			queryBody:    "?nama=Bordir+1+baris",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/bordir"+tt.queryBody, nil, &token)
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
				assert.NotEmpty(t, res[0]["harga"])
				assert.Equal(t, tt.expectedBody.Status, body.Status)
				switch tt.name {
				case "sukses with: filter nama":
					v, err := url.ParseQuery(tt.queryBody[1:])
					assert.NoError(t, err)
					assert.Contains(t, res[0]["nama"], v.Get("nama"))
				case "sukses limit 1":
					assert.Len(t, res, 1)
				}
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}

		})
	}
}

func TestBordirDelete(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			id:           idBordir,
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
			id:           idBordir + "123",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/bordir/"+tt.id, nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}

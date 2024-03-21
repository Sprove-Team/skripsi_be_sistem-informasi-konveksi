package test_user

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_user "github.com/be-sistem-informasi-konveksi/common/request/user"
	req_user_jenis_spv "github.com/be-sistem-informasi-konveksi/common/request/user/jenis_spv"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var userRoute route.UserRoute
var dbt *gorm.DB
var token string
var app = fiber.New()

func TestMain(m *testing.M) {
	dbt = test.GetDB()
	userH := handler_init.NewUserHandlerInit(dbt, test.Validator, test.UlidPkg, test.Encryptor)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	userRoute = route.NewUserRoute(userH, authMid)

	token = test.GetToken(dbt, authMid)

	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	userGroup := v1.Group("/user")
	userGroup.Route("/jenis_spv", userRoute.JenisSpv)
	userGroup.Route("/", userRoute.User)
	// Run tests
	exitVal := m.Run()
	dbt.Unscoped().Where("1 = 1").Delete(&entity.JenisSpv{})
	dbt.Unscoped().Where("1 = 1").Delete(&entity.User{})
	os.Exit(exitVal)
}

// ? JENIS SPV

func TestUserCreateSpv(t *testing.T) {
	tests := []struct {
		name         string
		payload      req_user_jenis_spv.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_user_jenis_spv.Create{
				Nama: "belanja",
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name: "err: conflict",
			payload: req_user_jenis_spv.Create{
				Nama: "belanja",
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name: "err: wajib diisi",
			payload: req_user_jenis_spv.Create{
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
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/user/jenis_spv", tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			if len(tt.expectedBody.ErrorsMessages) > 0 {
				for _, v := range tt.expectedBody.ErrorsMessages {
					assert.Contains(t, body.ErrorsMessages, v)
				}
				assert.Equal(t, tt.expectedBody.Status, body.Status)
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

var idSpv string

func TestUserUpdateSpv(t *testing.T) {
	spv := new(entity.JenisSpv)
	err := dbt.Select("id").First(spv).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	spvConflict := &entity.JenisSpv{
		Base: entity.Base{
			ID: test.UlidPkg.MakeUlid().String(),
		},
		Nama: "conflict",
	}
	err = dbt.Create(spvConflict).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	idSpv = spv.ID
	tests := []struct {
		name         string
		payload      req_user_jenis_spv.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_user_jenis_spv.Update{
				ID:   spv.ID,
				Nama: "bordir",
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name: "err: conflict",
			payload: req_user_jenis_spv.Update{
				ID:   spv.ID,
				Nama: spvConflict.Nama,
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name: "err: tidak ditemukan",
			payload: req_user_jenis_spv.Update{
				ID:   "01HM4B8QBH7MWAVAYP10WN6PKB",
				Nama: "bordir2",
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name: "err: ulid tidak valid",
			payload: req_user_jenis_spv.Update{
				ID:   spv.ID + "123",
				Nama: "bordir3",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/user/jenis_spv/"+tt.payload.ID, tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			if len(tt.expectedBody.ErrorsMessages) > 0 {
				for _, v := range tt.expectedBody.ErrorsMessages {
					assert.Contains(t, body.ErrorsMessages, v)
				}
				assert.Equal(t, tt.expectedBody.Status, body.Status)
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

func TestUserGetAllSpv(t *testing.T) {
	tests := []struct {
		name         string
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/user/jenis_spv", nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			var res []any
			err = mapstructure.Decode(body.Data, &res)
			assert.NoError(t, err)
			assert.Greater(t, len(res), 0)
			assert.NotEmpty(t, res[0])
			for _, v := range res {
				assert.NotEmpty(t, v.(map[string]any)["id"])
				assert.NotEmpty(t, v.(map[string]any)["created_at"])
				assert.NotEmpty(t, v.(map[string]any)["nama"])
			}
			assert.Equal(t, tt.expectedBody.Status, body.Status)

		})
	}
}

// ? USER
func TestUserCreate(t *testing.T) {

	tests := []struct {
		name         string
		payload      req_user.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses: create not spv user",
			payload: req_user.Create{
				Nama:     "manager_produksi",
				Username: "manager_produksi",
				Password: "manager_produksi",
				Role:     "MANAJER_PRODUKSI",
				Alamat:   "test1234",
				NoTelp:   "+62895397290606",
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name: "sukses: create spv user",
			payload: req_user.Create{
				Nama:       "supervisorbordir",
				Username:   "supervisorbordir",
				Password:   "supervisorbordir",
				Role:       "SUPERVISOR",
				Alamat:     "test1234",
				NoTelp:     "+62895897290606",
				JenisSpvID: idSpv,
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name: "err: conflict",
			payload: req_user.Create{
				Nama:     "manager_produksi",
				Username: "manager_produksi",
				Password: "manager_produksi",
				Role:     "MANAJER_PRODUKSI",
				Alamat:   "test1234",
				NoTelp:   "+62895397290606",
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name: "err: wajib diisi",
			payload: req_user.Create{
				Nama:     "",
				Username: "",
				Password: "",
				Role:     "",
				Alamat:   "",
				NoTelp:   "",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status: fiber.ErrBadRequest.Message,
				Code:   400,
				ErrorsMessages: []string{
					"nama wajib diisi",
					"role wajib diisi",
					"username wajib diisi",
					"password wajib diisi",
					"no telp wajib diisi",
					"alamat wajib diisi",
				},
			},
		},
		{
			name: "err: role harus berupa [DIREKTUR,ADMIN,BENDAHARA,MANAJER_PRODUKSI,SUPERVISOR]",
			payload: req_user.Create{
				Nama:     "manager_produksi2",
				Username: "manager_produksi2",
				Password: "manager_produksi2",
				Role:     "asdf",
				Alamat:   "test1234",
				NoTelp:   "+62895397290606",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"role harus berupa salah satu dari [DIREKTUR,ADMIN,BENDAHARA,MANAJER_PRODUKSI,SUPERVISOR]"},
			},
		},
		{
			name: "err: no telp harus berformat e164",
			payload: req_user.Create{
				Nama:     "manager_produksi3",
				Username: "manager_produksi3",
				Password: "manager_produksi3",
				Role:     "MANAJER_PRODUKSI",
				Alamat:   "test1234",
				NoTelp:   "0895397290606",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"no telp harus berformat e164"},
			},
		},
		{
			name: "err: panjang minimal password adalah 6 karakter",
			payload: req_user.Create{
				Nama:     "manager_produksi4",
				Username: "manager_produksi4",
				Password: "a4",
				Role:     "MANAJER_PRODUKSI",
				Alamat:   "test1234",
				NoTelp:   "+62895397290606",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"panjang minimal password adalah 6 karakter"},
			},
		},
		{
			name: "err: jenis spv id wajib diisi ketika role supervisor",
			payload: req_user.Create{
				Nama:     "supervisorbordir2",
				Username: "supervisorbordir2",
				Password: "supervisorbordir2",
				Role:     "SUPERVISOR",
				Alamat:   "test1234",
				NoTelp:   "+62895397290606",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"jenis spv id wajib diisi ketika role supervisor"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/user", tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			if len(tt.expectedBody.ErrorsMessages) > 0 {
				for _, v := range tt.expectedBody.ErrorsMessages {
					assert.Contains(t, body.ErrorsMessages, v)
				}
				assert.Equal(t, tt.expectedBody.Status, body.Status)
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

var idUser string
var idUserSpv string
var idSpv2 string

func TestUserUpdate(t *testing.T) {

	user := new(entity.User)
	err := dbt.Select("id").First(user, "ROLE NOT IN (?) ", []string{entity.RolesById[1], entity.RolesById[5]}).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	userSpv := new(entity.User)
	err = dbt.Select("id").First(userSpv, "jenis_spv_id = ?", idSpv).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	spv2 := &entity.JenisSpv{
		Base: entity.Base{
			ID: test.UlidPkg.MakeUlid().String(),
		},
		Nama: "belanja",
	}
	err = dbt.Create(spv2).Error
	if err != nil {
		helper.LogsError(err)
		return
	}

	idUser = user.ID
	idUserSpv = userSpv.ID
	idSpv2 = spv2.ID

	tests := []struct {
		name         string
		payload      req_user.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses: update not spv user",
			payload: req_user.Update{
				ID:       idUser,
				Nama:     "admin2",
				Username: "admin2",
				Password: "admin22",
				Role:     "ADMIN",
				Alamat:   "test1234",
				NoTelp:   "+62898397290606",
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name: "sukses: update spv user",
			payload: req_user.Update{
				ID:         idUserSpv,
				Nama:       "supervisorbelanja",
				Username:   "supervisorbelanja",
				Password:   "supervisorbelanja",
				Role:       "SUPERVISOR",
				Alamat:     "test1234",
				NoTelp:     "+62895897290606",
				JenisSpvID: idSpv2,
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name: "err: conflict",
			payload: req_user.Update{
				ID:       idUser,
				Nama:     "admin2",
				Username: "direktur",
				Password: "admin22",
				Role:     "ADMIN",
				Alamat:   "test1234",
				NoTelp:   "+62898397290606",
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name: "err: role harus berupa [DIREKTUR,ADMIN,BENDAHARA,MANAJER_PRODUKSI,SUPERVISOR]",
			payload: req_user.Update{
				ID:       idSpv,
				Nama:     "manager_produksi2",
				Username: "manager_produksi2",
				Password: "manager_produksi2",
				Role:     "asdf",
				Alamat:   "test1234",
				NoTelp:   "+62895397290606",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"role harus berupa salah satu dari [DIREKTUR,ADMIN,BENDAHARA,MANAJER_PRODUKSI,SUPERVISOR]"},
			},
		},
		{
			name: "err: no telp harus berformat e164",
			payload: req_user.Update{
				ID:       idSpv,
				Nama:     "manager_produksi3",
				Username: "manager_produksi3",
				Password: "manager_produksi3",
				Role:     "MANAJER_PRODUKSI",
				Alamat:   "test1234",
				NoTelp:   "0895397290606",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"no telp harus berformat e164"},
			},
		},
		{
			name: "err: panjang minimal password adalah 6 karakter",
			payload: req_user.Update{
				ID:       idSpv,
				Nama:     "manager_produksi4",
				Username: "manager_produksi4",
				Password: "a4",
				Role:     "MANAJER_PRODUKSI",
				Alamat:   "test1234",
				NoTelp:   "+62895397290606",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"panjang minimal password adalah 6 karakter"},
			},
		},
		{
			name: "err: jenis spv id wajib diisi ketika role supervisor",
			payload: req_user.Update{
				ID:       idUser,
				Nama:     "supervisorbelanja2",
				Username: "supervisorbelanja2",
				Password: "supervisorbelanja2",
				Role:     "SUPERVISOR",
				Alamat:   "test1234",
				NoTelp:   "+62895397290606",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"jenis spv id wajib diisi ketika role supervisor"},
			},
		},
		{
			name: "err: ulid tidak valid",
			payload: req_user.Update{
				ID:         idUser + "123",
				Nama:       "supervisorbelanja2",
				Username:   "supervisorbelanja2",
				Password:   "supervisorbelanja2",
				Role:       "SUPERVISOR",
				Alamat:     "test1234",
				NoTelp:     "+62895397290606",
				JenisSpvID: idSpv + "123",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid", "jenis spv id tidak berupa ulid yang valid"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/user/"+tt.payload.ID, tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			if len(tt.expectedBody.ErrorsMessages) > 0 {
				for _, v := range tt.expectedBody.ErrorsMessages {
					assert.Contains(t, body.ErrorsMessages, v)
				}
				assert.Equal(t, tt.expectedBody.Status, body.Status)
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

func TestUserGetAll(t *testing.T) {
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
			name:         "sukses with filter search",
			queryBody:    fmt.Sprintf("?search[nama]=supervisorbelanja&search[jenis_spv_id]=%s&search[role]=SUPERVISOR&search[username]=supervisorbelanja&search[alamat]=test123&search[no_telp]=%s6289589729", idSpv2, "%2B"),
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses with next",
			queryBody:    "?next=" + idUser,
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
			queryBody:    "?next=01HQVTTJ1S2606JGTYYZ5NDKNR123&search[jenis_spv_id]=ABCDSI123124AASDDC",
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"next tidak berupa ulid yang valid", "jenis spv id tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "err: role harus berupa [DIREKTUR,ADMIN,BENDAHARA,MANAJER_PRODUKSI,SUPERVISOR]",
			expectedCode: 400,
			queryBody:    "?search[role]=abcd",
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"role harus berupa salah satu dari [DIREKTUR,ADMIN,BENDAHARA,MANAJER_PRODUKSI,SUPERVISOR]"},
			},
		},
		{
			name:         "err: no telp harus berformat e164",
			queryBody:    "?search[no_telp]=08912312313",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"no telp harus berformat e164"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/user"+tt.queryBody, nil, &token)
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
				assert.NotEmpty(t, res[0]["role"])
				assert.NotEmpty(t, res[0]["username"])
				assert.NotEmpty(t, res[0]["no_telp"])
				assert.NotEmpty(t, res[0]["alamat"])
				assert.Equal(t, tt.expectedBody.Status, body.Status)
				switch tt.name {
				case "sukses with filter search":
					v, err := url.ParseQuery(tt.queryBody[1:])
					assert.NoError(t, err)
					assert.Contains(t, res[0]["nama"], v.Get("search[nama]"))
					assert.Equal(t, res[0]["role"], v.Get("search[role]"))
					assert.Contains(t, res[0]["username"], v.Get("search[username]"))
					assert.Contains(t, res[0]["alamat"], v.Get("search[alamat]"))
					assert.Contains(t, res[0]["no_telp"], v.Get("search[no_telp]"))
					jenisSpv := res[0]["jenis_spv"].(map[string]any)
					assert.NotEmpty(t, jenisSpv)
					assert.NotEmpty(t, jenisSpv["id"])
					assert.NotEmpty(t, jenisSpv["created_at"])
					assert.NotEmpty(t, jenisSpv["nama"])
					assert.Equal(t, jenisSpv["id"], v.Get("search[jenis_spv_id]"))
				case "sukses limit 1":
					assert.Len(t, res, 1)
				case "sukses with next":
					assert.NotEqual(t, idUser, res[0]["id"])
				}
			} else {
				if len(tt.expectedBody.ErrorsMessages) > 0 {
					for _, v := range tt.expectedBody.ErrorsMessages {
						assert.Contains(t, body.ErrorsMessages, v)
					}
					assert.Equal(t, tt.expectedBody.Status, body.Status)
				} else {
					assert.Equal(t, tt.expectedBody, body)
				}
			}

		})
	}
}

// ! ALL DELETE TEST

func TestUserDeleteSpv(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			id:           idSpv,
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
			id:           idSpv + "123",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/user/jenis_spv/"+tt.id, nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			if len(tt.expectedBody.ErrorsMessages) > 0 {
				for _, v := range tt.expectedBody.ErrorsMessages {
					assert.Contains(t, body.ErrorsMessages, v)
				}
				assert.Equal(t, tt.expectedBody.Status, body.Status)
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

func TestUserDelete(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			id:           idUser,
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
			id:           idUser + "123",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/user/"+tt.id, nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			if len(tt.expectedBody.ErrorsMessages) > 0 {
				for _, v := range tt.expectedBody.ErrorsMessages {
					assert.Contains(t, body.ErrorsMessages, v)
				}
				assert.Equal(t, tt.expectedBody.Status, body.Status)
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

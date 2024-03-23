package test_user

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_user "github.com/be-sistem-informasi-konveksi/common/request/user"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func UserCreate(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		payload      req_user.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses: create not spv user",
			token: tokens[entity.RolesById[1]],
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
			name:  "sukses: create spv user",
			token: tokens[entity.RolesById[1]],
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
			name:  "err: conflict",
			token: tokens[entity.RolesById[1]],
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
			name:  "err: wajib diisi",
			token: tokens[entity.RolesById[1]],
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
			name:  "err: role harus berupa [DIREKTUR,ADMIN,BENDAHARA,MANAJER_PRODUKSI,SUPERVISOR]",
			token: tokens[entity.RolesById[1]],
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
			name:  "err: no telp harus berformat e164",
			token: tokens[entity.RolesById[1]],
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
			name:  "err: panjang minimal password adalah 6 karakter",
			token: tokens[entity.RolesById[1]],
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
			name:  "err: jenis spv id wajib diisi ketika role supervisor",
			token: tokens[entity.RolesById[1]],
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
		{
			name:         "err: authorization " + entity.RolesById[2],
			payload:      req_user.Create{},
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			payload:      req_user.Create{},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			payload:      req_user.Create{},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			payload:      req_user.Create{},
			token:        tokens[entity.RolesById[5]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/user", tt.payload, &tt.token)
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

func UserUpdate(t *testing.T) {

	user := new(entity.User)
	err := dbt.Select("id").First(user, "ROLE NOT IN (?) AND id NOT IN (?)", []string{entity.RolesById[1], entity.RolesById[5]}, idsDefaultUser).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	userSpv := new(entity.User)
	err = dbt.Select("id").First(userSpv, "jenis_spv_id = ? AND id NOT IN (?)", idSpv, idsDefaultUser).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	spv2 := &entity.JenisSpv{
		Base: entity.Base{
			ID: test.UlidPkg.MakeUlid().String(),
		},
		Nama: "test_belanja2",
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
		token        string
		payload      req_user.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses: update not spv user",
			token: tokens[entity.RolesById[1]],
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
			name:  "sukses: update spv user",
			token: tokens[entity.RolesById[1]],
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
			name:  "err: conflict",
			token: tokens[entity.RolesById[1]],
			payload: req_user.Update{
				ID:       idUser,
				Nama:     static_data.DefaultUsers[0].Nama,
				Username: static_data.DefaultUsers[0].Username,
				Password: "admin22",
				Role:     "ADMIN",
				Alamat:   static_data.DefaultUsers[0].Alamat,
				NoTelp:   static_data.DefaultUsers[0].NoTelp,
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name:  "err: role harus berupa [DIREKTUR,ADMIN,BENDAHARA,MANAJER_PRODUKSI,SUPERVISOR]",
			token: tokens[entity.RolesById[1]],
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
			name:  "err: no telp harus berformat e164",
			token: tokens[entity.RolesById[1]],
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
			name:  "err: panjang minimal password adalah 6 karakter",
			token: tokens[entity.RolesById[1]],
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
			name:  "err: jenis spv id wajib diisi ketika role supervisor",
			token: tokens[entity.RolesById[1]],
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
			name:  "err: ulid tidak valid",
			token: tokens[entity.RolesById[1]],
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

		{
			name:         "err: authorization " + entity.RolesById[2],
			payload:      req_user.Update{},
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			payload:      req_user.Update{},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			payload:      req_user.Update{},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			payload:      req_user.Update{},
			token:        tokens[entity.RolesById[5]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/user/"+tt.payload.ID, tt.payload, &tt.token)
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

func UserGetAll(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		queryBody    string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses with filter search",
			token:        tokens[entity.RolesById[1]],
			queryBody:    fmt.Sprintf("?search[nama]=supervisorbelanja&search[jenis_spv_id]=%s&search[role]=SUPERVISOR&search[username]=supervisorbelanja&search[alamat]=test123&search[no_telp]=%s6289589729", idSpv2, "%2B"),
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses with next",
			token:        tokens[entity.RolesById[1]],
			queryBody:    "?next=" + idUser,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses limit 1",
			token:        tokens[entity.RolesById[1]],
			queryBody:    "?limit=1",
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: ulid tidak valid",
			token:        tokens[entity.RolesById[1]],
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
			token:        tokens[entity.RolesById[1]],
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
			token:        tokens[entity.RolesById[1]],
			queryBody:    "?search[no_telp]=08912312313",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"no telp harus berformat e164"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[2],
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			token:        tokens[entity.RolesById[5]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/user"+tt.queryBody, nil, &tt.token)
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

func UserGet(t *testing.T) {
	tests := []struct {
		id           string
		token        string
		name         string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			token:        tokens[entity.RolesById[1]],
			id:           idUser,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: ulid tidak valid",
			token:        tokens[entity.RolesById[1]],
			id:           "01HQVTTJ1S2606JGTYYZ5NDKNR123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[2],
			id:           idUser,
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			id:           idUser,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			id:           idUser,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idUser,
			token:        tokens[entity.RolesById[5]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/user/"+tt.id, nil, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)

			var res map[string]any
			if tt.name == "sukses" {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.Greater(t, len(res), 0)
				assert.NotEmpty(t, res)
				assert.NotEmpty(t, res["id"])
				assert.NotEmpty(t, res["created_at"])
				assert.NotEmpty(t, res["nama"])
				assert.NotEmpty(t, res["role"])
				assert.NotEmpty(t, res["username"])
				assert.NotEmpty(t, res["no_telp"])
				assert.NotEmpty(t, res["alamat"])
				assert.Equal(t, tt.expectedBody.Status, body.Status)
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

func UserDelete(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		id           string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			token:        tokens[entity.RolesById[1]],
			id:           idUser,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: tidak ditemukan",
			token:        tokens[entity.RolesById[1]],
			id:           "01HM4B8QBH7MWAVAYP10WN6PKA",
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name:         "err: ulid tidak valid",
			token:        tokens[entity.RolesById[1]],
			id:           idUser + "123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[2],
			id:           idUser,
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			id:           idUser,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			id:           idUser,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idUser,
			token:        tokens[entity.RolesById[5]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/user/"+tt.id, nil, &tt.token)
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

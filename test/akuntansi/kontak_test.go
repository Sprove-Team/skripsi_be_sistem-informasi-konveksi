package test_akuntansi

import (
	"net/url"
	"strings"
	"testing"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_akuntansi_kontak "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/kontak"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func AkuntansiCreateKontak(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		payload      req_akuntansi_kontak.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_kontak.Create{
				Nama:       "megabaran",
				NoTelp:     "+628964123452",
				Alamat:     "jln. pahlawan no 3D",
				Keterangan: "kontak langganan",
				Email:      "megabaran@yahoo.com",
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name:  "err: validasi format email & no telp",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_kontak.Create{
				Nama:       "megabaran2",
				NoTelp:     "0828964123452",
				Alamat:     "jln. pahlawan no 3D2",
				Keterangan: "kontak langganan2",
				Email:      "megabaran2",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"no telp harus berformat e164", "email harus berupa alamat email yang valid"},
			},
		},
		{
			name:  "err: wajib diisi",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_kontak.Create{
				Nama:   "",
				NoTelp: "",
				Alamat: "",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"nama wajib diisi", "no telp wajib diisi", "alamat wajib diisi"},
			},
		},
		{
			name:         "authorization " + entity.RolesById[2] + " passed",
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			payload:      req_akuntansi_kontak.Create{},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			payload:      req_akuntansi_kontak.Create{},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			payload:      req_akuntansi_kontak.Create{},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/akuntansi/kontak", tt.payload, &tt.token)
			assert.NoError(t, err)
			if strings.Contains(tt.name, "passed") {
				assert.NotEqual(t, tt.expectedCode, code)
				assert.NotEqual(t, tt.expectedBody.Code, body.Code)
				assert.NotEqual(t, tt.expectedBody.Status, body.Status)
				return
			}
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

var idKontak string
var idKontak2 string

func AkuntansiUpdateKontak(t *testing.T) {
	kontak := new(entity.Kontak)
	err := dbt.Model(kontak).Order("id DESC").First(kontak).Error
	if err != nil {
		panic(helper.LogsError(err))
	}
	kontak2 := &entity.Kontak{
		BaseSoftDelete: entity.BaseSoftDelete{
			ID: test.UlidPkg.MakeUlid().String(),
		},
		Nama:       "test",
		NoTelp:     "+62828964123459",
		Alamat:     "jln. pahlawan no 3D2",
		Keterangan: "kontak langganan2",
		Email:      "megabaran2@gmail.com",
	}
	if err := dbt.Create(kontak2).Error; err != nil {
		panic(helper.LogsError(err))
	}
	idKontak = kontak.ID
	idKontak2 = kontak2.ID

	tests := []struct {
		name         string
		token        string
		payload      req_akuntansi_kontak.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_kontak.Update{
				ID:         idKontak,
				Nama:       "test update",
				NoTelp:     "+62828912323459",
				Alamat:     "jln. heroik no 3D2",
				Keterangan: "kontak langganan tetap",
				Email:      "megabaran_test@gmail.com",
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:  "err: tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_kontak.Update{
				ID:         "01HM4B8QBH7MWAVAYP10WN6PKZ",
				Nama:       "test update not found",
				NoTelp:     "+62828912323451",
				Alamat:     "jln. heroik no 3D1",
				Keterangan: "kontak langganan1",
				Email:      "megabaran_test1@gmail.com",
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name:  "err: ulid tidak valid",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_kontak.Update{
				ID:         idKontak + "123",
				Nama:       "test update ulid tidak valid",
				NoTelp:     "+62828912323458",
				Alamat:     "jln. heroik no 3D8",
				Keterangan: "kontak langganan8",
				Email:      "megabaran_test8@gmail.com",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
		{
			name:  "authorization " + entity.RolesById[2] + " passed",
			token: tokens[entity.RolesById[2]],
			payload: req_akuntansi_kontak.Update{
				ID: idAkun,
			},
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name: "err: authorization " + entity.RolesById[3],
			payload: req_akuntansi_kontak.Update{
				ID: idAkun,
			},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name: "err: authorization " + entity.RolesById[4],
			payload: req_akuntansi_kontak.Update{
				ID: idAkun,
			},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name: "err: authorization " + entity.RolesById[5],
			payload: req_akuntansi_kontak.Update{
				ID: idAkun,
			},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/akuntansi/kontak/"+tt.payload.ID, tt.payload, &tt.token)
			assert.NoError(t, err)
			if strings.Contains(tt.name, "passed") {
				assert.NotEqual(t, tt.expectedCode, code)
				assert.NotEqual(t, tt.expectedBody.Code, body.Code)
				assert.NotEqual(t, tt.expectedBody.Status, body.Status)
				return
			}
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

func AkuntansiGetAllKontak(t *testing.T) {
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
		{ // same data with idKontak2
			name:         "sukses dengan filter",
			token:        tokens[entity.RolesById[1]],
			queryBody:    "?nama=test&no_telp=%2B628289123&email=megabaran",
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
			name:         "sukses dengan next",
			token:        tokens[entity.RolesById[1]],
			queryBody:    "?next=" + idKontak,
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
			queryBody:    "?next=01HQVTTJ1S2606JGTYYZ5NDKNR123",
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"next tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "authorization " + entity.RolesById[2] + " passed",
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[3] + " passed",
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[4] + " passed",
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[5] + " passed",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/kontak"+tt.queryBody, nil, &tt.token)
			assert.NoError(t, err)
			if strings.Contains(tt.name, "passed") {
				assert.NotEqual(t, tt.expectedCode, code)
				assert.NotEqual(t, tt.expectedBody.Code, body.Code)
				assert.NotEqual(t, tt.expectedBody.Status, body.Status)
				return
			}
			assert.Equal(t, tt.expectedCode, code)

			var res []map[string]interface{}
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)

				assert.Greater(t, len(res), 0)
				if len(res) <= 0 {
					panic("res less than zero")
				}
				assert.NotEmpty(t, res[0])
				assert.NotEmpty(t, res[0]["id"])
				assert.NotEmpty(t, res[0]["nama"])
				assert.NotEmpty(t, res[0]["alamat"])
				assert.NotEmpty(t, res[0]["email"])
				assert.NotEmpty(t, res[0]["no_telp"])
				assert.NotEmpty(t, res[0]["keterangan"])
				assert.NotEmpty(t, res[0]["created_at"])
				assert.Equal(t, tt.expectedBody.Status, body.Status)
				switch tt.name {
				case "sukses dengan filter":
					v, err := url.ParseQuery(tt.queryBody[1:])
					assert.NoError(t, err)
					assert.NotEmpty(t, res)
					assert.Contains(t, res[0]["nama"], v.Get("nama"))
					assert.Contains(t, res[0]["email"], v.Get("email"))
					assert.Equal(t, res[0]["no_telp"].(string)[0:len(v.Get("no_telp"))], v.Get("no_telp"))

				case "sukses limit 1":
					assert.Len(t, res, 1)
				case "sukses dengan next":
					assert.NotEmpty(t, res[0])
					assert.NotEqual(t, idKontak, res[0]["id"])
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

func AkuntansiGetKontak(t *testing.T) {
	tests := []struct {
		id           string
		name         string
		token        string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			token:        tokens[entity.RolesById[1]],
			id:           idKontak,
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
			id:           "01HQVTTJ1S2606JGTYYZ5NDKNR123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "authorization " + entity.RolesById[2] + " passed",
			id:           idAkun,
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			id:           idKontak,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			id:           idKontak,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idKontak,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/kontak/"+tt.id, nil, &tt.token)
			assert.NoError(t, err)
			if strings.Contains(tt.name, "passed") {
				assert.NotEqual(t, tt.expectedCode, code)
				assert.NotEqual(t, tt.expectedBody.Code, body.Code)
				assert.NotEqual(t, tt.expectedBody.Status, body.Status)
				return
			}
			assert.Equal(t, tt.expectedCode, code)

			var res map[string]any
			if tt.name == "sukses" {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				assert.NotEmpty(t, res["id"])
				assert.NotEmpty(t, res["nama"])
				assert.NotEmpty(t, res["alamat"])
				assert.NotEmpty(t, res["email"])
				assert.NotEmpty(t, res["no_telp"])
				assert.NotEmpty(t, res["keterangan"])
				assert.NotEmpty(t, res["created_at"])
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

func AkuntansiDeleteKontak(t *testing.T) {
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
			id:           idKontak,
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
			id:           idKontak + "123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "authorization " + entity.RolesById[2] + " passed",
			token:        tokens[entity.RolesById[2]],
			id:           idAkun,
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			id:           idKontak,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			id:           idKontak,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idKontak,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/akuntansi/kontak/"+tt.id, nil, &tt.token)
			assert.NoError(t, err)
			if strings.Contains(tt.name, "passed") {
				assert.NotEqual(t, tt.expectedCode, code)
				assert.NotEqual(t, tt.expectedBody.Code, body.Code)
				assert.NotEqual(t, tt.expectedBody.Status, body.Status)
				return
			}
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

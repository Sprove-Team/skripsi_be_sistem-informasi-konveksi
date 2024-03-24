package test_akuntansi

import (
	"net/url"
	"strings"
	"testing"

	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_akuntansi_akun "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/akun"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func AkuntansiCreateAkun(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		payload      req_akuntansi_akun.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_akun.Create{
				Nama:           "akun test",
				Kode:           "4",
				KelompokAkunID: idKelompokAkun,
				Deskripsi:      "des akun test",
				SaldoNormal:    "DEBIT",
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
			payload: req_akuntansi_akun.Create{
				Nama:           "akun test",
				Kode:           "4",
				KelompokAkunID: idKelompokAkun,
				Deskripsi:      "des akun test",
				SaldoNormal:    "DEBIT",
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name:  "err: saldo normal harus berupa salah satu dari [DEBIT,KREDIT]",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_akun.Create{
				Nama:           "akun test2",
				Kode:           "124",
				KelompokAkunID: idKelompokAkun,
				Deskripsi:      "des akun test2",
				SaldoNormal:    "ABCD",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"saldo normal harus berupa salah satu dari [DEBIT,KREDIT]"},
			},
		},
		{
			name:  "err: wajib diisi",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_akun.Create{
				Nama:           "",
				Kode:           "",
				KelompokAkunID: "",
				SaldoNormal:    "",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"nama wajib diisi", "kode wajib diisi", "kelompok akun id wajib diisi", "saldo normal wajib diisi"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			payload:      req_akuntansi_akun.Create{},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			payload:      req_akuntansi_akun.Create{},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			payload:      req_akuntansi_akun.Create{},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/akuntansi/akun", tt.payload, &tt.token)
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

var idAkun string
var idAkun2 string

func AkuntansiUpdateAkun(t *testing.T) {
	akun := new(entity.Akun)
	err := dbt.Model(akun).Order("id DESC").Preload("KelompokAkun").First(akun).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	akun2 := &entity.Akun{
		Base: entity.Base{
			ID: test.UlidPkg.MakeUlid().String(),
		},
		Kode:           akun.KelompokAkun.Kode + "1",
		Nama:           "akun conflict",
		KelompokAkunID: akun.KelompokAkunID,
		SaldoNormal:    "DEBIT",
		Deskripsi:      "des akun conflict",
	}
	if err := dbt.Create(akun2).Error; err != nil {
		helper.LogsError(err)
		return
	}
	idAkun = akun.ID
	idAkun2 = akun2.ID

	tests := []struct {
		name         string
		token        string
		payload      req_akuntansi_akun.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_akun.Update{
				ID:             idAkun,
				Nama:           "akun update test",
				Kode:           "99",
				KelompokAkunID: idKelompokAkun2,
				Deskripsi:      "update akun",
				SaldoNormal:    "KREDIT",
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:  "err: can't update default data",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_akun.Update{
				ID:             static_data.DataAkun[0][0].ID,
				Nama:           "can't update akun test",
				Kode:           "19",
				KelompokAkunID: static_data.DataKelompokAkun[12].ID,
				Deskripsi:      "des can't update akun test",
				SaldoNormal:    "KREDIT",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.CantModifiedDefaultData},
			},
		},
		{
			name:  "err: conflict",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_akun.Update{
				ID:             idAkun,
				Nama:           akun2.Nama,
				Kode:           "1",
				KelompokAkunID: akun2.KelompokAkunID,
				Deskripsi:      "des update conflict",
				SaldoNormal:    "KREDIT",
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name:  "err: tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_akun.Update{
				ID:             "01HM4B8QBH7MWAVAYP10WN6PKA",
				Nama:           "update not found",
				Kode:           "13",
				KelompokAkunID: static_data.DataKelompokAkun[12].ID,
				Deskripsi:      "des update not found",
				SaldoNormal:    "KREDIT",
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name:  "err: saldo normal harus berupa salah satu dari [DEBIT,KREDIT]",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_akun.Update{
				ID:             idAkun,
				Nama:           "akun test2",
				Kode:           "14",
				KelompokAkunID: idKelompokAkun,
				Deskripsi:      "des akun test2",
				SaldoNormal:    "ABCD",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"saldo normal harus berupa salah satu dari [DEBIT,KREDIT]"},
			},
		},
		{
			name:  "err: ulid tidak valid",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_akun.Update{
				ID:             idAkun + "123",
				Nama:           "kelompok ulid tidak valid",
				Kode:           "15",
				KelompokAkunID: static_data.DataKelompokAkun[12].ID,
				Deskripsi:      "des update ulid tidak valid",
				SaldoNormal:    "KREDIT",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			payload:      req_akuntansi_akun.Update{},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			payload:      req_akuntansi_akun.Update{},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			payload:      req_akuntansi_akun.Update{},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/akuntansi/akun/"+tt.payload.ID, tt.payload, &tt.token)
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

func AkuntansiGetAllAkun(t *testing.T) {
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
		{ // same data with idAkun2
			name:         "sukses with filter",
			token:        tokens[entity.RolesById[1]],
			queryBody:    "?nama=akun+update+test&kode=1123499",
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
			name:         "sukses with next",
			token:        tokens[entity.RolesById[1]],
			queryBody:    "?next=" + idAkun,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/akun"+tt.queryBody, nil, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)

			var res []map[string]interface{}
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)

				assert.Greater(t, len(res), 0)
				assert.NotEmpty(t, res[0])
				assert.NotEmpty(t, res[0]["id"])
				assert.NotEmpty(t, res[0]["kode"])
				assert.NotEmpty(t, res[0]["nama"])
				assert.NotEmpty(t, res[0]["created_at"])
				assert.NotEmpty(t, res[0]["saldo_normal"])
				assert.NotEmpty(t, res[0]["deskripsi"])
				assert.NotEmpty(t, res[0]["kelompok_akun"])
				assert.NotEmpty(t, res[0]["kelompok_akun"].(map[string]any)["id"])
				assert.NotEmpty(t, res[0]["kelompok_akun"].(map[string]any)["kode"])
				assert.NotEmpty(t, res[0]["kelompok_akun"].(map[string]any)["nama"])
				assert.Equal(t, tt.expectedBody.Status, body.Status)
				switch tt.name {
				case "sukses with filter":
					v, err := url.ParseQuery(tt.queryBody[1:])
					assert.NoError(t, err)
					assert.NotEmpty(t, res)
					assert.Equal(t, res[0]["nama"], v.Get("nama"))
					assert.Equal(t, res[0]["kode"], v.Get("kode"))
				case "sukses limit 1":
					assert.Len(t, res, 1)
				case "sukses with next":
					assert.NotEmpty(t, res[0])
					assert.NotEqual(t, idAkun, res[0]["id"])
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

func AkuntansiGetAkun(t *testing.T) {
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
			id:           idAkun,
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
			name:         "err: authorization " + entity.RolesById[3],
			id:           idAkun,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			id:           idAkun,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idAkun,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/akun/"+tt.id, nil, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)

			var res map[string]any
			if tt.name == "sukses" {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				assert.NotEmpty(t, res["id"])
				assert.NotEmpty(t, res["created_at"])
				assert.NotEmpty(t, res["nama"])
				assert.NotEmpty(t, res["kode"])
				assert.NotEmpty(t, res["saldo_normal"])
				assert.NotEmpty(t, res["deskripsi"])
				assert.NotEmpty(t, res["kelompok_akun"])
				assert.NotEmpty(t, res["kelompok_akun"].(map[string]any)["id"])
				assert.NotEmpty(t, res["kelompok_akun"].(map[string]any)["kode"])
				assert.NotEmpty(t, res["kelompok_akun"].(map[string]any)["nama"])
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

func AkuntansiDeleteAkun(t *testing.T) {
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
			id:           idAkun,
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
			name:         "err: can't delete default data",
			token:        tokens[entity.RolesById[1]],
			id:           static_data.DataAkun[0][0].ID,
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.CantModifiedDefaultData},
			},
		},
		{
			name:         "err: ulid tidak valid",
			token:        tokens[entity.RolesById[1]],
			id:           idAkun + "123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			id:           idAkun,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			id:           idAkun,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idAkun,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/akuntansi/akun/"+tt.id, nil, &tt.token)
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

package test_akuntansi

import (
	"net/url"
	"strings"
	"testing"

	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_akuntansi_kelompok_akun "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/kelompok_akun"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func AkuntansiCreateKelompokAkun(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		payload      req_akuntansi_kelompok_akun.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_kelompok_akun.Create{
				Nama:         "kelompok test",
				Kode:         "11111111",
				KategoriAkun: entity.KategoriAkunByKode["1"],
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
			payload: req_akuntansi_kelompok_akun.Create{
				Nama:         static_data.DataKelompokAkun[0].Nama,
				Kode:         "1",
				KategoriAkun: static_data.DataKelompokAkun[0].KategoriAkun,
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name:  "err: kategori akun harus berupa salah satu dari [ASET,KEWAJIBAN,MODAL,PENDAPATAN,BEBAN]",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_kelompok_akun.Create{
				Nama:         "kelompok err kategori",
				Kode:         "212312",
				KategoriAkun: "asdf",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kategori akun harus berupa salah satu dari [ASET,KEWAJIBAN,MODAL,PENDAPATAN,BEBAN]"},
			},
		},
		{
			name:  "err: wajib diisi",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_kelompok_akun.Create{
				Nama:         "",
				Kode:         "",
				KategoriAkun: "",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"nama wajib diisi", "kode wajib diisi", "kategori akun wajib diisi"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			payload:      req_akuntansi_kelompok_akun.Create{},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			payload:      req_akuntansi_kelompok_akun.Create{},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			payload:      req_akuntansi_kelompok_akun.Create{},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/akuntansi/kelompok_akun", tt.payload, &tt.token)
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

var idKelompokAkun string
var idKelompokAkun2 string

func AkuntansiUpdateKelompokAkun(t *testing.T) {
	kelompokAkun := new(entity.KelompokAkun)
	err := dbt.Model(kelompokAkun).Order("id DESC").First(kelompokAkun).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	kelompokAkun2 := &entity.KelompokAkun{
		Base: entity.Base{
			ID: test.UlidPkg.MakeUlid().String(),
		},
		Kode:         entity.KategoriAkun[kelompokAkun.KategoriAkun] + "1234",
		Nama:         "kelompok akun conflict",
		KategoriAkun: kelompokAkun.KategoriAkun,
	}
	if err := dbt.Create(kelompokAkun2).Error; err != nil {
		helper.LogsError(err)
		return
	}
	idKelompokAkun = kelompokAkun.ID
	idKelompokAkun2 = kelompokAkun2.ID
	tests := []struct {
		name         string
		token        string
		payload      req_akuntansi_kelompok_akun.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_kelompok_akun.Update{
				ID:           idKelompokAkun,
				Nama:         "kelompok akun update",
				Kode:         "1112",
				KategoriAkun: entity.KategoriAkunByKode["2"],
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
			payload: req_akuntansi_kelompok_akun.Update{
				ID:           static_data.DataKelompokAkun[0].ID,
				Nama:         "can't update",
				Kode:         "11111",
				KategoriAkun: entity.KategoriAkunByKode["4"],
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
			payload: req_akuntansi_kelompok_akun.Update{
				ID:           idKelompokAkun,
				Nama:         "update conflict",
				Kode:         "1234",
				KategoriAkun: kelompokAkun2.KategoriAkun,
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
			payload: req_akuntansi_kelompok_akun.Update{
				ID:           "01HM4B8QBH7MWAVAYP10WN6PKA",
				Nama:         "update not found",
				Kode:         "111125",
				KategoriAkun: entity.KategoriAkunByKode["1"],
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name:  "err: kategori akun harus berupa salah satu dari [ASET,KEWAJIBAN,MODAL,PENDAPATAN,BEBAN]",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_kelompok_akun.Update{
				ID:           idKelompokAkun,
				Nama:         "kelompok err kategori",
				Kode:         "21231",
				KategoriAkun: "asdf",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kategori akun harus berupa salah satu dari [ASET,KEWAJIBAN,MODAL,PENDAPATAN,BEBAN]"},
			},
		},
		{
			name:  "err: ulid tidak valid",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_kelompok_akun.Update{
				ID:           idKelompokAkun + "123",
				Nama:         "kelompok ulid tidak valid",
				Kode:         "111124",
				KategoriAkun: entity.KategoriAkunByKode["4"],
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
			payload:      req_akuntansi_kelompok_akun.Update{},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			payload:      req_akuntansi_kelompok_akun.Update{},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			payload:      req_akuntansi_kelompok_akun.Update{},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/akuntansi/kelompok_akun/"+tt.payload.ID, tt.payload, &tt.token)
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

func AkuntansiGetAllKelompokAkun(t *testing.T) {
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
		{ // same data with idKelompokAkun2
			name:         "sukses with filter",
			token:        tokens[entity.RolesById[1]],
			queryBody:    "?nama=kelompok+akun+conflict&kode=11234&kategori_akun=" + entity.KategoriAkunByKode["1"],
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
			queryBody:    "?next=" + idKelompokAkun,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: kategori akun harus berupa salah satu dari [ASET,KEWAJIBAN,MODAL,PENDAPATAN,BEBAN]",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 400,
			queryBody:    "?kategori_akun=ABCD",
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kategori akun harus berupa salah satu dari [ASET,KEWAJIBAN,MODAL,PENDAPATAN,BEBAN]"},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/kelompok_akun"+tt.queryBody, nil, &tt.token)
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
				assert.NotEmpty(t, res[0]["kategori_akun"])
				assert.NotEmpty(t, res[0]["created_at"])
				assert.Equal(t, tt.expectedBody.Status, body.Status)
				switch tt.name {
				case "sukses with filter":
					v, err := url.ParseQuery(tt.queryBody[1:])
					assert.NoError(t, err)
					assert.NotEmpty(t, res)
					assert.Equal(t, res[0]["nama"], v.Get("nama"))
					assert.Equal(t, res[0]["kode"], v.Get("kode"))
					assert.Equal(t, res[0]["kategori_akun"], v.Get("kategori_akun"))
				case "sukses limit 1":
					assert.Len(t, res, 1)
				case "sukses with next":
					assert.NotEmpty(t, res[0])
					assert.NotEqual(t, idKelompokAkun, res[0]["id"])
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

func AkuntansiGetKelompokAkun(t *testing.T) {
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
			id:           idKelompokAkun,
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
			id:           idKelompokAkun,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			id:           idKelompokAkun,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idKelompokAkun,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/kelompok_akun/"+tt.id, nil, &tt.token)
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
				assert.NotEmpty(t, res["kategori_akun"])
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

func AkuntansiDeleteKelompokAkun(t *testing.T) {
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
			id:           idKelompokAkun,
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
			id:           static_data.DataKelompokAkun[0].ID,
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
			id:           idKelompokAkun + "123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			id:           idKelompokAkun,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			id:           idKelompokAkun,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idKelompokAkun,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/akuntansi/kelompok_akun/"+tt.id, nil, &tt.token)
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

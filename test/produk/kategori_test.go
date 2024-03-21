package test_produk

import (
	"strings"
	"testing"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_produk "github.com/be-sistem-informasi-konveksi/common/request/produk"
	produkKategori "github.com/be-sistem-informasi-konveksi/common/request/produk/kategori"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func ProdukCreateKategori(t *testing.T) {

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
			name: "err: wajib diisi",
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

func ProdukUpdateKategori(t *testing.T) {
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

func ProdukGetAllKategori(t *testing.T) {
	kategori := &entity.KategoriProduk{
		Base: entity.Base{
			ID: test.UlidPkg.MakeUlid().String(),
		},
		Nama: "test next",
	}
	if err := dbt.Create(kategori).Error; err != nil {
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
			name:         "sukses limit 1",
			queryBody:    "?limit=1",
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses with next",
			queryBody:    "?next=" + idProduk,
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
				case "sukses with next":
					assert.NotEmpty(t, res[0])
					assert.NotEqual(t, idProduk, res[0]["id"])
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

func ProdukGetKategori(t *testing.T) {
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

func ProdukDeleteKategori(t *testing.T) {
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

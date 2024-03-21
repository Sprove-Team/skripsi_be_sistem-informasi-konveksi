package test_produk

import (
	"strings"
	"testing"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_produk "github.com/be-sistem-informasi-konveksi/common/request/produk"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func ProdukCreate(t *testing.T) {
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

var idProduk string
var idKategori2 string

func ProdukUpdate(t *testing.T) {
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

var idProdukHasHargaDetail string

func ProdukGetAll(t *testing.T) {
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
			name:         "sukses with next",
			queryBody:    "?next=" + idProduk,
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

func ProdukGet(t *testing.T) {
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

func ProdukDelete(t *testing.T) {
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

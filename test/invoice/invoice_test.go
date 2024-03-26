package test_akuntansi

import (
	"testing"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_invoice "github.com/be-sistem-informasi-konveksi/common/request/invoice"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func InvoiceCreate(t *testing.T) {
	detailInvoice1 := []req_invoice.ReqDetailInvoice{
		{
			ProdukID:     produk[0].ID,
			BordirID:     bordir[0].ID,
			SablonID:     sablon[0].ID,
			GambarDesign: "img-design.webp",
			Qty:          10,
			Total:        (produk[0].HargaDetails[1].Harga + bordir[0].Harga + sablon[0].Harga) * 10,
		},
		{
			ProdukID:     produk[1].ID,
			BordirID:     bordir[1].ID,
			SablonID:     sablon[1].ID,
			GambarDesign: "img-design.webp",
			Qty:          1,
			Total:        (produk[1].HargaDetails[0].Harga + bordir[1].Harga + sablon[1].Harga) * 1,
		},
	}
	detailInvoice2 := []req_invoice.ReqDetailInvoice{
		{
			ProdukID:     produk[1].ID,
			BordirID:     bordir[1].ID,
			SablonID:     sablon[1].ID,
			GambarDesign: "img-design-2.webp",
			Qty:          1,
			Total:        (produk[1].HargaDetails[0].Harga + bordir[1].Harga + sablon[1].Harga) * 1,
		},
		{
			ProdukID:     produk[0].ID,
			BordirID:     bordir[0].ID,
			SablonID:     sablon[0].ID,
			GambarDesign: "img-design-2.webp",
			Qty:          10,
			Total:        (produk[0].HargaDetails[1].Harga + bordir[0].Harga + sablon[0].Harga) * 10,
		},
	}
	tests := []struct {
		name         string
		token        string
		payload      req_invoice.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				KontakID: kontak[0].ID,
				Bayar: req_invoice.ReqBayar{
					BuktiPembayaran: []string{"img-bukti.webp"},
					Keterangan:      "DP 1",
					AkunID:          "01HP7DVBGTC06PXWT6FD66VERN", // kas
					Total:           (detailInvoice1[0].Total + detailInvoice1[1].Total) / 2,
				},
				TanggalDeadline: "2024-10-15T12:00:00Z",
				TanggalKirim:    "2024-10-15T12:00:00Z",
				Keterangan:      "ket invoice 1",
				DetailInvoice:   detailInvoice1,
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name:  "sukses with new kontak",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				NewKontak: req_invoice.ReqNewKontak{
					Nama:   "joni",
					NoTelp: "+6281234567890",
					Alamat: "123 Main Street",
					Email:  "john.doe@example.com",
				},
				Bayar: req_invoice.ReqBayar{
					BuktiPembayaran: []string{"img-bukti-2.webp"},
					Keterangan:      "DP 2",
					AkunID:          "01HP7DVBGTC06PXWT6FD66VERN", // kas
					Total:           (detailInvoice2[0].Total + detailInvoice2[0].Total) / 2,
				},
				TanggalDeadline: "2024-10-15T12:00:00Z",
				TanggalKirim:    "2024-10-15T12:00:00Z",
				Keterangan:      "ket invoice 2",
				DetailInvoice:   detailInvoice2,
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name:  "err: wajib diisi",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				KontakID: "",
				Bayar: req_invoice.ReqBayar{
					Keterangan: "",
					AkunID:     "",
					Total:      0,
				},
				Keterangan:      "",
				TanggalDeadline: "",
				TanggalKirim:    "",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kontak id wajib diisi jika new kontak tidak diisi", "bukti pembayaran wajib diisi", "keterangan wajib diisi", "akun id wajib diisi", "total wajib diisi", "tanggal deadline wajib diisi", "tanggal kirim wajib diisi", "detail invoice wajib diisi", "bukti pembayaran wajib diisi"},
			},
		},
		{
			name:  "err: item lebih dari 0",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				Bayar: req_invoice.ReqBayar{
					BuktiPembayaran: entity.BuktiPembayaran{},
				},
				DetailInvoice: []req_invoice.ReqDetailInvoice{},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"detail invoice harus berisi lebih dari 0 item", "bukti pembayaran harus berisi lebih dari 0 item"},
			},
		},
		{
			name:  "err: jika data A diisi maka B harus kosong",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				KontakID: kontak[0].ID,
				NewKontak: req_invoice.ReqNewKontak{
					Nama: "test",
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"jika kontak id diisi maka new kontak harus kosong"},
			},
		},
		{
			name:  "err: format tanggal pengiriman and deadline",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				TanggalDeadline: "2006-01-20",
				TanggalKirim:    "2006-01-20",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"tanggal deadline harus berformat RFC3999", "tanggal kirim harus berformat RFC3999"},
			},
		},
		// {
		// 	name:         "err: conflict",
		// 	token:        tokens[entity.RolesById[1]],
		// 	payload:      req_invoice.Create{},
		// 	expectedCode: 409,
		// 	expectedBody: test.Response{
		// 		Status: fiber.ErrConflict.Message,
		// 		Code:   409,
		// 	},
		// },
		// {
		// 	name:  "err: wajib diisi",
		// 	token: tokens[entity.RolesById[1]],
		// 	payload: req_invoice.Create{
		// 		Nama:  "",
		// 		Harga: 0,
		// 	},
		// 	expectedCode: 400,
		// 	expectedBody: test.Response{
		// 		Status:         fiber.ErrBadRequest.Message,
		// 		Code:           400,
		// 		ErrorsMessages: []string{"nama wajib diisi", "harga wajib diisi"},
		// 	},
		// },
		{
			name:         "err: authorization " + entity.RolesById[2],
			payload:      req_invoice.Create{},
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			payload:      req_invoice.Create{},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			payload:      req_invoice.Create{},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/invoice", tt.payload, &tt.token)
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

// var idBordir string

// func BordirUpdate(t *testing.T) {
// 	bordir := new(entity.Bordir)
// 	err := dbt.Select("id").First(bordir).Error
// 	if err != nil {
// 		helper.LogsError(err)
// 		return
// 	}
// 	idBordir = bordir.ID
// 	tests := []struct {
// 		name         string
// 		token        string
// 		payload      req_bordir.Update
// 		expectedBody test.Response
// 		expectedCode int
// 	}{
// 		{
// 			name:  "sukses",
// 			token: tokens[entity.RolesById[1]],
// 			payload: req_bordir.Update{
// 				ID:    bordir.ID,
// 				Nama:  "Bordir 1 baris",
// 				Harga: 20000,
// 			},
// 			expectedCode: 200,
// 			expectedBody: test.Response{
// 				Status: message.OK,
// 				Code:   200,
// 			},
// 		},
// 		{
// 			name:  "err: tidak ditemukan",
// 			token: tokens[entity.RolesById[1]],
// 			payload: req_bordir.Update{
// 				ID:    "01HM4B8QBH7MWAVAYP10WN6PKA",
// 				Nama:  "Bordir 3 baris",
// 				Harga: 20001,
// 			},
// 			expectedCode: 404,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrNotFound.Message,
// 				Code:   404,
// 			},
// 		},
// 		{
// 			name:  "err: ulid tidak valid",
// 			token: tokens[entity.RolesById[1]],
// 			payload: req_bordir.Update{
// 				ID:    bordir.ID + "123",
// 				Nama:  "Bordir 4 baris",
// 				Harga: 20004,
// 			},
// 			expectedCode: 400,
// 			expectedBody: test.Response{
// 				Status:         fiber.ErrBadRequest.Message,
// 				Code:           400,
// 				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
// 			},
// 		},
// 		{
// 			name:         "err: authorization " + entity.RolesById[2],
// 			payload:      req_bordir.Update{},
// 			token:        tokens[entity.RolesById[2]],
// 			expectedCode: 401,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrUnauthorized.Message,
// 				Code:   401,
// 			},
// 		},
// 		{
// 			name:         "err: authorization " + entity.RolesById[3],
// 			payload:      req_bordir.Update{},
// 			token:        tokens[entity.RolesById[3]],
// 			expectedCode: 401,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrUnauthorized.Message,
// 				Code:   401,
// 			},
// 		},
// 		{
// 			name:         "err: authorization " + entity.RolesById[4],
// 			payload:      req_bordir.Update{},
// 			token:        tokens[entity.RolesById[4]],
// 			expectedCode: 401,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrUnauthorized.Message,
// 				Code:   401,
// 			},
// 		},
// 		{
// 			name:         "err: authorization " + entity.RolesById[5],
// 			payload:      req_bordir.Update{},
// 			token:        tokens[entity.RolesById[5]],
// 			expectedCode: 401,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrUnauthorized.Message,
// 				Code:   401,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/bordir/"+tt.payload.ID, tt.payload, &tt.token)
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expectedCode, code)
// 			if len(tt.expectedBody.ErrorsMessages) > 0 {
// 				for _, v := range tt.expectedBody.ErrorsMessages {
// 					assert.Contains(t, body.ErrorsMessages, v)
// 				}
// 				assert.Equal(t, tt.expectedBody.Status, body.Status)
// 			} else {
// 				assert.Equal(t, tt.expectedBody, body)
// 			}
// 		})
// 	}
// }

// func BordirGetAll(t *testing.T) {
// 	bordir := &entity.Bordir{
// 		Base: entity.Base{
// 			ID: test.UlidPkg.MakeUlid().String(),
// 		},
// 		Nama:  "test next",
// 		Harga: 20000,
// 	}
// 	if err := dbt.Create(bordir).Error; err != nil {
// 		helper.LogsError(err)
// 		return
// 	}
// 	tests := []struct {
// 		name         string
// 		token        string
// 		queryBody    string
// 		expectedBody test.Response
// 		expectedCode int
// 	}{
// 		{
// 			name:         "sukses",
// 			token:        tokens[entity.RolesById[1]],
// 			expectedCode: 200,
// 			expectedBody: test.Response{
// 				Status: message.OK,
// 				Code:   200,
// 			},
// 		},
// 		{
// 			name:         "sukses limit 1",
// 			token:        tokens[entity.RolesById[1]],
// 			queryBody:    "?limit=1",
// 			expectedCode: 200,
// 			expectedBody: test.Response{
// 				Status: message.OK,
// 				Code:   200,
// 			},
// 		},
// 		{
// 			name:         "sukses with next",
// 			token:        tokens[entity.RolesById[1]],
// 			queryBody:    "?next=" + idBordir,
// 			expectedCode: 200,
// 			expectedBody: test.Response{
// 				Status: message.OK,
// 				Code:   200,
// 			},
// 		},
// 		{
// 			name:         "sukses with: filter nama",
// 			token:        tokens[entity.RolesById[1]],
// 			queryBody:    "?nama=Bordir+1+baris",
// 			expectedCode: 200,
// 			expectedBody: test.Response{
// 				Status: message.OK,
// 				Code:   200,
// 			},
// 		},
// 		{
// 			name:         "err: ulid tidak valid",
// 			token:        tokens[entity.RolesById[1]],
// 			expectedCode: 400,
// 			queryBody:    "?next=01HQVTTJ1S2606JGTYYZ5NDKNR123",
// 			expectedBody: test.Response{
// 				Status:         fiber.ErrBadRequest.Message,
// 				Code:           400,
// 				ErrorsMessages: []string{"next tidak berupa ulid yang valid"},
// 			},
// 		},
// 		{
// 			name:         "err: authorization " + entity.RolesById[2],
// 			token:        tokens[entity.RolesById[2]],
// 			expectedCode: 401,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrUnauthorized.Message,
// 				Code:   401,
// 			},
// 		},
// 		{
// 			name:         "err: authorization " + entity.RolesById[3],
// 			token:        tokens[entity.RolesById[3]],
// 			expectedCode: 401,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrUnauthorized.Message,
// 				Code:   401,
// 			},
// 		},
// 		{
// 			name:         "err: authorization " + entity.RolesById[4],
// 			token:        tokens[entity.RolesById[4]],
// 			expectedCode: 401,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrUnauthorized.Message,
// 				Code:   401,
// 			},
// 		},
// 		{
// 			name:         "err: authorization " + entity.RolesById[5],
// 			token:        tokens[entity.RolesById[5]],
// 			expectedCode: 401,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrUnauthorized.Message,
// 				Code:   401,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/bordir"+tt.queryBody, nil, &tt.token)
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expectedCode, code)

// 			var res []map[string]interface{}
// 			if strings.Contains(tt.name, "sukses") {
// 				err = mapstructure.Decode(body.Data, &res)
// 				assert.NoError(t, err)
// 				assert.Greater(t, len(res), 0)
// 				assert.NotEmpty(t, res[0])
// 				assert.NotEmpty(t, res[0]["id"])
// 				assert.NotEmpty(t, res[0]["created_at"])
// 				assert.NotEmpty(t, res[0]["nama"])
// 				assert.NotEmpty(t, res[0]["harga"])
// 				assert.Equal(t, tt.expectedBody.Status, body.Status)
// 				switch tt.name {
// 				case "sukses with: filter nama":
// 					v, err := url.ParseQuery(tt.queryBody[1:])
// 					assert.NoError(t, err)
// 					assert.Contains(t, res[0]["nama"], v.Get("nama"))
// 				case "sukses limit 1":
// 					assert.Len(t, res, 1)
// 				case "sukses with next":
// 					assert.NotEmpty(t, res[0])
// 					assert.NotEqual(t, idBordir, res[0]["id"])
// 				}
// 			} else {
// 				if len(tt.expectedBody.ErrorsMessages) > 0 {
// 					for _, v := range tt.expectedBody.ErrorsMessages {
// 						assert.Contains(t, body.ErrorsMessages, v)
// 					}
// 					assert.Equal(t, tt.expectedBody.Status, body.Status)
// 				} else {
// 					assert.Equal(t, tt.expectedBody, body)
// 				}
// 			}

// 		})
// 	}
// }

// func BordirDelete(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		token        string
// 		id           string
// 		expectedBody test.Response
// 		expectedCode int
// 	}{
// 		{
// 			name:         "sukses",
// 			token:        tokens[entity.RolesById[1]],
// 			id:           idBordir,
// 			expectedCode: 200,
// 			expectedBody: test.Response{
// 				Status: message.OK,
// 				Code:   200,
// 			},
// 		},
// 		{
// 			name:         "err: tidak ditemukan",
// 			token:        tokens[entity.RolesById[1]],
// 			id:           "01HM4B8QBH7MWAVAYP10WN6PKA",
// 			expectedCode: 404,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrNotFound.Message,
// 				Code:   404,
// 			},
// 		},
// 		{
// 			name:         "err: ulid tidak valid",
// 			token:        tokens[entity.RolesById[1]],
// 			id:           idBordir + "123",
// 			expectedCode: 400,
// 			expectedBody: test.Response{
// 				Status:         fiber.ErrBadRequest.Message,
// 				Code:           400,
// 				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
// 			},
// 		},
// 		{
// 			name:         "err: authorization " + entity.RolesById[2],
// 			id:           idBordir,
// 			token:        tokens[entity.RolesById[2]],
// 			expectedCode: 401,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrUnauthorized.Message,
// 				Code:   401,
// 			},
// 		},
// 		{
// 			name:         "err: authorization " + entity.RolesById[3],
// 			id:           idBordir,
// 			token:        tokens[entity.RolesById[3]],
// 			expectedCode: 401,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrUnauthorized.Message,
// 				Code:   401,
// 			},
// 		},
// 		{
// 			name:         "err: authorization " + entity.RolesById[4],
// 			id:           idBordir,
// 			token:        tokens[entity.RolesById[4]],
// 			expectedCode: 401,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrUnauthorized.Message,
// 				Code:   401,
// 			},
// 		},
// 		{
// 			name:         "err: authorization " + entity.RolesById[5],
// 			id:           idBordir,
// 			token:        tokens[entity.RolesById[5]],
// 			expectedCode: 401,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrUnauthorized.Message,
// 				Code:   401,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/bordir/"+tt.id, nil, &tt.token)
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expectedCode, code)
// 			if len(tt.expectedBody.ErrorsMessages) > 0 {
// 				for _, v := range tt.expectedBody.ErrorsMessages {
// 					assert.Contains(t, body.ErrorsMessages, v)
// 				}
// 				assert.Equal(t, tt.expectedBody.Status, body.Status)
// 			} else {
// 				assert.Equal(t, tt.expectedBody, body)
// 			}
// 		})
// 	}
// }

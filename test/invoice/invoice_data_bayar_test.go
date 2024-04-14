package test_akuntansi

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_invoice_data_bayar "github.com/be-sistem-informasi-konveksi/common/request/invoice/data_bayar"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func InvoiceCreateDataBayarByInvoiceId(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		payload      req_invoice_data_bayar.CreateByInvoiceId
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice_data_bayar.CreateByInvoiceId{
				InvoiceID:       idInvoice,
				BuktiPembayaran: []string{"img.webp"},
				Keterangan:      "data bayar 1",
				AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
				Total:           1000,
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
			payload: req_invoice_data_bayar.CreateByInvoiceId{
				InvoiceID:  idInvoice,
				Keterangan: "",
				AkunID:     "",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"bukti pembayaran wajib diisi", "keterangan wajib diisi", "akun id wajib diisi", "total wajib diisi"},
			},
		},
		{
			name:  "err: ulid tidak valid",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice_data_bayar.CreateByInvoiceId{
				InvoiceID:       idInvoice + "123",
				Keterangan:      "qwe",
				AkunID:          idInvoice + "123",
				BuktiPembayaran: []string{"adf.jpg"},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"invoice id tidak berupa ulid yang valid", "akun id tidak berupa ulid yang valid"},
			},
		},
		{
			name:  "err: item harus berisi/lebih besar dari 0",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice_data_bayar.CreateByInvoiceId{
				InvoiceID:       idInvoice,
				Keterangan:      "qwe",
				AkunID:          idInvoice,
				BuktiPembayaran: []string{},
				Total:           -1,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"bukti pembayaran harus berisi lebih dari 0 item", "total harus lebih besar dari 0"},
			},
		},
		{
			name:  "err: total bayar harus kurang atau sama dengan sisa tagihan",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice_data_bayar.CreateByInvoiceId{
				InvoiceID:       idInvoice,
				BuktiPembayaran: []string{"img.webp"},
				Keterangan:      "data bayar 1",
				AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
				Total:           10000000,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.BayarMustLessThanSisaTagihan},
			},
		},
		{
			name:  "err: tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice_data_bayar.CreateByInvoiceId{
				InvoiceID:       produk[0].ID,
				BuktiPembayaran: []string{"img.webp"},
				Keterangan:      "data bayar 1",
				AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
				Total:           1000,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"invoice tidak ditemukan"},
			},
		},
		{
			name: "err: authorization " + entity.RolesById[2],
			payload: req_invoice_data_bayar.CreateByInvoiceId{
				InvoiceID: idInvoice,
			},
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name: "authorization " + entity.RolesById[3] + " passed",
			payload: req_invoice_data_bayar.CreateByInvoiceId{
				InvoiceID: idInvoice,
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
			payload: req_invoice_data_bayar.CreateByInvoiceId{
				InvoiceID: idInvoice,
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
			payload: req_invoice_data_bayar.CreateByInvoiceId{
				InvoiceID: idInvoice,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/invoice/"+tt.payload.InvoiceID+"/data_bayar", tt.payload, &tt.token)
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

var idDataBayar string
var dataBayar entity.DataBayarInvoice
var idDataBayarTerkonfirmasi string

func InvoiceUpdateDataBayar(t *testing.T) {
	invoiceDataBayar := []entity.DataBayarInvoice{}
	err := dbt.Preload("Invoice", func(db *gorm.DB)*gorm.DB {
		return db.Select("id", "kontak_id")
	}).Find(&invoiceDataBayar).Error
	assert.NoError(t, err)
	if err != nil {
		panic(helper.LogsError(err))
	}
	if length := len(invoiceDataBayar); length <= 3 {
		assert.LessOrEqual(t, length, 3)
		panic(helper.LogsError(errors.New("length not lest or equal than 3")))
	}
	idDataBayarTerkonfirmasi = invoiceDataBayar[0].ID
	idDataBayar = invoiceDataBayar[2].ID
	dataBayar = invoiceDataBayar[2]
	tests := []struct {
		name         string
		idInvoice    string
		token        string
		payload      req_invoice_data_bayar.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:      "sukses",
			token:     tokens[entity.RolesById[1]],
			idInvoice: invoiceDataBayar[0].InvoiceID,
			payload: req_invoice_data_bayar.Update{
				ID:              invoiceDataBayar[0].ID,
				Status:          "TERKONFIRMASI",
				BuktiPembayaran: []string{"img.jpg"},
				Keterangan:      "data bayar TERKONFIRMASI " + test.UlidPkg.MakeUlid().String(),
				AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
				Total:           invoiceDataBayar[0].Total + 1,
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:      "sukses update status data bayar untuk role bendahara",
			token:     tokens[entity.RolesById[2]],
			idInvoice: invoiceDataBayar[1].InvoiceID,
			payload: req_invoice_data_bayar.Update{
				ID:              invoiceDataBayar[1].ID,
				Status:          "TERKONFIRMASI",
				BuktiPembayaran: []string{"img2.jpg"},
				Keterangan:      "data bayar terkonfirmasi2 " + test.UlidPkg.MakeUlid().String(),
				AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
				Total:           invoiceDataBayar[1].Total + 1,
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:      "sukses akan tetapi status tidak terupdate karena role tidak diperbolehkan",
			token:     tokens[entity.RolesById[3]],
			idInvoice: invoiceDataBayar[2].InvoiceID,
			payload: req_invoice_data_bayar.Update{
				ID:              invoiceDataBayar[2].ID,
				Status:          "TERKONFIRMASI",
				BuktiPembayaran: []string{"img3.jpg"},
				Keterangan:      "data bayar TERKONFIRMASI3 " + test.UlidPkg.MakeUlid().String(), // harus unik dengan sukses lainnya
				AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
				Total:           invoiceDataBayar[2].Total + 1,
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
			payload: req_invoice_data_bayar.Update{
				ID: "01HM4B8QBH7MWAVAYP10WN6PKA",
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name:  "err: tidak bisa mengubah status yang bernilai terkonfirmasi",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice_data_bayar.Update{
				ID:              invoiceDataBayar[1].ID,
				Status:          "TERKONFIRMASI",
				BuktiPembayaran: []string{"img3.jpg"},
				Keterangan:      "data bayar TERKONFIRMASI3",
				AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
				Total:           50003,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.CannotModifiedTerkonfirmasiDataBayar},
			},
		},
		{
			name:  "err: format entire field",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice_data_bayar.Update{
				ID:              invoiceDataBayar[2].ID + "123",
				AkunID:          "123",
				Status:          "123",
				BuktiPembayaran: entity.BuktiPembayaran{},
				Total:           -1,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status: fiber.ErrBadRequest.Message,
				Code:   400,
				ErrorsMessages: []string{
					"id tidak berupa ulid yang valid",
					"status harus berupa salah satu dari [TERKONFIRMASI,BELUM_TERKONFIRMASI]",
					"bukti pembayaran harus berisi lebih dari 0 item",
					"akun id tidak berupa ulid yang valid",
					"total harus lebih besar dari 0",
				},
			},
		},
		{
			name: "err: authorization " + entity.RolesById[4],
			payload: req_invoice_data_bayar.Update{
				ID: idDataBayar,
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
			payload: req_invoice_data_bayar.Update{
				ID: idDataBayar,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/invoice/data_bayar/"+tt.payload.ID, tt.payload, &tt.token)
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
				if strings.Contains(tt.name, "sukses") {
					dbyrInvoice := new(entity.DataBayarInvoice)
					err = dbt.First(dbyrInvoice, "keterangan = ?", tt.payload.Keterangan).Error
					assert.NoError(t, err)
					assert.NotNil(t, dbyrInvoice)
					if err != nil {
						panic(helper.LogsError(err))
					}
					if tt.name != "sukses akan tetapi status tidak terupdate karena role tidak diperbolehkan" {
						assert.Equal(t, tt.payload.Status, dbyrInvoice.Status)
						// cek tr
						tr := new(entity.Transaksi)
						err := dbt.Preload("AyatJurnals").First(tr, "keterangan = ?", tt.payload.Keterangan).Error
						assert.NoError(t, err)
						assert.NotNil(t, tr)
						if err != nil {
							panic(helper.LogsError(err))
						}

						assert.Equal(t, tt.payload.Keterangan, tr.Keterangan)
						assert.Equal(t, tt.payload.Total, tr.Total)
						assert.Equal(t, tt.payload.BuktiPembayaran, tr.BuktiPembayaran)
						var akunExist bool
						for _, ay := range tr.AyatJurnals {
							if ay.AkunID == tt.payload.AkunID {
								akunExist = true
								break
							}
						}
						assert.True(t, akunExist)
					} else {
						assert.NotEqual(t, tt.payload.Status, dbyrInvoice.Status)
					}

				}
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

func InvoiceGetAllDataBayar(t *testing.T) {
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
			queryBody:    "?next=" + idDataBayarTerkonfirmasi,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses dengan: semua filter",
			token:        tokens[entity.RolesById[1]],
			queryBody:    fmt.Sprintf("?status=%s&kontak_id=%s", dataBayar.Status,dataBayar.Invoice.KontakID),
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/invoice/data_bayar"+tt.queryBody, nil, &tt.token)
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
				for _, v := range res {
					assert.NotEmpty(t, v)
					assert.NotEmpty(t, v["id"])
					assert.NotEmpty(t, v["created_at"])
					inv, ok := v["invoice"].(map[string]any)
					assert.True(t, ok)
					assert.NotEmpty(t, inv["id"])
					assert.NotEmpty(t, inv["nomor_referensi"])
					assert.NotEmpty(t, inv["created_at"])
					assert.NotEmpty(t, inv["kontak"])
					kontak, ok := inv["kontak"].(map[string]any)
					assert.True(t, ok)
					assert.NotEmpty(t, kontak["id"])
					assert.NotEmpty(t, kontak["nama"])
					assert.Equal(t, tt.expectedBody.Status, body.Status)
					switch tt.name {
					case "sukses dengan: semua filter":
						v2, err := url.ParseQuery(tt.queryBody[1:])
						assert.NoError(t, err)
						assert.Equal(t, v2.Get("status"), v["status"])
						assert.Equal(t, kontak["id"], v2.Get("kontak_id"))
					case "sukses limit 1":
						assert.Len(t, res, 1)
					case "sukses dengan next":
						assert.NotEmpty(t, v)
						assert.NotEqual(t, idDataBayarTerkonfirmasi, v["id"])
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

func InvoiceGetDataBayar(t *testing.T) {
	if idDataBayar == "" {
		panic(fmt.Sprintf("empty id data bayar %s", idDataBayar))
	}
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
			id:           idDataBayar,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "authorization " + entity.RolesById[2] + " passed",
			token:        tokens[entity.RolesById[2]],
			id:           idInvoice,
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[3] + " passed",
			token:        tokens[entity.RolesById[3]],
			id:           idInvoice,
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			token:        tokens[entity.RolesById[4]],
			id:           idInvoice,
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			token:        tokens[entity.RolesById[5]],
			id:           idInvoice,
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/invoice/data_bayar/"+tt.id, nil, &tt.token)
			assert.NoError(t, err)
			if strings.Contains(tt.name, "passed") {
				assert.NotEqual(t, tt.expectedCode, code)
				assert.NotEqual(t, tt.expectedBody.Code, body.Code)
				assert.NotEqual(t, tt.expectedBody.Status, body.Status)
				return
			}
			assert.Equal(t, tt.expectedCode, code)
			var r map[string]interface{}
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &r)
				assert.NoError(t, err)
				assert.NotEmpty(t, r)

				assert.NotEmpty(t, r)
				assert.NotEmpty(t, r["id"])
				assert.NotEmpty(t, r["created_at"])
				assert.NotEmpty(t, r["invoice_id"])
				assert.NotEmpty(t, r["akun"])
				assert.NotEmpty(t, r["keterangan"])
				assert.NotEmpty(t, r["bukti_pembayaran"])
				assert.NotEmpty(t, r["total"])
				assert.NotEmpty(t, r["status"])
				akun, ok := r["akun"].(map[string]any)
				assert.True(t, ok)
				assert.NotEmpty(t, akun["id"])
				assert.NotEmpty(t, akun["created_at"])
				assert.NotEmpty(t, akun["nama"])
				assert.NotEmpty(t, akun["kode"])
				assert.NotEmpty(t, akun["saldo_normal"])
				assert.NotEmpty(t, akun["deskripsi"])
				

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

func InvoiceGetAllDataBayarByInvoiceId(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		invoiceId    string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			invoiceId:    idInvoice,
			token:        tokens[entity.RolesById[1]],
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "authorization " + entity.RolesById[2] + " passed",
			invoiceId:    idInvoice,
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[3] + " passed",
			invoiceId:    idInvoice,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			invoiceId:    idInvoice,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			invoiceId:    idInvoice,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/invoice/"+tt.invoiceId+"/data_bayar", nil, &tt.token)
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
				for _, r := range res {
					assert.NotEmpty(t, r)
					assert.NotEmpty(t, r["id"])
					assert.NotEmpty(t, r["created_at"])
					assert.NotEmpty(t, r["invoice_id"])
					assert.NotEmpty(t, r["akun"])
					akun, ok := r["akun"].(map[string]any)
					assert.True(t, ok)
					assert.NotEmpty(t, akun["id"])
					assert.NotEmpty(t, akun["created_at"])
					assert.NotEmpty(t, akun["nama"])
					assert.NotEmpty(t, akun["kode"])
					assert.NotEmpty(t, akun["saldo_normal"])
					assert.NotEmpty(t, akun["deskripsi"])
					assert.NotEmpty(t, r["keterangan"])
					assert.NotEmpty(t, r["bukti_pembayaran"])
					assert.NotEmpty(t, r["total"])
					assert.NotEmpty(t, r["status"])

					assert.Equal(t, tt.expectedBody.Status, body.Status)
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

func InvoiceDeleteDataBayar(t *testing.T) {
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
			id:           idDataBayar,
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
			name:         "err: tidak bisa menghapus status yang bernilai terkonfirmasi",
			token:        tokens[entity.RolesById[1]],
			id:           idDataBayarTerkonfirmasi,
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.CannotModifiedTerkonfirmasiDataBayar},
			},
		},
		{
			name:         "err: ulid tidak valid",
			token:        tokens[entity.RolesById[1]],
			id:           idDataBayar + "123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "authorization " + entity.RolesById[2] + " passed",
			id:           idDataBayar,
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[3] + " passed",
			id:           idDataBayar,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			id:           idDataBayar,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idDataBayar,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/invoice/data_bayar/"+tt.id, nil, &tt.token)
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

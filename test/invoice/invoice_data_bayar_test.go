package test_akuntansi

import (
	"errors"
	"fmt"
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
)

func InvoiceCreateDataBayar(t *testing.T) {
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
			name: "authorization " + entity.RolesById[2] + " passed",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/invoice/data_bayar/"+tt.payload.InvoiceID, tt.payload, &tt.token)
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

func InvoiceUpdateDataBayar(t *testing.T) {
	invoiceDataBayar := []entity.DataBayarInvoice{}
	err := dbt.Select("id").Find(&invoiceDataBayar).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	if len(invoiceDataBayar) >= 3 {
		fmt.Println(len(invoiceDataBayar))
		helper.LogsError(errors.New("length not equal 3"))
		return
	}
	tests := []struct {
		name         string
		token        string
		payload      req_invoice_data_bayar.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice_data_bayar.Update{
				ID:              invoiceDataBayar[0].ID,
				Status:          "TERKONFIRMASI",
				BuktiPembayaran: []string{"img.jpg"},
				Keterangan:      "data bayar TERKONFIRMASI",
				AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
				Total:           50000,
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:  "sukses update status data bayar untuk role bendahara",
			token: tokens[entity.RolesById[2]],
			payload: req_invoice_data_bayar.Update{
				ID:              invoiceDataBayar[1].ID,
				Status:          "TERKONFIRMASI",
				BuktiPembayaran: []string{"img2.jpg"},
				Keterangan:      "data bayar belum terkonfirmasi2",
				AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
				Total:           50001,
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:  "sukses akan tetapi status tidak terupdate karena role tidak diperbolehkan",
			token: tokens[entity.RolesById[3]],
			payload: req_invoice_data_bayar.Update{
				ID:              invoiceDataBayar[2].ID,
				Status:          "TERKONFIRMASI",
				BuktiPembayaran: []string{"img3.jpg"},
				Keterangan:      "data bayar TERKONFIRMASI3",
				AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
				Total:           50002,
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
				ID:              invoiceDataBayar[2].ID,
				Status:          "TERKONFIRMASI",
				BuktiPembayaran: []string{"img3.jpg"},
				Keterangan:      "data bayar TERKONFIRMASI3",
				AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
				Total:           50002,
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
			name:         "err: authorization " + entity.RolesById[4],
			payload:      req_invoice_data_bayar.Update{},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			payload:      req_invoice_data_bayar.Update{},
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
			assert.Equal(t, tt.expectedCode, code)
			if len(tt.expectedBody.ErrorsMessages) > 0 {
				for _, v := range tt.expectedBody.ErrorsMessages {
					assert.Contains(t, body.ErrorsMessages, v)
				}
				assert.Equal(t, tt.expectedBody.Status, body.Status)
			} else {
				if strings.Contains(tt.name, "sukses") {
					dbyrHP := new(entity.DataBayarHutangPiutang)
					err := dbt.Preload("Transaksi").Preload("Transaksi.AyatJurnals").First(dbyrHP, "total = ?", tt.payload.Total).Error
					assert.NoError(t, err)
					assert.NotNil(t, dbyrHP)
					if err != nil {
						return
					}
					assert.Equal(t, tt.payload.Keterangan, dbyrHP.Transaksi.Keterangan)
					assert.Equal(t, tt.payload.Total, dbyrHP.Total)
					assert.Equal(t, tt.payload.BuktiPembayaran, dbyrHP.Transaksi.BuktiPembayaran)
					var akunExist bool
					for _, ay := range dbyrHP.Transaksi.AyatJurnals {
						if ay.AkunID == tt.payload.AkunID {
							akunExist = true
							break
						}
					}
					assert.True(t, akunExist)
					dbyrInvoice := new(entity.DataBayarInvoice)
					err = dbt.First(dbyrHP, "total = ?", tt.payload.Total).Error
					assert.NoError(t, err)
					assert.NotNil(t, dbyrInvoice)
					if err != nil {
						helper.LogsError(err)
						return
					}
					if tt.name != "sukses akan tetapi status tidak terupdate karena role tidak diperbolehkan" {
						assert.Equal(t, tt.payload.Status, dbyrInvoice.Status)
					} else {
						assert.NotEqual(t, tt.payload.Status, dbyrInvoice.Status)
					}

				}
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

func InvoiceGetAllByInvoiceIdDataBayar(t *testing.T) {
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/invoice/data_bayar/"+tt.invoiceId, nil, &tt.token)
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
					return
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

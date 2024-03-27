package test_akuntansi

import (
	"strings"
	"testing"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_invoice_data_bayar "github.com/be-sistem-informasi-konveksi/common/request/invoice/data_bayar"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
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

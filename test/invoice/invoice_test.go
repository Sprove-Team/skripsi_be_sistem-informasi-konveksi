package test_akuntansi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_invoice "github.com/be-sistem-informasi-konveksi/common/request/invoice"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
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
			Qty:          2,
			Total:        (produk[1].HargaDetails[0].Harga + bordir[1].Harga + sablon[1].Harga) * 2,
		},
	}
	detailInvoice2 := []req_invoice.ReqDetailInvoice{
		{
			ProdukID:     produk[1].ID,
			BordirID:     bordir[1].ID,
			SablonID:     sablon[1].ID,
			GambarDesign: "img-design-2.webp",
			Qty:          5,
			Total:        (produk[1].HargaDetails[0].Harga + bordir[1].Harga + sablon[1].Harga) * 5,
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
			name:  "err: ulid tidak valid",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				KontakID: "asdfsadfsd",
				Bayar: req_invoice.ReqBayar{
					AkunID: "asdfasdfsf",
					Total:  0,
				},
				DetailInvoice: []req_invoice.ReqDetailInvoice{
					{
						ProdukID: "Asdfsadf",
						BordirID: "Asdfasdf",
						SablonID: "Asdfsaf",
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status: fiber.ErrBadRequest.Message,
				Code:   400,
				ErrorsMessages: []string{
					"kontak id tidak berupa ulid yang valid",
					"akun id tidak berupa ulid yang valid",
					"produk id tidak berupa ulid yang valid",
					"bordir id tidak berupa ulid yang valid",
					"sablon id tidak berupa ulid yang valid",
				},
			},
		},
		{
			name:  "err: total dan qty harus lebih besar dari 0",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				KontakID: "asdfsadfsd",
				Bayar: req_invoice.ReqBayar{
					AkunID: "asdfasdfsf",
					Total:  -1,
				},
				DetailInvoice: []req_invoice.ReqDetailInvoice{
					{
						Total: -1,
						Qty:   -1,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status: fiber.ErrBadRequest.Message,
				Code:   400,
				ErrorsMessages: []string{
					"total harus lebih besar dari 0", "qty harus lebih besar dari 0",
				},
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
			name:  "err: jika kontak id diisi maka new kontak harus kosong",
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
		{
			name:  "err: kontak tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				KontakID: "01HP7DVBGTC06PXWT6FD66VERN",
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
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kontak tidak ditemukan"},
			},
		},
		{
			name:  "err: akun tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				KontakID: kontak[0].ID,
				Bayar: req_invoice.ReqBayar{
					BuktiPembayaran: []string{"img-bukti-2.webp"},
					Keterangan:      "DP 2",
					AkunID:          "01HP7DVBGTC06PXWT6FD66VERC", // kas
					Total:           (detailInvoice2[0].Total + detailInvoice2[0].Total) / 2,
				},
				TanggalDeadline: "2024-10-15T12:00:00Z",
				TanggalKirim:    "2024-10-15T12:00:00Z",
				Keterangan:      "ket invoice 2",
				DetailInvoice:   detailInvoice2,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"akun tidak ditemukan"},
			},
		},
		{
			name:  "err: produk tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				KontakID: kontak[0].ID,
				Bayar: req_invoice.ReqBayar{
					BuktiPembayaran: []string{"img-bukti-2.webp"},
					Keterangan:      "DP 2",
					AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
					Total:           (detailInvoice2[0].Total + detailInvoice2[0].Total) / 2,
				},
				TanggalDeadline: "2024-10-15T12:00:00Z",
				TanggalKirim:    "2024-10-15T12:00:00Z",
				Keterangan:      "ket invoice 2",
				DetailInvoice: []req_invoice.ReqDetailInvoice{
					{
						ProdukID:     "01HP7DVBGTC06PXWT6FD66VERC",
						BordirID:     bordir[1].ID,
						SablonID:     sablon[1].ID,
						GambarDesign: "img-design-2.webp",
						Qty:          5,
						Total:        (produk[1].HargaDetails[0].Harga + bordir[1].Harga + sablon[1].Harga) * 5,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"produk tidak ditemukan"},
			},
		},
		{
			name:  "err: bordir tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				KontakID: kontak[0].ID,
				Bayar: req_invoice.ReqBayar{
					BuktiPembayaran: []string{"img-bukti-2.webp"},
					Keterangan:      "DP 2",
					AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
					Total:           (detailInvoice2[0].Total + detailInvoice2[0].Total) / 2,
				},
				TanggalDeadline: "2024-10-15T12:00:00Z",
				TanggalKirim:    "2024-10-15T12:00:00Z",
				Keterangan:      "ket invoice 2",
				DetailInvoice: []req_invoice.ReqDetailInvoice{
					{
						ProdukID:     produk[1].ID,
						BordirID:     "01HP7DVBGTC06PXWT6FD66VERC",
						SablonID:     sablon[1].ID,
						GambarDesign: "img-design-2.webp",
						Qty:          5,
						Total:        (produk[1].HargaDetails[0].Harga + bordir[1].Harga + sablon[1].Harga) * 5,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"bordir tidak ditemukan"},
			},
		},
		{
			name:  "err: sablon tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Create{
				KontakID: kontak[0].ID,
				Bayar: req_invoice.ReqBayar{
					BuktiPembayaran: []string{"img-bukti-2.webp"},
					Keterangan:      "DP 2",
					AkunID:          "01HP7DVBGTC06PXWT6FD66VERN",
					Total:           (detailInvoice2[0].Total + detailInvoice2[0].Total) / 2,
				},
				TanggalDeadline: "2024-10-15T12:00:00Z",
				TanggalKirim:    "2024-10-15T12:00:00Z",
				Keterangan:      "ket invoice 2",
				DetailInvoice: []req_invoice.ReqDetailInvoice{
					{
						ProdukID:     produk[1].ID,
						BordirID:     bordir[1].ID,
						SablonID:     "01HP7DVBGTC06PXWT6FD66VERC",
						GambarDesign: "img-design-2.webp",
						Qty:          5,
						Total:        (produk[1].HargaDetails[0].Harga + bordir[1].Harga + sablon[1].Harga) * 5,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"sablon tidak ditemukan"},
			},
		},
		{
			name:         "authorization " + entity.RolesById[2] + " passed",
			payload:      req_invoice.Create{},
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[3] + " passed",
			payload:      req_invoice.Create{},
			token:        tokens[entity.RolesById[3]],
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

var dataInvoice entity.Invoice

// []string{"DIREKTUR", "ADMIN", "MANAJER_PRODUKSI", "BENDAHARA"} - allowed,
// []string{"DIREKTUR", "MANAJER_PRODUKSI"} - allowed to modified statusProduksi
func InvoiceUpdate(t *testing.T) {
	invoices := make([]entity.Invoice, 2)
	err := dbt.Preload("DetailInvoice").Order("id ASC").Find(&invoices).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	if len(invoices) != 2 {
		fmt.Println("inv -> ", invoices)
		helper.LogsError(errors.New("err: invoices less than 2"))
		return
	}
	dataInvoice = invoices[0]
	tests := []struct {
		name         string
		token        string
		payload      req_invoice.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Update{
				ID:              invoices[0].ID,
				StatusProduksi:  "DIPROSES",
				TanggalDeadline: "2024-10-11T12:00:00Z",
				TanggalKirim:    "2024-10-11T12:00:00Z",
				Keterangan:      "update invoice 1",
				DetailInvoice: []req_invoice.ReqUpdateDetailInvoice{
					{
						ID:    invoices[0].DetailInvoice[0].ID,
						Qty:   invoices[0].DetailInvoice[0].Qty + 5,
						Total: invoices[0].DetailInvoice[0].Total + ((produk[0].HargaDetails[0].Harga + bordir[0].Harga + sablon[0].Harga) * 5),
					},
				},
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:  "sukses update status produksi",
			token: tokens[entity.RolesById[4]],
			payload: req_invoice.Update{
				ID:             invoices[1].ID,
				StatusProduksi: "DIPROSES",
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:  "err: " + entity.RolesById[2] + " tidak diperbolehkan mengedit status produksi",
			token: tokens[entity.RolesById[2]],
			payload: req_invoice.Update{
				ID:             invoices[1].ID,
				StatusProduksi: "SELESAI",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{strings.ToLower(entity.RolesById[2]) + message.UserNotAllowedToModifiedStatusProdusi},
			},
		},
		{
			name:  "err: " + entity.RolesById[3] + " tidak diperbolehkan mengedit status produksi",
			token: tokens[entity.RolesById[3]],
			payload: req_invoice.Update{
				ID:             invoices[1].ID,
				StatusProduksi: "DIPROSES",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{strings.ToLower(entity.RolesById[3]) + message.UserNotAllowedToModifiedStatusProdusi},
			},
		},
		{
			name:  "err: tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Update{
				ID: "01HM4B8QBH7MWAVAYP10WN6PKA",
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name:  "err: status produksi",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Update{
				ID:             dataInvoice.ID,
				StatusProduksi: "asdfasdf",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"status produksi harus berupa salah satu dari [BELUM_DIKERJAKAN,DIPROSES,SELESAI]"},
			},
		},
		{
			name:  "err: format tanggal kirim & deadline",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Update{
				ID:              dataInvoice.ID,
				TanggalDeadline: "asdfasdf",
				TanggalKirim:    "asdfsadf",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"tanggal deadline harus berformat RFC3999", "tanggal kirim harus berformat RFC3999"},
			},
		},
		{
			name:  "err: detail invoice harus berisi lebih dari 0 item",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Update{
				ID:            dataInvoice.ID,
				DetailInvoice: []req_invoice.ReqUpdateDetailInvoice{},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"detail invoice harus berisi lebih dari 0 item"},
			},
		},
		{
			name:  "err: total dan qty lebih besar dari 0",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Update{
				ID: dataInvoice.ID,
				DetailInvoice: []req_invoice.ReqUpdateDetailInvoice{
					{
						ID:    invoices[0].DetailInvoice[0].ID,
						Total: -1,
						Qty:   -1,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"total harus lebih besar dari 0", "qty harus lebih besar dari 0"},
			},
		},
		{
			name:  "err: wajib diisi",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Update{
				ID: dataInvoice.ID,
				DetailInvoice: []req_invoice.ReqUpdateDetailInvoice{
					{
						Total: 1,
						Qty:   1,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id wajib diisi"},
			},
		},
		{
			name:  "err: ulid tidak valid",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Update{
				ID: dataInvoice.ID,
				DetailInvoice: []req_invoice.ReqUpdateDetailInvoice{
					{
						ID:    dataInvoice.ID + "123",
						Total: 1,
						Qty:   1,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
		{
			name:  "err: detail invoice tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_invoice.Update{
				ID: dataInvoice.ID,
				DetailInvoice: []req_invoice.ReqUpdateDetailInvoice{
					{
						ID:    dataInvoice.ID,
						Total: 1,
						Qty:   1,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"detail invoice tidak ditemukan"},
			},
		},
		{
			name: "authorization " + entity.RolesById[2] + " passed",
			payload: req_invoice.Update{
				ID: dataInvoice.ID + "123",
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
			payload: req_invoice.Update{
				ID: dataInvoice.ID + "123",
			},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name: "authorization " + entity.RolesById[4] + " passed",
			payload: req_invoice.Update{
				ID: dataInvoice.ID + "123",
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
			payload: req_invoice.Update{
				ID: dataInvoice.ID + "123",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/invoice/"+tt.payload.ID, tt.payload, &tt.token)
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

func InvoiceGetAll(t *testing.T) {
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
			name:  "sukses with filter tanggal deadline, kirim dan kontak id",
			token: tokens[entity.RolesById[1]],
			queryBody: fmt.Sprintf(
				"?tanggal_deadline=%s&tanggal_kirim=%s&kontak_id=%s",
				dataInvoice.TanggalDeadline.Format(time.RFC3339),
				dataInvoice.TanggalKirim.Format(time.RFC3339),
				dataInvoice.KontakID,
			),
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		// {
		// 	name:         "sukses with next",
		// 	token:        tokens[entity.RolesById[1]],
		// 	queryBody:    "?next=" + dataInvoice.ID,
		// 	expectedCode: 200,
		// 	expectedBody: test.Response{
		// 		Status: message.OK,
		// 		Code:   200,
		// 	},
		// },
		// {
		// 	name:         "sukses with: filter nama",
		// 	token:        tokens[entity.RolesById[1]],
		// 	queryBody:    "?nama=Bordir+1+baris",
		// 	expectedCode: 200,
		// 	expectedBody: test.Response{
		// 		Status: message.OK,
		// 		Code:   200,
		// 	},
		// },
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/invoice"+tt.queryBody, nil, &tt.token)
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
				d, _ := json.MarshalIndent(res, "", " ")
				fmt.Println("dat -> ", string(d))
				for _, r := range res {
					assert.NotEmpty(t, r)
					assert.NotEmpty(t, r["id"])
					assert.NotEmpty(t, r["created_at"])
					assert.NotEmpty(t, r["status_produksi"])
					assert.NotEmpty(t, r["nomor_referensi"])
					assert.NotEmpty(t, r["total_qty"])
					assert.NotEmpty(t, r["total_harga"])
					assert.NotEmpty(t, r["keterangan"])
					assert.NotEmpty(t, r["tanggal_deadline"])
					assert.NotEmpty(t, r["tanggal_kirim"])
					assert.NotEmpty(t, r["kontak"])
					assert.NotEmpty(t, r["user_editor"])

					kontak, ok := r["kontak"].(map[string]any)
					assert.True(t, ok)
					assert.NotEmpty(t, kontak["id"])
					assert.NotEmpty(t, kontak["nama"])
					assert.NotEmpty(t, kontak["no_telp"])
					assert.NotEmpty(t, kontak["alamat"])
					assert.NotEmpty(t, kontak["email"])
					assert.NotEmpty(t, kontak["keterangan"])
					user, ok := r["user_editor"].(map[string]any)
					assert.True(t, ok)
					assert.NotEmpty(t, user["id"])
					assert.NotEmpty(t, user["nama"])
					assert.NotEmpty(t, user["role"])
					assert.NotEmpty(t, user["username"])

					assert.Equal(t, tt.expectedBody.Status, body.Status)
					switch tt.name {
					case "sukses with: filter nama":
						v, err := url.ParseQuery(tt.queryBody[1:])
						assert.NoError(t, err)
						assert.Equal(t, v.Get("tanggal_deadline"), r["tanggal_deadline"])
						assert.Equal(t, v.Get("tanggal_kirim"), r["tanggal_kirim"])
						assert.Equal(t, v.Get("kontak_id"), r["kontak_id"])
					case "sukses with next":
						// assert.NotEmpty(t, r)
						// assert.NotEqual(t, idBordir, res[0]["id"])
					}
				}
				// case "sukses limit 1":
				// 		assert.Len(t, res, 1)
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

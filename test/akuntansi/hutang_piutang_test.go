package test_akuntansi

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_akuntansi_hp "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func AkuntansiCreateHutangPiutang(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		payload      req_akuntansi_hp.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses with jenis PIUTANG",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.Create{
				KontakID:   idKontak,
				Jenis:      "PIUTANG",
				Keterangan: "menggunakan jasa konveksi seharga 10.000",
				Transaksi: req_akuntansi_hp.ReqTransaksi{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{"bukti-pembayaran.webp"},
					AyatJurnal: []req_akuntansi_hp.ReqAyatJurnal{
						{
							AkunID: "01HP7DVBGWR5ZR6C13RFQKQ2A0", // pendapatan jasa
							Kredit: 10000,
						},
						{
							AkunID: "01HP7DVBGTC06PXWT6FF89WRAB", // piutang usaha
							Debit:  10000,
						},
					},
				},
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name:  "sukses with jenis HUTANG",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.Create{
				KontakID:   idKontak,
				Jenis:      "HUTANG",
				Keterangan: "perusahaan membeli peralatan dengan kredit sebesar 20.0000",
				Transaksi: req_akuntansi_hp.ReqTransaksi{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{"bukti-pembayaran.webp"},
					AyatJurnal: []req_akuntansi_hp.ReqAyatJurnal{
						{
							AkunID: "01HP7DVBGVHMSA4VHWXZR27J7C", // hutang usaha
							Kredit: 20000,
						},
						{
							AkunID: "01HP7DVBGVHMSA4VHWXPQ3635H", // peralatan
							Debit:  20000,
						},
					},
				},
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name:  "err: kredit wajib diisi jika debit tidak diisi, begitu sebaliknya",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.Create{
				KontakID:   idKontak,
				Jenis:      "PIUTANG",
				Keterangan: "menggunakan jasa konveksi seharga 10.000",
				Transaksi: req_akuntansi_hp.ReqTransaksi{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{"bukti-pembayaran.webp"},
					AyatJurnal: []req_akuntansi_hp.ReqAyatJurnal{
						{
							AkunID: "01HP7DVBGWR5ZR6C13RFQKQ2A0", // pendapatan jasa
						},
						{
							AkunID: "01HP7DVBGTC06PXWT6FF89WRAB", // piutang usaha
							Debit:  10000,
						},
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kredit wajib diisi jika debit tidak diisi", "debit wajib diisi jika kredit tidak diisi"},
			},
		},
		{
			name:  "err: peletakan total debit dan kredit untuk hutang piutang tidak benar",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.Create{
				KontakID:   idKontak,
				Jenis:      "PIUTANG",
				Keterangan: "menggunakan jasa konveksi seharga 10.000",
				Transaksi: req_akuntansi_hp.ReqTransaksi{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{"bukti-pembayaran.webp"},
					AyatJurnal: []req_akuntansi_hp.ReqAyatJurnal{
						{
							AkunID: "01HP7DVBGWR5ZR6C13RFQKQ2A0", // pendapatan jasa
							Debit:  10000,
						},
						{
							AkunID: "01HP7DVBGTC06PXWT6FF89WRAB", // piutang usaha
							Debit:  10000,
						},
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.IncorrectPlacementOfCreditAndDebit},
			},
		},
		{
			name:  "err: total debit dan kredit harus sama",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.Create{
				KontakID:   idKontak,
				Jenis:      "PIUTANG",
				Keterangan: "menggunakan jasa konveksi seharga 10.000",
				Transaksi: req_akuntansi_hp.ReqTransaksi{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{"bukti-pembayaran.webp"},
					AyatJurnal: []req_akuntansi_hp.ReqAyatJurnal{
						{
							AkunID: "01HP7DVBGWR5ZR6C13RFQKQ2A0", // pendapatan jasa
							Kredit: 100000,
						},
						{
							AkunID: "01HP7DVBGTC06PXWT6FF89WRAB", // piutang usaha
							Debit:  10000,
						},
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.CreditDebitNotSame},
			},
		},
		{
			name:  "err: wajib diisi akun id pada ayat jurnal",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.Create{
				KontakID:   idKontak,
				Jenis:      "PIUTANG",
				Keterangan: "menggunakan jasa konveksi seharga 10.000",
				Transaksi: req_akuntansi_hp.ReqTransaksi{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{"bukti-pembayaran.webp"},
					AyatJurnal: []req_akuntansi_hp.ReqAyatJurnal{
						{
							Debit: 10000,
						},
						{
							Kredit: 10000,
						},
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"akun id wajib diisi"},
			},
		},
		{
			name:  "err: wajib diisi di field lainnya",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.Create{
				KontakID:   "",
				Jenis:      "",
				Keterangan: "",
				Transaksi:  req_akuntansi_hp.ReqTransaksi{},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kontak id wajib diisi", "tanggal wajib diisi", "ayat jurnal wajib diisi", "jenis wajib diisi", "keterangan wajib diisi"},
			},
		},
		{
			name:  "err: wajib diisi di field lainnya",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.Create{
				KontakID:   "",
				Jenis:      "",
				Keterangan: "",
				Transaksi:  req_akuntansi_hp.ReqTransaksi{},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kontak id wajib diisi", "tanggal wajib diisi", "ayat jurnal wajib diisi", "jenis wajib diisi", "keterangan wajib diisi"},
			},
		},
		{
			name:  "err: ulid tidak valid",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.Create{
				KontakID:   idKontak + "124",
				Jenis:      "PIUTANG",
				Keterangan: "menggunakan jasa konveksi seharga 10.000",
				Transaksi: req_akuntansi_hp.ReqTransaksi{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{"bukti-pembayaran.webp"},
					AyatJurnal: []req_akuntansi_hp.ReqAyatJurnal{
						{
							AkunID: "01HP7DVBGWR5ZR6C13RFQKQ2A0124", // pendapatan jasa
							Kredit: 10000,
						},
						{
							AkunID: "01HP7DVBGTC06PXWT6FF89WRAB124", // piutang usaha
							Debit:  10000,
						},
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kontak id tidak berupa ulid yang valid", "akun id tidak berupa ulid yang valid"},
			},
		},
		{
			name:  "err: incorrect entry akun HP",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.Create{
				KontakID:   idKontak,
				Jenis:      "HUTANG",
				Keterangan: "menggunakan jasa konveksi seharga 10.000",
				Transaksi: req_akuntansi_hp.ReqTransaksi{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{"bukti-pembayaran.webp"},
					AyatJurnal: []req_akuntansi_hp.ReqAyatJurnal{
						{
							AkunID: "01HP7DVBGWR5ZR6C13RFQKQ2A0", // pendapatan jasa
							Kredit: 10000,
						},
						{
							AkunID: "01HP7DVBGTC06PXWT6FF89WRAB", // piutang usaha
							Debit:  10000,
						},
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.IncorrectEntryAkunHP},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			payload:      req_akuntansi_hp.Create{},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			payload:      req_akuntansi_hp.Create{},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			payload:      req_akuntansi_hp.Create{},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/akuntansi/hutang_piutang", tt.payload, &tt.token)
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

var idTransaksiWithHP string

func AkuntansiGetAllHutangPiutang(t *testing.T) {
	// create invoice for test data hp that invoice id
	{
		tt := time.Now().Add(time.Hour * 24)
		invoice := &entity.Invoice{
			Base: entity.Base{
				ID: test.UlidPkg.MakeUlid().String(),
			},
			NomorReferensi:  "001",
			TanggalDeadline: &tt,
			TanggalKirim:    &tt,
			Keterangan:      "test for hp",
			KontakID:        idKontak,
			TotalQty:        10,
			TotalHarga:      10000,
		}
		if err := dbt.Create(invoice).Error; err != nil {
			helper.LogsError(err)
			return
		}
		tr := &entity.Transaksi{
			Base: entity.Base{
				ID: test.UlidPkg.MakeUlid().String(),
			},
			Keterangan: "test for hp",
			KontakID:   idKontak,
			Total:      10000,
			Tanggal:    tt,
		}
		if err := dbt.Create(tr).Error; err != nil {
			helper.LogsError(err)
			return
		}
		idHp := test.UlidPkg.MakeUlid().String()
		hpWithInvoiceId := &entity.HutangPiutang{
			Base: entity.Base{
				ID: idHp,
			},
			InvoiceID:   invoice.ID,
			TransaksiID: tr.ID,
			Jenis:       "PIUTANG",
			Total:       10000,
			Sisa:        10000,
		}

		if err := dbt.Create(hpWithInvoiceId).Error; err != nil {
			helper.LogsError(err)
			return
		}
	}

	dataKontak := new(entity.Kontak)
	if err := dbt.First(dataKontak, "id = ?", idKontak).Error; err != nil {
		helper.LogsError(err)
		return
	}

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
			queryBody:    "?status=BELUM_LUNAS&jenis=HUTANG&kontak_id=" + dataKontak.ID,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		// {
		// 	name:         "sukses limit 1",
		// 	token:        tokens[entity.RolesById[1]],
		// 	queryBody:    "?limit=1",
		// 	expectedCode: 200,
		// 	expectedBody: test.Response{
		// 		Status: message.OK,
		// 		Code:   200,
		// 	},
		// },
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/hutang_piutang"+tt.queryBody, nil, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)

			var res []map[string]interface{}
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				length := len(res)
				assert.Greater(t, length, 0)
				if length <= 0 {
					return
				}
				assert.NotEmpty(t, res[0])
				assert.NotEmpty(t, res[0]["nama"])
				assert.NotEmpty(t, res[0]["total_piutang"])
				assert.NotEmpty(t, res[0]["total_hutang"])
				assert.NotEmpty(t, res[0]["sisa_piutang"])
				assert.NotEmpty(t, res[0]["sisa_hutang"])
				assert.NotEmpty(t, res[0]["hutang_piutang"])
				var invoiceIdEverExist bool
				for _, hp := range res[0]["hutang_piutang"].([]any) {
					hp2 := hp.(map[string]any)
					assert.NotEmpty(t, hp2["id"])
					assert.NotEmpty(t, hp2["jenis"])
					assert.NotEmpty(t, hp2["transaksi_id"])
					assert.NotEmpty(t, hp2["tanggal"])
					assert.NotEmpty(t, hp2["status"])
					assert.NotEmpty(t, hp2["total"])
					assert.NotEmpty(t, hp2["sisa"])
					if idTransaksiWithHP == "" {
						idTransaksiWithHP = hp2["transaksi_id"].(string)
					}
					_, ok := hp2["invoice_id"]
					if ok {
						invoiceIdEverExist = true
					}
				}

				assert.Equal(t, tt.expectedBody.Status, body.Status)
				switch tt.name {
				case "sukses":
					assert.True(t, invoiceIdEverExist)
				case "sukses with filter":
					v, err := url.ParseQuery(tt.queryBody[1:])
					assert.NoError(t, err)
					assert.Equal(t, res[0]["status"], v.Get("status"))
					assert.Equal(t, res[0]["jenis"], v.Get("jenis"))
					assert.Equal(t, res[0]["nama"], dataKontak.Nama)
				case "sukses limit 1":
					assert.Len(t, res, 1)
				case "sukses with next":
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

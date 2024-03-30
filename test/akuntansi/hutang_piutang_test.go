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
			name:  "sukses dengan jenis PIUTANG",
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
			name:  "sukses dengan jenis HUTANG",
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
			name:         "authorization " + entity.RolesById[2] + " passed",
			token:        tokens[entity.RolesById[2]],
			expectedCode: 400,
			expectedBody: test.Response{
				Status: fiber.ErrBadRequest.Message,
				Code:   400,
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
				if strings.Contains(tt.name, "passed") {
					return
				}
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

var idTransaksiWithHP string

func AkuntansiGetAllHutangPiutang(t *testing.T) {
	// create invoice for test data hp that invoice id

	tt := time.Now().Add(time.Hour * 24)
	invoice := &entity.Invoice{
		Base: entity.Base{
			ID: test.UlidPkg.MakeUlid().String(),
		},
		NomorReferensi:  "001",
		TanggalDeadline: &tt,
		TanggalKirim:    &tt,
		Keterangan:      "test for hp",
		KontakID:        idKontak2,
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
		KontakID:   idKontak2,
		Total:      10000,
		Tanggal:    tt,
	}
	if err := dbt.Create(tr).Error; err != nil {
		helper.LogsError(err)
		return
	}
	hpWithInvoiceId := &entity.HutangPiutang{
		Base: entity.Base{
			ID: test.UlidPkg.MakeUlid().String(),
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
			name:         "sukses dengan filter",
			token:        tokens[entity.RolesById[1]],
			queryBody:    "?status=BELUM_LUNAS&jenis=HUTANG&kontak_id=" + dataKontak.ID,
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
			queryBody:    "?next=" + idKontak,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "authorization " + entity.RolesById[2] + " passed",
			token:        tokens[entity.RolesById[2]],
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
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

				for _, v := range res {
					assert.NotEmpty(t, v)
					assert.NotEmpty(t, v["nama"])
					assert.NotEmpty(t, v["kontak_id"])
					assert.NotEmpty(t, v["hutang_piutang"])
					var totalPiutang, totalHutang, sisaPiutang, sisaHutang int
					for _, hp := range v["hutang_piutang"].([]any) {
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
						invID, ok := hp2["invoice_id"]
						if ok {
							assert.Equal(t, invID, invoice.ID)
						}
						if hp2["jenis"] == "PIUTANG" {
							totalPiutang += int(hp2["total"].(float64))
							sisaPiutang += int(hp2["sisa"].(float64))
						} else {
							totalHutang += int(hp2["total"].(float64))
							sisaHutang += int(hp2["sisa"].(float64))
						}
						if tt.name == "sukses dengan filter" {
							vq, err := url.ParseQuery(tt.queryBody[1:])
							assert.NoError(t, err)
							assert.Equal(t, hp2["status"], vq.Get("status"))
							assert.Equal(t, hp2["jenis"], vq.Get("jenis"))
						}
					}
					switch tt.name {
					case "sukses":
						if dat, ok := v["total_piutang"].(float64); ok && dat != 0 {
							assert.Greater(t, int(dat), 0)
						}
						if dat, ok := v["total_hutang"].(float64); ok && dat != 0 {
							assert.Greater(t, int(dat), 0)
						}
						if dat, ok := v["sisa_piutang"].(float64); ok && dat != 0 {
							assert.Greater(t, int(dat), 0)
						}
						if dat, ok := v["sisa_hutang"].(float64); ok && dat != 0 {
							assert.Greater(t, int(dat), 0)
						}
					case "sukses dengan filter":
						// karna jenis yg difilter HUTANG
						assert.Empty(t, v["total_piutang"])
						assert.Empty(t, v["sisa_piutang"])
						assert.NotEmpty(t, v["total_hutang"])
						assert.NotEmpty(t, v["sisa_hutang"])
						assert.Equal(t, v["nama"], dataKontak.Nama)
						assert.Equal(t, v["kontak_id"], dataKontak.ID)
					case "sukses dengan next":
						assert.NotEqual(t, dataKontak.ID, v["id"])
					}
				}

				if tt.name == "sukses limit 1" {
					assert.Len(t, res, 1)
				}
				assert.Equal(t, tt.expectedBody.Status, body.Status)

			} else {
				if len(tt.expectedBody.ErrorsMessages) > 0 {
					for _, v := range tt.expectedBody.ErrorsMessages {
						assert.Contains(t, body.ErrorsMessages, v)
					}
					assert.Equal(t, tt.expectedBody.Status, body.Status)
				} else {
					if strings.Contains(tt.name, "passed") {
						return
					}
					assert.Equal(t, tt.expectedBody, body)
				}
			}

		})
	}
	if err := dbt.Unscoped().Delete(hpWithInvoiceId).Error; err != nil {
		helper.LogsError(err)
		return
	}
	if err := dbt.Unscoped().Delete(invoice).Error; err != nil {
		helper.LogsError(err)
		return
	}
	if err := dbt.Unscoped().Delete(tr).Error; err != nil {
		helper.LogsError(err)
		return
	}
}

var idTrWithBayarHP string

func AkuntansiCreateBayarHP(t *testing.T) {
	hp := new(entity.HutangPiutang)
	if err := dbt.First(hp).Error; err != nil {
		helper.LogsError(err)
		return
	}
	tests := []struct {
		name         string
		token        string
		payload      req_akuntansi_hp.CreateBayar
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.CreateBayar{
				HutangPiutangID: hp.ID,
				ReqBayar: req_akuntansi_hp.ReqBayar{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{"bukti-pembayaran.jpg"},
					Keterangan:      "bayar lunas",
					AkunBayarID:     "01HP7DVBGTC06PXWT6FD66VERN",
					Total:           5000,
				},
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name:  "err: bayar too much",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.CreateBayar{
				HutangPiutangID: hp.ID,
				ReqBayar: req_akuntansi_hp.ReqBayar{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{"bukti-pembayaran.jpg"},
					Keterangan:      "bayar lunas",
					AkunBayarID:     "01HP7DVBGTC06PXWT6FD66VERN",
					Total:           100000,
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.BayarMustLessThanSisaTagihan},
			},
		},
		{
			name:  "err: wajib diisi",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.CreateBayar{
				HutangPiutangID: hp.ID,
				ReqBayar: req_akuntansi_hp.ReqBayar{
					Tanggal:     "",
					Keterangan:  "",
					AkunBayarID: "",
					Total:       0,
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"total wajib diisi", "tanggal wajib diisi", "bukti pembayaran wajib diisi", "keterangan wajib diisi", "akun bayar id wajib diisi"},
			},
		},
		{
			name:  "err: bukti pembayaran haru memliki item lebih dari 0",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.CreateBayar{
				HutangPiutangID: hp.ID,
				ReqBayar: req_akuntansi_hp.ReqBayar{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{},
					Keterangan:      "bayar lunas",
					AkunBayarID:     "01HP7DVBGTC06PXWT6FD66VERN",
					Total:           10000,
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"bukti pembayaran harus berisi lebih dari 0 item"},
			},
		},
		{
			name:  "err: format tanggal",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.CreateBayar{
				HutangPiutangID: hp.ID,
				ReqBayar: req_akuntansi_hp.ReqBayar{
					Tanggal:         "2024-01-01",
					BuktiPembayaran: []string{},
					Keterangan:      "bayar lunas",
					AkunBayarID:     "01HP7DVBGTC06PXWT6FD66VERN",
					Total:           10000,
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"tanggal harus berformat RFC3339"},
			},
		},
		{
			name:  "err: ulid tidak valid",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.CreateBayar{
				HutangPiutangID: hp.ID + "123",
				ReqBayar: req_akuntansi_hp.ReqBayar{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{},
					Keterangan:      "bayar lunas",
					AkunBayarID:     "01HP7DVBGTC06PXWT6FD66VERN123",
					Total:           10000,
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"hutang piutang id tidak berupa ulid yang valid", "akun bayar id tidak berupa ulid yang valid"},
			},
		},
		{
			name:  "err: total harus lebih dari 0",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_hp.CreateBayar{
				HutangPiutangID: hp.ID + "123",
				ReqBayar: req_akuntansi_hp.ReqBayar{
					Tanggal:         "2024-01-01T14:17:03.723Z",
					BuktiPembayaran: []string{"bukti-pembayaran.jpg"},
					Keterangan:      "bayar lunas",
					AkunBayarID:     "01HP7DVBGTC06PXWT6FD66VERN",
					Total:           -10000,
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"total harus lebih besar dari 0"},
			},
		},
		{
			name:  "authorization " + entity.RolesById[2] + " passed",
			token: tokens[entity.RolesById[2]],
			payload: req_akuntansi_hp.CreateBayar{
				HutangPiutangID: hp.ID + "123",
			},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/akuntansi/hutang_piutang/bayar/"+tt.payload.HutangPiutangID, tt.payload, &tt.token)
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
				if tt.name == "sukses" {
					datResAfterBayar := new(entity.HutangPiutang)
					if err := dbt.Preload("DataBayarHutangPiutang").First(datResAfterBayar, "id = ?", hp.ID).Error; err != nil {
						helper.LogsError(err)
						return
					}
					assert.NotEmpty(t, datResAfterBayar.DataBayarHutangPiutang)
					if len(datResAfterBayar.DataBayarHutangPiutang) > 0 {
						idTrWithBayarHP = datResAfterBayar.DataBayarHutangPiutang[0].TransaksiID
					}
					sisa := hp.Sisa - tt.payload.Total
					assert.Equal(t, datResAfterBayar.Sisa, sisa)
					if sisa <= 0 {
						assert.Equal(t, datResAfterBayar.Status, "LUNAS")
					} else {
						assert.Equal(t, datResAfterBayar.Status, "BELUM_LUNAS")
					}
				}
			}
		})
	}
}

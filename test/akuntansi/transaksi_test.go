package test_akuntansi

import (
	"math"
	"strings"
	"testing"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_akuntansi_transaksi "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/transaksi"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func AkuntansiCreateTransaksi(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		payload      req_akuntansi_transaksi.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses without kontak",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Create{
				BuktiPembayaran: []string{"bukti-pembayaran.webp"},
				Tanggal:         "2023-10-30T14:17:03.723Z",
				Keterangan:      "Membayar beban listrik dan air sebesar Rp. 10.000,-",
				KontakID:        "",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEQTSH53Y",
						Debit:  10000,
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN",
						Kredit: 10000,
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
			name:  "sukses with kontak",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Create{
				BuktiPembayaran: []string{"bukti-pembayaran2.webp"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "Membayar beban listrik dan air sebesar Rp. 10.000,- with kontak",
				KontakID:        idKontak,
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEQTSH53Y",
						Debit:  10000,
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN",
						Kredit: 10000,
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
			payload: req_akuntansi_transaksi.Create{
				BuktiPembayaran: []string{"bukti-pembayaran2.webp"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "Membayar beban listrik dan air sebesar Rp. 10.000,- with kontak",
				KontakID:        idKontak,
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEQTSH53Y",
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN",
						Kredit: 10000,
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
			name:  "err: total debit dan kredit harus sama",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Create{
				BuktiPembayaran: []string{"bukti-pembayaran2.webp"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "Membayar beban listrik dan air sebesar Rp. 10.000,- with kontak",
				KontakID:        idKontak,
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEQTSH53Y",
						Kredit: 10000,
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN",
						Kredit: 10000,
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
			name:  "err: duplicate akun",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Create{
				BuktiPembayaran: []string{"bukti-pembayaran2.webp"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "err",
				KontakID:        idKontak,
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEQTSH53Y",
						Kredit: 10000,
					},
					{
						AkunID: "01HP7DVBGX4JR0KETMEQTSH53Y",
						Debit:  10000,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.AkunCannotBeSame},
			},
		},
		{
			name:  "err: akun not found",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Create{
				BuktiPembayaran: []string{"bukti-pembayaran2.webp"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "err",
				KontakID:        idKontak,
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: test.UlidPkg.MakeUlid().String(),
						Kredit: 10000,
					},
					{
						AkunID: test.UlidPkg.MakeUlid().String(),
						Debit:  10000,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.AkunNotFound},
			},
		},
		{
			name:  "err: wajib diisi akun id pada ayat jurnal",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Create{
				BuktiPembayaran: []string{"bukti-pembayaran2.webp"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "err",
				KontakID:        idKontak,
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						Debit: 10000,
					},
					{
						Kredit: 10000,
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
			name:  "err: tanggal harus berformat RFC3999",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Create{
				BuktiPembayaran: []string{"bukti-pembayaran3.webp"},
				Tanggal:         "2023-10-28",
				Keterangan:      "err",
				KontakID:        idKontak,
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEQTSH53Y",
						Kredit: 10000,
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN",
						Kredit: 10000,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"tanggal harus berformat RFC3999"},
			},
		},
		{
			name:  "err: panjang minimal ayat jurnal adalah 2 item",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Create{
				BuktiPembayaran: []string{"bukti-pembayaran3.webp"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "err",
				KontakID:        idKontak,
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEQTSH53Y",
						Kredit: 10000,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"panjang minimal ayat jurnal adalah 2 item"},
			},
		},
		{
			name:  "err: ulid tidak valid",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Create{
				BuktiPembayaran: []string{"bukti-pembayaran3.webp"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "err",
				KontakID:        "ADFADCADFADFASDFASFSAFASFASFDASDFASDF123",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "ADFADCADFADFASDFASFSAFASFASFDASDFASDF123",
						Kredit: 10000,
					},
					{
						AkunID: "ADFADCADFADFASDFASFSAFASFASFDASDFASDF123",
						Kredit: 10000,
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
			name:  "err: wajib diisi",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Create{
				BuktiPembayaran: nil,
				Tanggal:         "",
				Keterangan:      "",
				KontakID:        "",
				AyatJurnal:      nil,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"tanggal wajib diisi", "keterangan wajib diisi", "ayat jurnal wajib diisi"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			payload:      req_akuntansi_transaksi.Create{},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			payload:      req_akuntansi_transaksi.Create{},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			payload:      req_akuntansi_transaksi.Create{},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/akuntansi/transaksi", tt.payload, &tt.token)
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

var idTransaksi string
var idTransaksiWithKontak string

func AkuntansiUpdateTransaksi(t *testing.T) {
	transaksi := new(entity.Transaksi)
	err := dbt.Model(transaksi).First(transaksi, "kontak_id IS NULL").Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	transaksiWithKontak := new(entity.Transaksi)
	err = dbt.Model(transaksi).First(&transaksiWithKontak, "kontak_id IS NOT NULL").Error
	if err != nil {
		helper.LogsError(err)
		return
	}

	idTransaksi = transaksi.ID
	idTransaksiWithKontak = transaksiWithKontak.ID

	tests := []struct {
		name         string
		token        string
		payload      req_akuntansi_transaksi.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              idTransaksi,
				BuktiPembayaran: []string{"bukti-update.webp"},
				Tanggal:         "2023-10-10T14:17:03.723Z",
				Keterangan:      "Membayar beban gaji Rp. 15.000,- update",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEJXMZCMF", // beban gaji
						Debit:  15000,
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
						Kredit: 15000,
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
			name:  "sukses update tr hp",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              idTransaksiWithHP, // jenis PIUTANG
				BuktiPembayaran: []string{"bukti-update-hp.webp"},
				Tanggal:         "2023-10-20T14:17:03.723Z",
				Keterangan:      "joni menggunakan jasa konveksi sebesar 20.000 dan dibayar bulan depan",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGTC06PXWT6FF89WRAB", // piutang usaha
						Debit:  20000,
					},
					{
						AkunID: "01HP7DVBGVHMSA4VHWXPQ3635H", // peralatan
						Kredit: 20000,
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
			name:  "sukses update tr bayar hp",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              idTrWithBayarHP, // tr bayar piutang, sisa hp jadi 15.000 karena di ubah di test pertama dan sudah dibayar 5.000
				BuktiPembayaran: []string{"bukti-update-hp.webp"},
				Tanggal:         "2023-10-20T14:17:03.723Z",
				Keterangan:      "joni jadinya melunasi utangnya",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGTC06PXWT6FF89WRAB", // piutang usaha
						Kredit: 20000,
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
						Debit:  20000,
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
			name:  "err: transaksi merupakan hutang piutang, akun harus sama dengan jenis hutang piutang",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              idTransaksiWithHP, // jenis PIUTANG
				BuktiPembayaran: []string{"bukti-update-hp.webp"},
				Tanggal:         "2023-10-20T14:17:03.723Z",
				Keterangan:      "joni menggunakan jasa konveksi sebesar 20.000 dan dibayar bulan depan",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGVHMSA4VHWXZR27J7C", // hutang usaha
						Debit:  20000,
					},
					{
						AkunID: "01HP7DVBGVHMSA4VHWXPQ3635H", // peralatan
						Kredit: 20000,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.AkunNotMatchWithJenisHPTr},
			},
		},
		{
			name:  "err: akun hutang piutang tidak ada",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              idTransaksiWithHP, // jenis PIUTANG
				BuktiPembayaran: []string{"bukti-update-hp.webp"},
				Tanggal:         "2023-10-20T14:17:03.723Z",
				Keterangan:      "joni menggunakan jasa konveksi sebesar 20.000 dan dibayar bulan depan",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
						Debit:  20000,
					},
					{
						AkunID: "01HP7DVBGVHMSA4VHWXPQ3635H", // peralatan
						Kredit: 20000,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.AkunHPDoesNotExist},
			},
		},
		{
			name:  "err: total hp harus lebih besar atau sama dengan total yang telah di bayar",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              idTransaksiWithHP, // jenis PIUTANG, total yg telah dibayar sebesar 5000 based on hp test
				BuktiPembayaran: []string{"bukti-update-hp.webp"},
				Tanggal:         "2023-10-20T14:17:03.723Z",
				Keterangan:      "joni menggunakan jasa konveksi sebesar 4000 dan dibayar bulan depan",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGWR5ZR6C13RFQKQ2A0", // pendapatan jasa
						Kredit: 4000,
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FF89WRAB", // piutang usaha
						Debit:  4000,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.TotalHPMustGeOrEqToTotalByr},
			},
		},
		{
			name:  "err: total yang dibayar harus lebih kecil atau sama dengan sisa tagihan",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              idTrWithBayarHP,
				BuktiPembayaran: []string{"bukti-update-hp.webp"},
				Tanggal:         "2023-10-20T14:17:03.723Z",
				Keterangan:      "joni membayar sebesar 10000000",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGTC06PXWT6FF89WRAB", // piutang usaha
						Kredit: 10000000,
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
						Debit:  10000000,
					},
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
			name:  "err: akun tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              idTransaksiWithHP, // jenis PIUTANG
				BuktiPembayaran: []string{"bukti-update-hp.webp"},
				Tanggal:         "2023-10-20T14:17:03.723Z",
				Keterangan:      "joni menggunakan jasa konveksi sebesar 20.000 dan dibayar bulan depan",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "02HP7DVBGTC06PXWT6FD66VERN", // akun not found
						Debit:  20000,
					},
					{
						AkunID: "01HP7DVBGVHMSA4VHWXPQ3635H", // peralatan
						Kredit: 20000,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.AkunNotFound},
			},
		},
		{
			name:  "err: tanggal harus berformat RFC3999",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              idTransaksi,
				BuktiPembayaran: []string{"bukti-pembayaran3.webp"},
				Tanggal:         "2023-10-28",
				Keterangan:      "Membayar beban listrik dan air sebesar Rp. 10.000,3",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEQTSH53Y",
						Kredit: 10000,
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN",
						Kredit: 10000,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"tanggal harus berformat RFC3999"},
			},
		},
		{
			name:  "err: total debit dan kredit harus sama",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              idTransaksi,
				BuktiPembayaran: []string{"bukti-pembayaran2.webp"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "err",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEQTSH53Y",
						Kredit: 10000,
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN",
						Kredit: 10000,
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
			payload: req_akuntansi_transaksi.Update{
				ID:              idTransaksi,
				BuktiPembayaran: []string{"bukti-pembayaran2.webp"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "err",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						Debit: 10000,
					},
					{
						Kredit: 10000,
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
			name:  "err: panjang minimal ayat jurnal adalah 2 item",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              idTransaksi,
				BuktiPembayaran: []string{"bukti-pembayaran3.webp"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "err",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEQTSH53Y",
						Kredit: 10000,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"panjang minimal ayat jurnal adalah 2 item"},
			},
		},
		{
			name:  "err: tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              "01HM4B8QBH7MWAVAYP10WN6PKA",
				BuktiPembayaran: []string{"bukti-update.webp"},
				Tanggal:         "2023-10-10T14:17:03.723Z",
				Keterangan:      "err",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEJXMZCMF", // beban gaji
						Debit:  15000,
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
						Kredit: 15000,
					},
				},
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name:  "err: ulid tidak valid",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              idAkun + "123",
				BuktiPembayaran: []string{"bukti-update.webp"},
				Tanggal:         "2023-10-10T14:17:03.723Z",
				Keterangan:      "err",
				AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
					{
						AkunID: "01HP7DVBGX4JR0KETMEJXMZCMF1234", // beban gaji
						Debit:  15000,
					},
					{
						AkunID: "01HP7DVBGTC06PXWT6FD66VERN1234", // kas
						Kredit: 15000,
					},
				},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid", "akun id tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			payload:      req_akuntansi_transaksi.Update{},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			payload:      req_akuntansi_transaksi.Update{},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			payload:      req_akuntansi_transaksi.Update{},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/akuntansi/transaksi/"+tt.payload.ID, tt.payload, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			if len(tt.expectedBody.ErrorsMessages) > 0 {
				for _, v := range tt.expectedBody.ErrorsMessages {
					assert.Contains(t, body.ErrorsMessages, v)
				}
				assert.Equal(t, tt.expectedBody.Status, body.Status)
			} else {
				if strings.Contains(tt.name, "sukses") {
					hp := new(entity.HutangPiutang)
					if err := dbt.First(hp, "transaksi_id = ?", idTransaksiWithHP).Error; err != nil {
						helper.LogsError(err)
						return
					}

					totalUpdate := math.Abs(tt.payload.AyatJurnal[0].Debit - tt.payload.AyatJurnal[0].Kredit)
					switch tt.name {
					case "sukses update tr hp":
						totalByrOld := (hp.Total - hp.Sisa)
						assert.Equal(t, hp.Total, totalUpdate)
						assert.Equal(t, hp.Sisa, totalUpdate-totalByrOld)
					case "sukses update tr bayar hp":
						sisaCurrent := hp.Total - totalUpdate
						assert.Equal(t, hp.Sisa, sisaCurrent)
					}

					if hp.Sisa <= 0 {
						assert.Equal(t, hp.Status, "LUNAS")
					} else {
						assert.Equal(t, hp.Status, "BELUM_LUNAS")
					}
				}
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

func AkuntansiGetAllTransaksi(t *testing.T) {
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
			queryBody:    "?start_date=2023-10-01&end_date=2024-12-31",
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: start_date & end_date wajib diisi",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"start date wajib diisi", "end date wajib diisi"},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/transaksi"+tt.queryBody, nil, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)

			var res []map[string]interface{}
			if tt.name == "sukses" {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)

				// data without kontak
				assert.Greater(t, len(res), 0)
				for _, v := range res {
					assert.NotEmpty(t, v)
					assert.NotEmpty(t, v["id"])
					assert.NotEmpty(t, v["created_at"])
					assert.NotEmpty(t, v["keterangan"])
					assert.NotEmpty(t, v["total"])
					assert.NotEmpty(t, v["tanggal"])
					if v["id"] == idTransaksiWithKontak {
						assert.NotEmpty(t, v["kontak"])
						assert.NotEmpty(t, v["kontak"].(map[string]any)["id"])
						assert.NotEmpty(t, v["kontak"].(map[string]any)["created_at"])
						assert.NotEmpty(t, v["kontak"].(map[string]any)["nama"])
						assert.NotEmpty(t, v["kontak"].(map[string]any)["no_telp"])
						assert.NotEmpty(t, v["kontak"].(map[string]any)["alamat"])
						assert.NotEmpty(t, v["kontak"].(map[string]any)["email"])
						assert.NotEmpty(t, v["kontak"].(map[string]any)["keterangan"])
					}
				}

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

func AkuntansiGetTransaksi(t *testing.T) {
	tests := []struct {
		id           string
		name         string
		token        string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses without kontak",
			token:        tokens[entity.RolesById[1]],
			id:           idTransaksi,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses with kontak",
			token:        tokens[entity.RolesById[1]],
			id:           idTransaksiWithKontak,
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
			id:           idTransaksi,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			id:           idTransaksi,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idTransaksi,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/transaksi/"+tt.id, nil, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)

			var res map[string]any
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				assert.NotEmpty(t, res)
				assert.NotEmpty(t, res["id"])
				assert.NotEmpty(t, res["created_at"])
				assert.NotEmpty(t, res["keterangan"])
				assert.NotEmpty(t, res["total"])
				assert.NotEmpty(t, res["tanggal"])
				assert.Equal(t, tt.expectedBody.Status, body.Status)
				switch tt.name {
				case "sukses without kontak":
					assert.Empty(t, res["kontak"])
				case "sukses with kontak":
					assert.NotEmpty(t, res["kontak"])
					assert.NotEmpty(t, res["kontak"].(map[string]any)["id"])
					assert.NotEmpty(t, res["kontak"].(map[string]any)["created_at"])
					assert.NotEmpty(t, res["kontak"].(map[string]any)["nama"])
					assert.NotEmpty(t, res["kontak"].(map[string]any)["no_telp"])
					assert.NotEmpty(t, res["kontak"].(map[string]any)["alamat"])
					assert.NotEmpty(t, res["kontak"].(map[string]any)["email"])
					assert.NotEmpty(t, res["kontak"].(map[string]any)["keterangan"])
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

func AkuntansiDeleteTransaksi(t *testing.T) {
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
			id:           idTransaksi,
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
			id:           idTransaksi + "123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			id:           idTransaksi,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[4],
			id:           idTransaksi,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idTransaksi,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/akuntansi/transaksi/"+tt.id, nil, &tt.token)
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

func AkuntansiGetHistoryTransaksi(t *testing.T) {
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
			queryBody:    "?start_date=2023-10-01&end_date=2024-12-31",
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: start_date & end_date wajib diisi",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"start date wajib diisi", "end date wajib diisi"},
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/transaksi"+tt.queryBody, nil, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)

			var res []map[string]interface{}
			if tt.name == "sukses" {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.Greater(t, len(res), 0)
				for _, v := range res {
					assert.NotEmpty(t, v)
					assert.NotEmpty(t, v["id"])
					assert.NotEmpty(t, v["created_at"])
					assert.NotEmpty(t, v["keterangan"])
					assert.NotEmpty(t, v["total"])
					assert.NotEmpty(t, v["tanggal"])
				}
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

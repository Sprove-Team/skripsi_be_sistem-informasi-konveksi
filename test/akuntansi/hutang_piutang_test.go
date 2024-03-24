//go:build test_exclude

package test_akuntansi

import (
	"testing"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_akuntansi_hp "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
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
		// {
		// 	name:  "err: kategori akun harus berupa salah satu dari [ASET,KEWAJIBAN,MODAL,PENDAPATAN,BEBAN]",
		// 	token: tokens[entity.RolesById[1]],
		// 	payload: req_akuntansi_hp.Create{
		// 		Nama:         "kelompok err kategori",
		// 		Kode:         "212312",
		// 		KategoriAkun: "asdf",
		// 	},
		// 	expectedCode: 400,
		// 	expectedBody: test.Response{
		// 		Status:         fiber.ErrBadRequest.Message,
		// 		Code:           400,
		// 		ErrorsMessages: []string{"kategori akun harus berupa salah satu dari [ASET,KEWAJIBAN,MODAL,PENDAPATAN,BEBAN]"},
		// 	},
		// },
		// {
		// 	name:  "err: wajib diisi",
		// 	token: tokens[entity.RolesById[1]],
		// 	payload: req_akuntansi_hp.Create{
		// 		Nama:         "",
		// 		Kode:         "",
		// 		KategoriAkun: "",
		// 	},
		// 	expectedCode: 400,
		// 	expectedBody: test.Response{
		// 		Status:         fiber.ErrBadRequest.Message,
		// 		Code:           400,
		// 		ErrorsMessages: []string{"nama wajib diisi", "kode wajib diisi", "kategori akun wajib diisi"},
		// 	},
		// },
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

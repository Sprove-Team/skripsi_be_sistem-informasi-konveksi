package test_akuntansi

import (
	"testing"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_akuntansi_transaksi "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/transaksi"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
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
				BuktiPembayaran: []string{"http://bukti-pembayaran.jpg"},
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
				BuktiPembayaran: []string{"http://bukti-pembayaran2.jpg"},
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
				BuktiPembayaran: []string{"http://bukti-pembayaran2.jpg"},
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
				BuktiPembayaran: []string{"http://bukti-pembayaran2.jpg"},
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
			name:  "err: wajib diisi akun id pada ayat jurnal",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Create{
				BuktiPembayaran: []string{"http://bukti-pembayaran2.jpg"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "Membayar beban listrik dan air sebesar Rp. 10.000,- with kontak",
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
				BuktiPembayaran: []string{"http://bukti-pembayaran3.jpg"},
				Tanggal:         "2023-10-28",
				Keterangan:      "Membayar beban listrik dan air sebesar Rp. 10.000,3",
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
				BuktiPembayaran: []string{"http://bukti-pembayaran3.jpg"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "Membayar beban listrik dan air sebesar Rp. 10.000,3",
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
				BuktiPembayaran: []string{"http://bukti-pembayaran3.jpg"},
				Tanggal:         "2023-10-28T14:17:03.723Z",
				Keterangan:      "Membayar beban listrik dan air sebesar Rp. 10.000,3",
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

func AkuntansiUpdateTransaksi(t *testing.T) {
	transaksi := new(entity.Transaksi)
	err := dbt.Model(transaksi).Order("id DESC").First(transaksi).Error
	if err != nil {
		helper.LogsError(err)
		return
	}

	idTransaksi = transaksi.ID

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
				BuktiPembayaran: []string{"http://bukti-update.jpg"},
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
			name:  "err: tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_akuntansi_transaksi.Update{
				ID:              "01HM4B8QBH7MWAVAYP10WN6PKA",
				BuktiPembayaran: []string{"http://bukti-update.jpg"},
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
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		// {
		// 	name: "err: ulid tidak valid",
		// token: tokens[entity.RolesById[1]],
		// 	payload: req_akuntansi_akun.Update{
		// 		ID:             idAkun + "123",
		// 		Nama:           "kelompok ulid tidak valid",
		// 		Kode:           "15",
		// 		KelompokAkunID: static_data.DataKelompokAkun[12].ID,
		// 		Deskripsi:      "des update ulid tidak valid",
		// 		SaldoNormal:    "KREDIT",
		// 	},
		// 	expectedCode: 400,
		// 	expectedBody: test.Response{
		// 		Status:         fiber.ErrBadRequest.Message,
		// 		Code:           400,
		// 		ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
		// 	},
		// },
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
				assert.Equal(t, tt.expectedBody, body)
			}
		})
	}
}

// func AkuntansiGetAllAkun(t *testing.T) {
// 	tests := []struct {
// 		name         string
// token string
// 		queryBody    string
// 		expectedBody test.Response
// 		expectedCode int
// 	}{
// 		{
// 			name:         "sukses",
// token: tokens[entity.RolesById[1]],
// 			expectedCode: 200,
// 			expectedBody: test.Response{
// 				Status: message.OK,
// 				Code:   200,
// 			},
// 		},
// 		{ // same data with idAkun2
// 			name:         "sukses with filter",
// token: tokens[entity.RolesById[1]],
// 			queryBody:    "?nama=akun+update+test&kode=1123499",
// 			expectedCode: 200,
// 			expectedBody: test.Response{
// 				Status: message.OK,
// 				Code:   200,
// 			},
// 		},
// 		{
// 			name:         "sukses limit 1",
// token: tokens[entity.RolesById[1]],
// 			queryBody:    "?limit=1",
// 			expectedCode: 200,
// 			expectedBody: test.Response{
// 				Status: message.OK,
// 				Code:   200,
// 			},
// 		},
// 		{
// 			name:         "sukses with next",
// token: tokens[entity.RolesById[1]],
// 			queryBody:    "?next=" + idAkun,
// 			expectedCode: 200,
// 			expectedBody: test.Response{
// 				Status: message.OK,
// 				Code:   200,
// 			},
// 		},
// 		{
// 			name:         "err: ulid tidak valid",
// token: tokens[entity.RolesById[1]],
// 			expectedCode: 400,
// 			queryBody:    "?next=01HQVTTJ1S2606JGTYYZ5NDKNR123",
// 			expectedBody: test.Response{
// 				Status:         fiber.ErrBadRequest.Message,
// 				Code:           400,
// 				ErrorsMessages: []string{"next tidak berupa ulid yang valid"},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/akun"+tt.queryBody, nil, &token)
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expectedCode, code)

// 			var res []map[string]interface{}
// 			if strings.Contains(tt.name, "sukses") {
// 				err = mapstructure.Decode(body.Data, &res)
// 				assert.NoError(t, err)

// 				assert.Greater(t, len(res), 0)
// 				assert.NotEmpty(t, res[0])
// 				assert.NotEmpty(t, res[0]["id"])
// 				assert.NotEmpty(t, res[0]["kode"])
// 				assert.NotEmpty(t, res[0]["nama"])
// 				assert.NotEmpty(t, res[0]["created_at"])
// 				assert.NotEmpty(t, res[0]["saldo_normal"])
// 				assert.NotEmpty(t, res[0]["deskripsi"])
// 				assert.NotEmpty(t, res[0]["kelompok_akun"])
// 				assert.NotEmpty(t, res[0]["kelompok_akun"].(map[string]any)["id"])
// 				assert.NotEmpty(t, res[0]["kelompok_akun"].(map[string]any)["kode"])
// 				assert.NotEmpty(t, res[0]["kelompok_akun"].(map[string]any)["nama"])
// 				assert.Equal(t, tt.expectedBody.Status, body.Status)
// 				switch tt.name {
// 				case "sukses with filter":
// 					v, err := url.ParseQuery(tt.queryBody[1:])
// 					assert.NoError(t, err)
// 					assert.NotEmpty(t, res)
// 					assert.Equal(t, res[0]["nama"], v.Get("nama"))
// 					assert.Equal(t, res[0]["kode"], v.Get("kode"))
// 				case "sukses limit 1":
// 					assert.Len(t, res, 1)
// 				case "sukses with next":
// 					assert.NotEmpty(t, res[0])
// 					assert.NotEqual(t, idAkun, res[0]["id"])
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

// func AkuntansiGetAkun(t *testing.T) {
// 	tests := []struct {
// 		id           string
// 		name         string
// token string
// 		expectedBody test.Response
// 		expectedCode int
// 	}{
// 		{
// 			name:         "sukses",
// token: tokens[entity.RolesById[1]],
// 			id:           idAkun,
// 			expectedCode: 200,
// 			expectedBody: test.Response{
// 				Status: message.OK,
// 				Code:   200,
// 			},
// 		},
// 		{
// 			name:         "err: tidak ditemukan",
// token: tokens[entity.RolesById[1]],
// 			id:           "01HM4B8QBH7MWAVAYP10WN6PKA",
// 			expectedCode: 404,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrNotFound.Message,
// 				Code:   404,
// 			},
// 		},
// 		{
// 			name:         "err: ulid tidak valid",
// token: tokens[entity.RolesById[1]],
// 			id:           "01HQVTTJ1S2606JGTYYZ5NDKNR123",
// 			expectedCode: 400,
// 			expectedBody: test.Response{
// 				Status:         fiber.ErrBadRequest.Message,
// 				Code:           400,
// 				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/akun/"+tt.id, nil, &token)
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expectedCode, code)

// 			var res map[string]any
// 			if tt.name == "sukses" {
// 				err = mapstructure.Decode(body.Data, &res)
// 				assert.NoError(t, err)
// 				assert.Greater(t, len(res), 0)
// 				assert.NotEmpty(t, res)
// 				assert.NotEmpty(t, res["id"])
// 				assert.NotEmpty(t, res["created_at"])
// 				assert.NotEmpty(t, res["nama"])
// 				assert.NotEmpty(t, res["kode"])
// 				assert.NotEmpty(t, res["saldo_normal"])
// 				assert.NotEmpty(t, res["deskripsi"])
// 				assert.NotEmpty(t, res["kelompok_akun"])
// 				assert.NotEmpty(t, res["kelompok_akun"].(map[string]any)["id"])
// 				assert.NotEmpty(t, res["kelompok_akun"].(map[string]any)["kode"])
// 				assert.NotEmpty(t, res["kelompok_akun"].(map[string]any)["nama"])
// 				assert.Equal(t, tt.expectedBody.Status, body.Status)
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

// func AkuntansiDeleteAkun(t *testing.T) {
// 	tests := []struct {
// 		name         string
// token string
// 		id           string
// 		expectedBody test.Response
// 		expectedCode int
// 	}{
// 		{
// 			name:         "sukses",
// token: tokens[entity.RolesById[1]],
// 			id:           idAkun,
// 			expectedCode: 200,
// 			expectedBody: test.Response{
// 				Status: message.OK,
// 				Code:   200,
// 			},
// 		},
// 		{
// 			name:         "err: tidak ditemukan",
// token: tokens[entity.RolesById[1]],
// 			id:           "01HM4B8QBH7MWAVAYP10WN6PKA",
// 			expectedCode: 404,
// 			expectedBody: test.Response{
// 				Status: fiber.ErrNotFound.Message,
// 				Code:   404,
// 			},
// 		},
// 		{
// 			name:         "err: can't delete default data",
// token: tokens[entity.RolesById[1]],
// 			id:           static_data.DataAkun[0][0].ID,
// 			expectedCode: 400,
// 			expectedBody: test.Response{
// 				Status:         fiber.ErrBadRequest.Message,
// 				Code:           400,
// 				ErrorsMessages: []string{message.CantModifiedDefaultData},
// 			},
// 		},
// 		{
// 			name:         "err: ulid tidak valid",
// token: tokens[entity.RolesById[1]],
// 			id:           idAkun + "123",
// 			expectedCode: 400,
// 			expectedBody: test.Response{
// 				Status:         fiber.ErrBadRequest.Message,
// 				Code:           400,
// 				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/akuntansi/akun/"+tt.id, nil, &token)
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

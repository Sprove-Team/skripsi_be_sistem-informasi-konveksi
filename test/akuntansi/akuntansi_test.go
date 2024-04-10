package test_akuntansi

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/be-sistem-informasi-konveksi/common/message"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func AkuntansiGetJU(t *testing.T) {
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
			queryBody:    fmt.Sprintf("?start_date=%s&end_date=%s", ttTrStartDateAkuntansi, ttTrEndDateAkuntansi),
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses download excel",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 200,
			queryBody:    fmt.Sprintf("?start_date=%s&end_date=%s&download=1", ttTrStartDateAkuntansi, ttTrEndDateAkuntansi),
		},
		{
			name:         "err: format start date & end date",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 400,
			queryBody:    "?start_date=2023-01-01T22:17:03.723+08:00&end_date=2024-01-01T22:17:03.723+08:00",
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"start date harus berformat Tahun-Bulan-Tanggal", "end date harus berformat Tahun-Bulan-Tanggal"},
			},
		},
		{
			name:         "err: wajib diisi",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"start date wajib diisi", "end date wajib diisi"},
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
			if tt.name == "sukses download excel" {
				code, err := test.GetAttachTestRequestResponse(app, "GET", "/api/v1/akuntansi/jurnal_umum"+tt.queryBody, nil, "ju.xlsx", &tt.token)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, code)
				info, err := os.Stat("ju.xlsx")
				assert.False(t, os.IsNotExist(err))
				assert.Greater(t, info.Size(), int64(0))
				err = os.Remove("ju.xlsx")
				assert.NoError(t, err)
				return
			}
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/jurnal_umum"+tt.queryBody, nil, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			var res map[string]interface{}
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res["total_kredit"])
				assert.NotEmpty(t, res["total_debit"])
				assert.NotEmpty(t, res["transaksi"])
				tr, ok := res["transaksi"].([]any)
				assert.True(t, ok)
				if len(tr) <= 0 {
					return
				}

				var totalDebit float64
				var totalKredit float64
				for _, v := range tr {
					v2 := v.(map[string]any)
					assert.NotEmpty(t, v2["tanggal"])
					assert.NotEmpty(t, v2["transaksi_id"])
					assert.NotEmpty(t, v2["keterangan"])
					assert.NotEmpty(t, v2["ayat_jurnal"])
					ays, ok := v2["ayat_jurnal"].([]any)
					assert.True(t, ok)
					if len(ays) <= 0 {
						return
					}
					assert.GreaterOrEqual(t, len(ays), 2)
					for _, ay := range ays {
						ay2 := ay.(map[string]any)
						assert.NotEmpty(t, ay2["id"])
						assert.NotEmpty(t, ay2["akun_id"])
						assert.NotEmpty(t, ay2["kode_akun"])
						assert.NotEmpty(t, ay2["nama_akun"])
						if dat, ok := ay2["debit"].(float64); ok && dat != 0 {
							totalDebit += dat
							assert.Greater(t, dat, float64(0))
						}
						if dat, ok := ay2["kredit"].(float64); ok && dat != 0 {
							totalKredit += dat
							assert.Greater(t, dat, float64(0))
						}
					}
				}

				assert.Equal(t, totalDebit, float64(45050000))
				assert.Equal(t, totalKredit, float64(45050000))
				assert.Equal(t, totalDebit, totalKredit)
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
}

func AkuntansiGetBB(t *testing.T) {
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
			queryBody:    "?start_date=2023-01-20&end_date=2024-12-30",
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses dengan filter",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 200,
			queryBody:    fmt.Sprintf("?start_date=%s&end_date=%s&akun_id=01HP7DVBGTC06PXWT6FD66VERN%s01HP7DVBGTC06PXWT6FF89WRAB", ttTrStartDateAkuntansi, ttTrEndDateAkuntansi, "%2C"),
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses download excel",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 200,
			queryBody:    fmt.Sprintf("?start_date=%s&end_date=%s&akun_id=01HP7DVBGTC06PXWT6FD66VERN%s01HP7DVBGTC06PXWT6FF89WRAB&download=1", ttTrStartDateAkuntansi, ttTrEndDateAkuntansi, "%2C"),
		},
		{
			name:         "err: format start date & end date",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 400,
			queryBody:    "?start_date=2023-01-01T22:17:03.723+08:00&end_date=2024-01-01T22:17:03.723+08:00",
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"start date harus berformat Tahun-Bulan-Tanggal", "end date harus berformat Tahun-Bulan-Tanggal"},
			},
		},
		{
			name:         "err: wajib diisi",
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
			name:         "authorization " + entity.RolesById[2] + " passed",
			token:        tokens[entity.RolesById[2]],
			expectedCode: 400,
			expectedBody: test.Response{
				Status: fiber.ErrBadRequest.Message,
				Code:   400,
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
			if tt.name == "sukses download excel" {
				code, err := test.GetAttachTestRequestResponse(app, "GET", "/api/v1/akuntansi/buku_besar"+tt.queryBody, nil, "bb.xlsx", &tt.token)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, code)
				info, err := os.Stat("bb.xlsx")
				assert.False(t, os.IsNotExist(err))
				assert.Greater(t, info.Size(), int64(0))
				err = os.Remove("bb.xlsx")
				assert.NoError(t, err)
				return
			}
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/buku_besar"+tt.queryBody, nil, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			var res []map[string]any
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				if len(res) <= 0 {
					return
				}

				for _, v := range res {
					assert.NotEmpty(t, v["kode_akun"])
					assert.NotEmpty(t, v["nama_akun"])
					assert.NotEmpty(t, v["saldo_normal"])
					var totalDebit, totalKredit, totalSaldo float64
					if ays, ok := v["ayat_jurnal"].([]any); ok && len(ays) > 0 {
						for _, ay := range ays {
							ay2 := ay.(map[string]any)
							assert.NotEmpty(t, ay2["tanggal"])
							assert.NotEmpty(t, ay2["keterangan"])
							if ay2["keterangan"] != "saldo awal" {
								assert.NotEmpty(t, ay2["transaksi_id"])
							}
							totalDebit += ay2["debit"].(float64)
							totalKredit += ay2["kredit"].(float64)
							if v["saldo_normal"] == "DEBIT" {
								totalSaldo = totalDebit - totalKredit
							} else {
								totalSaldo = totalKredit - totalDebit
							}
							if dat, ok := ay2["saldo"].(float64); ok {
								assert.Equal(t, totalSaldo, dat)
							}
						}
					}
					if v["total_debit"].(float64) != 0 {
						assert.Equal(t, totalDebit, v["total_debit"].(float64))
					}
					if v["total_kredit"].(float64) != 0 {
						assert.Equal(t, totalKredit, v["total_kredit"].(float64))
					}
					if v["total_saldo"].(float64) != 0 {
						assert.Equal(t, totalSaldo, v["total_saldo"].(float64))
					}
					if tt.name == "sukses dengan filter" {
						akunNama := []string{"kas", "piutang usaha"} // fit with the query akun
						assert.Contains(t, akunNama, v["nama_akun"])
					}

					// check akuntansi sesuai dengan soal
					if tt.name == "sukses" {
						switch v["kode_akun"].(string) {
						case "111":
							assert.Equal(t, float64(34100000), totalSaldo)
						case "121":
							assert.Equal(t, float64(1050000), totalSaldo)
						case "131":
							assert.Equal(t, float64(350000), totalSaldo)
						case "141":
							assert.Equal(t, float64(5500000), totalSaldo)
						case "271":
							assert.Equal(t, float64(5500000), totalSaldo)
						case "3101":
							assert.Equal(t, float64(35000000), totalSaldo)
						case "5141":
							assert.Equal(t, float64(1000000), totalSaldo)
						case "5144":
							assert.Equal(t, float64(250000), totalSaldo)
						case "5142":
							assert.Equal(t, float64(600000), totalSaldo)
						}
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
					if strings.Contains(tt.name, "passed") {
						return
					}
					assert.Equal(t, tt.expectedBody, body)
				}
			}
		})
	}
}

func AkuntansiGetNC(t *testing.T) {
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
			queryBody:    fmt.Sprintf("?date=%s", ttTrYearMonthAkuntansi),
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses download excel",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 200,
			queryBody:    fmt.Sprintf("?date=%s&download=?", ttTrYearMonthAkuntansi),
		},
		{
			name:         "err: format date",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 400,
			queryBody:    "?date=2023-01-01T22:17:03.723+08:00",
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"date harus berformat Tahun-Bulan"},
			},
		},
		{
			name:         "err: wajib diisi",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"date wajib diisi"},
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
			if tt.name == "sukses download excel" {
				code, err := test.GetAttachTestRequestResponse(app, "GET", "/api/v1/akuntansi/neraca_saldo"+tt.queryBody, nil, "nc.xlsx", &tt.token)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, code)
				info, err := os.Stat("nc.xlsx")
				assert.False(t, os.IsNotExist(err))
				assert.Greater(t, info.Size(), int64(0))
				err = os.Remove("nc.xlsx")
				assert.NoError(t, err)
				return
			}
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/neraca_saldo"+tt.queryBody, nil, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			var res map[string]any
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				if len(res) <= 0 {
					return
				}

				var totalDebit, totalKredit float64
				if ays, ok := res["saldo_akun"].([]any); ok && len(ays) > 0 {
					for _, ay := range ays {
						ay2 := ay.(map[string]any)
						assert.NotEmpty(t, ay2["kode_akun"])
						assert.NotEmpty(t, ay2["nama_akun"])
						totalDebit += ay2["saldo_debit"].(float64)
						totalKredit += ay2["saldo_kredit"].(float64)
					}
				}
				assert.Equal(t, totalDebit, res["total_debit"].(float64))
				assert.Equal(t, totalKredit, res["total_kredit"].(float64))
				assert.Equal(t, res["total_debit"], res["total_kredit"])
				// check sesuai dengan soal akuntansi
				assert.Equal(t, float64(42850000), totalDebit)
				assert.Equal(t, float64(42850000), totalKredit)

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
}

func AkuntansiGetLB(t *testing.T) {
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
			queryBody:    "?start_date=2023-01-20&end_date=2024-12-30",
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses download excel",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 200,
			queryBody:    "?start_date=2023-01-20&end_date=2024-12-30&download=1",
		},
		{
			name:         "err: format start date & end date",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 400,
			queryBody:    "?start_date=2023-01-01T22:17:03.723+08:00&end_date=2024-01-01T22:17:03.723+08:00",
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"start date harus berformat Tahun-Bulan-Tanggal", "end date harus berformat Tahun-Bulan-Tanggal"},
			},
		},
		{
			name:         "err: wajib diisi",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"start date wajib diisi", "end date wajib diisi"},
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
			if tt.name == "sukses download excel" {
				code, err := test.GetAttachTestRequestResponse(app, "GET", "/api/v1/akuntansi/laba_rugi"+tt.queryBody, nil, "lb.xlsx", &tt.token)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, code)
				info, err := os.Stat("lb.xlsx")
				assert.False(t, os.IsNotExist(err))
				assert.Greater(t, info.Size(), int64(0))
				err = os.Remove("lb.xlsx")
				assert.NoError(t, err)
				return
			}
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/akuntansi/laba_rugi"+tt.queryBody, nil, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			var res []map[string]any
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				if len(res) <= 0 {
					return
				}

				var labaRugi float64
				for _, v := range res {
					assert.NotEmpty(t, v["kategori_akun"])
					assert.NotEmpty(t, v["akun"])
					var total float64
					if ays, ok := v["akun"].([]any); ok && len(ays) > 0 {
						for _, ay := range ays {
							ay2 := ay.(map[string]any)
							assert.NotEmpty(t, ay2["nama_akun"])
							assert.NotEmpty(t, ay2["kode_akun"])
							total += ay2["total"].(float64)
						}
					}
					if v["total"].(float64) != 0 {
						assert.Equal(t, total, v["total"].(float64))
					}
					if v["kategori_akun"] == "PENDAPATAN" {
						labaRugi += total
					} else {
						labaRugi -= total
					}
				}
				assert.Equal(t, float64(500000), labaRugi)
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
}

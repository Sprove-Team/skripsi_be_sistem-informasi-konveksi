package test_tugas

import (
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_tugas "github.com/be-sistem-informasi-konveksi/common/request/tugas"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TugasCreate(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		payload      req_tugas.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_tugas.Create{
				InvoiceID:       invoiceId[0],
				JenisSpvID:      spvId,
				TanggalDeadline: "2024-10-20",
				UserID:          userId,
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name:  "err: id user bukan termasuk jenis dari spv",
			token: tokens[entity.RolesById[1]],
			payload: req_tugas.Create{
				InvoiceID:       invoiceId[0],
				JenisSpvID:      spvId,
				TanggalDeadline: "2024-10-20",
				UserID:          []string{static_data.DefaultUsers[4].ID},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.UserNotFoundOrNotSpv},
			},
		},
		{
			name:  "err: id user tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_tugas.Create{
				InvoiceID:       invoiceId[0],
				JenisSpvID:      spvId,
				TanggalDeadline: "2024-10-20",
				UserID:          []string{test.UlidPkg.MakeUlid().String()},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.UserNotFoundOrNotSpv},
			},
		},
		{
			name:         "err: wajib diisi",
			token:        tokens[entity.RolesById[1]],
			payload:      req_tugas.Create{},
			expectedCode: 400,
			expectedBody: test.Response{
				Status: fiber.ErrBadRequest.Message,
				Code:   400,
				ErrorsMessages: []string{
					"invoice id wajib diisi",
					"jenis spv id wajib diisi",
					"tanggal deadline wajib diisi",
					"user id wajib diisi",
				},
			},
		},
		{
			name:  "err: data user id harus lebih dari 0",
			token: tokens[entity.RolesById[1]],
			payload: req_tugas.Create{
				UserID: []string{},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"user id harus berisi lebih dari 0 item"},
			},
		},
		{
			name:  "err: format tidak valid untuk semua field",
			token: tokens[entity.RolesById[1]],
			payload: req_tugas.Create{
				InvoiceID:       "asdf",
				JenisSpvID:      "asdf",
				UserID:          []string{"asdf"},
				TanggalDeadline: "20-08",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"invoice id tidak berupa ulid yang valid", "jenis spv id tidak berupa ulid yang valid", "tanggal deadline harus berformat Tahun-Bulan-Tanggal", "user id tidak berupa ulid yang valid"},
			},
		},
		{
			name: "err: authorization " + entity.RolesById[2],
			payload: req_tugas.Create{
				InvoiceID:       "asdf",
				JenisSpvID:      "asdf",
				UserID:          []string{"asdf"},
				TanggalDeadline: "20-08",
			},
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name: "err: authorization " + entity.RolesById[3],
			payload: req_tugas.Create{
				InvoiceID:       "asdf",
				JenisSpvID:      "asdf",
				UserID:          []string{"asdf"},
				TanggalDeadline: "20-08",
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
			payload: req_tugas.Create{
				InvoiceID:       "asdf",
				JenisSpvID:      "asdf",
				UserID:          []string{"asdf"},
				TanggalDeadline: "20-08",
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
			payload: req_tugas.Create{
				InvoiceID:       "asdf",
				JenisSpvID:      "asdf",
				UserID:          []string{"asdf"},
				TanggalDeadline: "20-08",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/tugas", tt.payload, &tt.token)
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

var idTugas string

func TugasUpdate(t *testing.T) {
	tugas := new(entity.Tugas)
	if err := dbt.Select("id").First(tugas).Error; err != nil {
		panic(helper.LogsError(err))
	}
	idTugas = tugas.ID

	tests := []struct {
		name         string
		token        string
		payload      req_tugas.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_tugas.Update{
				ID:              idTugas,
				JenisSpvID:      spvId,
				TanggalDeadline: "2024-04-20",
				UserID:          userId,
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:  "err: id user bukan termasuk jenis dari spv",
			token: tokens[entity.RolesById[1]],
			payload: req_tugas.Update{
				ID:              idTugas,
				JenisSpvID:      spvId,
				TanggalDeadline: "2024-04-20",
				UserID:          []string{static_data.DefaultUsers[4].ID},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.UserNotFoundOrNotSpv},
			},
		},
		{
			name:  "err: id user tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_tugas.Update{
				ID:              idTugas,
				JenisSpvID:      spvId,
				TanggalDeadline: "2024-04-20",
				UserID:          []string{test.UlidPkg.MakeUlid().String()},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.UserNotFoundOrNotSpv},
			},
		},
		{
			name:  "err: data user id harus lebih dari 0",
			token: tokens[entity.RolesById[1]],
			payload: req_tugas.Update{
				ID:     idTugas,
				UserID: []string{},
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"user id harus berisi lebih dari 0 item"},
			},
		},
		{
			name:  "err: format tidak valid untuk semua field",
			token: tokens[entity.RolesById[1]],
			payload: req_tugas.Update{
				ID:              "123",
				JenisSpvID:      "asdf",
				UserID:          []string{"asdf"},
				TanggalDeadline: "20-08",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid", "jenis spv id tidak berupa ulid yang valid", "tanggal deadline harus berformat Tahun-Bulan-Tanggal", "user id tidak berupa ulid yang valid"},
			},
		},
		{
			name: "err: authorization " + entity.RolesById[2],
			payload: req_tugas.Update{
				ID:              idTugas,
				JenisSpvID:      "asdf",
				UserID:          []string{"asdf"},
				TanggalDeadline: "20-08",
			},
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name: "err: authorization " + entity.RolesById[3],
			payload: req_tugas.Update{
				ID:              idTugas,
				JenisSpvID:      "asdf",
				UserID:          []string{"asdf"},
				TanggalDeadline: "20-08",
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
			payload: req_tugas.Update{
				ID:              idTugas,
				JenisSpvID:      "asdf",
				UserID:          []string{"asdf"},
				TanggalDeadline: "20-08",
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
			payload: req_tugas.Update{
				ID:              idTugas,
				JenisSpvID:      "asdf",
				UserID:          []string{"asdf"},
				TanggalDeadline: "20-08",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/tugas/"+tt.payload.ID, tt.payload, &tt.token)
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

func TugasGetAll(t *testing.T) {
	tests := []struct {
		name         string
		queryBody    string
		token        string
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
			name:         "sukses dengan filter",
			queryBody:    "?bulan=4&tahun=2024",
			token:        tokens[entity.RolesById[1]],
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: format semua field query body",
			token:        tokens[entity.RolesById[1]],
			queryBody:    "?bulan=100&tahun=1999",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"tahun harus 2006 atau lebih besar", "bulan harus 12 atau kurang"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[2],
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
			name:         "authorization " + entity.RolesById[4] + " passed",
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[5] + " passed",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/tugas"+tt.queryBody, nil, &tt.token)
			assert.NoError(t, err)
			if strings.Contains(tt.name, "passed") {
				assert.NotEqual(t, tt.expectedCode, code)
				assert.NotEqual(t, tt.expectedBody.Code, body.Code)
				assert.NotEqual(t, tt.expectedBody.Status, body.Status)
				return
			}
			assert.Equal(t, tt.expectedCode, code)
			var res []map[string]any
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.Greater(t, len(res), 0)
				for _, vres := range res {
					assert.NotEmpty(t, vres)
					assert.NotEmpty(t, vres["id"])
					assert.NotEmpty(t, vres["created_at"])
					assert.NotEmpty(t, vres["jenis_spv"])
					jenisSpv, ok := vres["jenis_spv"].(map[string]any)
					assert.True(t, ok)
					assert.NotEmpty(t, jenisSpv["id"])
					assert.NotEmpty(t, jenisSpv["created_at"])
					assert.NotEmpty(t, jenisSpv["nama"])

					assert.NotEmpty(t, vres["tanggal_deadline"])
					assert.NotEmpty(t, vres["users"])
					for _, v := range vres["users"].([]any) {
						v2 := v.(map[string]any)
						assert.NotEmpty(t, v2["id"])
						assert.NotEmpty(t, v2["nama"])
						assert.NotEmpty(t, v2["role"])
						assert.NotEmpty(t, v2["username"])
					}
					if tt.name == "sukses dengan filter" {
						parse, err := time.Parse(time.RFC3339, vres["tanggal_deadline"].(string))
						assert.NoError(t, err)
						v, err := url.ParseQuery(tt.queryBody[1:])
						assert.NoError(t, err)
						assert.Equal(t, v.Get("tahun"), strconv.Itoa(parse.Year()))
						assert.Equal(t, v.Get("bulan"), strconv.Itoa(int(parse.Month())))
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

func TugasGetByInvoiceId(t *testing.T) {
	tests := []struct {
		name         string
		invoiceId    string
		token        string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			token:        tokens[entity.RolesById[1]],
			invoiceId:    invoiceId[0],
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[2],
			invoiceId:    invoiceId[0],
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			invoiceId:    invoiceId[0],
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[4] + " passed",
			invoiceId:    invoiceId[0],
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[5] + " passed",
			invoiceId:    invoiceId[0],
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/tugas/invoice/"+tt.invoiceId, nil, &tt.token)
			assert.NoError(t, err)
			if strings.Contains(tt.name, "passed") {
				assert.NotEqual(t, tt.expectedCode, code)
				assert.NotEqual(t, tt.expectedBody.Code, body.Code)
				assert.NotEqual(t, tt.expectedBody.Status, body.Status)
				return
			}
			assert.Equal(t, tt.expectedCode, code)
			var res []map[string]any
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.Greater(t, len(res), 0)
				if len(res) <= 0 {
					panic("res les than 0")
				}
				assert.NotEmpty(t, res[0])
				assert.NotEmpty(t, res[0]["id"])
				assert.NotEmpty(t, res[0]["created_at"])
				assert.NotEmpty(t, res[0]["jenis_spv"])
				jenisSpv, ok := res[0]["jenis_spv"].(map[string]any)
				assert.True(t, ok)
				assert.NotEmpty(t, jenisSpv["id"])
				assert.NotEmpty(t, jenisSpv["created_at"])
				assert.NotEmpty(t, jenisSpv["nama"])

				assert.NotEmpty(t, res[0]["tanggal_deadline"])
				assert.NotEmpty(t, res[0]["users"])
				for _, v := range res[0]["users"].([]any) {
					v2 := v.(map[string]any)
					assert.NotEmpty(t, v2["id"])
					assert.NotEmpty(t, v2["nama"])
					assert.NotEmpty(t, v2["role"])
					assert.NotEmpty(t, v2["username"])
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

func TugasGet(t *testing.T) {
	subTugas := entity.SubTugas{
		Base: entity.Base{
			ID: test.UlidPkg.MakeUlid().String(),
		},
		Nama:      "beli bordir di a",
		Deskripsi: "beli bordir toko a di alamat ini",
		Status:    "BELUM_DIKERJAKAN",
		TugasID:   idTugas,
	}
	if err := dbt.Create(&subTugas).Error; err != nil {
		panic(helper.LogsError(err))
	}

	tests := []struct {
		name         string
		id           string
		token        string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			token:        tokens[entity.RolesById[1]],
			id:           idTugas,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[2],
			id:           idTugas,
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			id:           idTugas,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[4] + " passed",
			id:           idTugas,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[5] + " passed",
			id:           idTugas,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/tugas/"+tt.id, nil, &tt.token)
			assert.NoError(t, err)
			if strings.Contains(tt.name, "passed") {
				assert.NotEqual(t, tt.expectedCode, code)
				assert.NotEqual(t, tt.expectedBody.Code, body.Code)
				assert.NotEqual(t, tt.expectedBody.Status, body.Status)
				return
			}
			assert.Equal(t, tt.expectedCode, code)
			var res map[string]any
			if strings.Contains(tt.name, "sukses") {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.Greater(t, len(res), 0)
				if len(res) <= 0 {
					panic("res les than 0")
				}
				assert.NotEmpty(t, res)
				assert.NotEmpty(t, res["id"])
				assert.NotEmpty(t, res["created_at"])
				assert.NotEmpty(t, res["jenis_spv"])
				jenisSpv, ok := res["jenis_spv"].(map[string]any)
				assert.True(t, ok)
				assert.NotEmpty(t, jenisSpv["id"])
				assert.NotEmpty(t, jenisSpv["created_at"])
				assert.NotEmpty(t, jenisSpv["nama"])

				invoice, ok := res["invoice"].(map[string]any)
				assert.True(t, ok)
				assert.NotEmpty(t, invoice["id"])
				assert.NotEmpty(t, invoice["created_at"])
				assert.NotEmpty(t, invoice["status_produksi"])
				assert.NotEmpty(t, invoice["nomor_referensi"])
				assert.NotEmpty(t, invoice["total_qty"])
				assert.NotEmpty(t, invoice["total_harga"])
				assert.NotEmpty(t, invoice["keterangan"])
				assert.NotEmpty(t, invoice["tanggal_kirim"])

				assert.NotEmpty(t, res["tanggal_deadline"])
				assert.NotEmpty(t, res["users"])
				for _, v := range res["users"].([]any) {
					v2 := v.(map[string]any)
					assert.NotEmpty(t, v2["id"])
					assert.NotEmpty(t, v2["nama"])
					assert.NotEmpty(t, v2["role"])
					assert.NotEmpty(t, v2["username"])
				}

				assert.NotEmpty(t, res["sub_tugas"])
				for _, v := range res["sub_tugas"].([]any) {
					v2 := v.(map[string]any)
					assert.NotEmpty(t, v2["id"])
					assert.NotEmpty(t, v2["created_at"])
					assert.NotEmpty(t, v2["nama"])
					assert.NotEmpty(t, v2["status"])
					assert.NotEmpty(t, v2["deskripsi"])
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

func TugasDelete(t *testing.T) {
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
			id:           idTugas,
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
			id:           idTugas + "123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[2],
			id:           idTugas,
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			id:           idTugas,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[4] + " passed",
			id:           idTugas,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idTugas,
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
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/tugas/"+tt.id, nil, &tt.token)
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

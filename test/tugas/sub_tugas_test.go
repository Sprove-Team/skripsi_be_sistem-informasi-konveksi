package test_tugas

import (
	"strings"
	"testing"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_sub_tugas "github.com/be-sistem-informasi-konveksi/common/request/tugas/sub_tugas"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func SubTugasCreate(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		payload      req_sub_tugas.CreateByTugasId
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_sub_tugas.CreateByTugasId{
				TugasID:   idTugas,
				Nama:      "tugas testing",
				Status:    "BELUM_DIKERJAKAN",
				Deskripsi: "test",
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name:  "err: tidak ditemukan",
			token: tokens[entity.RolesById[1]],
			payload: req_sub_tugas.CreateByTugasId{
				TugasID:   spvId,
				Nama:      "tugas testing",
				Status:    "BELUM_DIKERJAKAN",
				Deskripsi: "test",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.TugasNotFound},
			},
		},
		{
			name:  "err: wajib diisi",
			token: tokens[entity.RolesById[1]],
			payload: req_sub_tugas.CreateByTugasId{
				TugasID: idTugas,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"nama wajib diisi"},
			},
		},
		{
			name:  "err: format semua field",
			token: tokens[entity.RolesById[1]],
			payload: req_sub_tugas.CreateByTugasId{
				TugasID:   idTugas + "123",
				Nama:      "asdf",
				Status:    "asdf",
				Deskripsi: "123",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"tugas id tidak berupa ulid yang valid", "status harus berupa salah satu dari [BELUM_DIKERJAKAN,DIPROSES,SELESAI]"},
			},
		},
		{
			name: "err: authorization " + entity.RolesById[2],
			payload: req_sub_tugas.CreateByTugasId{
				TugasID: idTugas,
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
			payload: req_sub_tugas.CreateByTugasId{
				TugasID: idTugas,
			},
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:  "authorization " + entity.RolesById[4] + " passed",
			token: tokens[entity.RolesById[4]],
			payload: req_sub_tugas.CreateByTugasId{
				TugasID: idTugas,
			},
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:  "err: authorization " + entity.RolesById[5],
			token: tokens[entity.RolesById[5]],
			payload: req_sub_tugas.CreateByTugasId{
				TugasID: idTugas,
			},
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/tugas/sub_tugas/"+tt.payload.TugasID, tt.payload, &tt.token)
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

var idSubTugas string

func SubTugasUpdate(t *testing.T) {
	subTugas := new(entity.SubTugas)
	if err := dbt.First(subTugas).Error; err != nil {
		panic(helper.LogsError(err))
		return
	}
	idSubTugas = subTugas.ID
	tests := []struct {
		name         string
		token        string
		payload      req_sub_tugas.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_sub_tugas.Update{
				ID:        idSubTugas,
				Nama:      "tugas testing",
				Status:    "SELESAI",
				Deskripsi: "test",
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:  "sukses spv",
			token: tokens[entity.RolesById[5]],
			payload: req_sub_tugas.Update{
				ID:        idSubTugas,
				Nama:      "tugas testing2",
				Status:    "DIPROSES",
				Deskripsi: "test2",
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
			payload: req_sub_tugas.Update{
				ID:        spvId,
				Nama:      "tugas testing",
				Status:    "BELUM_DIKERJAKAN",
				Deskripsi: "test",
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name:  "err: format status harus salah satu dari ini [BELUM_DIKERJAKAN,DIPROSES,SELESAI]",
			token: tokens[entity.RolesById[1]],
			payload: req_sub_tugas.Update{
				ID:        idSubTugas,
				Nama:      "tugas testing",
				Status:    "ASDFA123",
				Deskripsi: "test",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"status harus berupa salah satu dari [BELUM_DIKERJAKAN,DIPROSES,SELESAI]"},
			},
		},
		{
			name: "err: authorization " + entity.RolesById[2],
			payload: req_sub_tugas.Update{
				ID: idSubTugas,
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
			payload: req_sub_tugas.Update{
				ID: idSubTugas,
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
			payload: req_sub_tugas.Update{
				ID: idSubTugas,
			},
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/tugas/sub_tugas/"+tt.payload.ID, tt.payload, &tt.token)
			assert.NoError(t, err)
			if strings.Contains(tt.name, "passed") {
				assert.NotEqual(t, tt.expectedCode, code)
				assert.NotEqual(t, tt.expectedBody.Code, body.Code)
				assert.NotEqual(t, tt.expectedBody.Status, body.Status)
				return
			}
			if tt.name == "sukses spv" {
				sbtugas := new(entity.SubTugas)
				if err := dbt.First(sbtugas, "id = ?", idSubTugas).Error; err != nil {
					panic(helper.LogsError(err))
					return
				}
				assert.NotEqual(t, tt.payload.Nama, sbtugas.Nama)
				assert.NotEqual(t, tt.payload.Deskripsi, sbtugas.Deskripsi)
				assert.Equal(t, tt.payload.Status, sbtugas.Status)
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

func SubTugasDelete(t *testing.T) {
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
			id:           idSubTugas,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: tidak ditemukan",
			token:        tokens[entity.RolesById[1]],
			id:           idSubTugas,
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[2],
			id:           idSubTugas,
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[3],
			id:           idSubTugas,
			token:        tokens[entity.RolesById[3]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "err: authorization " + entity.RolesById[5],
			id:           idSubTugas,
			token:        tokens[entity.RolesById[5]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[4] + " passed",
			id:           idSubTugas,
			token:        tokens[entity.RolesById[4]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/tugas/sub_tugas/"+tt.id, nil, &tt.token)
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

package test_user

import (
	"testing"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req_user_jenis_spv "github.com/be-sistem-informasi-konveksi/common/request/user/jenis_spv"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func UserCreateSpv(t *testing.T) {
	tests := []struct {
		name         string
		payload      req_user_jenis_spv.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_user_jenis_spv.Create{
				Nama: "belanja",
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name: "err: conflict",
			payload: req_user_jenis_spv.Create{
				Nama: "belanja",
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name: "err: wajib diisi",
			payload: req_user_jenis_spv.Create{
				Nama: "",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"nama wajib diisi"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/user/jenis_spv", tt.payload, &token)
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

var idSpv string

func UserUpdateSpv(t *testing.T) {
	spv := new(entity.JenisSpv)
	err := dbt.Select("id").First(spv).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	spvConflict := &entity.JenisSpv{
		Base: entity.Base{
			ID: test.UlidPkg.MakeUlid().String(),
		},
		Nama: "conflict",
	}
	err = dbt.Create(spvConflict).Error
	if err != nil {
		helper.LogsError(err)
		return
	}
	idSpv = spv.ID
	tests := []struct {
		name         string
		payload      req_user_jenis_spv.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_user_jenis_spv.Update{
				ID:   spv.ID,
				Nama: "bordir",
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name: "err: conflict",
			payload: req_user_jenis_spv.Update{
				ID:   spv.ID,
				Nama: spvConflict.Nama,
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name: "err: tidak ditemukan",
			payload: req_user_jenis_spv.Update{
				ID:   "01HM4B8QBH7MWAVAYP10WN6PKB",
				Nama: "bordir2",
			},
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name: "err: ulid tidak valid",
			payload: req_user_jenis_spv.Update{
				ID:   spv.ID + "123",
				Nama: "bordir3",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/user/jenis_spv/"+tt.payload.ID, tt.payload, &token)
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

func UserGetAllSpv(t *testing.T) {
	tests := []struct {
		name         string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/user/jenis_spv", nil, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			var res []any
			err = mapstructure.Decode(body.Data, &res)
			assert.NoError(t, err)
			assert.Greater(t, len(res), 0)
			assert.NotEmpty(t, res[0])
			for _, v := range res {
				assert.NotEmpty(t, v.(map[string]any)["id"])
				assert.NotEmpty(t, v.(map[string]any)["created_at"])
				assert.NotEmpty(t, v.(map[string]any)["nama"])
			}
			assert.Equal(t, tt.expectedBody.Status, body.Status)

		})
	}
}

func UserDeleteSpv(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:         "sukses",
			id:           idSpv,
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: tidak ditemukan",
			id:           "01HM4B8QBH7MWAVAYP10WN6PKA",
			expectedCode: 404,
			expectedBody: test.Response{
				Status: fiber.ErrNotFound.Message,
				Code:   404,
			},
		},
		{
			name:         "err: ulid tidak valid",
			id:           idSpv + "123",
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"id tidak berupa ulid yang valid"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "DELETE", "/api/v1/user/jenis_spv/"+tt.id, nil, &token)
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

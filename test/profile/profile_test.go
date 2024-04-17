package test_profile

import (
	"strings"
	"testing"

	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_profile "github.com/be-sistem-informasi-konveksi/common/request/profile"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func ProfileGet(t *testing.T) {
	tests := []struct {
		name         string
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
			name:         "authorization " + entity.RolesById[2] + " passed",
			token:        tokens[entity.RolesById[2]],
			expectedCode: 401,
			expectedBody: test.Response{
				Status: fiber.ErrUnauthorized.Message,
				Code:   401,
			},
		},
		{
			name:         "authorization " + entity.RolesById[3] + " passed",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/profile", nil, &tt.token)
			assert.NoError(t, err)
			if strings.Contains(tt.name, "passed") {
				assert.NotEqual(t, tt.expectedCode, code)
				assert.NotEqual(t, tt.expectedBody.Code, body.Code)
				assert.NotEqual(t, tt.expectedBody.Status, body.Status)
				return
			}
			assert.Equal(t, tt.expectedCode, code)
			var res map[string]interface{}
			if tt.name == "sukses" {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
				dataUserDefault := static_data.DefaultUsers[0]
				assert.Equal(t, res["id"], dataUserDefault.ID)
				assert.Equal(t, res["nama"], dataUserDefault.Nama)
				assert.Equal(t, res["role"], dataUserDefault.Role)
				assert.Equal(t, res["username"], dataUserDefault.Username)
				assert.Equal(t, res["no_telp"], dataUserDefault.NoTelp)
				assert.Equal(t, res["alamat"], dataUserDefault.Alamat)
				assert.NotEmpty(t, res["created_at"])
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

func ProfileUpdate(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		payload      req_profile.Update
		expectedBody test.Response
		expectedCode int
	}{
		{
			name:  "sukses",
			token: tokens[entity.RolesById[1]],
			payload: req_profile.Update{
				Nama:        "direktur2",
				Username:    static_data.CredentialUsers[entity.RolesById[1]].Username + "123",
				OldPassword: static_data.CredentialUsers[entity.RolesById[1]].Password,
				NewPassword: static_data.CredentialUsers[entity.RolesById[1]].Password + "123",
				NoTelp:      "+62895397290999",
				Alamat:      "jln. botoijo heaven 1",
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:  "err: format field old, new password dan no telp",
			token: tokens[entity.RolesById[2]],
			payload: req_profile.Update{
				Nama:        "bendahara",
				Username:    static_data.CredentialUsers[entity.RolesById[2]].Username + "123",
				OldPassword: static_data.CredentialUsers[entity.RolesById[2]].Password + "123",
				NewPassword: static_data.CredentialUsers[entity.RolesById[2]].Password + "123",
				NoTelp:      "082895397290991",
				Alamat:      "jln. botoijo heaven 1",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"new password tidak boleh sama dengan old password", "no telp harus berformat e164"},
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
			name:         "authorization " + entity.RolesById[3] + " passed",
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
			code, body, err := test.GetJsonTestRequestResponse(app, "PUT", "/api/v1/profile", tt.payload, &tt.token)
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

package test_auth

import (
	"testing"

	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_auth "github.com/be-sistem-informasi-konveksi/common/request/auth"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

var refreshToken string

func AuthLogin(t *testing.T) {
	tests := []struct {
		name         string
		payload      req_auth.Login
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_auth.Login{
				Username: static_data.CredentialUsers[entity.RolesById[1]].Username,
				Password: static_data.CredentialUsers[entity.RolesById[1]].Password,
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name: "err: username atau password tidak valid",
			payload: req_auth.Login{
				Username: static_data.CredentialUsers[entity.RolesById[1]].Username + "123",
				Password: static_data.CredentialUsers[entity.RolesById[1]].Password + "123",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.InvalidUsernameOrPassword},
			},
		},
		{
			name:         "err: wajib diisi",
			payload:      req_auth.Login{},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"username wajib diisi", "password wajib diisi"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/auth/login", tt.payload, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			var res map[string]any
			if tt.name == "sukses" {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res["token"])
				assert.NotEmpty(t, res["refresh_token"])
				refreshToken = res["refresh_token"].(string)
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

func AuthRefreshToken(t *testing.T) {
	tests := []struct {
		name         string
		payload      req_auth.GetNewToken
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_auth.GetNewToken{
				RefreshToken: refreshToken,
			},
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name: "err: token tidak valid untuk id yang tidak ditemukan",
			payload: req_auth.GetNewToken{
				RefreshToken: refreshTokenUserDoesntExits,
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.InvalidRefreshToken},
			},
		},
		{
			name: "err: token tidak valid untuk signature yang tidak valid",
			payload: req_auth.GetNewToken{
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJyZWZyZXNoX3Rva2VuIiwiZXhwIjoxNzEyNDcwMjYxLCJpZCI6IjAxSFNOOEMxS1YwVFpCTUhQVFBUSk5UWTg4In0.CBFBWuN_yAmfZYqDrBk7mVLkKqiXSABZ5EJOCJ6B9yA",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.InvalidRefreshToken},
			},
		},
		{
			name: "err: refresh token expired",
			payload: req_auth.GetNewToken{
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJyZWZyZXNoX3Rva2VuIiwiZXhwIjoxMDAwMDAwMDAwLCJpZCI6IjAxSFNOOEMxS1YwVFpCTUhQVFBUSk5UWTg4In0.709E6kJosXvEyhknlIL6wTOrrSrA2P9DSAyDelrE57k",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{message.RefreshTokenExpired},
			},
		},
		{
			name:         "err: wajib diisi",
			payload:      req_auth.GetNewToken{},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"refresh token wajib diisi"},
			},
		},
		{
			name: "err: format tidak jwt",
			payload: req_auth.GetNewToken{
				RefreshToken: "asdfasdf",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"refresh token harus berformat jwt"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/auth/refresh_token", tt.payload, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			var res map[string]any
			if tt.name == "sukses" {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res["token"])
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

func AuthWhoAmI(t *testing.T) {
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
			name:         "sukses",
			token:        tokens[entity.RolesById[2]],
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses",
			token:        tokens[entity.RolesById[3]],
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses",
			token:        tokens[entity.RolesById[4]],
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "sukses",
			token:        tokens[entity.RolesById[5]],
			expectedCode: 200,
			expectedBody: test.Response{
				Status: message.OK,
				Code:   200,
			},
		},
		{
			name:         "err: unauthorized",
			token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhY2Nlc3NfdG9rZW4iLCJleHAiOjE3MDAwMDAwMDAsImlkIjoiMDFIU044QzFLVjBUWkJNSFBUUFRKTlRZODgiLCJuYW1hIjoiYWt1bl9kaXJla3R1ciIsInVzZXJuYW1lIjoiYWt1bl9kaXJla3R1ciIsInJvbGUiOiJESVJFS1RVUiJ9.fIGDTc-JKISFqzfFPoPXat1dcxM04dvoNMs6q_xNSgM",
			expectedCode: 401,
			expectedBody: test.Response{
				Status:         fiber.ErrUnauthorized.Message,
				Code:           401,
				ErrorsMessages: []string{message.UnauthInvalidToken},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "GET", "/api/v1/auth/whoami", nil, &tt.token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			var res map[string]any
			if tt.name == "sukses" {
				err = mapstructure.Decode(body.Data, &res)
				assert.NoError(t, err)
				parse, err := jwtPkg.ParseToken(false, tt.token, &pkg.Claims{})
				assert.NoError(t, err)
				assert.NotEmpty(t, parse)
				claims, ok := parse.Claims.(*pkg.Claims)
				assert.True(t, ok)

				assert.Equal(t, float64(claims.ExpiresAt.Unix()), res["exp"])
				assert.Equal(t, claims.Subject, res["sub"])
				assert.Equal(t, claims.ID, res["id"])
				assert.Equal(t, claims.Nama, res["nama"])
				assert.Equal(t, claims.Username, res["username"])
				assert.Equal(t, claims.Role, res["role"])
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

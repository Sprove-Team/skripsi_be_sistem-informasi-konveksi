package test_akuntansi

import (
	"fmt"
	"testing"

	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_akuntansi_kelompok_akun "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/kelompok_akun"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func AkuntansiCreateKelompokAkun(t *testing.T) {
	tests := []struct {
		name         string
		payload      req_akuntansi_kelompok_akun.Create
		expectedBody test.Response
		expectedCode int
	}{
		{
			name: "sukses",
			payload: req_akuntansi_kelompok_akun.Create{
				Nama:         "kelompok test",
				Kode:         "11111111",
				KategoriAkun: entity.KategoriAkunByKode["1"],
			},
			expectedCode: 201,
			expectedBody: test.Response{
				Status: message.Created,
				Code:   201,
			},
		},
		{
			name: "err: conflict",
			payload: req_akuntansi_kelompok_akun.Create{
				Nama:         static_data.DataKelompokAkun[0].Nama,
				Kode:         "1",
				KategoriAkun: static_data.DataKelompokAkun[0].KategoriAkun,
			},
			expectedCode: 409,
			expectedBody: test.Response{
				Status: fiber.ErrConflict.Message,
				Code:   409,
			},
		},
		{
			name: "err: kategori akun harus berupa salah satu dari [ASET,KEWAJIBAN,MODAL,PENDAPATAN,BEBAN]",
			payload: req_akuntansi_kelompok_akun.Create{
				Nama:         "kelompok err kategori",
				Kode:         "2",
				KategoriAkun: "asdf",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"kategori akun harus berupa salah satu dari [ASET,KEWAJIBAN,MODAL,PENDAPATAN,BEBAN]"},
			},
		},
		{
			name: "err: wajib diisi",
			payload: req_akuntansi_kelompok_akun.Create{
				Nama:         "",
				Kode:         "",
				KategoriAkun: "",
			},
			expectedCode: 400,
			expectedBody: test.Response{
				Status:         fiber.ErrBadRequest.Message,
				Code:           400,
				ErrorsMessages: []string{"nama wajib diisi", "kode wajib diisi", "kategori akun wajib diisi"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/akuntansi/kelompok_akun", tt.payload, &token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, code)
			fmt.Println("err->", body.ErrorsMessages)
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

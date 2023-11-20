package produk_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handler "github.com/be-sistem-informasi-konveksi/api/handler/produk/harga_detail"
	mockUC "github.com/be-sistem-informasi-konveksi/api/usecase/produk/harga_detail/mock"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk/harga_detail"
	resGlobal "github.com/be-sistem-informasi-konveksi/common/response/global"
	"github.com/be-sistem-informasi-konveksi/pkg"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	mockCtrl := gomock.NewController(t)
	mockUC := mockUC.NewMockHargaDetailProdukUsecase(mockCtrl)
	defer mockCtrl.Finish()
	validator := pkg.NewValidator()
	h := handler.NewHargaDetailProdukHandler(mockUC, validator)
	
	tests := []struct {
		name         string
		reqData      interface{}
		expectedBody interface{}
		expectedCode int
		wantErr      bool
	}{
		{
			name: "Success: 201",
			reqData: req.Create{
				ProdukId: "0b7561357efa11eea0e95efc22537c19",
				HargaDetail: []req.HargaDetail{
					{
						QTY: 1,
						Harga: 50000,
					},
					{
						QTY: 2,
						Harga: 10000,
					},
				},
			},
			expectedBody: resGlobal.SuccessResWithoutData("C"),
			expectedCode: fiber.StatusCreated,
			wantErr:      false,
		},
		{
			name: "Failed: 400 ~ validator",
			reqData: req.Create{
				ProdukId: "",
				HargaDetail: []req.HargaDetail{},
			},
			expectedCode: fiber.StatusBadRequest,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			reqBody, err := json.Marshal(tt.reqData)
			assert.NoError(t, err)
			
			reqT := httptest.NewRequest(http.MethodPost, "/api/v1/direktur/produk/harga_detail", bytes.NewReader(reqBody))

			reqT.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			
			
			mockUC.EXPECT().Create(reqT.Context(), tt.reqData.(req.Create)).Return(nil).AnyTimes()
			
			app := fiber.New()

			app.Post("/api/v1/direktur/produk/harga_detail", h.Create)

			resP, _ :=  app.Test(reqT, -1)

			assert.Equal(t,resP.StatusCode, tt.expectedCode, "unexpected status code value")

			var body struct {
				Data    interface{} `json:"data,omitempty"`
				Message string      `json:"message"`
				Status  int         `json:"status"`
			}

			json.NewDecoder(resP.Body).Decode(&body)

			assert.Equal(t,body.Status, tt.expectedCode, "unexpected status value")
			assert.NotEmpty(t, body.Message, "message should not empty")

			if tt.wantErr {
				if strings.Contains(tt.name, "validator") {
					assert.NotEmpty(t, body.Data, "data should not empty")
					assert.NotNil(t, body.Data.(map[string]interface{})["error_field"], "error_field should not nil")
				}
			}else {
				assert.Nil(t, body.Data)
			}
		})
	}
	
}

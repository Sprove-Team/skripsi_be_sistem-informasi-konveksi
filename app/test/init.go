package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status         string   `json:"status"`
	Code           int      `json:"code"`
	Data           any      `json:"data,omitempty"`
	ErrorsMessages []string `json:"errors_messages,omitempty"`
}

func GetJsonTestRequestResponse(app *fiber.App, method string, url string, reqBody any) (code int, respBody Response, err error) {
	bodyJson := []byte("")
	if reqBody != nil {
		bodyJson, _ = json.Marshal(reqBody)
	}

	req := httptest.NewRequest(method, url, bytes.NewReader(bodyJson))
	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := app.Test(req, 50)
	code = resp.StatusCode
	// If error we're done
	if err != nil {
		return
	}
	// If no body content, we're done
	if resp.ContentLength == 0 {
		return
	}
	bodyData := make([]byte, resp.ContentLength)
	_, _ = resp.Body.Read(bodyData)
	err = json.Unmarshal(bodyData, &respBody)

	return
}

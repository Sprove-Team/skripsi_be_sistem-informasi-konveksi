package global

import (
	response "github.com/be-sistem-informasi-konveksi/common/reponse"
)

func CustomRes(status int, message string, data *map[string]interface{}) *response.BaseFormatRes {
	if data != nil {
		return &response.BaseFormatRes{
			Status:  status,
			Data:    *data,
			Message: message,
		}
	}
	return &response.BaseFormatRes{
		Status:  status,
		Message: message,
	}
}

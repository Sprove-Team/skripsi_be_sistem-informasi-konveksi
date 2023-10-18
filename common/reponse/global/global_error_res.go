package global

import (
	"github.com/be-sistem-informasi-konveksi/common/message"
	response "github.com/be-sistem-informasi-konveksi/common/reponse"
)

func getErrStatusMessage(status int) string {
	switch status {
	case 400:
		return message.BadRequest
  case 409:
    return message.Conflict
	case 404:
		return message.NotFound
	default:
		return message.InternalServerError
	}
}

func ErrorResWithoutData(status int) *response.BaseFormatRes {
	return &response.BaseFormatRes{
		Message: getErrStatusMessage(status),
		Status:  status,
	}
}

func ErrorResWithData(data interface{}, status int) *response.BaseFormatRes {
	return &response.BaseFormatRes{
		Message: getErrStatusMessage(status),
		Status:  status,
		Data: map[string]interface{}{
			"errors": data,
		},
	}
}

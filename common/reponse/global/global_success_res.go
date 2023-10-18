package global

import (
	"github.com/be-sistem-informasi-konveksi/common/message"
	response "github.com/be-sistem-informasi-konveksi/common/reponse"
)

func getTypeMessage(char string) (int, string) {
	switch char {
	case "C":
		return 201, message.CreateDataOK
	case "R":
		return 200, message.GetDataOK
	case "U":
		return 200, message.UpdateDataOK
	case "D":
		return 200, message.DeleteDataOK
	default:
		return 200, message.OK
	}
}

func SuccessResWithData(data map[string]interface{}, chars ...string) *response.BaseFormatRes {
	status, msg := getTypeMessage(chars[0])
	return &response.BaseFormatRes{
		Status:  status,
		Message: msg,
		Data:    data,
	}
}

func SuccessResWithoutData(chars ...string) *response.BaseFormatRes {
	status, msg := getTypeMessage(chars[0])
	return &response.BaseFormatRes{
		Status:  status,
		Message: msg,
	}
}

package response

// type Erro

type BaseFormat struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

type BaseFormatSuccessRes struct {
	BaseFormat
	Data interface{} `json:"data,omitempty"`
}

type BaseFormatError struct {
	BaseFormat
	Errors map[string][]string `json:"errors,omitempty"`
}

// func BaseSuccessRes

func BaseRes(code int, status string) *BaseFormat {
	return &BaseFormat{
		Code:   code,
		Status: status,
	}
}

func ErrorRes(code int, status string, errors map[string][]string) *BaseFormatError {
	return &BaseFormatError{
		BaseFormat: BaseFormat{
			Code:   code,
			Status: status,
		},
		Errors: errors,
	}
}

func SuccessRes(code int, status string, data interface{}) *BaseFormatSuccessRes {
	return &BaseFormatSuccessRes{
		BaseFormat: BaseFormat{
			Code:   code,
			Status: status,
		},
		Data: data,
	}
}

// type BaseFormatRes struct {
// 	Data    interface{} `json:"data,omitempty"`
// 	Message string      `json:"message"`
// 	Status  int         `json:"status"`
// }

// FieldName string `json:"field_name"`
// Message   string `json:"message"`

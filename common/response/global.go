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
	ErrorsMessages []string `json:"errors_messages,omitempty"`
}

type BaseFormatInterErr struct {
	BaseFormat
	Message string `json:"message"`
}

// func BaseSuccessRes

func BaseRes(code int, status string) *BaseFormat {
	return &BaseFormat{
		Code:   code,
		Status: status,
	}
}

func ErrorRes(code int, status string, errors []string) *BaseFormatError {
	return &BaseFormatError{
		BaseFormat: BaseFormat{
			Code:   code,
			Status: status,
		},
		ErrorsMessages: errors,
	}
}

func ErrorInterWithMessageRes(code int, status, message string) *BaseFormatInterErr {
	return &BaseFormatInterErr{
		BaseFormat: BaseFormat{
			Code:   code,
			Status: status,
		},
		Message: message,
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

package response

type BaseFormatRes struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

type BaseFormatError struct {
	ValueInput   interface{} `json:"value_input"`
	ErrorMessage string      `json:"error_message"`  
}

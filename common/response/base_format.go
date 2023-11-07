package response

type BaseFormatRes struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

type BaseFormatError struct {
	FieldName string `json:"field_name"`
	Message   string `json:"message"`
}

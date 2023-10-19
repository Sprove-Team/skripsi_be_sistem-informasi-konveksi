package response

type BaseFormatRes struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

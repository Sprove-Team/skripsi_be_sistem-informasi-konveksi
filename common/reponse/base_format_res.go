package response

type BaseFormatRes struct {
	Data    map[string]interface{} `json:"data,omitempty"`
	Message string                 `json:"message"`
	Status  int                    `json:"status"`
}

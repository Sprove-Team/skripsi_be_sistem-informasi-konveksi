package sablon

type Create struct {
	Nama  string  `json:"nama" validate:"required,printascii"`
	Harga float64 `json:"harga" validate:"required,number"`
}

type Update struct {
	ID    string  `params:"id" validate:"required,ulid"`
	Nama  string  `json:"nama" validate:"omitempty,printascii"`
	Harga float64 `json:"harga" validate:"omitempty,number"`
}

type GetAll struct {
	Nama string `query:"nama" validate:"omitempty,printascii"`
	// Page   int                `query:"page" validate:"omitempty,number"`
	Next  string `query:"next" validate:"omitempty"`
	Limit int    `query:"limit" validate:"omitempty,number"`
}

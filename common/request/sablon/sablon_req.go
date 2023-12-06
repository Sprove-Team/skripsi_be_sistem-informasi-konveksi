package sablon

type CreateSablon struct {
	Nama  string  `json:"nama" validate:"required,printascii"`
	Harga float64 `json:"harga" validate:"required,number"`
}

type UpdateSablon struct {
	ID    string  `params:"id" validate:"required,ulid"`
	Nama  string  `json:"nama" validate:"omitempty,printascii"`
	Harga float64 `json:"harga" validate:"omitempty,number"`
}

type SearchFilterSablon struct {
	Nama string `json:"nama" validate:"omitempty,printascii"`
}

type GetAllSablon struct {
	Search SearchFilterSablon `json:"search" validate:"omitempty"`
	// Page   int                `query:"page" validate:"omitempty,number"`
	Next  string `query:"next" validate:"omitempty"`
	Limit int    `query:"limit" validate:"omitempty,number"`
}

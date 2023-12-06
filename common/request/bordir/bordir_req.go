package bordir

type CreateBordir struct {
	Nama  string  `json:"nama" validate:"required,printascii"`
	Harga float64 `json:"harga" validate:"required,number"`
}

type UpdateBordir struct {
	ID    string  `params:"id" validate:"required,ulid"`
	Nama  string  `json:"nama" validate:"omitempty,printascii"`
	Harga float64 `json:"harga" validate:"omitempty,number"`
}

type SearchFilterBordir struct {
	Nama string `json:"nama" validate:"omitempty,printascii"`
}

type GetAllBordir struct {
	Search SearchFilterBordir `json:"search" validate:"omitempty"`
	// Page   int                `query:"page" validate:"omitempty,number"`
	Next  string `query:"page" validate:"omitempty"`
	Limit int    `query:"limit" validate:"omitempty,number"`
}

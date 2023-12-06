package produk

type Create struct {
	Nama string `json:"nama" validate:"required,printascii"`
}

type Update struct {
	Nama string `json:"nama" validate:"omitempty,printascii"`
	ID   string `json:"id" validate:"required,ulid"`
}

type SearchFilter struct {
	Nama string `json:"nama" validate:"omitempty,printascii"`
}

type GetAll struct {
	Search SearchFilter `json:"search" validate:"omitempty"`
	// Page   int          `query:"page" validate:"omitempty,number"`
	Next  string `query:"next" validate:"omitempty"`
	Limit int    `query:"limit" validate:"omitempty,number"`
}

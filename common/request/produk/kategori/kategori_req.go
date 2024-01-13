package produk

type Create struct {
	Nama string `json:"nama" validate:"required,printascii"`
}

type Update struct {
	Nama string `json:"nama" validate:"omitempty,printascii"`
	ID   string `json:"id" validate:"required,ulid"`
}

type GetAll struct {
	Nama  string `query:"nama" validate:"omitempty,printascii"`
	Next  string `query:"next" validate:"omitempty"`
	Limit int    `query:"limit" validate:"omitempty,number"`
}

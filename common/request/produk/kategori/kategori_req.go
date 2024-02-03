package produk

type Create struct {
	Nama string `json:"nama" validate:"required"`
}

type Update struct {
	ID   string `params:"id" validate:"required,ulid"`
	Nama string `json:"nama" validate:"omitempty"`
}

type GetAll struct {
	Nama  string `query:"nama" validate:"omitempty"`
	Next  string `query:"next" validate:"omitempty"`
	Limit int    `query:"limit" validate:"omitempty,number"`
}

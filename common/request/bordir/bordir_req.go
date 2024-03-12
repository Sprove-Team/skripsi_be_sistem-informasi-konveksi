package req_bordir

type Create struct {
	Nama  string  `json:"nama" validate:"required"`
	Harga float64 `json:"harga" validate:"required,number"`
}

type Update struct {
	ID    string  `params:"id" validate:"required,ulid"`
	Nama  string  `json:"nama" validate:"omitempty"`
	Harga float64 `json:"harga" validate:"omitempty,number"`
}

type GetAll struct {
	Nama  string `query:"nama" validate:"omitempty"`
	Next  string `query:"page" validate:"omitempty"`
	Limit int    `query:"limit" validate:"omitempty,number"`
}

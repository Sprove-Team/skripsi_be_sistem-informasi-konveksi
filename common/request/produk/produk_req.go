package produk

type Create struct {
	Nama       string `json:"nama" validate:"required,printascii"`
	KategoriID string `json:"kategori_id" validate:"required,ulid"`
}

type Update struct {
	ID         string `params:"id" validate:"required,ulid"`
	Nama       string `json:"nama" validate:"omitempty,printascii"`
	KategoriID string `json:"kategori_id" validate:"omitempty,ulid"`
}

// type SearchFilter struct {
// 	Nama        string `json:"nama" validate:"omitempty,printascii"`
// 	KategoriID  string `json:"kategori_id" validate:"omitempty,ulid"`
// 	HargaDetail string `json:"harga_detail" validate:"omitempty,oneof=EMPTY NOT_EMPTY"`
// }

type GetAll struct {
	Nama        string `query:"nama" validate:"omitempty,printascii"`
	KategoriID  string `query:"kategori_id" validate:"omitempty,ulid"`
	HargaDetail string `query:"harga_detail" validate:"omitempty,oneof=EMPTY NOT_EMPTY"`

	// Page   int                `query:"page" validate:"omitempty,number"`
	Next  string `query:"next" validate:"omitempty"`
	Limit int    `query:"limit" validate:"omitempty,number"`
}

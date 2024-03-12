package req_produk

type Create struct {
	Nama       string `json:"nama" validate:"required"`
	KategoriID string `json:"kategori_id" validate:"required,ulid"`
}

type Update struct {
	ID         string `params:"id" validate:"required,ulid"`
	Nama       string `json:"nama" validate:"omitempty"`
	KategoriID string `json:"kategori_id" validate:"omitempty,ulid"`
}

type GetAll struct {
	Nama        string `query:"nama" validate:"omitempty"`
	KategoriID  string `query:"kategori_id" validate:"omitempty,ulid"`
	HargaDetail string `query:"harga_detail" validate:"omitempty,oneof=EMPTY NOT_EMPTY"`
	Next        string `query:"next" validate:"omitempty"`
	Limit       int    `query:"limit" validate:"omitempty,number"`
}

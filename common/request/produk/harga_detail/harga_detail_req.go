package produk

type Create struct {
	ProdukId string  `json:"produk_id" validate:"required,ulid"`
	QTY      uint    `json:"qty" validate:"required,number,gt=0"`
	Harga    float64 `json:"harga" validate:"required,number"`
}

type Update struct {
	ID    string  `params:"id" validate:"required,ulid"`
	QTY   uint    `json:"qty" validate:"omitempty,number,gt=0"`
	Harga float64 `json:"harga" validate:"omitempty,number"`
}

type GetByProdukId struct {
	ProdukId string `params:"produk_id" validate:"required,ulid"`
}

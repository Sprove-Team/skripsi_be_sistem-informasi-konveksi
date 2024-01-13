package produk

type HargaDetail struct {
	QTY   uint    `json:"qty" validate:"required,number"`
	Harga float64 `json:"harga" validate:"required,number"`
}

type Create struct {
	ProdukId    string        `json:"produk_id" validate:"required,ulid"`
	HargaDetail []HargaDetail `json:"harga_detail" validate:"gt=0,dive,required"`
}

type Update struct {
	ID    string  `params:"id" validate:"required,ulid"`
	QTY   uint    `json:"qty" validate:"omitempty,number"`
	Harga float64 `json:"harga" validate:"omitempty,number"`
}

type DeleteByProdukId struct {
	ProdukId string `params:"produk_id" validate:"required,ulid"`
}

type GetByProdukId struct {
	ProdukId string `params:"produk_id" validate:"required,ulid"`
}

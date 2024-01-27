package produk

type HargaDetail struct {
	QTY   uint    `json:"qty" validate:"required,number,gt=1"`
	Harga float64 `json:"harga" validate:"required,number"`
}

type CreateByProdukId struct {
	ProdukId    string        `params:"produk_id" validate:"required,ulid"`
	HargaDetail []HargaDetail `json:"harga_detail" validate:"gt=0,dive,required"`
}

type HargaDetailWithId struct {
	ID    string  `json:"id" validate:"required,ulid"`
	QTY   uint    `json:"qty" validate:"required,number"`
	Harga float64 `json:"harga" validate:"required,number"`
}

type UpdateByProdukId struct {
	ProdukId    string              `params:"produk_id" validate:"required,ulid"`
	HargaDetail []HargaDetailWithId `json:"harga_detail" validate:"gt=0,dive,required"`
}

type DeleteByProdukId struct {
	ProdukId string `params:"produk_id" validate:"required,ulid"`
}

type GetByProdukId struct {
	ProdukId string `params:"produk_id" validate:"required,ulid"`
}

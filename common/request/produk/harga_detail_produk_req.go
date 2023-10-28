package produk

type hargaProduk struct {
	QTY   uint    `json:"qty" validate:"required,number"`
	Harga float64 `json:"harga" validate:"required,number"`
}

type CreateHargaDetailProduk struct {
	ProdukId     string        `json:"produk_id" validate:"required,uuidv4_no_hyphens"`
	HargaProduks []hargaProduk `json:"harga_produks" validate:"gt=0,dive,required"`
}

type UpdateHargaDetailProdukById struct {
	ProdukId string  `json:"produk_id" validate:"omitempty,uuidv4_no_hyphens"`
	ID       string  `json:"id" validate:"required,uuidv4_no_hyphens"`
	QTY      uint    `json:"qty" validate:"omitempty,number"`
	Harga    float64 `json:"harga" validate:"omitempty,number"`
}

type DeleteByProdukId struct {
	ProdukId string `params:"produk_id" validate:"required,uuidv4_no_hyphens"`
}

type GetByProdukId struct {
	ProdukId string `params:"produk_id" validate:"required,uuidv4_no_hyphens"`
}

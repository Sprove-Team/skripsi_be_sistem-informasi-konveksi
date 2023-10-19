package produk

type hargaProduk struct {
	QTY   uint    `json:"qty" validate:"required,number"`
	Harga float64 `json:"harga" validate:"required,number"`
}

type CreateHargaDetailProduk struct {
	ProdukId     string        `json:"produk_id" validate:"required,uuidv4_no_hyphens"`
	HargaProduks []hargaProduk `json:"harga_produks" validate:"gt=0,dive,required"`
}

package direktur

type CreateHargaDetailProduk struct {
	ProdukId string  `json:"produk_id" validate:"required,uuid4"`
	QTY      uint    `json:"qty" validate:"required,number"`
	Harga    float64 `json:"harga" validate:"required"`
}

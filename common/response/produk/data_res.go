package produk

type HargaDetailsRes struct {
	ID    uint    `json:"id"`
	QTY   uint    `json:"qty"`
	Harga float64 `json:"harga"`
}

type DataProdukRes struct {
	ID          string            `json:"id"`
	Nama        string            `json:"nama"`
	HargaDetail []HargaDetailsRes `json:"harga_detail"`
	KategoriId  uint              `json:"kategori_id"`
}

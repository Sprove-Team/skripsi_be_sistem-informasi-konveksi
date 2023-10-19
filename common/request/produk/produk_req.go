package produk

type CreateProduk struct {
	Nama       string `json:"nama" validate:"required,ascii"`
	IDKategori uint   `json:"id_kategori" validate:"required,number"`
}

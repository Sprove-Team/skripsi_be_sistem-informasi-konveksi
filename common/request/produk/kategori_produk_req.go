package produk

type CreateKategoriProduk struct {
	Nama string `json:"nama" validate:"required,ascii"`
}

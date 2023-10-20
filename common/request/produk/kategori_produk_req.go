package produk

type CreateKategoriProduk struct {
	Nama string `json:"nama" validate:"required,ascii"`
}

type UpdateKategoriProduk struct {
	Nama string `json:"nama" validate:"required,ascii"`
	ID   uint   `params:"id" validate:"required,number"`
}

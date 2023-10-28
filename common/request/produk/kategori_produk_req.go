package produk

type CreateKategoriProduk struct {
	Nama string `json:"nama" validate:"required,ascii"`
}

type UpdateKategoriProduk struct {
	Nama string `json:"nama" validate:"omitempty,ascii"`
	ID   string `json:"id" validate:"required,uuidv4_no_hyphens"`
}

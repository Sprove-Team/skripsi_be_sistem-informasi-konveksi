package produk

type CreateKategoriProduk struct {
	Nama string `json:"nama" validate:"required,ascii"`
}

type UpdateKategoriProduk struct {
	Nama string `json:"nama" validate:"omitempty,ascii"`
	ID   string `json:"id" validate:"required,uuidv4_no_hyphens"`
}

type SearchFilterKategoriProduk struct {
	Nama string `json:"nama" validate:"omitempty,ascii"`
}

type GetAllKategoriProduk struct {
	Search SearchFilterKategoriProduk `json:"search" validate:"omitempty"`
	Page   int                        `query:"page" validate:"omitempty,number"`
	Limit  int                        `query:"limit" validate:"omitempty,number"`
}
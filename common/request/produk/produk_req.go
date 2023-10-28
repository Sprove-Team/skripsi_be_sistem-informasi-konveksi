package produk

type CreateProduk struct {
	Nama       string `json:"nama" validate:"required,ascii"`
	KategoriID string `json:"kategori_id" validate:"required,uuidv4_no_hyphens"`
}

type UpdateProduk struct {
	ID         string `json:"produk_id" validate:"required,uuidv4_no_hyphens"`
	Nama       string `json:"nama" validate:"omitempty,ascii"`
	KategoriID string `json:"kategori_id" validate:"omitempty,uuidv4_no_hyphens"`
}

type SearchFilter struct {
	Nama       string `json:"nama" validate:"omitempty,ascii"`
	KategoriID string `json:"kategori_id" validate:"omitempty,uuidv4_no_hyphens"`
}

type GetAll struct {
	Search SearchFilter `json:"search" validate:"omitempty"`
	Page   int          `query:"page" validate:"omitempty,number"`
	Limit  int          `query:"limit" validate:"omitempty,number"`
}
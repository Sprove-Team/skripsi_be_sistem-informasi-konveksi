package direktur

type CreateKategoriProduk struct {
	Nama string `json:"nama" validate:"required,ascii"`
}

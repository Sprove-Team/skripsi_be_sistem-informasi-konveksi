package akuntansi

type GetAll struct {
	Nama  string `query:"nama" validate:"omitempty,printascii"`
	Kode  string `query:"kode" validate:"omitempty"`
	Next  string `query:"Next" validate:"omitempty,ulid"`
	Limit int    `query:"limit" validate:"omitempty,number"`
}

type Create struct {
	Nama         string `json:"nama" validate:"required,printascii"`
	Kode         string `json:"kode" validate:"required"`
	KategoriAkun string `json:"kategori_akun" validate:"required,oneof=ASET KEWAJIBAN MODAL PENDAPATAN BEBAN"`
}

type Update struct {
	ID           string `param:"id" validate:"required,ulid"`
	Nama         string `json:"nama" validate:"omitempty,printascii"`
	Kode         string `json:"kode" validate:"omitempty"`
	KategoriAkun string `json:"kategori_akun" validate:"omitempty,oneof=ASET KEWAJIBAN MODAL PENDAPATAN BEBAN"`
}

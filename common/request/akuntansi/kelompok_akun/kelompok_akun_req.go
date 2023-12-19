package akuntansi

type Create struct {
	Nama string `json:"nama" validate:"required,printascii"`
	Kode string `json:"kode" validate:"required"`
}

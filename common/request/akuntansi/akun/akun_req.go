package akuntansi

type Create struct {
	GolonganAkunID string `json:"golongan_akun_id" validate:"required,ulid"`
	Nama           string `json:"nama" validate:"required,printascii"`
	Kode           string `json:"kode" validate:"required"`
	Deskripsi      string `json:"deskripsi" validate:"omitempty"`
	SaldoNormal    string `json:"saldo_normal" validate:"required,oneof=DEBIT KREDIT"`
}

type GetAll struct {
	Nama  string `query:"nama" validate:"omitempty,printascii"`
	Kode  string `query:"kode" validate:"omitempty"`
	Next  string `query:"Next" validate:"omitempty,ulid"`
	Limit int    `query:"limit" validate:"omitempty,number"`
}

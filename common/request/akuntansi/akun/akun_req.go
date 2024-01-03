package akuntansi

type Create struct {
	KelompokAkunID string `json:"kelompok_akun_id" validate:"required,ulid"`
	Nama           string `json:"nama" validate:"required,printascii"`
	Kode           string `json:"kode" validate:"required"`
	Deskripsi      string `json:"deskripsi" validate:"omitempty"`
	SaldoNormal    string `json:"saldo_normal" validate:"required,oneof=DEBIT KREDIT"`
}

type Update struct {
	ID             string `param:"id" validate:"required,ulid"`
	KelompokAkunID string `json:"kelompok_akun_id" validate:"omitempty,ulid"`
	Nama           string `json:"nama" validate:"omitempty,printascii"`
	Kode           string `json:"kode" validate:"omitempty"`
	Deskripsi      string `json:"deskripsi" validate:"omitempty"`
	SaldoNormal    string `json:"saldo_normal" validate:"omitempty,oneof=DEBIT KREDIT"`
}

type GetAll struct {
	Nama  string `query:"nama" validate:"omitempty,printascii"`
	Kode  string `query:"kode" validate:"omitempty"`
	Next  string `query:"Next" validate:"omitempty,ulid"`
	Limit int    `query:"limit" validate:"omitempty,number"`
}

package akuntansi

type GetAll struct {
	Nama   string `query:"nama" validate:"omitempty,printascii"`
	NoTelp string `query:"no_telp" validate:"omitempty,e164"`
	Email  string `query:"email" validate:"omitempty,email"`
	Next   string `query:"next" validate:"omitempty,ulid"`
	Limit  int    `query:"limit" validate:"omitempty,number"`
}

type Create struct {
	Nama       string `json:"nama" validate:"required,printascii"`
	NoTelp     string `json:"no_telp" validate:"required,e164"`
	Alamat     string `json:"alamat" validate:"required"`
	Keterangan string `json:"keterangan" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
}

type Update struct {
	ID         string `params:"id" validate:"required,ulid"`
	Nama       string `json:"nama" validate:"omitempty,printascii"`
	NoTelp     string `json:"no_telp" validate:"omitempty,e164"`
	Alamat     string `json:"alamat" validate:"omitempty"`
	Keterangan string `json:"keterangan" validate:"omitempty"`
	Email      string `json:"email" validate:"omitempty,email"`
}

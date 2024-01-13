package user

type Create struct {
	Nama       string `json:"nama" validate:"required,printascii"`
	Role       string `json:"role" validate:"required,oneof=DIREKTUR ADMIN BENDAHARA MANAJER_PRODUKSI SUPERVISOR"`
	Username   string `json:"username" validate:"required,printascii"`
	Password   string `json:"password" validate:"required,printascii,min=6"`
	NoTelp     string `json:"no_telp" validate:"required,e164"`
	Alamat     string `json:"alamat" validate:"required,printascii"`
	JenisSpvID string `json:"jenis_spv_id" validate:"required_if=Role SUPERVISOR,omitempty,ulid"`
}

type Update struct {
	ID         string `params:"id" validate:"required,ulid"`
	Nama       string `json:"nama" validate:"omitempty,printascii"`
	Role       string `json:"role" validate:"omitempty,oneof=DIREKTUR ADMIN BENDAHARA MANAJER_PRODUKSI SUPERVISOR"`
	Username   string `json:"username" validate:"omitempty,printascii"`
	Password   string `json:"password" validate:"omitempty,printascii,min=6"`
	NoTelp     string `json:"no_telp" validate:"omitempty,e164"`
	Alamat     string `json:"alamat" validate:"omitempty,printascii"`
	JenisSpvID string `json:"jenis_spv_id" validate:"required_if=Role SUPERVISOR,omitempty,ulid"`
}

type SearchFilter struct {
	Nama        string `json:"nama" validate:"omitempty,printascii"`
	Role        string `json:"role" validate:"omitempty,oneof=DIREKTUR ADMIN BENDAHARA MANAJER_PRODUKSI SUPERVISOR"`
	Username    string `json:"username" validate:"omitempty,printascii"`
	Alamat      string `json:"alamat" validate:"omitempty,printascii"`
	NoTelp      string `json:"no_telp" validate:"omitempty,e164"`
	JenisSpvID  string `json:"jenis_spv_id" validate:"required_if=AllJenisSpv false,omitempty,ulid"`
	AllJenisSpv bool   `json:"all_jenis_spv" validate:"omitempty,boolean"`
}

type GetAll struct {
	Search SearchFilter `json:"search" validate:"omitempty"`
	Next   string       `query:"page" validate:"omitempty"`
	Limit  int          `query:"limit" validate:"omitempty,number"`
}

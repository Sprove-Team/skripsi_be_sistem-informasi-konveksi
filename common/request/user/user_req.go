package user

type Create struct {
	Nama       string `json:"nama" validate:"required"`
	Role       string `json:"role" validate:"required,oneof=DIREKTUR ADMIN BENDAHARA MANAJER_PRODUKSI SUPERVISOR"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required,min=6"`
	NoTelp     string `json:"no_telp" validate:"required,e164"`
	Alamat     string `json:"alamat" validate:"required"`
	JenisSpvID string `json:"jenis_spv_id" validate:"required_if=Role SUPERVISOR,omitempty,ulid"`
}

type Update struct {
	ID         string `params:"id" validate:"required,ulid"`
	Nama       string `json:"nama" validate:"omitempty"`
	Role       string `json:"role" validate:"omitempty,oneof=DIREKTUR ADMIN BENDAHARA MANAJER_PRODUKSI SUPERVISOR"`
	Username   string `json:"username" validate:"omitempty"`
	Password   string `json:"password" validate:"omitempty,min=6"`
	NoTelp     string `json:"no_telp" validate:"omitempty,e164"`
	Alamat     string `json:"alamat" validate:"omitempty"`
	JenisSpvID string `json:"jenis_spv_id" validate:"required_if=Role SUPERVISOR,omitempty,ulid"`
}

type SearchFilter struct {
	Nama        string `query:"nama" validate:"omitempty"`
	Role        string `query:"role" validate:"omitempty,oneof=DIREKTUR ADMIN BENDAHARA MANAJER_PRODUKSI SUPERVISOR"`
	Username    string `query:"username" validate:"omitempty"`
	Alamat      string `query:"alamat" validate:"omitempty"`
	NoTelp      string `query:"no_telp" validate:"omitempty,e164"`
	JenisSpvID  string `query:"jenis_spv_id" validate:"required_if=AllJenisSpv false,omitempty,ulid"`
	AllJenisSpv bool   `query:"all_jenis_spv" validate:"omitempty,boolean"`
}

type GetAll struct {
	Search SearchFilter `query:"search" validate:"omitempty"`
	Next   string       `query:"page" validate:"omitempty"`
	Limit  int          `query:"limit" validate:"omitempty,number"`
}

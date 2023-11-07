package user

type CreateUser struct {
	Nama       string `json:"nama" validate:"required,printascii"`
	Role       string `json:"role" validate:"required,oneof=DIREKTUR ADMIN BENDAHARA MANAJER_PRODUKSI SUPERVISOR"`
	Username   string `json:"username" validate:"required,printascii"`
	Password   string `json:"password" validate:"required,printascii,min=6"`
	NoTelp     string `json:"no_telp" validate:"required,e164"`
	Alamat     string `json:"alamat" validate:"required,printascii"`
	JenisSpvID string `json:"jenis_spv_id" validate:"required_if=Role SUPERVISOR,omitempty,uuidv4_no_hyphens"`
}

type UpdateUser struct {
	ID         string `params:"id" validate:"required,uuidv4_no_hyphens"`
	Nama       string `json:"nama" validate:"omitempty,printascii"`
	Role       string `json:"role" validate:"omitempty,oneof=DIREKTUR ADMIN BENDAHARA MANAJER_PRODUKSI SUPERVISOR"`
	Username   string `json:"username" validate:"omitempty,printascii"`
	Password   string `json:"password" validate:"omitempty,printascii,min=6"`
	NoTelp     string `json:"no_telp" validate:"omitempty,e164"`
	Alamat     string `json:"alamat" validate:"omitempty,printascii"`
	JenisSpvID string `json:"jenis_spv_id" validate:"required_if=Role SUPERVISOR,omitempty,uuidv4_no_hyphens"`
}

type SearchFilterUser struct {
	Nama           string `json:"nama" validate:"omitempty,printascii"`
	Role           string `json:"role" validate:"omitempty,oneof=DIREKTUR ADMIN BENDAHARA MANAJER_PRODUKSI SUPERVISOR"`
	Username       string `json:"username" validate:"omitempty,printascii"`
	Alamat         string `json:"alamat" validate:"omitempty,printascii"`
	NoTelp         string `json:"no_telp" validate:"omitempty,e164"`
	JenisSpvID     string `json:"jenis_spv_id" validate:"required_if=GetAllJenisSpv true,omitempty,uuidv4_no_hyphens"`
	GetAllJenisSpv bool   `json:"jenis_spv" validate:"omitempty,boolean"`
}

type GetAllUser struct {
	Search SearchFilterUser `json:"search" validate:"omitempty"`
	Page   int              `query:"page" validate:"omitempty,number"`
	Limit  int              `query:"limit" validate:"omitempty,number"`
}

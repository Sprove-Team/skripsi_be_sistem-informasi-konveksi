package req_user_jenis_spv

type Create struct {
	Nama string `json:"nama" validate:"required"`
}

type Update struct {
	ID   string `prams:"id" vaidate:"required,ulid"`
	Nama string `json:"nama" validate:"omitempty"`
}

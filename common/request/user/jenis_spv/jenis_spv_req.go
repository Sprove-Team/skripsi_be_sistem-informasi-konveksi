package user

type Create struct {
	Nama string `json:"nama" validate:"required,printascii"`
}

type Update struct {
	ID   string `prams:"id" vaidate:"required,uuidv4_no_hyphens"`
	Nama string `json:"nama" validate:"omitempty,printascii"`
}

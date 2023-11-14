package user

type CreateJenisSpv struct {
	Nama string `json:"nama" validate:"required,printascii"`
}

package req_profile

type Update struct {
	Nama         string `json:"nama" validate:"omitempty"`
	Username     string `json:"username" validate:"omitempty"`
	PasswordLama string `json:"old_password" validate:"required_with=PasswordBaru,omitempty,min=6"`
	PasswordBaru string `json:"new_password" validate:"required_with=PasswordLama,omitempty,nefield=PasswordLama,min=6"`
	NoTelp       string `json:"no_telp" validate:"omitempty,e164"`
	Alamat       string `json:"alamat" validate:"omitempty"`
}

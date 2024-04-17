package req_profile

type Update struct {
	Nama        string `json:"nama" validate:"omitempty"`
	Username    string `json:"username" validate:"omitempty"`
	OldPassword string `json:"old_password" validate:"required_with=NewPassword,omitempty,min=6"`
	NewPassword string `json:"new_password" validate:"required_with=OldPassword,omitempty,nefield=OldPassword,min=6"`
	NoTelp      string `json:"no_telp" validate:"omitempty,e164"`
	Alamat      string `json:"alamat" validate:"omitempty"`
}

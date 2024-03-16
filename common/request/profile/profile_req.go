package req_profile

type Update struct {
	Nama        string `json:"nama" validate:"omitempty"`
	Username    string `json:"username" validate:"omitempty"`
	OldPassword string `json:"old_password" validate:"omitempty,required_with=NewPassword,nefield=NewPassword"`
	NewPassword string `json:"new_password" validate:"omitempty,required_with=OldPassword,nefield=OldPassword"`
	NoTelp      string `json:"no_telp" validate:"omitempty,e164"`
	Alamat      string `json:"alamat" validate:"omitempty"`
}

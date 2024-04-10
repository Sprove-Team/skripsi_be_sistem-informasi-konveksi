package res_user

type DataGetAllUserRes struct {
	ID       string      `json:"id"`
	Nama     string      `json:"nama"`
	Role     string      `json:"role"`
	Username string      `json:"username"`
	NoTelp   string      `json:"no_telp"`
	Alamat   string      `json:"alamat"`
	JenisSpv interface{} `json:"jenis_spv,omitempty"`
}

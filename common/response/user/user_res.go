package res_user

import "github.com/be-sistem-informasi-konveksi/entity"

type DataGetUserRes struct {
	entity.User
	TotalTugas uint `json:"total_tugas,omitempty"`
}

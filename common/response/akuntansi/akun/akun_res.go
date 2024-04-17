package res_akun

import "github.com/be-sistem-informasi-konveksi/entity"

type GetAll struct {
	entity.Akun
	Default bool `json:"default,omitempty"`
}

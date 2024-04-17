package res_kelompok_akun

import "github.com/be-sistem-informasi-konveksi/entity"

type GetAll struct {
	entity.KelompokAkun
	Default bool `json:"default,omitempty"`
}
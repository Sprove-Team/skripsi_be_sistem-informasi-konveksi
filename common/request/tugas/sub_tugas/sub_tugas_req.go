package req_sub_tugas

type CreateByTugasId struct {
	TugasID   string `params:"tugas_id" validate:"required,ulid"`
	Nama      string `json:"nama" validate:"required"`
	Status    string `json:"status" validate:"omitempty,oneof=BELUM_DIKERJAKAN DIPROSES SELESAI"`
	Deskripsi string `json:"deskripsi" validate:"omitempty"`
}

type Update struct {
	ID        string `params:"id" validate:"omitempty"`
	Nama      string `json:"nama" validate:"omitempty"`
	Status    string `json:"status" validate:"omitempty,oneof=BELUM_DIKERJAKAN DIPROSES SELESAI"`
	Deskripsi string `json:"deskripsi" validate:"omitempty"`
}

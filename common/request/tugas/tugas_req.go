package req_tugas

type Create struct {
	InvoiceID       string   `json:"invoice_id" validate:"required,ulid"`
	JenisSpvID      string   `json:"jenis_spv_id" validate:"required,ulid"`
	TanggalDeadline string   `json:"tanggal_deadline" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	UserID          []string `json:"user_id" validate:"gt=0,dive,ulid"`
}

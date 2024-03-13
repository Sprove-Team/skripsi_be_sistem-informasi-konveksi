package req_tugas

type Create struct {
	InvoiceID       string   `json:"invoice_id" validate:"required,ulid"`
	JenisSpvID      string   `json:"jenis_spv_id" validate:"required,ulid"`
	TanggalDeadline string   `json:"tanggal_deadline" validate:"required,datetime=2006-01-02"`
	UserID          []string `json:"user_id" validate:"required,min=1,dive,ulid"`
}

type GetByInvoiceId struct {
	InvoiceID string `params:"invoice_id" validate:"required,ulid"`
}

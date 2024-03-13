package entity

import "time"

type Tugas struct {
	Base
	InvoiceID       string     `gorm:"type:varchar(26);index:idx_invoice_id;default:null" json:"-"`
	JenisSpvID      string     `gorm:"type:varchar(26);index:idx_jenis_spv_id;default:null" json:"-"`
	Invoice         *Invoice   `json:"invoice,omitempty"`
	JenisSpv        *JenisSpv  `json:"jenis_spv"`
	TanggalDeadline *time.Time `gorm:"type:datetime(3)" json:"tanggal_deadline"`
	Users           []User     `gorm:"many2many:user_tugas" json:"users"`
}

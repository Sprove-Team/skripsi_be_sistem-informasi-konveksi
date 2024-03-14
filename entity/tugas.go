package entity

import "time"

type Tugas struct {
	Base
	InvoiceID       string     `gorm:"type:varchar(26);index:idx_invoice_id;default:null" json:"-"`
	JenisSpvID      string     `gorm:"type:varchar(26);index:idx_jenis_spv_id;default:null" json:"-"`
	Invoice         *Invoice   `json:"invoice,omitempty"`
	JenisSpv        *JenisSpv  `json:"jenis_spv,omitempty"`
	TanggalDeadline *time.Time `gorm:"type:datetime(3)" json:"tanggal_deadline,omitempty"`
	Users           []User     `gorm:"many2many:user_tugas;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"users,omitempty"`
	SubTugas        []SubTugas `gorm:"foreignKey:TugasID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sub_tugas,omitempty"`
}

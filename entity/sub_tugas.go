package entity

type SubTugas struct {
	Base
	TugasID   string `gorm:"type:varchar(26);index:idx_invoice_id;default:null" json:"-"`
	Tugas     *Tugas `json:"tugas,omitempty"`
	Nama      string `gorm:"type:varchar(150)" json:"nama"`
	Deskripsi string `gorm:"type:TEXT;default:null" json:"deskripsi,omitempty"`
	Status    string `gorm:"type:enum('menunggu','diproses','selesai')" json:"status"`
}

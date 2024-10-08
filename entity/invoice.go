package entity

import "time"

type Invoice struct {
	Base
	UserID           string             `gorm:"type:varchar(26);index:idx_user_id;default:null" json:"-"`
	KontakID         string             `gorm:"type:varchar(26);index:idx_kontak_id;default:null" json:"-"`
	StatusProduksi   string             `gorm:"type:enum('BELUM_DIKERJAKAN','DIPROSES','SELESAI');default:'BELUM_DIKERJAKAN'" json:"status_produksi,omitempty"`
	NomorReferensi   string             `gorm:"type:varchar(5);uniqueIndex;not null" json:"nomor_referensi,omitempty"`
	TotalQty         int                `gorm:"type:MEDIUMINT unsigned" json:"total_qty,omitempty"`
	TotalHarga       float64            `gorm:"type:decimal(10,2)" json:"total_harga,omitempty"`
	Keterangan       string             `gorm:"type:longtext" json:"keterangan,omitempty"`
	TanggalDeadline  *time.Time         `gorm:"type:datetime(3)" json:"tanggal_deadline,omitempty"`
	TanggalKirim     *time.Time         `gorm:"type:datetime(3)" json:"tanggal_kirim,omitempty"`
	HutangPiutang    *HutangPiutang      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"hutang_piutang,omitempty"`
	Kontak           *Kontak            `json:"kontak,omitempty"`
	User             *User              `json:"user_editor,omitempty"`
	DetailInvoice    []DetailInvoice    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"detail_invoice,omitempty"`
	DataBayarInvoice []DataBayarInvoice `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"data_bayar,omitempty"`
	Tugas            []Tugas            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

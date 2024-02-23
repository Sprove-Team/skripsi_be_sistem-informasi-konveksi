package entity

import "time"

type Invoice struct {
	Base
	UserID           string             `gorm:"type:varchar(26);index:idx_user_id;default:null" json:"-"`
	KontakID         string             `gorm:"type:varchar(26);index:idx_kontak_id;default:null" json:"-"`
	StatusProduksi   string             `gorm:"type:enum('BELUM_DIKERJAKAN','DIPROSES','SELESAI');default:'BELUM_DIKERJAKAN'" json:"status_produksi"`
	NomorReferensi   string             `gorm:"type:varchar(5);uniqueIndex;not null" json:"nomor_referensi"`
	TotalQty         int                `gorm:"type:MEDIUMINT unsigned" json:"total_qty"`
	TotalHarga       float64            `gorm:"type:decimal(10,2)" json:"total_harga"`
	Keterangan       string             `gorm:"type:longtext" json:"keterangan,omitempty"`
	TanggalDeadline  *time.Time         `gorm:"type:datetime(3)" json:"tanggal_deadline"`
	TanggalKirim     *time.Time         `gorm:"type:datetime(3)" json:"tanggal_kirim"`
	HutangPiutang    HutangPiutang      `gorm:"foreignKey:InvoiceID;references:ID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE" json:"-"`
	Kontak           *Kontak            `json:"kontak"`
	User             *User              `json:"user_editor"`
	DetailInvoice    []DetailInvoice    `gorm:"foreignKey:InvoiceID;references:ID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE" json:"detail_invoice,omitempty"`
	DataBayarInvoice []DataBayarInvoice `gorm:"foreignKey:InvoiceID;references:ID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE" json:"data_bayar,omitempty"`
}

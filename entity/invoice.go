package entity

import "time"

type Invoice struct {
	Base
	UserID           string          `gorm:"type:varchar(26);index:idx_user_id;default:null" json:"user_id"`
	StatusProduksiID string          `gorm:"type:varchar(26);index:idx_status_produksi_id;default:null" json:"status_produksi_id"`
	NomorReferensi   string          `gorm:"type:varchar(5);uniqueIndex;not null" json:"nomor_referensi"`
	Kepada           string          `gorm:"type:varchar(150);not null;default:''" json:"kepada"`
	NoTelp           string          `gorm:"type:varchar(50)" json:"no_telp,omitempty"`
	Alamat           string          `gorm:"type:varchar(225)" json:"alamat,omitempty"`
	StatusBayar      string          `gorm:"type:enum('LUNAS','BELUM_LUNAS');default:'BELUM_LUNAS'" json:"status_bayar"`
	TotalQty         int             `gorm:"type:MEDIUMINT unsigned" json:"total_qty"`
	TotalHarga       float64         `gorm:"type:decimal(10,2)" json:"total_harga"`
	SisaTagihan      float64         `gorm:"type:decimal(10,2)" json:"sisa_tagihan"`
	Keterangan       string          `gorm:"type:longtext" json:"keterangan,omitempty"`
	BuktiPembayaran  string          `gorm:"type:varchar(225)" json:"bukti_pembayaran"`
	TanggalDeadline  *time.Time      `gorm:"type:datetime(3)" json:"tanggal_deadline"`
	TanggalKirim     *time.Time      `gorm:"type:datetime(3)" json:"tanggal_kirim"`
	DetailInvoices   []DetailInvoice `gorm:"foreignKey:InvoiceID;references:ID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE" json:"detail_invoice,omitempty"`
}

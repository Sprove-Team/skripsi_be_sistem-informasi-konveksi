package entity

type DataBayarInvoice struct {
	Base
	InvoiceID       string          `gorm:"type:varchar(26);not null;index:idx_invoice_id" json:"invoice_id,omitempty"`
	Invoice         *Invoice        `json:"invoice,omitempty"`
	AkunID          string          `gorm:"type:varchar(26);not null;index:idx_akun_id" json:"-"`
	Akun            *Akun           `json:"akun,omitempty"`
	Keterangan      string          `gorm:"type:longtext" json:"keterangan,omitempty"`
	BuktiPembayaran BuktiPembayaran `gorm:"serializer:json" json:"bukti_pembayaran,omitempty"`
	Total           float64         `gorm:"type:decimal(10,2);default:0" json:"total,omitempty"`
	Status          string          `gorm:"type:enum('TERKONFIRMASI','BELUM_TERKONFIRMASI');default:'BELUM_TERKONFIRMASI'" json:"status,omitempty"`
}

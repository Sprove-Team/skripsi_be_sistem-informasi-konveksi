package entity

type DataBayarInvoice struct {
	Base
	InvoiceID       string          `gorm:"type:varchar(26);not null;index:idx_invoice_id" json:"invoice_id"`
	AkunID          string          `gorm:"type:varchar(26);not null;index:idx_akun_id" json:"akun_id"`
	Akun            *Akun           `json:"akun"`
	Keterangan      string          `gorm:"type:longtext" json:"keterangan"`
	BuktiPembayaran BuktiPembayaran `gorm:"serializer:json" json:"bukti_pembayaran"`
	Total           float64         `gorm:"type:decimal(10,2);default:0" json:"total"`
	Status          string          `gorm:"type('TERKONFIRMASI','BELUM_TERKONFIRMASI');default:'BELUM_TERKONFIRMASI'" json:"status"`
}

func (DataBayarInvoice) TableName() string {
	return "data_bayar_invoice"
}

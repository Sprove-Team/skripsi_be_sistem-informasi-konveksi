package entity

type DetailInvoice struct {
	Base
	ProdukID     string  `gorm:"type:varchar(26);index:idx_produk_id;default:null" json:"produk_id"`
	InvoiceID    string  `gorm:"type:varchar(26);index:idx_invoice_id;default:null" json:"invoice_id"`
	BordirID     string  `gorm:"type:varchar(26);index:idx_bordir_id;default:null" json:"bordir_id"`
	SablonID     string  `gorm:"type:varchar(26);index:idx_sablon_id;default:null" json:"sablon_id"`
	GambarDesign string  `gorm:"type:varchar(225)" json:"gambar_design"`
	Qty          int     `gorm:"type:MEDIUMINT" json:"qty"`
	Total        float64 `gorm:"type:decimal(10,2)" json:"total"`
}

func (DetailInvoice) TableName() string {
	return "detail_invoice"
}

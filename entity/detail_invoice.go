package entity

type DetailInvoice struct {
	Base
	ProdukID     string  `gorm:"type:varchar(26);index:idx_produk_id;default:null" json:"-"`
	BordirID     string  `gorm:"type:varchar(26);index:idx_bordir_id;default:null" json:"-"`
	SablonID     string  `gorm:"type:varchar(26);index:idx_sablon_id;default:null" json:"-"`
	InvoiceID    string  `gorm:"type:varchar(26);index:idx_invoice_id;default:null" json:"invoice_id,omitempty"`
	Produk       *Produk `json:"produk,omitempty"`
	Sablon       *Sablon `json:"sablon,omitempty"`
	Bordir       *Bordir `json:"bordir,omitempty"`
	GambarDesign string  `gorm:"type:varchar(225)" json:"gambar_design,omitempty"`
	Qty          int     `gorm:"type:MEDIUMINT" json:"qty,omitempty"`
	Total        float64 `gorm:"type:decimal(10,2)" json:"total,omitempty"`
}

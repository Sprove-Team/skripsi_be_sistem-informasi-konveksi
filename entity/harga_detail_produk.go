package entity

type HargaDetailProduk struct {
	Base
	ProdukID string  `gorm:"type:varchar(26);index:idx_produk_id;not null" json:"produk_id,omitempty"`
	QTY      uint    `gorm:"type:int(10) unsigned;not null;" json:"qty,omitempty"`
	Harga    float64 `gorm:"type:decimal(10, 2);" json:"harga,omitempty"`
}

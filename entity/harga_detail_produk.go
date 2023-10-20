package entity

type HargaDetailProduk struct {
	ProdukID string  `gorm:"type:varchar(32);primaryKey;uniqueIndex;not null" json:"produk_id,omitempty"`
	ID       uint    `gorm:"primaryKey;not null" json:"id"`
	QTY      uint    `gorm:"type:int(10) unsigned;uniqueIndex;not null;" json:"qty"`
	Harga    float64 `gorm:"type:decimal(10, 2);" json:"harga"`
}

func (HargaDetailProduk) TableName() string {
	return "harga_detail_produk"
}

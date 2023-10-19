package entity

type HargaDetailProduk struct {
	ProdukID string  `gorm:"type:varchar(36)" json:"produk_id"`
	ID       uint    `gorm:"primaryKey;not null" json:"id"`
	QTY      uint    `gorm:"default:0;uniqueIndex" json:"qty"`
	Harga    float64 `gorm:"type:decimal(10, 2);" json:"harga"`
}

func (HargaDetailProduk) TableName() string {
	return "harga_detail_produk"
}

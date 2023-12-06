package entity

type HargaDetailProduk struct {
	ProdukID string  `gorm:"type:varchar(26);index:idx_produk_id;not null" json:"produk_id"`
	ID       string  `gorm:"type:varchar(26);primaryKey;uniqueIndex;not null" json:"id"`
	QTY      uint    `gorm:"type:int(10) unsigned;not null;" json:"qty"`
	Harga    float64 `gorm:"type:decimal(10, 2);" json:"harga"`
}

func (HargaDetailProduk) TableName() string {
	return "harga_detail_produk"
}

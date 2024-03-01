package entity

type HargaDetailProduk struct {
	Base
	ProdukID string  `gorm:"type:varchar(26);index:idx_produk_id;not null" json:"produk_id"`
	QTY      uint    `gorm:"type:int(10) unsigned;uniqueIndex;not null;" json:"qty"`
	Harga    float64 `gorm:"type:decimal(10, 2);" json:"harga"`
}

func (HargaDetailProduk) TableName() string {
	return "harga_detail_produk"
}

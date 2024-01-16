package entity

type Produk struct {
	Base
	KategoriProdukID string              `gorm:"type:varchar(26);index:idx_kategori_produk_id;not null" json:"kategori_id"`
	Nama             string              `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama"`
	HargaDetails     []HargaDetailProduk `gorm:"foreignKey:ProdukID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"harga_detail,omitempty"`
}

func (Produk) TableName() string {
	return "produk"
}

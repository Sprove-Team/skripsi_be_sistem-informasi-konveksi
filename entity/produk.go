package entity

type Produk struct {
	ID               string              `gorm:"type:varchar(32);primaryKey;uniqueindex;not null" json:"id"`
	Nama             string              `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama"`
	HargaDetails     []HargaDetailProduk `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"harga_details"`
	KategoriProdukID uint                `json:"kategori_id"`
}

func (Produk) TableName() string {
	return "produk"
}

package entity

type KategoriProduk struct {
	Nama    string   `gorm:"type:varchar(150);unique;not null"`
	Produks []Produk `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"produk"`
	ID      uint     `gorm:"primaryKey;not null"`
}

func (KategoriProduk) TableName() string {
	return "kategori_produk"
}

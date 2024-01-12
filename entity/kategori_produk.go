package entity

type KategoriProduk struct {
	Base
	Nama    string   `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama"`
	Produks []Produk `gorm:"foreignKey:KategoriProdukID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (KategoriProduk) TableName() string {
	return "kategori_produk"
}

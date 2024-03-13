package entity

type KategoriProduk struct {
	Base
	Nama    string   `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama,omitempty"`
	Produks []Produk `gorm:"foreignKey:KategoriProdukID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

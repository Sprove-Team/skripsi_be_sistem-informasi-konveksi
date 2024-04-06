package entity

type KategoriProduk struct {
	Base
	Nama    string   `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama,omitempty"`
	Produks []Produk `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

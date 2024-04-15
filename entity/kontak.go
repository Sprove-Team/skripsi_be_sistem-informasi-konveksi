package entity

type Kontak struct {
	Base
	Nama       string      `gorm:"type:varchar(150);not null" json:"nama,omitempty"`
	NoTelp     string      `gorm:"type:varchar(50)" json:"no_telp,omitempty"`
	Alamat     string      `gorm:"type:varchar(255)" json:"alamat,omitempty"`
	Email      string      `gorm:"type:varchar(255);default=''" json:"email,omitempty"`
	Keterangan string      `gorm:"type:longtext" json:"keterangan,omitempty"`
	Transaksi  []Transaksi `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"transaksi,omitempty"`
	Invoice    []Invoice   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"invoice,omitempty"`
}

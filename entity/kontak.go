package entity

type Kontak struct {
	Base
	Nama       string      `gorm:"type:varchar(150);not null" json:"nama,omitempty"`
	NoTelp     string      `gorm:"type:varchar(50)" json:"no_telp,omitempty"`
	Alamat     string      `gorm:"type:varchar(225)" json:"alamat,omitempty"`
	Email      string      `gorm:"type:varchar(225);default=''" json:"email,omitempty"`
	Keterangan string      `gorm:"type:longtext" json:"keterangan,omitempty"`
	Transaksi  []Transaksi `gorm:"foreignKey:KontakID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"transaksi,omitempty"`
}

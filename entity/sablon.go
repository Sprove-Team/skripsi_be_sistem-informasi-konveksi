package entity

type Sablon struct {
	ID    string  `gorm:"type:varchar(26);primaryKey;index:idx_sablon_id;not null" json:"id"`
	Nama  string  `gorm:"type:varchar(150);unique;not null" json:"nama"`
	Harga float64 `gorm:"type:decimal(10,2);not null" harga:"harga"`
}
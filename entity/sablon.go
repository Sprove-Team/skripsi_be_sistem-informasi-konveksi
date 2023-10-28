package entity

type Sablon struct {
	ID    string  `gorm:"primaryKey;index:idx_sablon_id;not null"`
	Nama  string  `gorm:"type:varchar(32);unique;not null"`
	Harga float64 `gorm:"type:decimal(10,2);not null"`
}
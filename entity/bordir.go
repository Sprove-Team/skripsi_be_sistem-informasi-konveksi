package entity

type Bordir struct {
	ID    string  `gorm:"type:varchar(26);primaryKey;index:idx_bordir_id;not null" json:"id"`
	Nama  string  `gorm:"type:varchar(32);unique;not null" json:"nama"`
	Harga float64 `gorm:"type:decimal(10,2);not null" json:"harga"`
}

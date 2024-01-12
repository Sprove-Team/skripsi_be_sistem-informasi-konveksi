package entity

type Bordir struct {
	Base
	Nama           string          `gorm:"type:varchar(32);unique;not null" json:"nama"`
	Harga          float64         `gorm:"type:decimal(10,2);not null" json:"harga"`
	DetailInvoices []DetailInvoice `gorm:"foreignKey:BordirID;references:ID"`
}

func (Bordir) TableName() string {
	return "bordir"
}

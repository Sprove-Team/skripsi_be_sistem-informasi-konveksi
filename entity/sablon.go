package entity

type Sablon struct {
	Base
	Nama          string          `gorm:"type:varchar(150);unique;not null" json:"nama"`
	Harga         float64         `gorm:"type:decimal(10,2);not null" json:"harga"`
	DetailInvoice []DetailInvoice `gorm:"foreignKey:SablonID;references:ID" json:"detail_invoice,omitempty"`
}

func (Sablon) TableName() string {
	return "sablon"
}

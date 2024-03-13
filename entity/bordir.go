package entity

type Bordir struct {
	Base
	Nama          string          `gorm:"type:varchar(32);unique;not null" json:"nama,omitempty"`
	Harga         float64         `gorm:"type:decimal(10,2);not null" json:"harga,omitempty"`
	DetailInvoice []DetailInvoice `gorm:"foreignKey:BordirID;references:ID" json:"detail_invoice,omitempty"`
}

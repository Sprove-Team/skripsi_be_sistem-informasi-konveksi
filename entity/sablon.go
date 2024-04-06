package entity

type Sablon struct {
	Base
	Nama          string          `gorm:"type:varchar(150);unique;not null" json:"nama,omitempty"`
	Harga         float64         `gorm:"type:decimal(10,2);not null" json:"harga,omitempty"`
	DetailInvoice []DetailInvoice `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"detail_invoice,omitempty"`
}

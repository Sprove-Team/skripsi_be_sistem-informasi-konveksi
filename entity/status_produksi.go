package entity

type StatusProduksi struct {
	Base
	Nama     string    `gorm:"type:varchar(50)" json:"nama"`
	Invoices []Invoice `gorm:"foreignKey:StatusProduksiID;references:ID" json:"invoice"`
}

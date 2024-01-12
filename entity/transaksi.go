package entity

import (
	"time"
)

type Transaksi struct {
	Base
	Keterangan      string       `gorm:"type:longtext" json:"keterangan,omitempty"`
	BuktiPembayaran string       `gorm:"type:varchar(255)" json:"bukti_pembayaran,omitempty"`
	Total           float64      `gorm:"type:decimal(10,2);default:0" json:"total"`
	Tanggal         time.Time    `gorm:"type:datetime(3)" json:"tanggal"`
	AyatJurnals     []AyatJurnal `gorm:"foreignKey:TransaksiID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ayat_jurnal,omitempty"`
}

func (Transaksi) TableName() string {
	return "transaksi"
}

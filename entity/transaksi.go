package entity

import (
	"time"
)

// BuktiPembayaran string       `gorm:"type:varchar(255)" json:"bukti_pembayaran,omitempty"`
type Transaksi struct {
	Base
	Keterangan      string          `gorm:"type:longtext" json:"keterangan,omitempty"`
	BuktiPembayaran BuktiPembayaran `gorm:"serializer:json" json:"bukti_pembayaran"`
	Total           float64         `gorm:"type:decimal(10,2);default:0" json:"total"`
	Tanggal         time.Time       `gorm:"type:datetime(3)" json:"tanggal"`
	KontakID        string          `gorm:"type:varchar(26);index:idx_kontak_id;default:null" json:"kontak_id,omitempty"`
	Kontak          *Kontak         `json:"kontak,omitempty"`
	AyatJurnals     []AyatJurnal    `gorm:"foreignKey:TransaksiID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ayat_jurnal,omitempty"`
}

func (Transaksi) TableName() string {
	return "transaksi"
}

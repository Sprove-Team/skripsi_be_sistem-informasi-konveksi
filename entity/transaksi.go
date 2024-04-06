package entity

import (
	"time"
)

// note: siapa yang megang foreign key itu child, sedangkan yang tidak itu parentnya

type Transaksi struct {
	BaseSoftDelete
	Keterangan             string                  `gorm:"type:longtext" json:"keterangan,omitempty"`
	BuktiPembayaran        BuktiPembayaran         `gorm:"serializer:json" json:"bukti_pembayaran,omitempty"`
	Total                  float64                 `gorm:"type:decimal(10,2);default:0" json:"total,omitempty"`
	Tanggal                time.Time               `gorm:"type:datetime(3)" json:"tanggal,omitempty"`
	KontakID               string                  `gorm:"type:varchar(26);index:idx_kontak_id;default:null" json:"-"`
	Kontak                 *Kontak                 `json:"kontak,omitempty"`
	HutangPiutang          *HutangPiutang          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	DataBayarHutangPiutang *DataBayarHutangPiutang `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	AyatJurnals            []AyatJurnal            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"ayat_jurnal,omitempty"`
}

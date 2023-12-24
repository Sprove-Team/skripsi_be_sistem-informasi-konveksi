package entity

import (
	"time"

	"gorm.io/gorm"
)

type Transaksi struct {
	ID              string          `gorm:"type:varchar(26);primaryKey;index:idx_akun_id;not null" json:"id"`
	Keterangan      string          `gorm:"type:longtext" json:"keterangan"`
	BuktiPembayaran string          `gorm:"type:varchar(255)" json:"bukti_pembayaran"`
	Total           float64         `gorm:"type:decimal(10,2);default:0" json:"total"`
	Tanggal         time.Time       `gorm:"type:datetime(3)" json:"tanggal"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"-"`
	DeletedAt       *gorm.DeletedAt `gorm:"index" json:"-"`
	AyatJurnals     []AyatJurnal    `gorm:"foreignKey:TransaksiID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ayat_jurnal"`
}

func (Transaksi) TableName() string {
	return "transaksi"
}

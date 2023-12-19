package entity

import (
	"time"

	"gorm.io/gorm"
)

type Transaksi struct {
	ID              string          `gorm:"type:varchar(26);primaryKey;index:idx_akun_id;not null" json:"id"`
	Tanggal         time.Time       `gorm:"type:date"`
	Keterangan      string          `gorm:"type:longtext" json:"keterangan"`
	BuktiPembayaran string          `gorm:"type:varchar(255)" json:"bukti_pembayaran"`
	CreatedAt       time.Time       `json:"dibuat_pada"`
	UpdatedAt       time.Time       `json:"diedit_pada"`
	DeletedAt       *gorm.DeletedAt `gorm:"index" json:"dihapus_pada"`
	AyatJurnals     []AyatJurnal    `gorm:"foreignKey:TransaksiID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ayat_jurnal"`
}

func (Transaksi) TableName() string {
	return "transaksi"
}

package entity

import (
	"time"

	"gorm.io/gorm"
)

type AyatJurnal struct {
	ID          string          `gorm:"type:varchar(26);primaryKey;index:idx_akun_id;not null" json:"id"`
	TransaksiID string          `gorm:"type:varchar(26);index:idx_transaksi_id;not null" json:"transaksi_id"`
	AkunID      string          `gorm:"type:varchar(26);index:idx_akun_id;not null" json:"akun_id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Debit       float64         `gorm:"type:decimal(10,2);default:0" json:"debit"`
	Saldo       float64         `gorm:"type:decimal(10,2);default:0" json:"saldo"`
	Kredit      float64         `gorm:"type:decimal(10,2);default:0" json:"kredit"`
}

func (AyatJurnal) TableName() string {
	return "ayat_jurnal"
}

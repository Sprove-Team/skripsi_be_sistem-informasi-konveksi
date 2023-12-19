package entity

import (
	"time"

	"gorm.io/gorm"
)

var SaldoNormal = map[string]bool{
	"DEBIT":  true,
	"KREDIT": true,
}

type Akun struct {
	ID             string          `gorm:"type:varchar(26);primaryKey;index:idx_akun_id;not null" json:"id"`
	GolonganAkunID string          `gorm:"type:varchar(26);index:idx_golongan_akun_id;not null" json:"golongan_akun_id"`
	Nama           string          `gorm:"type:varchar(150);not null" json:"nama"`
	Kode           string          `gorm:"type:varchar(10);uniqueIndex;not null" json:"kode"`
	SaldoNormal    string          `gorm:"type:enum('DEBIT','KREDIT');default:'DEBIT';not null" json:"saldo_normal"`
	SaldoAkhir     float64         `gorm:"type:decimal(10,2);default:0" json:"saldo_akhir"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Akun) TableName() string {
	return "akun"
}

package entity

import (
	"time"

	"gorm.io/gorm"
)

type GolonganAkun struct {
	ID             string          `gorm:"type:varchar(26);primaryKey;index:idx_akun_id;not null" json:"id"`
	KelompokAkunID string          `gorm:"type:varchar(26);index:idx_kelompok_akun_id;not null" json:"kelompok_akun_id"`
	Kode           string          `gorm:"type:varchar(10);uniqueIndex;not null" json:"kode"`
	Nama           string          `gorm:"type:varchar(150);not null" json:"nama"`
	CreatedAt      time.Time       `json:"dibuat_pada"`
	UpdatedAt      time.Time       `json:"diedit_pada"`
	DeletedAt      *gorm.DeletedAt `gorm:"index" json:"dihapus_pada"`
	Akuns          []Akun          `gorm:"foreignKey:GolonganAkunID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"akun"`
}

func (GolonganAkun) TableName() string {
	return "golongan_akun"
}

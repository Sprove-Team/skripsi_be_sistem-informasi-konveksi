package entity

import (
	"time"

	"gorm.io/gorm"
)

type KelompokAkun struct {
	ID            string          `gorm:"type:varchar(26);primaryKey;index:idx_akun_id;not null" json:"id"`
	Kode          string          `gorm:"type:varchar(10);uniqueIndex;not null" json:"kode,omitempty"`
	Nama          string          `gorm:"type:varchar(150);not null" json:"nama"`
	CreatedAt     *time.Time      `json:"created_at,omitempty"`
	UpdatedAt     *time.Time      `json:"updated_at,omitempty"`
	DeletedAt     *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	GolonganAkuns []GolonganAkun  `gorm:"foreignKey:KelompokAkunID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"golongan_akun,omitempty"`
}

func (KelompokAkun) TableName() string {
	return "kelompok_akun"
}

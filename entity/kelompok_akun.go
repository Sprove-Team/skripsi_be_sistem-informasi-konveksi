package entity

import (
	"time"

	"gorm.io/gorm"
)

type KelompokAkun struct {
	ID            string          `gorm:"type:varchar(26);primaryKey;index:idx_akun_id;not null" json:"id"`
	Kode          string          `gorm:"type:varchar(10);uniqueIndex;not null" json:"kode"`
	Nama          string          `gorm:"type:varchar(150);not null" json:"nama"`
	CreatedAt     time.Time       `json:"dibuat_pada"`
	UpdatedAt     time.Time       `json:"diedit_pada"`
	DeletedAt     *gorm.DeletedAt `gorm:"index" json:"dihapus_pada"`
	GolonganAkuns []GolonganAkun  `gorm:"foreignKey:KelompokAkunID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"golongan_akun"`
}

func (KelompokAkun) TableName() string {
	return "kelompok_akun"
}

package entity

import (
	"time"

	"gorm.io/gorm"
)

type GolonganAkun struct {
	ID             string          `gorm:"type:varchar(26);primaryKey;index:idx_akun_id;not null" json:"id"`
	KelompokAkunID string          `gorm:"type:varchar(26);index:idx_kelompok_akun_id;not null" json:"-"`
	KelompokAkun   *KelompokAkun   `json:"kelompok_akun,omitempty"`
	Kode           string          `gorm:"type:varchar(10);uniqueIndex;not null" json:"kode,omitempty"`
	Nama           string          `gorm:"type:varchar(150);not null" json:"nama,omitempty"`
	CreatedAt      *time.Time      `json:"created_at,omitempty"`
	UpdatedAt      *time.Time      `json:"updated_at,omitempty"`
	DeletedAt      *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Akuns          []Akun          `gorm:"foreignKey:GolonganAkunID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"akun,omitempty"`
}

func (GolonganAkun) TableName() string {
	return "golongan_akun"
}

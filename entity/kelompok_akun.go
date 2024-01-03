package entity

import (
	"time"

	"gorm.io/gorm"
)

var KategoriAkun = map[string]string{
	"ASET":       "1",
	"KEWAJIBAN":  "2",
	"MODAL":      "3",
	"PENDAPATAN": "4",
	"BEBAN":      "5",
}

var KategoriAkunByKode = map[string]string{
	"1": "ASET",
	"2": "KEWAJIBAN",
	"3": "MODAL",
	"4": "PENDAPATAN",
	"5": "BEBAN",
}

type KelompokAkun struct {
	ID           string          `gorm:"type:varchar(26);primaryKey;index:idx_akun_id;not null" json:"id"`
	Kode         string          `gorm:"type:varchar(10);uniqueIndex;not null" json:"kode,omitempty"`
	Nama         string          `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama"`
	KategoriAkun string          `gorm:"type:enum('ASET','KEWAJIBAN','MODAL','PENDAPATAN','BEBAN')" json:"kategori_akun,omitempty"`
	CreatedAt    *time.Time      `json:"created_at,omitempty"`
	UpdatedAt    *time.Time      `json:"updated_at,omitempty"`
	DeletedAt    *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Akuns        []Akun          `gorm:"foreignKey:KelompokAkunID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"akun,omitempty"`
}

func (KelompokAkun) TableName() string {
	return "kelompok_akun"
}

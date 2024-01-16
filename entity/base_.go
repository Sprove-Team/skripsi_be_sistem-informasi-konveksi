package entity

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        string          `gorm:"type:varchar(26);primaryKey;index:idx_akun_id;not null" json:"id"`
	CreatedAt *time.Time      `json:"-"`
	UpdatedAt *time.Time      `json:"-"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`
}

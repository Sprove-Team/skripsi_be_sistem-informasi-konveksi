package entity

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        string          `gorm:"type:varchar(26);primaryKey;index:idx_akun_id;not null" json:"id"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

package entity

import (
	"time"

	"gorm.io/gorm"
)

type BuktiPembayaran []string
type Base struct {
	ID        string          `gorm:"type:varchar(26);primaryKey;index:idx_id;not null" json:"id"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"-"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`
}

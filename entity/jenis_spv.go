package entity

type JenisSpv struct {
	Base
	Nama  string `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama,omitempty"`
	Tugas Tugas  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Users []User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

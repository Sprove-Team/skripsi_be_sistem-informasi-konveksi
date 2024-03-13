package entity

type JenisSpv struct {
	Base
	Nama  string `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama,omitempty"`
	Tugas Tugas  `gorm:"foreignKey:JenisSpvID;references:ID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE" json:"-"`
	Users []User `gorm:"foreignKey:JenisSpvID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

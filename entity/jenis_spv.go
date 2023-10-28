package entity

type JenisSpv struct {
	ID    string `gorm:"type:varchar(32);primaryKey;uniqueIndex;not null" json:"id"`
	Nama  string `gorm:"type:varchar(150);not null" json:"nama"`
	Users []User `gorm:"foreignKey:JenisSpvID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (JenisSpv) TableName() string {
	return "jenis_spv"
}

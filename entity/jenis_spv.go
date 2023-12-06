package entity

type JenisSpv struct {
	ID    string `gorm:"type:varchar(26);primaryKey;uniqueIndex;not null" json:"id"`
	Nama  string `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama"`
	Users []User `gorm:"foreignKey:JenisSpvID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (JenisSpv) TableName() string {
	return "jenis_spv"
}

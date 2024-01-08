package entity

type Invoice struct {
	ID             string `gorm:"type:varchar(26);primaryKey;uniqueIndex;not null" json:"id"`
	NomorReferensi string `gorm:"type:varchar(5)"`
}

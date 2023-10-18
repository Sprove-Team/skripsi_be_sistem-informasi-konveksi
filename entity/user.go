package entity

type User struct {
	ID       string `gorm:"type:varchar(36);primaryKey;uniqueindex;not null"`
	Nama     string `gorm:"type:varchar(150);not null"`
	Username string `gorm:"type:varchar(150);uniqueindex;not null"`
	Password string `gorm:"type:varchar(64)"`
	NoTelp   string `gorm:"type:varchar(20)"`
	Alamat   string `gorm:"type:varchar(64)"`
}

func (User) TableName() string {
	return "user"
}

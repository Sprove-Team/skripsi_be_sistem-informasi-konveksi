package user

type DataGetAllUserRes struct {
	ID       string `gorm:"type:varchar(36);primaryKey;uniqueIndex;not null" json:"id"`
	Nama     string `gorm:"type:varchar(150);not null" json:"nama"`
	Role     string `gorm:"type:enum('DIREKTUR','BENDAHARA','ADMIN','MANAJER_PRODUKSI','SUPERVISOR')" json:"role"`
	Username string `gorm:"type:varchar(150);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(100);not null" json:"-"`
	NoTelp   string `gorm:"type:varchar(20)" json:"no_telp"`
	Alamat   string `gorm:"type:varchar(150)" json:"alamat"`
}

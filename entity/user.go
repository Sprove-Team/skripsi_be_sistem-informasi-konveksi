package entity

var Roles = map[string]bool{
	"DIREKTUR":         true,
	"BENDAHARA":        true,
	"ADMIN":            true,
	"MANAJER_PRODUKSI": true,
	"SUPERVISOR":       true,
}

type User struct {
	Base
	Nama       string    `gorm:"type:varchar(150);not null" json:"nama"`
	Role       string    `gorm:"type:enum('DIREKTUR','BENDAHARA','ADMIN','MANAJER_PRODUKSI','SUPERVISOR')" json:"role"`
	Username   string    `gorm:"type:varchar(150);uniqueIndex;not null" json:"username"`
	Password   string    `gorm:"type:varchar(100);not null" json:"-"`
	NoTelp     string    `gorm:"type:varchar(20)" json:"no_telp"`
	Alamat     string    `gorm:"type:varchar(150)" json:"alamat"`
	JenisSpvID string    `gorm:"type:varchar(26);index:idx_jenis_spv_id;default:null" json:"jenis_spv_id"`
	JenisSpv   *JenisSpv `json:"jenis_spv,omitempty"`
}

func (User) TableName() string {
	return "user"
}

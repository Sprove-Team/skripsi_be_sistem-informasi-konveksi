package entity

// var SaldoNormal = map[string]bool{
// 	"DEBIT":  true,
// 	"KREDIT": true,
// }

type Akun struct {
	Base
	KelompokAkunID string        `gorm:"type:varchar(26);index:idx_kelompok_akun_id;not null" json:"-"`
	KelompokAkun   *KelompokAkun `json:"kelompok_akun,omitempty"`
	Nama           string        `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama"`
	Kode           string        `gorm:"type:varchar(10);uniqueIndex;not null" json:"kode"`
	SaldoNormal    string        `gorm:"type:enum('DEBIT','KREDIT');default:'DEBIT';not null" json:"saldo_normal"`
	Saldo          float64       `gorm:"type:decimal(10,2);default:0" json:"saldo"`
	Deskripsi      string        `gorm:"type:TEXT;default:null" json:"deskripsi,omitempty"`
	AyatJurnals    []AyatJurnal  `gorm:"foreignKey:AkunID;references:ID" json:"ayat_jurnal,omitempty"`
}

func (Akun) TableName() string {
	return "akun"
}

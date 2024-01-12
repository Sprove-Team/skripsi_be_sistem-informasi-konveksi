package entity

var KategoriAkun = map[string]string{
	"ASET":       "1",
	"KEWAJIBAN":  "2",
	"MODAL":      "3",
	"PENDAPATAN": "4",
	"BEBAN":      "5",
}

var KategoriAkunByKode = map[string]string{
	"1": "ASET",
	"2": "KEWAJIBAN",
	"3": "MODAL",
	"4": "PENDAPATAN",
	"5": "BEBAN",
}

type KelompokAkun struct {
	Base
	Kode         string `gorm:"type:varchar(10);uniqueIndex;not null" json:"kode,omitempty"`
	Nama         string `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama"`
	KategoriAkun string `gorm:"type:enum('ASET','KEWAJIBAN','MODAL','PENDAPATAN','BEBAN')" json:"kategori_akun,omitempty"`
	Akuns        []Akun `gorm:"foreignKey:KelompokAkunID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"akun,omitempty"`
}

func (KelompokAkun) TableName() string {
	return "kelompok_akun"
}

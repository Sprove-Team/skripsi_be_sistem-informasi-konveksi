package entity

type Akun struct {
	Base
	KelompokAkunID   string             `gorm:"type:varchar(26);index:idx_kelompok_akun_id;not null" json:"-"`
	KelompokAkun     *KelompokAkun      `json:"kelompok_akun,omitempty"`
	Nama             string             `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama,omitempty"`
	Kode             string             `gorm:"type:varchar(10);uniqueIndex;not null" json:"kode,omitempty"`
	SaldoNormal      string             `gorm:"type:enum('DEBIT','KREDIT');default:'DEBIT';not null" json:"saldo_normal,omitempty"`
	Deskripsi        string             `gorm:"type:TEXT;default:null" json:"deskripsi,omitempty"`
	AyatJurnal       []AyatJurnal       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"ayat_jurnal,omitempty"`
	DataBayarInvoice []DataBayarInvoice `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"data_bayar_invoice,omitempty"`
}

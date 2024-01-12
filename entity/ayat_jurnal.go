package entity

type AyatJurnal struct {
	Base
	TransaksiID string  `gorm:"type:varchar(26);index:idx_transaksi_id;not null" json:"transaksi_id"`
	AkunID      string  `gorm:"type:varchar(26);index:idx_akun_id;not null" json:"-"`
	Akun        *Akun   `json:"akun"`
	Debit       float64 `gorm:"type:decimal(10,2);default:0" json:"debit"`
	Saldo       float64 `gorm:"type:decimal(10,2);default:0" json:"saldo"`
	Kredit      float64 `gorm:"type:decimal(10,2);default:0" json:"kredit"`
}

func (AyatJurnal) TableName() string {
	return "ayat_jurnal"
}

package entity

type AyatJurnal struct {
	Base
	TransaksiID string  `gorm:"type:varchar(26);index:idx_transaksi_id;not null" json:"transaksi_id,omitempty"`
	AkunID      string  `gorm:"type:varchar(26);index:idx_akun_id;not null" json:"-"`
	Akun        *Akun   `json:"akun,omitempty"`
	Debit       float64 `gorm:"type:decimal(10,2);default:0" json:"debit,omitempty"`
	Saldo       float64 `gorm:"type:decimal(10,2);default:0" json:"saldo,omitempty"`
	Kredit      float64 `gorm:"type:decimal(10,2);default:0" json:"kredit,omitempty"`
}

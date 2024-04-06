package entity

type DataBayarHutangPiutang struct {
	Base
	HutangPiutangID string        `gorm:"type:varchar(26);not null;index:idx_hutang_piutang_id" json:"hutang_piutang_id,omitempty"`
	HutangPiutang   HutangPiutang `json:"-"`
	Total           float64       `gorm:"type:decimal(10,2);default:0" json:"total,omitempty"`
	TransaksiID     string        `gorm:"type:varchar(26);not null;uniqueIndex;" json:"transaksi_id,omitempty"`
	Transaksi       Transaksi     `json:"-"`
}

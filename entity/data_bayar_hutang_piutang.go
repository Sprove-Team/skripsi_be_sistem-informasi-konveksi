package entity

type DataBayarHutangPiutang struct {
	Base
	HutangPiutangID string        `gorm:"type:varchar(26);not null;uniqueIndex:idx_hutang_piutang_id" json:"hutang_piutang_id"`
	HutangPiutang   HutangPiutang `json:"-"`
	TransaksiID     string        `gorm:"type:varchar(26);not null;uniqueIndex:idx_hutang_piutang_id" json:"transaksi_id"`
	Transaksi       Transaksi     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

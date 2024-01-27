package entity

type HutangPiutang struct {
	Base
	InvoiceSlug            string                   `gorm:"type:varchar(26);default:''" json:"invoice_slug,omitempty"`
	Jenis                  string                   `gorm:"type:enum('PIUTANG','UTANG')" json:"jenis"`
	TransaksiID            string                   `json:"transaksi_id"`
	Transaksi              Transaksi                `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Status                 string                   `gorm:"type:enum('BELUM_LUNAS','LUNAS');default:'BELUM_LUNAS'" json:"status"`
	DataBayarHutangPiutang []DataBayarHutangPiutang `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"data_bayar,omitempty"`
}

package entity

type HutangPiutang struct {
	Base
	InvoiceID              string                   `gorm:"type:varchar(26);index:idx_invoice_id,unique;default:null" json:"invoice_id,omitempty"`
	TransaksiID            string                   `json:"transaksi_id"`
	Jenis                  string                   `gorm:"type:enum('PIUTANG','HUTANG')" json:"jenis"`
	Transaksi              Transaksi                `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Status                 string                   `gorm:"type:enum('BELUM_LUNAS','LUNAS');default:'BELUM_LUNAS'" json:"status"`
	Total                  float64                  `gorm:"type:decimal(10,2);default:0" json:"total"`
	Sisa                   float64                  `gorm:"type:decimal(10,2);default:0" json:"sisa"`
	DataBayarHutangPiutang []DataBayarHutangPiutang `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"data_bayar,omitempty"`
}

package entity

import (
	"time"

	"gorm.io/gorm"
)

// BuktiPembayaran string       `gorm:"type:varchar(255)" json:"bukti_pembayaran,omitempty"`
type Transaksi struct {
	Base
	Keterangan             string                  `gorm:"type:longtext" json:"keterangan,omitempty"`
	BuktiPembayaran        BuktiPembayaran         `gorm:"serializer:json" json:"bukti_pembayaran"`
	Total                  float64                 `gorm:"type:decimal(10,2);default:0" json:"total"`
	Tanggal                time.Time               `gorm:"type:datetime(3)" json:"tanggal"`
	KontakID               string                  `gorm:"type:varchar(26);index:idx_kontak_id;default:null" json:"kontak_id,omitempty"`
	Kontak                 *Kontak                 `json:"kontak,omitempty"`
	HutangPiutang          *HutangPiutang          `json:"-"`
	DataBayarHutangPiutang *DataBayarHutangPiutang `json:"-"`
	AyatJurnals            []AyatJurnal            `gorm:"foreignKey:TransaksiID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ayat_jurnal,omitempty"`
}

func (t *Transaksi) BeforeDelete(tx *gorm.DB) (err error) {
	// jika tr termasuk tr bayar, update data hp
	var dhp DataBayarHutangPiutang
	if err := tx.First(&dhp, "transaksi_id = ?", t.ID).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}
	if dhp.ID == "" {
		return
	}

	var hp HutangPiutang
	if err := tx.First(&hp, "id = ?", dhp.HutangPiutangID).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	hp.Sisa += t.Total
	if hp.Sisa > 0 {
		hp.Status = "BELUM_LUNAS"
	}
	if err := tx.Select("sisa", "status").Updates(&hp).Error; err != nil {
		return err
	}

	return
}

func (Transaksi) TableName() string {
	return "transaksi"
}

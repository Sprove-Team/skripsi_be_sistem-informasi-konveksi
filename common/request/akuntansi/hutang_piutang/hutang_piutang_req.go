package req_akuntansi_hp

import (
	"github.com/be-sistem-informasi-konveksi/entity"
)

type ReqAyatJurnal struct {
	AkunID string  `json:"akun_id" validate:"required,ulid"`
	Kredit float64 `json:"kredit" validate:"number,gte=0,required_without=Debit,excluded_with=Debit"`
	Debit  float64 `json:"debit" validate:"number,gte=0,required_without=Kredit,excluded_with=Kredit"`
}

type ReqTransaksi struct {
	Tanggal         string                 `json:"tanggal" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	BuktiPembayaran entity.BuktiPembayaran `json:"bukti_pembayaran" validate:"omitempty"`
	AyatJurnal      []ReqAyatJurnal        `json:"ayat_jurnal" validate:"required,eq=2,dive"`
}

type ReqBayar struct {
	Tanggal         string                 `json:"tanggal" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	BuktiPembayaran entity.BuktiPembayaran `json:"bukti_pembayaran" validate:"required,gt=0"`
	Keterangan      string                 `json:"keterangan" validate:"required"`
	AkunBayarID     string                 `json:"akun_bayar_id" validate:"required,ulid"`
	Total           float64                `json:"total" validate:"required,number,gt=0"`
}

type Create struct {
	KontakID   string       `json:"kontak_id" validate:"required,omitempty,ulid"`
	Jenis      string       `json:"jenis" validate:"required,oneof=PIUTANG HUTANG"`
	Keterangan string       `json:"keterangan" validate:"required"`
	Transaksi  ReqTransaksi `json:"transaksi" validate:"required"`
}

type CreateBayar struct {
	HutangPiutangID string `params:"hutang_piutang_id" validate:"required,ulid"`
	ReqBayar
}

type GetAll struct {
	KontakID string `query:"kontak_id" validate:"omitempty,ulid"`
	Jenis    string `query:"jenis" validate:"omitempty,oneof=PIUTANG HUTANG ALL"`
	Status   string `query:"status" validate:"omitempty,oneof=BELUM_LUNAS LUNAS ALL"`
	Next     string `query:"next" validate:"omitempty,ulid"`
	Limit    int    `query:"limit" validate:"omitempty"`
}

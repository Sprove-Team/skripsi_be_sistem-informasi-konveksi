package req_akuntansi_transaksi

import "github.com/be-sistem-informasi-konveksi/entity"

type ReqAyatJurnal struct {
	AkunID string  `json:"akun_id" validate:"required,ulid"`
	Kredit float64 `json:"kredit" validate:"number,gte=0,required_without=Debit,excluded_with=Debit"`
	Debit  float64 `json:"debit" validate:"number,gte=0,required_without=Kredit,excluded_with=Kredit"`
}

type Create struct {
	BuktiPembayaran entity.BuktiPembayaran `json:"bukti_pembayaran" validate:"omitempty"`
	Tanggal         string                 `json:"tanggal" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Keterangan      string                 `json:"keterangan" validate:"required"`
	KontakID        string                 `json:"kontak_id" validate:"omitempty,ulid"`
	AyatJurnal      []ReqAyatJurnal        `json:"ayat_jurnal" validate:"required,min=2,dive"`
}

type Update struct {
	ID              string                 `params:"id" validate:"required,ulid"`
	BuktiPembayaran entity.BuktiPembayaran `json:"bukti_pembayaran" validate:"omitempty"`
	Tanggal         string                 `json:"tanggal" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Keterangan      string                 `json:"keterangan" validate:"omitempty"`
	AyatJurnal      []ReqAyatJurnal        `json:"ayat_jurnal" validate:"omitempty,min=2,dive"`
}

type GetAll struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"end_date" validate:"required,datetime=2006-01-02"`
}

type GetHistory struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"end_date" validate:"required,datetime=2006-01-02"`
}

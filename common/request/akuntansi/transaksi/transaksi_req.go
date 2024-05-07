package req_akuntansi_transaksi

import (
	"mime/multipart"
)

type ReqAyatJurnal struct {
	AkunID string  `json:"akun_id" validate:"required,ulid"`
	Kredit float64 `json:"kredit" validate:"number,gte=0,required_without=Debit,excluded_with=Debit"`
	Debit  float64 `json:"debit" validate:"number,gte=0,required_without=Kredit,excluded_with=Kredit"`
}

type Create struct {
	BuktiPembayaran []*multipart.FileHeader `form:"bukti_pembayaran" validate:"omitempty"`
	Data            string                  `form:"data"`
	Tanggal         string                  `json:"tanggal" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Keterangan      string                  `json:"keterangan" validate:"required"`
	KontakID        string                  `json:"kontak_id" validate:"omitempty,ulid"`
	AyatJurnal      []ReqAyatJurnal         `json:"ayat_jurnal" validate:"required,min=2,dive"`
}

type Update struct {
	ID              string                  `params:"id" validate:"required,ulid"`
	BuktiPembayaran []*multipart.FileHeader `form:"bukti_pembayaran" validate:"omitempty"`
	Data            string                  `form:"data"`
	Tanggal         string                  `json:"tanggal" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Keterangan      string                  `json:"keterangan" validate:"omitempty"`
	AyatJurnal      []ReqAyatJurnal         `json:"ayat_jurnal" validate:"omitempty,min=2,dive"`
}

type GetAll struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"end_date" validate:"required,datetime=2006-01-02"`
	TimeZone  string `query:"time_zone" validate:"omitempty,timezone"`
}

type GetHistory struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"end_date" validate:"required,datetime=2006-01-02"`
	TimeZone  string `query:"time_zone" validate:"omitempty,timezone"`
}

package invoice

import "github.com/be-sistem-informasi-konveksi/entity"

type ReqDetailInvoice struct {
	ProdukID     string  `json:"produk_id" validate:"required,ulid"`
	BordirID     string  `json:"bordir_id" validate:"required,ulid"`
	SablonID     string  `json:"sablon_id" validate:"required,ulid"`
	GambarDesign string  `json:"gambar_design" validate:"omitempty"`
	Total        float64 `json:"total" validate:"required,number"`
	Qty          int     `json:"qty" validate:"required,number,min=1"`
}

type ReqNewKontak struct {
	Nama   string `json:"nama" validate:"required"`
	NoTelp string `json:"no_telp" validate:"required,e164"`
	Alamat string `json:"alamat" validate:"required"`
	Email  string `json:"email" validate:"omitempty,email"`
}

type ReqBayar struct {
	BuktiPembayaran entity.BuktiPembayaran `json:"bukti_pembayaran" validate:"required,dive,url"`
	Keterangan      string                 `json:"keterangan" validate:"required"`
	AkunID          string                 `json:"akun_id" validate:"required,ulid"`
	Total           float64                `json:"total" validate:"required,number,gt=0"`
}

type Create struct {
	UserID          string             `json:"user_id" validate:"required,ulid"`
	StatusProduksi  string             `json:"status_produksi" validate:"omitempty,oneof=BELUM_DIKERJAKAN DIPROSES SELESAI"`
	KontakID        string             `json:"kontak_id" validate:"required_without=NewKontak,excluded_with=NewKontak"`
	NewKontak       ReqNewKontak       `json:"new_kontak" validate:"omitempty"`
	Bayar           ReqBayar           `json:"bayar" validate:"required"`
	TotalHarga      float64            `json:"total_harga" validate:"required,number"`
	TotalQty        int                `json:"total_qty" validate:"required,number,min=1"`
	TanggalDeadline string             `json:"tanggal_deadline" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	TanggalKirim    string             `json:"tanggal_kirim" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Keterangan      string             `json:"keterangan" validate:"required"`
	DetailInvoice   []ReqDetailInvoice `json:"detail_invoice" validate:"gt=0,dive,required"`
}

type Update struct {
	ID              string             `params:"id" validate:"required"`
	StatusProduksi  string             `json:"status_produksi" validate:"omitempty,oneof=BELUM_DIKERJAKAN DIPROSES SELESAI"`
	KontakID        string             `json:"kontak_id" validate:"omitempty"`
	TotalHarga      float64            `json:"total_harga" validate:"required_with=DetailInvoice,number"`
	TotalQty        int                `json:"total_qty" validate:"required_with=DetailInvoice,number,min=1"`
	TanggalDeadline string             `json:"tanggal_deadline" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	TanggalKirim    string             `json:"tanggal_kirim" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Keterangan      string             `json:"keterangan" validate:"omitempty"`
	DetailInvoice   []ReqDetailInvoice `json:"detail_invoice" validate:"gt=0,dive,omitempty"`
}

type GetAll struct {
	StatusProduksi  string `query:"status_produksi" validate:"omitempty,oneof=BELUM_DIKERJAKAN DIPROSES SELESAI"`
	KontakID        string `query:"kontak_id" validate:"omitempty"`
	TanggalDeadline string `query:"tanggal_deadline" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	TanggalKirim    string `query:"tanggal_kirim" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	SortBy          string `query:"sort_by" validate:"omitempty,oneof=TANGGAL_DEADLINE TANGGAL_KIRIM"`
	OrderBy         string `query:"order_by" validate:"omitempty,oneof=ASC DESC"`
	Next            string `query:"next" validate:"omitempty"`
	Limit           int    `query:"limit" validate:"omitempty,number"`
}

package req_invoice

import "mime/multipart"

type ReqNewKontak struct {
	NamaKontak   string `json:"nama" validate:"required"`
	NoTelpKontak string `json:"no_telp" validate:"required,e164"`
	AlamatKontak string `json:"alamat" validate:"required"`
	EmailKontak  string `json:"email" validate:"omitempty,email"`
}

type ReqBayar struct {
	KeteranganPembayaran string  `json:"keterangan" validate:"required"`
	MetodePembayaran     string  `json:"akun_id" validate:"required,ulid"`
	TotalBayar           float64 `json:"total" validate:"required,number,gt=0"`
}

type ReqDetailInvoice struct {
	ProdukID     string  `json:"produk_id" validate:"required,ulid"`
	BordirID     string  `json:"bordir_id" validate:"omitempty,ulid"`
	SablonID     string  `json:"sablon_id" validate:"omitempty,ulid"`
	TotalPesanan float64 `json:"total" validate:"required,number,gt=0"`
	QtyPesanan   int     `json:"qty" validate:"required,number,gt=0"`
}

type Create struct {
	BuktiPembayaran   []*multipart.FileHeader `form:"bukti_pembayaran" validate:"required"`
	GambarDesign      []*multipart.FileHeader `form:"gambar_design" validate:"required,equalLengthWithField=DetailInvoice"`
	Data              string                  `form:"data"`
	KontakID          string                  `json:"kontak_id" validate:"required,ulid"`
	Bayar             ReqBayar                `json:"bayar" validate:"required"`
	TanggalDeadline   string                  `json:"tanggal_deadline" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	TanggalKirim      string                  `json:"tanggal_kirim" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	KeteranganPesanan string                  `json:"keterangan" validate:"required"`
	DetailInvoice     []ReqDetailInvoice      `json:"detail_invoice" validate:"required,gt=0,dive"`
}

type ReqUpdateDetailInvoice struct {
	ID    string  `json:"id" validate:"required,ulid"`
	Total float64 `json:"total" validate:"omitempty,number,gt=0"`
	Qty   int     `json:"qty" validate:"omitempty,number,gt=0"`
}

type Update struct {
	ID              string                   `params:"id" validate:"required"`
	StatusProduksi  string                   `json:"status_produksi" validate:"omitempty,oneof=BELUM_DIKERJAKAN DIPROSES SELESAI"`
	TanggalDeadline string                   `json:"tanggal_deadline" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	TanggalKirim    string                   `json:"tanggal_kirim" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Keterangan      string                   `json:"keterangan" validate:"omitempty"`
	DetailInvoice   []ReqUpdateDetailInvoice `json:"detail_invoice" validate:"omitempty,gt=0,dive"`
}

type GetAll struct {
	StatusProduksi  string `query:"status_produksi" validate:"omitempty,oneof=BELUM_DIKERJAKAN DIPROSES SELESAI"`
	KontakID        string `query:"kontak_id" validate:"omitempty,ulid"`
	TanggalDeadline string `query:"tanggal_deadline" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	TanggalKirim    string `query:"tanggal_kirim" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	TanggalDipesan  string `query:"tanggal_dipesan" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	SortBy  string `query:"sort_by" validate:"omitempty,oneof=TANGGAL_DEADLINE TANGGAL_KIRIM TANGGAL_DIPESAN"`
	OrderBy string `query:"order_by" validate:"omitempty,oneof=ASC DESC"`
	Next    string `query:"next" validate:"omitempty,ulid"`
	Limit   int    `query:"limit" validate:"omitempty,number"`
}

package invoice

type ReqDetailInvoice struct {
	ProdukID      string  `json:"produk_id" validate:"required,ulid"`
	HargaProdukID string  `json:"harga_produk_id" validate:"required,ulid"`
	BordirID      string  `json:"bordir_id" validate:"required,ulid"`
	SablonID      string  `json:"sablon_id" validate:"required,ulid"`
	GambarDesign  string  `json:"gambar_design" validate:"omitempty,url_cloud_storage"`
	Qty           float64 `json:"qty" validate:"required,number"`
}

type Create struct {
	UserID         string  `json:"user_id" validate:"required,ulid"`
	StatusProduksi string  `json:"status_produksi_id" validate:"required,oneof=BELUM_DIKERJAKAN DIPROSES SELESAI"`
	Kepada         string  `json:"kepada" validate:"required"`
	NoTelp         string  `json:"no_telp" validate:"required"`
	Alamat         string  `json:"alamat" validate:"required"`
	Bayar          float64 `json:"bayar" validate:"required,number"`
	Keterangan     string  `json:"keterangan" validate:"required"`
	// BuktiPembayaran  string             `json:"bukti_pembayaran" validate:"required,url_cloud_storage"`
	BuktiPembayaran string             `json:"bukti_pembayaran" validate:"omitempty,url_cloud_storage"`
	TanggalDeadline string             `json:"tanggal_deadline" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	TanggalKirim    string             `json:"tanggal_kirim" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	DetailInvoice   []ReqDetailInvoice `json:"detail_invoice" validate:"gt=0,dive,required"`
}

type GetAll struct {
	Nama           string `query:"nama" validate:"omitempty"`
	StatusProduksi string `query:"status_produksi" validate:"omitempty,ulid"`
	Next           string `query:"next" validate:"omitempty"`
	Limit          int    `query:"limit" validate:"omitempty,number"`
}

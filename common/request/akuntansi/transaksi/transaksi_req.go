package akuntansi

type ReqAyatJurnal struct {
	AkunID string  `json:"akun_id" validate:"required,ulid"`
	Kredit float64 `json:"kredit" validate:"required_unless=Debit "`
	Debit  float64 `json:"debit" validate:"required_unless=Kredit "`
}

type Create struct {
	BuktiPembayaran string          `json:"bukti_pembayaran" validate:"omitempty,url_cloud_storage"`
	Tanggal         string          `json:"tanggal" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Keterangan      string          `json:"keterangan" validate:"required"`
	AyatJurnals     []ReqAyatJurnal `json:"ayat_jurnal" validate:"required,min=2"`
	// BuktiPembayaran string `json:"bukti_pembayaran"`
}

type GetAll struct {
	StartDate string `query:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `query:"end_date" validate:"required,datetime=2006-01-02"`
}

package res_akuntansi_transaksi

type (
	DataAyatJurnalTR struct {
		AyatJurnalID string  `json:"ayat_jurnal_id"`
		AkunID       string  `json:"akun_id"`
		Debit        float64 `json:"debit"`
		Saldo        float64 `json:"saldo"`
		Kredit       float64 `json:"kredit"`
	}
	DataGetAllTransaksi struct {
		ID              string             `json:"id"`
		Keterangan      string             `json:"keterangan"`
		BuktiPembayaran string             `json:"bukti_pembayaran"`
		Tanggal         string             `json:"tanggal"`
		TotalKredit     float64            `json:"total_kredit"`
		TotalDebit      float64            `json:"total_debit"`
		AyatJurnals     []DataAyatJurnalTR `json:"ayat_jurnal"`
	}
)

package akuntansi

type (
	DataAyatJurnalJU struct {
		AyatJurnalID string  `json:"ayat_jurnal_id"`
		AkunID       string  `json:"akun_id"`
		KodeAkun     string  `json:"kode_akun"`
		NamaAkun     string  `json:"nama_akun"`
		Debit        float64 `json:"debit"`
		Kredit       float64 `json:"kredit"`
	}
	DataTransaksiJU struct {
		Tanggal     string             `json:"tanggal"`
		TransaksiID string             `json:"transaksi_id"`
		Keterangan  string             `json:"keterangan"`
		AyatJurnal  []DataAyatJurnalJU `json:"ayat_jurnal"`
	}

	JurnalUmumRes struct {
		TotalDebit  float64           `json:"total_debit"`
		TotalKredit float64           `json:"total_kredit"`
		Transaksi   []DataTransaksiJU `json:"transaksi"`
	}
)

type (
	DataAyatJurnalBB struct {
		TransaksiID string  `json:"transaksi_id"`
		Tanggal     string  `json:"tanggal"`
		Keterangan  string  `json:"keterangan"`
		Debit       float64 `json:"debit"`
		Kredit      float64 `json:"kredit"`
		Saldo       float64 `json:"saldo"`
	}
	BukuBesarRes struct {
		KodeAkun    string             `json:"kode_akun"`
		NamaAkun    string             `json:"nama_akun"`
		SaldoNormal string             `json:"saldo_normal"`
		TotalSaldo  float64            `json:"total_saldo"`
		TotalDebit  float64            `json:"total_debit"`
		TotalKredit float64            `json:"total_kredit"`
		AyatJurnals []DataAyatJurnalBB `json:"ayat_jurnal"`
	}
)

type (
	DataSaldoAkun struct {
		KodeAkun    string  `json:"kode_akun"`
		NamaAkun    string  `json:"nama_akun"`
		SaldoDebit  float64 `json:"saldo_debit"`
		SaldoKredit float64 `json:"saldo_kredit"`
	}
	NeracaSaldoRes struct {
		TotalKredit    float64         `json:"total_kredit"`
		TotalDebit     float64         `json:"total_debit"`
		DataSaldoAkuns []DataSaldoAkun `json:"saldo_akun"`
	}
)

type (
	// DataSaldoAkunLBR struct {
	//   NamaAkun string
	//   KodeAkun string
	//   SaldoDebit  float64 `json:"saldo_debit"`
	//   SaldoKredit float64 `json:"saldo_kredit"`
	// }
	// DataAkunLBR struct {
	// 	NamaAkun string
	// 	Saldo    float64
	// }
	DataSaldoLBR struct {
		NamaKelompokAkun string          `json:"nama_kelompok_akun"`
		Total            float64         `json:"total"`
		DataSaldoAkuns   []DataSaldoAkun `json:"saldo_akun"`
	}
	LabaRugiRes struct {
		SaldoLabaRugi []DataSaldoLBR `json:"saldo_laba_rugi"`
		Laba          float64
		Rugi          float64
	}
)

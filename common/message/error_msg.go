package message

// global
var (
	BadRequest          = "request tidak valid"
	InternalServerError = "terjadi kesalahan pada server"
	NotFound            = "data tidak ditemukan"
	Conflict            = "data telah telah ditambahkan"
	RequestTimeout      = "request timeout"
)

// middleware auth

var (
	UnauthInvalidToken   = "token tidak valid"
	UnauthTokenExpired   = "token telah kadaluarsa"
	UnauthUserNotAllowed = "pengguna tidak diizinkan untuk mengakses halaman ini"
	UnauthUserNotFound   = "pengguna telah dihapus atau hubungi direktur"
)

// misc
var InvalidImageFormat = "format gambar tidak valid"

// produk
var (
	// data
	ProdukNotFound = "produk tidak ditemukan"
	// kategori
	KategoriProdukNotFound = "kategori produk tidak ditemukan"
	// harga detail
	HargaDetailProdukNotFound              = "harga detail produk tidak ditemukan"
	HargaDetailProdukEmpty                 = "harga detail produk pada produk ini masih kosong"
	HargaDetailProdukNotFoundOrNotAddedYet = "harga detail produk dengan qty ini tidak ditemukan atau belum ditambahkan, hubungi direktur untuk informasi lebih lanjut"
)

// bordir
var (
	BordirNotFound = "bordir tidak ditemukan"
)

// sablon
var (
	SablonNotFound = "sablon tidak ditemukan"
)

// user
var (
	PasswordIsNotStrong = "password setidaknya harus berisi angka dan huruf besar"
	UserNotFound        = "user tidak ditemukan"
)

// misc
var (
	CantDeleteDefaultData = "tidak dapat menghapus atau mengubah data default"
)

// akuntansi
var (
	// transaksi
	CreditDebitNotSame        = "total debit dan kredit harus sama"
	AkunCannotBeSame          = "akun tidak boleh sama"
	AkunCannotDeleted         = "tidak bisa menghapus akun karena akun masih digunakan pada ayat jurnal"
	AkunHutangPiutangNotEq2   = "transaksi merupakan hutang piutang, data ayat_jurnal harus berjumlah 2"
	AkunNotMatchWithJenisHPTr = "transaksi merupakan hutang piutang, akun harus sama dengan jenis hutang piutang"
	// kategori akun
	KategoriAkunNotFound = "kategori akun tidak ditemukan"
	// kelompok akun
	KelompokAkunNotFound = "kelompok akun tidak ditemukan"
	// akun
	AkunNotFound = "akun tidak ditemukan"
	// kontak
	KontakNotFound = "kontak tidak ditemukan"
	// hutang piutang
	HutangPiutangNotFound              = "hutang piutang tidak di temukan"
	AkunNotMatchWithJenisHP            = "akun harus sama dengan jenis hutang piutang"
	InvalidAkunHutangPiutang           = "akun bukan merupakan akun hutang piutang"
	IncorrectPlacementOfCreditAndDebit = "peletakan total debit dan kredit untuk hutang piutang tidak benar"
	TotalHPMustGeOrEqToTotalByr        = "jumlah hutang piutang harus lebih besar atau sama dengan total hutang piutang yang telah dibayar"
	// bayar hutang piutang
	BayarMustLessThanSisaTagihan = "jumlah yang dibayar harus kurang atau sama dengan sisa tagihan"
	InvalidAkunBayar             = "akun untuk bayar hutang piutang tidak valid"
)

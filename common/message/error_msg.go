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

// auth
var (
	InvalidUsernameOrPassword = "username atau password tidak valid"
	RefreshTokenExpired       = "refresh token telah kadaluarsa"
	InvalidRefreshToken       = "refresh token tidak valid"
	UserNotFoundOrDeleted     = "user tidak ditemukan atau telah dihapus"
)

// misc
var (
	InvalidImageFormat = "format gambar tidak valid"
)

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

// profile
var (
	NotFitOldPassword = "old password tidak cocok"
	UsernameConflict  = "username telah digunakan"
)

// tugas
var (
	UserNotFoundOrNotSpv = "user tidak ditemukan atau jenis supervisor tidak sesuai dengan user"
	TugasNotFound        = "tugas tidak ditemukan"
)

// misc
var (
	CantModifiedDefaultData = "tidak dapat menghapus atau mengubah data default"
)

// akuntansi
var (
	// transaksi
	CreditDebitNotSame                         = "total debit dan kredit harus sama"
	AkunCannotBeSame                           = "akun tidak boleh sama"
	AkunCannotDeleted                          = "tidak bisa menghapus akun karena akun masih digunakan pada ayat jurnal"
	AkunHutangPiutangNotEq2                    = "transaksi merupakan hutang piutang, data ayat_jurnal harus berjumlah 2"
	AkunNotMatchWithJenisHPTr                  = "transaksi merupakan hutang piutang, akun harus sama dengan jenis hutang piutang"
	CantDeleteTrIfDataByrBlmTerkonfirmasiExist = "transaksi tidak dapat di hapus bila masih terdapat data bayar yang berstatus belum terkonfirmasi"
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
	IncorrectEntryAkunHP               = "transaksi menyebabkan pengurangan atau penambahan pada akun hutang piutang yang salah"
	AkunHPDoesNotExist                 = "akun hutang piutang tidak ada pada ayat jurnal"
	IncorrectPlacementOfCreditAndDebit = "peletakan total debit dan kredit untuk hutang piutang tidak benar"
	TotalHPMustGeOrEqToTotalByr        = "jumlah hutang piutang harus lebih besar atau sama dengan total hutang piutang yang telah dibayar"
	// bayar hutang piutang
	BayarMustLessThanSisaTagihan = "jumlah yang dibayar harus kurang atau sama dengan sisa tagihan"
	InvalidAkunBayar             = "akun untuk bayar hutang piutang tidak valid"
)

// invoice
var (
	TotalBayarMustGeOrEqToTotalByr     = "total bayar invoice harus lebih besar atau sama dengan total yang telah di bayar"
	BayarMustLessThanTotalHargaInvoice = "jumlah yang dibayar harus kurang atau sama dengan total harga invoice"
	InvoiceNotFound                    = "invoice tidak ditemukan"
	DetailInvoiceNotFound              = "detail invoice tidak ditemukan"
	//data bayar
	CannotModifiedTerkonfirmasiDataBayar = "tidak bisa mengubah atau menghapus data bayar invoice yang berstatus `TERKONFIRMASI`"
	// user not allowed
	UserNotAllowedToModifiedStatusProdusi = " tidak diizinkan untuk mengedit status produksi"
)

package message

// global
var (
	BadRequest          = "request tidak valid"
	InternalServerError = "terjadi kesalahan pada server"
	NotFound            = "data tidak tidak ditemukan"
	Conflict            = "data telah telah ditambahkan"
	RequestTimeout      = "request timeout"
)

// middleware auth

var (
	UnauthInvalidFormatToken = "format token jwt tidak valid"
	UnauthTokenExpired       = "token telah kadaluarsa"
	UnauthUserNotAllowed     = "pengguna tidak diizinkan untuk mengakses halaman ini"
)

// produk
var (
	// data
	// ProdukConflict = "produk telah ditambahkan"
	ProdukNotFound = "produk tidak ditemukan"
	// kategori
	// KategoriConflict = "kategori produk telah ditambahkan"
	KategoriNotFound = "kategori produk tidak ditemukan"
	// harga detail
	// HargaDetailConflict = "harga detail produk telah ditambahkan"
	HargaDetailNotFound = "harga detail produk tidak ditemukan"
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
)

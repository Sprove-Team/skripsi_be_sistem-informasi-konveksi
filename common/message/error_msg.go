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

// produk
var (
	// data
	ProdukNotFound = "produk tidak ditemukan"
	// kategori
	KategoriNotFound = "kategori produk tidak ditemukan"
	// harga detail
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

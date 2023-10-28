package message

// global
var (
	BadRequest          = "request tidak valid"
	InternalServerError = "terjadi kesalahan pada server"
	NotFound            = "data tidak tidak ditemukan"
	Conflict            = "data telah telah ditambahkan"
	RequestTimeout      = "request timeout"
)

// produk
var (
  // data
	ProdukConflict = "produk telah ditambahkan"
	ProdukNotFound = "produk tidak ditemukan"
  // kategori
  KategoriConflict = "kategori produk telah ditambahkan"
  KategoriNotFound = "kategori produk tidak ditemukan"
  // harga detail
  HargaDetailConflict = "harga detail produk telah ditambahkan"
  HargaDetailNotFound = "harga detail produk tidak ditemukan"

)

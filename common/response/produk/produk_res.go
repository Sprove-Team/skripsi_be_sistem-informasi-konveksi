package produk

import "github.com/be-sistem-informasi-konveksi/entity"

type ProdukRes struct {
	ID          string           `json:"id"`
	Nama        string           `json:"nama"`
	HargaDetail []entity.HargaDetailProduk `json:"harga_detail"`
	Kategori    entity.KategoriProduk      `json:"kategori"`
}

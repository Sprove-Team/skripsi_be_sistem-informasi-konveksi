package entity

type Produk struct {
	Base
	KategoriProdukID string              `gorm:"type:varchar(26);index:idx_kategori_produk_id;not null" json:"-"`
	KategoriProduk   *KategoriProduk     `json:"kategori,omitempty"`
	DetailInvoice    *DetailInvoice      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Nama             string              `gorm:"type:varchar(150);uniqueIndex;not null" json:"nama,omitempty"`
	HargaDetails     []HargaDetailProduk `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"harga_detail,omitempty"`
}

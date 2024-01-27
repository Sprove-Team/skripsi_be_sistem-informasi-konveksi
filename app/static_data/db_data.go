package static_data

import (
	"fmt"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

// akuntansi

var DataKelompokAkun = []entity.KelompokAkun{
	// Aset
	{Nama: "kas & bank", KategoriAkun: "ASET"},
	{Nama: "piutang", KategoriAkun: "ASET"},
	{Nama: "aset lancar lainnya", KategoriAkun: "ASET"},
	{Nama: "aset tetap berwujud", KategoriAkun: "ASET"},
	{Nama: "aset tetap tidak berwujud", KategoriAkun: "ASET"},
	{Nama: "aset tidak lancar lainnya", KategoriAkun: "ASET"},
	// Kewajiban
	{Nama: "hutang", KategoriAkun: "KEWAJIBAN"},
	{Nama: "kewajiban lancar lainnya", KategoriAkun: "KEWAJIBAN"},
	{Nama: "kewajiban jangka panjang", KategoriAkun: "KEWAJIBAN"},
	// Modal
	{Nama: "modal", KategoriAkun: "MODAL"},
	// Pendapatan
	{Nama: "pendapatan usaha", KategoriAkun: "PENDAPATAN"},
	{Nama: "pendapatan di luar usaha", KategoriAkun: "PENDAPATAN"},
	// Beban
	{Nama: "harga pokok penjualan", KategoriAkun: "BEBAN"},
	{Nama: "beban usaha lainnya", KategoriAkun: "BEBAN"},
	{Nama: "beban diluar usaha lainnya", KategoriAkun: "BEBAN"},
}

var DataAkun = [][]entity.Akun{
	// Kas & Bank
	{
		{Nama: "kas", SaldoNormal: "DEBIT", Deskripsi: "Kas dari usaha"},
	},
	// Piutang
	{
		{Nama: "piutang usaha", SaldoNormal: "DEBIT", Deskripsi: "Piutang dari penjualan barang atau jasa kepada pelanggan"},
	},
	// Aset Lancar Lainnya
	{
		{Nama: "perlengkapan", SaldoNormal: "DEBIT", Deskripsi: "Perlengkapan yang digunakan dalam usaha"},
	},
	// Aset Tetap Berwujud
	{
		{Nama: "peralatan", SaldoNormal: "DEBIT", Deskripsi: "Peralatan yang digunakan dalam usaha"},
		{Nama: "mesin-mesin", SaldoNormal: "DEBIT", Deskripsi: "Mesin-mesin yang digunakan dalam usaha"},
	},
	// Aset Tetap Tidak Berwujud
	{
		{Nama: "hak cipta", SaldoNormal: "DEBIT", Deskripsi: "Hak cipta dari suatu karya atau lainnya"},
	},
	// Aset Tidak Lancar Lainnya
	{
		{Nama: "aset tidak lancar lainnya", SaldoNormal: "DEBIT", Deskripsi: "Aset tidak lancar lainnya"},
	},
	// Hutang
	{
		{Nama: "hutang usaha", SaldoNormal: "KREDIT", Deskripsi: "Hutang dari pembelian barang atau jasa yang belum dibayar"},
	},
	// Kewajiban Lancar Lainnya
	{
		{Nama: "hutang sewa", SaldoNormal: "KREDIT", Deskripsi: "Hutang sewa yang merupakan kewajiban lancar"},
		{Nama: "hutang gaji", SaldoNormal: "KREDIT", Deskripsi: "Hutang gaji kepada karyawan yang belum dibayar"},
		{Nama: "hutang pajak", SaldoNormal: "KREDIT", Deskripsi: "Hutang  pajak yang harus dibayar dalam waktu dekat"},
	},
	// Kewajiban Jangka Panjang
	{
		{Nama: "hutang bank", SaldoNormal: "KREDIT", Deskripsi: "Hutang kepada bank yang jatuh tempo dalam waktu lebih dari satu tahun"},
		{Nama: "hutang obligasi", SaldoNormal: "KREDIT", Deskripsi: "Hutang berupa obligasi yang jatuh tempo dalam waktu lebih dari satu tahun"},
	},
	// Modal
	{
		{Nama: "modal pribadi", SaldoNormal: "KREDIT", Deskripsi: "Modal yang ditanamkan oleh pemilik secara pribadi"},
		{Nama: "modal saham", SaldoNormal: "KREDIT", Deskripsi: "Modal yang didapat dari pendapatan usaha"},
		{Nama: "prive", SaldoNormal: "DEBIT", Deskripsi: "Pengambilan uang oleh pemilik untuk keperluan pribadi"},
	},
	// Pendapatan Usaha
	{
		{Nama: "pendapatan jasa", SaldoNormal: "KREDIT", Deskripsi: "Pendapatan dari jasa yang diberikan"},
	},
	// Pendapatan Diluar Usaha
	{
		{Nama: "pendapatan bunga", SaldoNormal: "KREDIT", Deskripsi: "Pendapatan dari bunga yang diterima"},
	},
	// Harga Pokok Penjualan
	{
		{Nama: "harga pokok pendapatan", SaldoNormal: "DEBIT", Deskripsi: "Harga pokok pendapatan usaha"},
	},
	// Beban Usaha Lainnya
	{
		{Nama: "beban gaji", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk membayar gaji karyawan"},
		{Nama: "beban sewa", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk sewa tempat atau fasilitas"},
		{Nama: "beban asuransi", SaldoNormal: "DEBIT", Deskripsi: "Beban asuransi yang dibayarkan"},
		{Nama: "beban air, listrik dan telepon", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk utilitas"},
		{Nama: "beban perlengkapan", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk membeli perlengkapan"},
		{Nama: "beban penyusutan bangunan", SaldoNormal: "DEBIT", Deskripsi: "Beban penyusutan bangunan"},
		{Nama: "beban penyusutan kendaraan", SaldoNormal: "DEBIT", Deskripsi: "Beban penyusutan kendaraan"},
		{Nama: "beban penyusutan peralatan", SaldoNormal: "DEBIT", Deskripsi: "Beban penyusutan peralatan"},
		{Nama: "beban administrasi lainnya", SaldoNormal: "DEBIT", Deskripsi: "Beban administrasi lainnya"},
		{Nama: "beban iklan", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk biaya iklan"},
		{Nama: "pajak penghasilan", SaldoNormal: "DEBIT", Deskripsi: "Beban pajak penghasilan yang harus dibayar"},
	},
	// Beban Diluar Usaha Lainnya
	{
		{Nama: "beban bunga", SaldoNormal: "DEBIT", Deskripsi: "Beban bunga yang harus dibayar"},
	},
}

var (
	KelompokAkuns = func(ulid pkg.UlidPkg) (datas []entity.KelompokAkun) {

		datas = make([]entity.KelompokAkun, 0, len(DataKelompokAkun))

		for i, kelompokAkun := range DataKelompokAkun {
			datas = append(datas, entity.KelompokAkun{
				Base: entity.Base{
					ID: ulid.MakeUlid().String(),
				},
				Nama:         kelompokAkun.Nama,
				Kode:         fmt.Sprintf("%s%d", entity.KategoriAkun[kelompokAkun.KategoriAkun], i+1),
				KategoriAkun: kelompokAkun.KategoriAkun,
			})
		}
		return datas
	}

	Akuns = func(kelompokAkuns []entity.KelompokAkun, ulid pkg.UlidPkg) (datas []entity.Akun) {

		datas = make([]entity.Akun, 0, 34)

		for i, kelompokAkun := range kelompokAkuns {
			for j, akun := range DataAkun[i] {
				datas = append(datas, entity.Akun{
					Base: entity.Base{
						ID: ulid.MakeUlid().String(),
					},
					KelompokAkunID: kelompokAkun.ID,
					Nama:           akun.Nama,
					Deskripsi:      akun.Deskripsi,
					SaldoNormal:    akun.SaldoNormal,
					Kode:           fmt.Sprintf("%s%d", kelompokAkun.Kode, j+1),
				})
			}
		}

		return datas
	}
)

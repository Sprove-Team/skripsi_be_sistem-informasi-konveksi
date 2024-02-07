package static_data

import (
	"fmt"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

// akuntansi

var DataKelompokAkun = []entity.KelompokAkun{
	// Aset
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMWDJTQNGS"}, Nama: "kas & bank", KategoriAkun: "ASET"},
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMWFGPJ0NC"}, Nama: "piutang", KategoriAkun: "ASET"},
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMWKA5M0WJ"}, Nama: "aset lancar lainnya", KategoriAkun: "ASET"},
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMWNKW2MY2"}, Nama: "aset tetap berwujud", KategoriAkun: "ASET"},
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMWR3BQCJE"}, Nama: "aset tetap tidak berwujud", KategoriAkun: "ASET"},
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMWVJP5V5X"}, Nama: "aset tidak lancar lainnya", KategoriAkun: "ASET"},
	// Kewajiban
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMWXFG0ETJ"}, Nama: "hutang", KategoriAkun: "KEWAJIBAN"},
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMWXZHT3YH"}, Nama: "kewajiban lancar lainnya", KategoriAkun: "KEWAJIBAN"},
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMWY91025S"}, Nama: "kewajiban jangka panjang", KategoriAkun: "KEWAJIBAN"},
	// Modal
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMX1N3XR54"}, Nama: "modal", KategoriAkun: "MODAL"},
	// Pendapatan
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMX56J933W"}, Nama: "pendapatan usaha", KategoriAkun: "PENDAPATAN"},
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMX94R3KZV"}, Nama: "pendapatan di luar usaha", KategoriAkun: "PENDAPATAN"},
	// Beban
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXBWMC2GK"}, Nama: "harga pokok penjualan", KategoriAkun: "BEBAN"},
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXDG18Q11"}, Nama: "beban usaha lainnya", KategoriAkun: "BEBAN"},
	{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXEJ8RX88"}, Nama: "beban diluar usaha lainnya", KategoriAkun: "BEBAN"},
}

var DataAkun = [][]entity.Akun{
	// Kas & Bank
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXGPC0KSP"}, Nama: "kas", SaldoNormal: "DEBIT", Deskripsi: "Kas dari usaha"},
	},
	// Piutang
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXKJAR066"}, Nama: "piutang usaha", SaldoNormal: "DEBIT", Deskripsi: "Piutang dari penjualan barang atau jasa kepada pelanggan"},
	},
	// Aset Lancar Lainnya
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXMGHHDYC"}, Nama: "perlengkapan", SaldoNormal: "DEBIT", Deskripsi: "Perlengkapan yang digunakan dalam usaha"},
	},
	// Aset Tetap Berwujud
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXNF1JSDW"}, Nama: "peralatan", SaldoNormal: "DEBIT", Deskripsi: "Peralatan yang digunakan dalam usaha"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXNF6PZNS"}, Nama: "mesin-mesin", SaldoNormal: "DEBIT", Deskripsi: "Mesin-mesin yang digunakan dalam usaha"},
	},
	// Aset Tetap Tidak Berwujud
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXSAWKPPX"}, Nama: "hak cipta", SaldoNormal: "DEBIT", Deskripsi: "Hak cipta dari suatu karya atau lainnya"},
	},
	// Aset Tidak Lancar Lainnya
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXWBEQR81"}, Nama: "aset tidak lancar lainnya", SaldoNormal: "DEBIT", Deskripsi: "Aset tidak lancar lainnya"},
	},
	// Hutang
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMY04Y7H4W"}, Nama: "hutang usaha", SaldoNormal: "KREDIT", Deskripsi: "Hutang dari pembelian barang atau jasa yang belum dibayar"},
	},
	// Kewajiban Lancar Lainnya
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMY29AFFR3"}, Nama: "hutang sewa", SaldoNormal: "KREDIT", Deskripsi: "Hutang sewa yang merupakan kewajiban lancar"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMY4Z6BSAA"}, Nama: "hutang gaji", SaldoNormal: "KREDIT", Deskripsi: "Hutang gaji kepada karyawan yang belum dibayar"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMY8S1KX96"}, Nama: "hutang pajak", SaldoNormal: "KREDIT", Deskripsi: "Hutang pajak yang harus dibayar dalam waktu dekat"},
	},
	// Kewajiban Jangka Panjang
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMYB48T5H1"}, Nama: "hutang bank", SaldoNormal: "KREDIT", Deskripsi: "Hutang kepada bank yang jatuh tempo dalam waktu lebih dari satu tahun"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMYEGQQFGW"}, Nama: "hutang obligasi", SaldoNormal: "KREDIT", Deskripsi: "Hutang berupa obligasi yang jatuh tempo dalam waktu lebih dari satu tahun"},
	},
	// Modal
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMYHH0FBDA"}, Nama: "modal pribadi", SaldoNormal: "KREDIT", Deskripsi: "Modal yang ditanamkan oleh pemilik secara pribadi"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMYHVFHV2D"}, Nama: "modal saham", SaldoNormal: "KREDIT", Deskripsi: "Modal yang didapat dari pendapatan usaha"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMYK2WKWPD"}, Nama: "prive", SaldoNormal: "DEBIT", Deskripsi: "Pengambilan uang oleh pemilik untuk keperluan pribadi"},
	},
	// Pendapatan Usaha
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMYKJ8YQ1M"}, Nama: "pendapatan jasa", SaldoNormal: "KREDIT", Deskripsi: "Pendapatan dari jasa yang diberikan"},
	},
	// Pendapatan Diluar Usaha
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMYQ0EMTYV"}, Nama: "pendapatan bunga", SaldoNormal: "KREDIT", Deskripsi: "Pendapatan dari bunga yang diterima"},
	},
	// Harga Pokok Penjualan
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMYSCADKDZ"}, Nama: "harga pokok pendapatan", SaldoNormal: "DEBIT", Deskripsi: "Harga pokok pendapatan usaha"},
	},
	// Beban Usaha Lainnya
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMYWK9P4WK"}, Nama: "beban gaji", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk membayar gaji karyawan"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMYZ0QR6R6"}, Nama: "beban sewa", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk sewa tempat atau fasilitas"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMYZYC1866"}, Nama: "beban asuransi", SaldoNormal: "DEBIT", Deskripsi: "Beban asuransi yang dibayarkan"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMZ3AEMKJ4"}, Nama: "beban air, listrik dan telepon", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk utilitas"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMZ6QX9BP4"}, Nama: "beban perlengkapan", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk membeli perlengkapan"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMZ962N99Q"}, Nama: "beban penyusutan bangunan", SaldoNormal: "DEBIT", Deskripsi: "Beban penyusutan bangunan"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXGPC0KSP"}, Nama: "beban penyusutan kendaraan", SaldoNormal: "DEBIT", Deskripsi: "Beban penyusutan kendaraan"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXKJAR066"}, Nama: "beban penyusutan peralatan", SaldoNormal: "DEBIT", Deskripsi: "Beban penyusutan peralatan"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXMGHHDYC"}, Nama: "beban administrasi lainnya", SaldoNormal: "DEBIT", Deskripsi: "Beban administrasi lainnya"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXNF1JSDW"}, Nama: "beban iklan", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk biaya iklan"},
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXNF6PZNS"}, Nama: "pajak penghasilan", SaldoNormal: "DEBIT", Deskripsi: "Beban pajak penghasilan yang harus dibayar"},
	},
	// Beban Diluar Usaha Lainnya
	{
		{Base: entity.Base{ID: "01HN592DA7QWCV1SFMXSAWKPPX"}, Nama: "beban bunga", SaldoNormal: "DEBIT", Deskripsi: "Beban bunga yang harus dibayar"},
	},
}

var DefaultKodeAkunNKelompokAkun = make(map[string]bool)
var (
	KelompokAkuns = func(ulid pkg.UlidPkg) (datas []entity.KelompokAkun) {

		datas = make([]entity.KelompokAkun, 0, len(DataKelompokAkun))

		for i, kelompokAkun := range DataKelompokAkun {
			kode := fmt.Sprintf("%s%d", entity.KategoriAkun[kelompokAkun.KategoriAkun], i+1)
			datas = append(datas, entity.KelompokAkun{
				Base: entity.Base{
					ID: kelompokAkun.ID,
				},
				Nama:         kelompokAkun.Nama,
				Kode:         kode,
				KategoriAkun: kelompokAkun.KategoriAkun,
			})
			DefaultKodeAkunNKelompokAkun[kelompokAkun.ID] = true
		}
		return datas
	}

	Akuns = func(kelompokAkuns []entity.KelompokAkun, ulid pkg.UlidPkg) (datas []entity.Akun) {

		datas = make([]entity.Akun, 0, 34)

		for i, kelompokAkun := range kelompokAkuns {
			for j, akun := range DataAkun[i] {
				kode := fmt.Sprintf("%s%d", kelompokAkun.Kode, j+1)
				datas = append(datas, entity.Akun{
					Base: entity.Base{
						ID: akun.ID,
					},
					KelompokAkunID: kelompokAkun.ID,
					Nama:           akun.Nama,
					Deskripsi:      akun.Deskripsi,
					SaldoNormal:    akun.SaldoNormal,
					Kode:           kode,
				})
				DefaultKodeAkunNKelompokAkun[akun.ID] = true
			}
		}

		return datas
	}
)

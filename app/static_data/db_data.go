package static_data

import (
	"fmt"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

// akuntansi
var (
	KelompokAkuns = func(ulid pkg.UlidPkg) (datas []entity.KelompokAkun) {
		dataKelompokAkun := []entity.KelompokAkun{
			// Aset
			{Nama: "aset lancar", KategoriAkun: "ASET"},
			{Nama: "aset tidak lancar", KategoriAkun: "ASET"},
			// Kewajiban
			{Nama: "kewajiban lancar", KategoriAkun: "KEWAJIBAN"},
			{Nama: "kewajiban jangka panjang", KategoriAkun: "KEWAJIBAN"},
			// Modal
			{Nama: "modal", KategoriAkun: "MODAL"},
			// Pendapatan
			{Nama: "pendapatan usaha", KategoriAkun: "PENDAPATAN"},
			{Nama: "pendapatan di luar usaha", KategoriAkun: "PENDAPATAN"},
			// Beban
			{Nama: "beban usaha", KategoriAkun: "BEBAN"},
			{Nama: "beban diluar usaha", KategoriAkun: "BEBAN"},
		}

		datas = make([]entity.KelompokAkun, 0, 9)

		for i, kelompokAkun := range dataKelompokAkun {
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
		dataAkun := [][]entity.Akun{
			// Aset Lancar
			{
				{Nama: "kas", SaldoNormal: "DEBIT", Deskripsi: "Kas sebagai aset lancar"},
				{Nama: "piutang usaha", SaldoNormal: "DEBIT", Deskripsi: "Piutang dari penjualan barang atau jasa kepada pelanggan"},
				{Nama: "piutang lain-lain", SaldoNormal: "DEBIT", Deskripsi: "Piutang lainnya yang merupakan aset lancar"},
				{Nama: "perlengkapan", SaldoNormal: "DEBIT", Deskripsi: "Persediaan perlengkapan sebagai aset lancar"},
			},
			// Aset Tidak Lancar
			{
				{Nama: "mesin-mesin", SaldoNormal: "DEBIT", Deskripsi: "Mesin-mesin sebagai aset tetap"},
				{Nama: "peralatan", SaldoNormal: "DEBIT", Deskripsi: "Peralatan sebagai aset tetap"},
			},
			// Kewajiban Lancar
			{
				{Nama: "utang usaha", SaldoNormal: "KREDIT", Deskripsi: "Utang dari pembelian barang atau jasa yang belum dibayar"},
				{Nama: "utang sewa", SaldoNormal: "KREDIT", Deskripsi: "Utang sewa yang merupakan kewajiban lancar"},
				{Nama: "utang gaji", SaldoNormal: "KREDIT", Deskripsi: "Utang gaji kepada karyawan yang belum dibayar"},
				{Nama: "utang pajak", SaldoNormal: "KREDIT", Deskripsi: "Utang pajak yang harus dibayar dalam waktu dekat"},
			},
			// Kewajiban Jangka Panjang
			{
				{Nama: "utang bank", SaldoNormal: "KREDIT", Deskripsi: "Utang kepada bank yang jatuh tempo dalam waktu lebih dari satu tahun"},
				{Nama: "utang obligasi", SaldoNormal: "KREDIT", Deskripsi: "Utang berupa obligasi yang jatuh tempo dalam waktu lebih dari satu tahun"},
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
			// Beban Usaha
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
			},
			// Beban Diluar Usaha
			{
				{Nama: "beban bunga", SaldoNormal: "DEBIT", Deskripsi: "Beban bunga yang harus dibayar"},
				{Nama: "pajak penghasilan", SaldoNormal: "DEBIT", Deskripsi: "Beban pajak penghasilan yang harus dibayar"},
			},
		}

		datas = make([]entity.Akun, 0, 34)

		for i, kelompokAkun := range kelompokAkuns {
			for j, akun := range dataAkun[i] {
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

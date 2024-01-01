package static_data

import (
	"fmt"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

// akuntansi
var (
	// KelompokAkun = func(ulid pkg.UlidPkg) (datas []entity.KelompokAkun) {
	// 	datas = []entity.KelompokAkun{
	// 		{
	// 			ID:           ulid.MakeUlid().String(),
	// 			Nama:         "aset",
	// 			Kode:         "1",
	// 			KategoriAkun: "ASET",
	// 			{
	// 				ID:           ulid.MakeUlid().String(),
	// 				Nama:         "kewajiban",
	// 				Kode:         "2",
	// 				KategoriAkun: "KEWAJIBAN",
	// 			},
	// 			{
	// 				ID:           ulid.MakeUlid().String(),
	// 				Nama:         "modal (ekuitas)",
	// 				Kode:         "3",
	// 				KategoriAkun: "MODAL",
	// 			},
	//
	// 			{
	// 				ID:           ulid.MakeUlid().String(),
	// 				Nama:         "pendapatan",
	// 				Kode:         "4",
	// 				KategoriAkun: "PENDAPATAN",
	// 			},
	// 			{
	// 				ID:           ulid.MakeUlid().String(),
	// 				Nama:         "beban",
	// 				Kode:         "5",
	// 				KategoriAkun: "BEBAN",
	// 			},
	// 		},
	// 	}
	// return datas
	// }
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
				ID:           ulid.MakeUlid().String(),
				Nama:         kelompokAkun.Nama,
				Kode:         fmt.Sprintf("%s%d", entity.KategoriAkun[kelompokAkun.KategoriAkun], i+1),
				KategoriAkun: kelompokAkun.KategoriAkun,
			})
		}
		return datas
	}

	Akuns = func(kelompokAkuns []entity.KelompokAkun, ulid pkg.UlidPkg) (datas []entity.Akun) {
		// dataAkun := []entity.Akun{
		//  {
		//     ID: ulid.MakeUlid().String(),
		//     GolonganAkunID: golonganAkuns[0].ID,
		//     Kode: "111",
		//     Nama: "kas",
		//     Deskripsi: "kas usaha",
		//  },
		// }
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
				{Nama: "prive", SaldoNormal: "DEBIT", Deskripsi: "Pengambilan uang oleh pemilik untuk keperluan pribadi"},
			},
			// Pendapatan Usaha
			{
				{Nama: "penjualan", SaldoNormal: "KREDIT", Deskripsi: "Pendapatan dari penjualan barang atau jasa kepada pelanggan"},
				{Nama: "return penjualan", SaldoNormal: "DEBIT", Deskripsi: "Pengurangan pendapatan karena retur penjualan"},
				{Nama: "potongan penjualan", SaldoNormal: "DEBIT", Deskripsi: "Pengurangan pendapatan karena potongan penjualan"},
			},
			// Pendapatan Diluar Usaha
			{
				{Nama: "pendapatan bunga", SaldoNormal: "KREDIT", Deskripsi: "Pendapatan dari bunga yang diterima"},
				{Nama: "laba penjualan aset tetap", SaldoNormal: "KREDIT", Deskripsi: "Laba dari penjualan aset tetap"},
				{Nama: "laba penjualan surat berharga", SaldoNormal: "KREDIT", Deskripsi: "Laba dari penjualan surat berharga"},
				{Nama: "pendapatan lainnya", SaldoNormal: "KREDIT", Deskripsi: "Pendapatan lainnya diluar usaha utama"},
			},
			// Beban Usaha
			{
				{Nama: "beban gaji", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk membayar gaji karyawan"},
				{Nama: "beban sewa", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk sewa tempat atau fasilitas"},
				{Nama: "beban asuransi", SaldoNormal: "DEBIT", Deskripsi: "Beban asuransi yang dibayarkan"},
				{Nama: "beban penyusutan aset tetap", SaldoNormal: "DEBIT", Deskripsi: "Beban penyusutan untuk aset tetap"},
				{Nama: "beban air, listrik dan telepon", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk utilitas"},
				{Nama: "beban perlengkapan", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk membeli perlengkapan"},
				{Nama: "beban administrasi lainnya", SaldoNormal: "DEBIT", Deskripsi: "Beban administrasi lainnya"},
				{Nama: "beban iklan", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk biaya iklan"},
				{Nama: "kerugian piutang", SaldoNormal: "DEBIT", Deskripsi: "Kerugian dari piutang yang tidak tertagih"},
				{Nama: "beban pemasaran lainnya", SaldoNormal: "DEBIT", Deskripsi: "Beban pemasaran lainnya"},
			},
			// Beban Diluar Usaha
			{
				{Nama: "beban bunga", SaldoNormal: "DEBIT", Deskripsi: "Beban bunga yang harus dibayar"},
				{Nama: "rugi penjualan aset tetap", SaldoNormal: "DEBIT", Deskripsi: "Rugi dari penjualan aset tetap"},
				{Nama: "rugi penjualan surat berharga", SaldoNormal: "DEBIT", Deskripsi: "Rugi dari penjualan surat berharga"},
			},
		}
		// dataSaldoNormalAkun := [][]string{
		//   {"DEBIT", "DEBIT", "DEBIT", "DEBIT"},
		//   {"DEBIT", ""}
		// }
		//

		datas = make([]entity.Akun, 0, 34)

		for i, kelompokAkun := range kelompokAkuns {
			for j, akun := range dataAkun[i] {
				datas = append(datas, entity.Akun{
					ID:             ulid.MakeUlid().String(),
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

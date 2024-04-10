package static_data

import (
	"fmt"

	req_auth "github.com/be-sistem-informasi-konveksi/common/request/auth"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

// akuntansi

var DataKelompokAkun = []entity.KelompokAkun{
	// Aset
	{Base: entity.Base{ID: "01HP7DVBGQB9D016VRJN4HEZXM"}, Nama: "kas & bank", KategoriAkun: "ASET"},
	{Base: entity.Base{ID: "01HP7DVBGQB9D016VRJQNE3SW3"}, Nama: "piutang", KategoriAkun: "ASET"},
	{Base: entity.Base{ID: "01HP7DVBGRW9YXSCGCYZNBRNEN"}, Nama: "aset lancar lainnya", KategoriAkun: "ASET"},
	{Base: entity.Base{ID: "01HP7DVBGRW9YXSCGCZ0RK7RJQ"}, Nama: "aset tetap berwujud", KategoriAkun: "ASET"},
	{Base: entity.Base{ID: "01HP7DVBGRW9YXSCGCZ3WQPGJE"}, Nama: "aset tetap tidak berwujud", KategoriAkun: "ASET"},
	{Base: entity.Base{ID: "01HP7DVBGRW9YXSCGCZ6ARQK30"}, Nama: "aset tidak lancar lainnya", KategoriAkun: "ASET"},
	// Kewajiban
	{Base: entity.Base{ID: "01HP7DVBGRW9YXSCGCZ75DFW6Z"}, Nama: "hutang", KategoriAkun: "KEWAJIBAN"},
	{Base: entity.Base{ID: "01HP7DVBGRW9YXSCGCZ9KY625R"}, Nama: "kewajiban lancar lainnya", KategoriAkun: "KEWAJIBAN"},
	{Base: entity.Base{ID: "01HP7DVBGSM7P4AVN1TEFNC2G6"}, Nama: "kewajiban jangka panjang", KategoriAkun: "KEWAJIBAN"},
	// Modal
	{Base: entity.Base{ID: "01HP7DVBGSM7P4AVN1TG9179N2"}, Nama: "modal", KategoriAkun: "MODAL"},
	// Pendapatan
	{Base: entity.Base{ID: "01HP7DVBGSM7P4AVN1TJEYQZBQ"}, Nama: "pendapatan usaha", KategoriAkun: "PENDAPATAN"},
	{Base: entity.Base{ID: "01HP7DVBGTC06PXWT6F2T4K9CB"}, Nama: "pendapatan di luar usaha", KategoriAkun: "PENDAPATAN"},
	// Beban
	{Base: entity.Base{ID: "01HP7DVBGTC06PXWT6F6FB9DJF"}, Nama: "harga pokok penjualan", KategoriAkun: "BEBAN"},
	{Base: entity.Base{ID: "01HP7DVBGTC06PXWT6F9DES6BC"}, Nama: "beban usaha lainnya", KategoriAkun: "BEBAN"},
	{Base: entity.Base{ID: "01HP7DVBGTC06PXWT6FD1PCEWR"}, Nama: "beban diluar usaha lainnya", KategoriAkun: "BEBAN"},
}

var DataAkun = [][]entity.Akun{
	// Kas & Bank
	{
		{Base: entity.Base{ID: "01HP7DVBGTC06PXWT6FD66VERN"}, Nama: "kas", SaldoNormal: "DEBIT", Deskripsi: "Kas dari usaha"},
	},
	// Piutang
	{
		{Base: entity.Base{ID: "01HP7DVBGTC06PXWT6FF89WRAB"}, Nama: "piutang usaha", SaldoNormal: "DEBIT", Deskripsi: "Piutang dari penjualan barang atau jasa kepada pelanggan"},
	},
	// Aset Lancar Lainnya
	{
		{Base: entity.Base{ID: "01HP7DVBGTC06PXWT6FFSJ2FEA"}, Nama: "perlengkapan", SaldoNormal: "DEBIT", Deskripsi: "Perlengkapan yang digunakan dalam usaha"},
	},
	// Aset Tetap Berwujud
	{
		{Base: entity.Base{ID: "01HP7DVBGVHMSA4VHWXPQ3635H"}, Nama: "peralatan", SaldoNormal: "DEBIT", Deskripsi: "Peralatan yang digunakan dalam usaha"},
		{Base: entity.Base{ID: "01HP7DVBGVHMSA4VHWXRNS06R3"}, Nama: "mesin-mesin", SaldoNormal: "DEBIT", Deskripsi: "Mesin-mesin yang digunakan dalam usaha"},
	},
	// Aset Tetap Tidak Berwujud
	{
		{Base: entity.Base{ID: "01HP7DVBGVHMSA4VHWXW9B76CM"}, Nama: "hak cipta", SaldoNormal: "DEBIT", Deskripsi: "Hak cipta dari suatu karya atau lainnya"},
	},
	// Aset Tidak Lancar Lainnya
	{
		{Base: entity.Base{ID: "01HP7DVBGVHMSA4VHWXX1XTCG1"}, Nama: "aset tidak lancar lainnya", SaldoNormal: "DEBIT", Deskripsi: "Aset tidak lancar lainnya"},
	},
	// Hutang
	{
		{Base: entity.Base{ID: "01HP7DVBGVHMSA4VHWXZR27J7C"}, Nama: "hutang usaha", SaldoNormal: "KREDIT", Deskripsi: "Hutang dari pembelian barang atau jasa yang belum dibayar"},
	},
	// Kewajiban Lancar Lainnya
	{
		{Base: entity.Base{ID: "01HP7DVBGVHMSA4VHWY0RN5RJ6"}, Nama: "hutang sewa", SaldoNormal: "KREDIT", Deskripsi: "Hutang sewa yang merupakan kewajiban lancar"},
		{Base: entity.Base{ID: "01HP7DVBGVHMSA4VHWY1DQW3JK"}, Nama: "hutang gaji", SaldoNormal: "KREDIT", Deskripsi: "Hutang gaji kepada karyawan yang belum dibayar"},
		{Base: entity.Base{ID: "01HP7DVBGWR5ZR6C13R8HM48D5"}, Nama: "hutang pajak", SaldoNormal: "KREDIT", Deskripsi: "Hutang pajak yang harus dibayar dalam waktu dekat"},
	},
	// Kewajiban Jangka Panjang
	{
		{Base: entity.Base{ID: "01HP7DVBGWR5ZR6C13RBFX5S8K"}, Nama: "hutang bank", SaldoNormal: "KREDIT", Deskripsi: "Hutang kepada bank yang jatuh tempo dalam waktu lebih dari satu tahun"},
		{Base: entity.Base{ID: "01HP7DVBGWR5ZR6C13RDW0Z9H6"}, Nama: "hutang obligasi", SaldoNormal: "KREDIT", Deskripsi: "Hutang berupa obligasi yang jatuh tempo dalam waktu lebih dari satu tahun"},
	},
	// Modal
	{
		{Base: entity.Base{ID: "01HP7DVBGWR5ZR6C13RF3VG3XF"}, Nama: "modal pribadi", SaldoNormal: "KREDIT", Deskripsi: "Modal yang ditanamkan oleh pemilik secara pribadi"},
		{Base: entity.Base{ID: "01HP7DVBGWR5ZR6C13RFG776Z0"}, Nama: "modal saham", SaldoNormal: "KREDIT", Deskripsi: "Modal yang didapat dari pendapatan usaha"},
		{Base: entity.Base{ID: "01HP7DVBGWR5ZR6C13RFMCTDG1"}, Nama: "prive", SaldoNormal: "DEBIT", Deskripsi: "Pengambilan uang oleh pemilik untuk keperluan pribadi"},
	},
	// Pendapatan Usaha
	{
		{Base: entity.Base{ID: "01HP7DVBGWR5ZR6C13RFQKQ2A0"}, Nama: "pendapatan jasa", SaldoNormal: "KREDIT", Deskripsi: "Pendapatan dari jasa yang diberikan"},
	},
	// Pendapatan Diluar Usaha
	{
		{Base: entity.Base{ID: "01HP7DVBGWR5ZR6C13RJE693S0"}, Nama: "pendapatan bunga", SaldoNormal: "KREDIT", Deskripsi: "Pendapatan dari bunga yang diterima"},
	},
	// Harga Pokok Penjualan
	{
		{Base: entity.Base{ID: "01HP7DVBGX4JR0KETMEFRXQARB"}, Nama: "harga pokok pendapatan", SaldoNormal: "DEBIT", Deskripsi: "Harga pokok pendapatan usaha"},
	},
	// Beban Usaha Lainnya
	{
		{Base: entity.Base{ID: "01HP7DVBGX4JR0KETMEJXMZCMF"}, Nama: "beban gaji", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk membayar gaji karyawan"},
		{Base: entity.Base{ID: "01HP7DVBGX4JR0KETMEKJ1VJ5J"}, Nama: "beban sewa", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk sewa tempat atau fasilitas"},
		{Base: entity.Base{ID: "01HP7DVBGX4JR0KETMEKZ1SXFH"}, Nama: "beban asuransi", SaldoNormal: "DEBIT", Deskripsi: "Beban asuransi yang dibayarkan"},
		{Base: entity.Base{ID: "01HP7DVBGX4JR0KETMEQTSH53Y"}, Nama: "beban air, listrik dan telepon", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk utilitas"},
		{Base: entity.Base{ID: "01HP7DVBGX4JR0KETMER295ZPJ"}, Nama: "beban perlengkapan", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk membeli perlengkapan"},
		{Base: entity.Base{ID: "01HP7DVBGYA7XXS371H7YKKFEG"}, Nama: "beban penyusutan bangunan", SaldoNormal: "DEBIT", Deskripsi: "Beban penyusutan bangunan"},
		{Base: entity.Base{ID: "01HP7DVBGYA7XXS371HAH90EXW"}, Nama: "beban penyusutan kendaraan", SaldoNormal: "DEBIT", Deskripsi: "Beban penyusutan kendaraan"},
		{Base: entity.Base{ID: "01HP7DVBGYA7XXS371HE8QG3VE"}, Nama: "beban penyusutan peralatan", SaldoNormal: "DEBIT", Deskripsi: "Beban penyusutan peralatan"},
		{Base: entity.Base{ID: "01HP7DVBGZH26AQAGYEHS9FD91"}, Nama: "beban administrasi lainnya", SaldoNormal: "DEBIT", Deskripsi: "Beban administrasi lainnya"},
		{Base: entity.Base{ID: "01HP7DVBGZH26AQAGYEJQ4AWDB"}, Nama: "beban iklan", SaldoNormal: "DEBIT", Deskripsi: "Beban untuk biaya iklan"},
		{Base: entity.Base{ID: "01HP7DVBGZH26AQAGYENCBQDD5"}, Nama: "pajak penghasilan", SaldoNormal: "DEBIT", Deskripsi: "Beban pajak penghasilan yang harus dibayar"},
	},
	// Beban Diluar Usaha Lainnya
	{
		{Base: entity.Base{ID: "01HP7DVBGZH26AQAGYESCA4Z7M"}, Nama: "beban bunga", SaldoNormal: "DEBIT", Deskripsi: "Beban bunga yang harus dibayar"},
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

// user
var getId = func() string {
	return pkg.NewUlidPkg().MakeUlid().String()
}
var CredentialUsers = map[string]req_auth.Login{
	entity.RolesById[1]: {
		Username: "akun_direktur",
		Password: "direktur123456",
	},
	entity.RolesById[2]: {
		Username: "akun_bendahara",
		Password: "bendahara123456",
	},
	entity.RolesById[3]: {
		Username: "akun_admin",
		Password: "admin123456",
	},
	entity.RolesById[4]: {
		Username: "akun_manajerproduksi",
		Password: "manajerproduksi123456",
	},
	entity.RolesById[5]: {
		Username: "akun_supervisorbelanja",
		Password: "supservisorbelanja123456",
	},
}

var DefaultSupervisor = []entity.JenisSpv{
	{
		Base: entity.Base{
			ID: "01HSN8C1KV0TZBMHPTPR8AAZT8",
		},
		Nama: "belanja",
	},
}

var DefaultUsers = []entity.User{
	{
		Base: entity.Base{
			ID: "01HSN8C1KV0TZBMHPTPTJNTY88",
		},
		Nama:     "akun_direktur",
		Role:     entity.RolesById[1],
		Username: CredentialUsers[entity.RolesById[1]].Username,
		Password: func() string {
			pass, _ := helper.NewEncryptor().HashPassword(CredentialUsers[entity.RolesById[1]].Password)
			return pass
		}(),
		NoTelp: "+62895397290606",
		Alamat: "angantaka",
	},
	{
		Base: entity.Base{
			ID: "01HSN8C1PK0NNVS3AD3PHJ4D54",
		},
		Nama:     "akun_bendahara",
		Role:     entity.RolesById[2],
		Username: CredentialUsers[entity.RolesById[2]].Username,
		Password: func() string {
			pass, _ := helper.NewEncryptor().HashPassword(CredentialUsers[entity.RolesById[2]].Password)
			return pass
		}(),
		NoTelp: "+62895397290607",
		Alamat: "angantaka2",
	},
	{
		Base: entity.Base{
			ID: "01HSN8C1SYNMCFP6RZMDK3DDGE",
		},
		Nama:     "akun_admin",
		Role:     entity.RolesById[3],
		Username: CredentialUsers[entity.RolesById[3]].Username,
		Password: func() string {
			pass, _ := helper.NewEncryptor().HashPassword(CredentialUsers[entity.RolesById[3]].Password)
			return pass
		}(),
		NoTelp: "+62895397290608",
		Alamat: "angantaka3",
	},
	{
		Base: entity.Base{
			ID: "01HSN8C1WYR5CTWGQKMQ6JTFAR",
		},
		Nama:     "akun_manajerproduksi",
		Role:     entity.RolesById[4],
		Username: CredentialUsers[entity.RolesById[4]].Username,
		Password: func() string {
			pass, _ := helper.NewEncryptor().HashPassword(CredentialUsers[entity.RolesById[4]].Password)
			return pass
		}(),
		NoTelp: "+62895397290609",
		Alamat: "angantaka4",
	},
	{
		Base: entity.Base{
			ID: "01HSN8C1ZW2W5R29PYQBEXZACP",
		},
		Nama:     "akun_supervisorbelanja",
		Role:     entity.RolesById[5],
		Username: CredentialUsers[entity.RolesById[5]].Username,
		Password: func() string {
			pass, _ := helper.NewEncryptor().HashPassword(CredentialUsers[entity.RolesById[5]].Password)
			return pass
		}(),
		NoTelp:     "+62895397290610",
		Alamat:     "angantaka5",
		JenisSpvID: DefaultSupervisor[0].ID,
	},
}

package test_akuntansi

import (
	"os"
	"testing"
	"time"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/app/static_data"
	req_akuntansi_transaksi "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/transaksi"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

var dbt *gorm.DB
var tokens map[string]string
var app = fiber.New()

func cleanUp() {
	var ids []string
	for _, kelompok := range static_data.DataKelompokAkun {
		ids = append(ids, kelompok.ID)
	}
	if err := dbt.Unscoped().Delete(&entity.KelompokAkun{}, "id NOT IN (?)", ids).Error; err != nil {
		panic(helper.LogsError(err))
	}
	if err := dbt.Unscoped().Where("id = ?", idDataBayarInvoice).Delete(&entity.DataBayarInvoice{}).Error; err != nil {
		panic(helper.LogsError(err))
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.HutangPiutang{}).Error; err != nil {
		panic(helper.LogsError(err))
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.Transaksi{}).Error; err != nil {
		panic(helper.LogsError(err))
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.Kontak{}).Error; err != nil {
		panic(helper.LogsError(err))
	}
}

func TestMain(m *testing.M) {
	test.GetDB()
	dbt = test.DBT
	akuntansiH := handler_init.NewAkuntansiHandlerInit(dbt, test.Validator, test.UlidPkg, test.Excelize)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	akuntansiRoute := route.NewAkuntansiRoute(akuntansiH, authMid)

	test.GetTokens(dbt, authMid)
	tokens = test.Tokens
	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	akuntansiGroup := v1.Group("/akuntansi")
	akuntansiGroup.Route("/akun", akuntansiRoute.Akun)
	akuntansiGroup.Route("/kontak", akuntansiRoute.Kontak)
	akuntansiGroup.Route("/kelompok_akun", akuntansiRoute.KelompokAkun)
	akuntansiGroup.Route("/transaksi", akuntansiRoute.Transaksi)
	akuntansiGroup.Route("/hutang_piutang", akuntansiRoute.HutangPiutang)
	akuntansiGroup.Route("", akuntansiRoute.Akuntansi)
	// Run tests
	exitVal := m.Run()
	cleanUp()
	os.Exit(exitVal)
}

var ttTrStartDateAkuntansi string
var ttTrEndDateAkuntansi string
var ttTrYearMonthAkuntansi string

func setUpAkuntansiTr() {
	tt := time.Now()
	payloads := []req_akuntansi_transaksi.Create{
		{
			Tanggal:         tt.Format(time.RFC3339),
			BuktiPembayaran: []string{"1"},
			Keterangan:      "setoran modal pribadi",
			AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
				{
					AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
					Debit:  35000000,
				},
				{
					AkunID: "01HP7DVBGWR5ZR6C13RF3VG3XF", // modal pribadi
					Kredit: 35000000,
				},
			},
		},
		{
			Tanggal:         tt.Add(time.Hour * 24).Format(time.RFC3339),
			BuktiPembayaran: []string{"2"},
			Keterangan:      "membeli perlatan secara kredit",
			AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
				{
					AkunID: "01HP7DVBGVHMSA4VHWXPQ3635H", // peralatan
					Debit:  5500000,
				},
				{
					AkunID: "01HP7DVBGVHMSA4VHWXZR27J7C", // hutang usaha
					Kredit: 5500000,
				},
			},
		},
		{
			Tanggal:         tt.Add(time.Hour * 24 * 2).Format(time.RFC3339),
			BuktiPembayaran: []string{"3"},
			Keterangan:      "membeli perlengkapan secara tunai",
			AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
				{
					AkunID: "01HP7DVBGTC06PXWT6FFSJ2FEA", // perlengkapan
					Debit:  350000,
				},
				{
					AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
					Kredit: 350000,
				},
			},
		},
		{
			Tanggal:         tt.Add(time.Hour * 24 * 3).Format(time.RFC3339),
			BuktiPembayaran: []string{"4"},
			Keterangan:      "memperoleh pendapatan jasa",
			AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
				{
					AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
					Debit:  450000,
				},
				{
					AkunID: "01HP7DVBGWR5ZR6C13RFQKQ2A0", // pendapatan jasa
					Kredit: 450000,
				},
			},
		},
		{
			Tanggal:         tt.Add(time.Hour * 24 * 4).Format(time.RFC3339),
			BuktiPembayaran: []string{"4"},
			Keterangan:      "memperoleh pendapatan jasa",
			AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
				{
					AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
					Debit:  250000,
				},
				{
					AkunID: "01HP7DVBGTC06PXWT6FF89WRAB", // piutang usaha
					Debit:  450000,
				},
				{
					AkunID: "01HP7DVBGWR5ZR6C13RFQKQ2A0", // pendapatan jasa
					Kredit: 700000,
				},
			},
		},
		{
			Tanggal:         tt.Add(time.Hour * 24 * 5).Format(time.RFC3339),
			BuktiPembayaran: []string{"6"},
			Keterangan:      "pembayaran beban",
			AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
				{
					AkunID: "01HP7DVBGX4JR0KETMEQTSH53Y", // beban air, listrik dan telepon
					Debit:  250000,
				},
				{
					AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
					Kredit: 250000,
				},
			},
		},
		{
			Tanggal:         tt.Add(time.Hour * 24 * 6).Format(time.RFC3339),
			BuktiPembayaran: []string{"7"},
			Keterangan:      "memperoleh pendapatan jasa",
			AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
				{
					AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
					Debit:  600000,
				},
				{
					AkunID: "01HP7DVBGWR5ZR6C13RFQKQ2A0", // pendapatan jasa
					Kredit: 600000,
				},
			},
		},
		{
			Tanggal:         tt.Add(time.Hour * 24 * 7).Format(time.RFC3339),
			BuktiPembayaran: []string{"8"},
			Keterangan:      "pembayaran beban",
			AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
				{
					AkunID: "01HP7DVBGX4JR0KETMEKJ1VJ5J", // beban sewa
					Debit:  600000,
				},
				{
					AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
					Kredit: 600000,
				},
			},
		},
		{
			Tanggal:         tt.Add(time.Hour * 24 * 8).Format(time.RFC3339),
			BuktiPembayaran: []string{"9"},
			Keterangan:      "memperoleh pendapatan jasa",
			AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
				{
					AkunID: "01HP7DVBGTC06PXWT6FF89WRAB", // piutang usaha
					Debit:  600000,
				},
				{
					AkunID: "01HP7DVBGWR5ZR6C13RFQKQ2A0", // pendapatan jasa
					Kredit: 600000,
				},
			},
		},
		{
			Tanggal:         tt.Add(time.Hour * 24 * 9).Format(time.RFC3339),
			BuktiPembayaran: []string{"10"},
			Keterangan:      "pembayaran gaji",
			AyatJurnal: []req_akuntansi_transaksi.ReqAyatJurnal{
				{
					AkunID: "01HP7DVBGX4JR0KETMEJXMZCMF", // beban gaji
					Debit:  1000000,
				},
				{
					AkunID: "01HP7DVBGTC06PXWT6FD66VERN", // kas
					Kredit: 1000000,
				},
			},
		},
	}
	token := tokens[entity.RolesById[1]]
	for _, payload := range payloads {
		code, body, err := test.GetJsonTestRequestResponse(app, "POST", "/api/v1/akuntansi/transaksi", payload, &token)
		if err != nil {
			panic(err)
		}
		if body.Code != 201 {
			panic(body)
		}
		if code != 201 {
			panic(code)
		}
	}
	ts := tt.AddDate(0, 0, -1)
	te := tt.AddDate(0, 0, 20)
	ttTrStartDateAkuntansi = ts.Format(time.DateOnly)
	ttTrEndDateAkuntansi = te.Format(time.DateOnly)
	ttTrYearMonthAkuntansi = tt.Format("2006-01")
}

func TestEndPointAkuntansi(t *testing.T) {
	// kelompok akun
	AkuntansiCreateKelompokAkun(t)
	AkuntansiUpdateKelompokAkun(t)
	AkuntansiGetAllKelompokAkun(t)
	AkuntansiGetKelompokAkun(t)

	// akun
	AkuntansiCreateAkun(t)
	AkuntansiUpdateAkun(t)
	AkuntansiGetAllAkun(t)
	AkuntansiGetAkun(t)

	// kontak
	AkuntansiCreateKontak(t)
	AkuntansiUpdateKontak(t)
	AkuntansiGetAllKontak(t)

	// hutang piutang
	AkuntansiCreateHutangPiutang(t)
	AkuntansiGetAllHutangPiutang(t)
	AkuntansiCreateBayarHP(t)

	// transaksi
	AkuntansiCreateTransaksi(t)
	AkuntansiUpdateTransaksi(t)
	AkuntansiGetAllTransaksi(t)
	AkuntansiGetTransaksi(t)
	AkuntansiGetHistoryTransaksi(t)

	// delete all
	AkuntansiDeleteKelompokAkun(t)
	AkuntansiDeleteAkun(t)
	AkuntansiDeleteTransaksi(t)
	AkuntansiDeleteKontak(t)

	// setup akuntansi case
	cleanUp()
	setUpAkuntansiTr()
	// akuntansi
	AkuntansiGetJU(t)
	AkuntansiGetBB(t)
	AkuntansiGetNC(t)
	AkuntansiGetLB(t)

}

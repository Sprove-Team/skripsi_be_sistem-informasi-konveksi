package test_akuntansi

import (
	"os"
	"testing"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
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

var produk []entity.Produk
var bordir []entity.Bordir
var sablon []entity.Sablon
var kontak []entity.Kontak
var kategori []entity.KategoriProduk

func setUpData() {
	kategori = []entity.KategoriProduk{
		{
			Base: entity.Base{
				ID: test.UlidPkg.MakeUlid().String(),
			},
			Nama: "kategori 1",
		},
		{
			Base: entity.Base{
				ID: test.UlidPkg.MakeUlid().String(),
			},
			Nama: "kategori 2",
		},
	}
	produk = []entity.Produk{
		{
			Base: entity.Base{
				ID: test.UlidPkg.MakeUlid().String(),
			},
			Nama:             "produk 1",
			KategoriProdukID: kategori[0].ID,
			HargaDetails: []entity.HargaDetailProduk{
				{
					Base: entity.Base{
						ID: test.UlidPkg.MakeUlid().String(),
					},
					QTY:   1,
					Harga: 20000,
				},
				{
					Base: entity.Base{
						ID: test.UlidPkg.MakeUlid().String(),
					},
					QTY:   10,
					Harga: 15000,
				},
			},
		},
		{
			Base: entity.Base{
				ID: test.UlidPkg.MakeUlid().String(),
			},
			Nama:             "produk 2",
			KategoriProdukID: kategori[1].ID,
			HargaDetails: []entity.HargaDetailProduk{{
				Base: entity.Base{
					ID: test.UlidPkg.MakeUlid().String(),
				},
				QTY:   1,
				Harga: 10000,
			},
				{
					Base: entity.Base{
						ID: test.UlidPkg.MakeUlid().String(),
					},
					QTY:   10,
					Harga: 9000,
				}},
		},
	}
	bordir = []entity.Bordir{
		{
			Base: entity.Base{
				ID: test.UlidPkg.MakeUlid().String(),
			},
			Nama:  "bordir 1",
			Harga: 25000,
		},
		{
			Base: entity.Base{
				ID: test.UlidPkg.MakeUlid().String(),
			},
			Nama:  "bordir 2",
			Harga: 15000,
		},
	}
	sablon = []entity.Sablon{
		{
			Base: entity.Base{
				ID: test.UlidPkg.MakeUlid().String(),
			},
			Nama:  "sablon 1",
			Harga: 22000,
		},
		{
			Base: entity.Base{
				ID: test.UlidPkg.MakeUlid().String(),
			},
			Nama:  "sablon 2",
			Harga: 20000,
		},
	}
	kontak = []entity.Kontak{
		{
			Base: entity.Base{
				ID: "01E9CXBFVX2VH4P9V2FJJXVCZG",
			},
			Nama:       "John Doe",
			NoTelp:     "+6281234567890",
			Alamat:     "123 Main Street",
			Email:      "john.doe@example.com",
			Keterangan: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		},
		{
			Base: entity.Base{
				ID: "01E9CXBFVX2VH4P9V2FJJXVCZH",
			},
			Nama:       "Jane Smith",
			NoTelp:     "+6282345678901",
			Alamat:     "456 Elm Street",
			Email:      "jane.smith@example.com",
			Keterangan: "Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		},
	}

	if err := dbt.Create(kategori).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Create(produk).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Create(bordir).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Create(sablon).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Create(kontak).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
}

func cleanUp() {
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.HutangPiutang{}).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.Transaksi{}).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.DataBayarInvoice{}).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.DetailInvoice{}).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.Invoice{}).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.Kontak{}).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.KategoriProduk{}).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.HargaDetailProduk{}).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.Produk{}).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.Bordir{}).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.Sablon{}).Error; err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	test.GetDB()
	dbt = test.DBT
	setUpData()
	invoiceH := handler_init.NewInvoiceHandlerInit(dbt, test.Validator, test.UlidPkg, test.Encryptor)
	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	invoiceRoute := route.NewInvoiceRoute(invoiceH, authMid)

	test.GetTokens(dbt, authMid)
	tokens = test.Tokens
	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	invoiceGroup := v1.Group("/invoice")
	invoiceGroup.Route("/data_bayar", invoiceRoute.DataBayarInvoice)
	invoiceGroup.Route("/", invoiceRoute.Invoice)
	// Run tests
	exitVal := m.Run()
	cleanUp()
	os.Exit(exitVal)
}

func TestEndPointInvoice(t *testing.T) {
	InvoiceCreate(t)
}

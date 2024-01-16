package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"

	"github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	timeoutMid "github.com/be-sistem-informasi-konveksi/api/middleware/timeout"
	user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/config"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	helper "github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

func main() {
	// pkg
	validator := pkg.NewValidator()
	ulidPkg := pkg.NewUlidPkg()

	status_app := os.Getenv("APP_STATUS")
	var dbGormConf config.DBGorm
	if status_app == "PRODUCTION" {
		dbGormConf = config.DBGorm{
			DB_Username: os.Getenv("DB_USERNAME_PRODUCTION"),
			DB_Password: os.Getenv("DB_PASSWORD_PRODUCTION"),
			DB_Name:     os.Getenv("DB_NAME_PRODUCTION"),
			DB_Port:     os.Getenv("DB_PORT_PRODUCTION"),
			DB_Host:     os.Getenv("DB_HOST_PRODUCTION"),
		}
	} else {
		dbGormConf = config.DBGorm{
			DB_Username: os.Getenv("DB_USERNAME"),
			DB_Password: os.Getenv("DB_PASSWORD"),
			DB_Name:     os.Getenv("DB_NAME"),
			DB_Port:     os.Getenv("DB_PORT"),
			DB_Host:     os.Getenv("DB_HOST"),
		}
	}

	dbGorm := dbGormConf.InitDBGorm(ulidPkg)
	app := fiber.New()

	// helper
	encryptor := helper.NewEncryptor()

	// middleware
	userRepo := user.NewUserRepo(dbGorm)
	authMid := auth.NewAuthMiddleware(userRepo)
	timeoutMid := timeoutMid.NewTimeoutMiddleware()

	// domain
	api := app.Group("/api")
	v1 := api.Group("/v1", timeoutMid.Timeout(nil))
	// special metrics
	// prometheus := fiberprometheus.New("matrics-konveksi")
	// prometheus.RegisterAt(app, "/metrics")
	// app.Use(prometheus.Middleware)
	// v1.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	{
		// produk
		produkHandler := handler_init.NewProdukHandlerInit(dbGorm, validator, ulidPkg)
		produkRoute := route.NewProdukRoute(produkHandler, authMid)
		produkGroup := v1.Group("/produk")
		{
			produkGroup.Route("/harga_detail", produkRoute.HargaDetailProduk)
			produkGroup.Route("/kategori", produkRoute.KategoriProduk)
			produkGroup.Route("/", produkRoute.Produk)
		}

		// bordir
		bordirHandler := handler_init.NewBordirHandlerInit(dbGorm, validator, ulidPkg)
		bordirRoute := route.NewBordirRoute(bordirHandler, authMid)
		bordirGroup := v1.Group("/bordir")
		{
			bordirGroup.Route("/", bordirRoute.Bordir)
		}

		// sablon
		sablonHandler := handler_init.NewSablonHandlerInit(dbGorm, validator, ulidPkg)
		sablonRoute := route.NewSablonRoute(sablonHandler, authMid)
		sablonGroup := v1.Group("/sablon")
		{
			sablonGroup.Route("/", sablonRoute.Sablon)
		}

		// user
		userHandler := handler_init.NewUserHandlerInit(dbGorm, validator, ulidPkg, encryptor)
		userRoute := route.NewUserRoute(userHandler, authMid)
		userGroup := v1.Group("/user")
		{
			userGroup.Route("/jenis_spv", userRoute.JenisSpv)
			userGroup.Route("/", userRoute.User)
		}

		// akuntansi
		akuntansiHandler := handler_init.NewAkuntansiHandlerInit(dbGorm, validator, ulidPkg)
		akuntansiRoute := route.NewAkuntansiRoute(akuntansiHandler, authMid)
		akuntansiGroup := v1.Group("/akuntansi")
		{
			akuntansiGroup.Route("", akuntansiRoute.Akuntansi) // pelaporan akuntansi
			akuntansiGroup.Route("/akun", akuntansiRoute.Akun)
			akuntansiGroup.Route("/kelompok_akun", akuntansiRoute.KelompokAkun)
			akuntansiGroup.Route("/transaksi", akuntansiRoute.Transaksi)
		}

		// invoice
		invoiceHandler := handler_init.NewInvoiceHandlerInit(dbGorm, validator, ulidPkg)
		invoiceRoute := route.NewInvoiceRoute(invoiceHandler, authMid)
		invoiceGroup := v1.Group("/invoice")
		{
			invoiceGroup.Route("", invoiceRoute.Invoice)
		}

	}
	if os.Getenv("APP_STATUS") == "PRODUCTION" {
		app.Listen(os.Getenv("APP_PORT_PRODUCTION"))
	} else {
		app.Listen(os.Getenv("APP_PORT"))
	}
}

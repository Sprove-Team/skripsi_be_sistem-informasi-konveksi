package main

import (
	"os"

	"github.com/gofiber/fiber/v2"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	corsMid "github.com/be-sistem-informasi-konveksi/api/middleware/cors"
	timeoutMid "github.com/be-sistem-informasi-konveksi/api/middleware/timeout"
	userRepo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/config"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	helper "github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// pkg
	validator := pkg.NewValidator()
	ulidPkg := pkg.NewUlidPkg()
	jwtPkg := pkg.NewJwt(os.Getenv("JWT_TOKEN"), os.Getenv("JWT_REFTOKEN"))
	dbGormConf := config.DBGorm{
		DB_Username: os.Getenv("DB_USERNAME"),
		DB_Password: os.Getenv("DB_PASSWORD"),
		DB_Name:     os.Getenv("DB_NAME"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_Port:     os.Getenv("DB_PORT"),
	}

	dbGorm := dbGormConf.InitDBGorm(ulidPkg)
	app := fiber.New()

	// helper
	encryptor := helper.NewEncryptor()

	// middleware
	userRepo := userRepo.NewUserRepo(dbGorm)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	timeoutMid := timeoutMid.NewTimeoutMiddleware()
	corsMid := corsMid.NewCorsMiddleware()

	// domain
	app.Use(corsMid.Cors())
	api := app.Group("/api")
	v1 := api.Group("/v1", timeoutMid.Timeout(nil))
	// special metrics
	// prometheus := fiberprometheus.New("matrics-konveksi")
	// prometheus.RegisterAt(app, "/metrics")
	// app.Use(prometheus.Middleware)
	// v1.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	{
		// auth
		authHandler := handler_init.NewAuthHandlerInit(dbGorm, jwtPkg, validator, encryptor)
		authRoute := route.NewAuthRoute(authHandler)
		authGroup := v1.Group("/auth")
		{
			authGroup.Route("", authRoute.Auth)
		}
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
			akuntansiGroup.Route("/kontak", akuntansiRoute.Kontak)
			akuntansiGroup.Route("/kelompok_akun", akuntansiRoute.KelompokAkun)
			akuntansiGroup.Route("/transaksi", akuntansiRoute.Transaksi)
			akuntansiGroup.Route("/hutang_piutang", akuntansiRoute.HutangPiutang)
		}

		// invoice
		invoiceHandler := handler_init.NewInvoiceHandlerInit(dbGorm, validator, ulidPkg, encryptor)
		invoiceRoute := route.NewInvoiceRoute(invoiceHandler, authMid)
		invoiceGroup := v1.Group("/invoice")
		{
			invoiceGroup.Route("", invoiceRoute.Invoice)
			invoiceGroup.Route("/data_bayar", invoiceRoute.DataBayarInvoice)
		}

	}

	app.Listen(":" + os.Getenv("APP_PORT"))

}

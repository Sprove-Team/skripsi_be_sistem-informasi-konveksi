package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/logger"

	_ "time/tzdata"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	corsMid "github.com/be-sistem-informasi-konveksi/api/middleware/cors"
	timeoutMid "github.com/be-sistem-informasi-konveksi/api/middleware/timeout"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
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
	excelize := pkg.NewExcelizePkg()
	jwtPkg := pkg.NewJwt(os.Getenv("JWT_TOKEN"), os.Getenv("JWT_REFTOKEN"))
	dbGormConf := config.DBGorm{
		DB_Username: os.Getenv("DB_USERNAME"),
		DB_Password: os.Getenv("DB_PASSWORD"),
		DB_Name:     os.Getenv("DB_NAME"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_Port:     os.Getenv("DB_PORT"),
		LogLevel:    logger.Info,
	}
	// dbGormConf := config.DBGorm{
	// 	DB_Username: "root",
	// 	DB_Password: os.Getenv("DB_PASSWORD_ROOT"),
	// 	DB_Name:     os.Getenv("DB_NAME_TEST"),
	// 	DB_HOST:     os.Getenv("DB_HOST_TEST"),
	// 	DB_Port:     os.Getenv("DB_PORT_TEST"),
	// 	LogLevel:    logger.Info,
	// }
	if os.Getenv("ENVIRONMENT") == "DEVELOPMENT" {
		dbGormConf.DB_HOST = "localhost"
		dbGormConf.LogLevel = logger.Info
	}
	dbGorm := dbGormConf.InitDBGorm(ulidPkg)
	app := fiber.New()

	// helper
	encryptor := helper.NewEncryptor()

	// middleware
	repoUser := repo_user.NewUserRepo(dbGorm)
	authMid := middleware_auth.NewAuthMiddleware(repoUser)
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
		authHandler := handler_init.NewAuthHandlerInit(dbGorm, jwtPkg, validator, ulidPkg, encryptor)
		authRoute := route.NewAuthRoute(authHandler, authMid)
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
		akuntansiHandler := handler_init.NewAkuntansiHandlerInit(dbGorm, validator, ulidPkg, excelize)
		akuntansiRoute := route.NewAkuntansiRoute(akuntansiHandler, authMid)
		akuntansiGroup := v1.Group("/akuntansi")
		{
			akuntansiGroup.Route("/akun", akuntansiRoute.Akun)
			akuntansiGroup.Route("/kontak", akuntansiRoute.Kontak)
			akuntansiGroup.Route("/kelompok_akun", akuntansiRoute.KelompokAkun)
			akuntansiGroup.Route("/transaksi", akuntansiRoute.Transaksi)
			akuntansiGroup.Route("/hutang_piutang", akuntansiRoute.HutangPiutang)
			akuntansiGroup.Route("", akuntansiRoute.Akuntansi) // pelaporan akuntansi
		}

		// invoice
		invoiceHandler := handler_init.NewInvoiceHandlerInit(dbGorm, validator, ulidPkg, encryptor)
		invoiceRoute := route.NewInvoiceRoute(invoiceHandler, authMid)
		invoiceGroup := v1.Group("/invoice")
		{
			invoiceGroup.Route("/:id/data_bayar", invoiceRoute.DataBayarByInvoiceId)
			invoiceGroup.Route("/data_bayar", invoiceRoute.DataBayar)
			invoiceGroup.Route("", invoiceRoute.Invoice)
		}

		// tugas
		tugasHandler := handler_init.NewTugasHandlerInit(dbGorm, validator, ulidPkg)
		tugasRoute := route.NewTugasRoute(tugasHandler, authMid)
		tugasGroup := v1.Group("/tugas")
		{
			tugasGroup.Route("/sub_tugas", tugasRoute.SubTugas)
			tugasGroup.Route("", tugasRoute.Tugas)
		}

		// profile
		profileHandler := handler_init.NewProfileHandlerInit(dbGorm, validator, ulidPkg, encryptor)
		profileRoute := route.NewProfileRoute(profileHandler, authMid)
		profileGroup := v1.Group("/profile")
		{
			profileGroup.Route("", profileRoute.Profile)
		}

	}

	app.Listen(":" + os.Getenv("APP_PORT"))

}

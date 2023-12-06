package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"

	timeoutMid "github.com/be-sistem-informasi-konveksi/api/middleware/timeout"
	"github.com/be-sistem-informasi-konveksi/app/config"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	helper "github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

func main() {

	dbGormConf := config.DBGormConf{
		DB_Username: os.Getenv("DB_USERNAME"),
		DB_Password: os.Getenv("DB_PASSWORD"),
		DB_Name:     os.Getenv("DB_NAME"),
		DB_Port:     os.Getenv("DB_PORT"),
		DB_Host:     os.Getenv("DB_HOST"),
	}
	dbGorm := dbGormConf.InitDBGormConf()
	app := fiber.New()

	// pkg
	validator := pkg.NewValidator()
	// uuidGen := pkg.NewGoogleUUID()
	ulidPkg := pkg.NewUlidPkg()

	// helper
	paginate := helper.NewPaginate()
	encryptor := helper.NewEncryptor()
	
	// middleware
	// userRepo := userRepo.NewUserRepo(dbGorm)
	// authMid := _authMid.NewAuthMiddleware(userRepo)
	timeoutMid := timeoutMid.NewTimeoutMiddleware()

	// domain

	api := app.Group("/api")
	v1 := api.Group("/v1")
	{
		// domain direktur
		direktur := v1.Group("/direktur", timeoutMid.Timeout(nil))	
		{
			// produk
			produkHandler := handler_init.NewProdukHandlerInit(dbGorm, validator, ulidPkg, paginate)
			produkRoute := route.NewProdukRoute(produkHandler)
			produkGroup := direktur.Group("/produk")
			{
				produkGroup.Route("/harga_detail", produkRoute.HargaDetailProduk)
				produkGroup.Route("/kategori", produkRoute.KategoriProduk)
				produkGroup.Route("/", produkRoute.Produk)
			}
	
			// bordir
			bordirHandler := handler_init.NewBordirHandlerInit(dbGorm, validator, ulidPkg, paginate)
			bordirRoute := route.NewBordirRoute(bordirHandler)
			bordirGroup := direktur.Group("/bordir")
			{
				bordirGroup.Route("/", bordirRoute.Bordir)
			}
	
			// sablon
			sablonHandler := handler_init.NewSablonHandlerInit(dbGorm, validator, ulidPkg, paginate)
			sablonRoute := route.NewSablonRoute(sablonHandler)
			sablonGroup := direktur.Group("/sablon")
			{
				sablonGroup.Route("/", sablonRoute.Sablon)
			}
	
			//user
			userHandler := handler_init.NewUserHandlerInit(dbGorm, validator, ulidPkg, paginate, encryptor)
			userRoute := route.NewUserRoute(userHandler)
			userGroup := direktur.Group("/user")
			{
				userGroup.Route("/jenis_spv", userRoute.JenisSpv)
				userGroup.Route("/", userRoute.User)
			}
		}
	}


	app.Listen(":8000")
}

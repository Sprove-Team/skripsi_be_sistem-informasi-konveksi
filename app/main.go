package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"

	_timeoutMid "github.com/be-sistem-informasi-konveksi/api/middleware/timeout"
	"github.com/be-sistem-informasi-konveksi/app/config"
	_handler_init "github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	helper "github.com/be-sistem-informasi-konveksi/helper"
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

	// helper
	validator := helper.NewValidator()
	uuidGen := helper.NewGoogleUUID()
	paginate := helper.NewPaginate()
	encryptor := helper.NewEncryptor()
	
	// middleware
	// userRepo := userRepo.NewUserRepo(dbGorm)
	// authMid := _authMid.NewAuthMiddleware(userRepo)
	timeoutMid := _timeoutMid.NewTimeoutMiddleware()

	// domain
	api := app.Group("/api")
	v1 := api.Group("/v1")
	{
		// domain direktur
		direktur := v1.Group("/direktur", timeoutMid.Timeout(nil))	
		{
			// produk
			produkHandler := _handler_init.NewProdukHandlerInit(dbGorm, validator, uuidGen, paginate)
			produkRoute := route.NewProdukRoute(produkHandler)
			produkGroup := direktur.Group("/produk")
			{
				produkGroup.Route("/harga_detail", produkRoute.HargaDetailProduk)
				produkGroup.Route("/kategori", produkRoute.KategoriProduk)
				produkGroup.Route("/", produkRoute.Produk)
			}
	
			// bordir
			bordirHandler := _handler_init.NewBordirHandlerInit(dbGorm, validator, uuidGen, paginate)
			bordirRoute := route.NewBordirRoute(bordirHandler)
			bordirGroup := direktur.Group("/bordir")
			{
				bordirGroup.Route("/", bordirRoute.Bordir)
			}
	
			// sablon
			sablonHandler := _handler_init.NewSablonHandlerInit(dbGorm, validator, uuidGen, paginate)
			sablonRoute := route.NewSablonRoute(sablonHandler)
			sablonGroup := direktur.Group("/sablon")
			{
				sablonGroup.Route("/", sablonRoute.Sablon)
			}
	
			//user
			userHandler := _handler_init.NewUserHandlerInit(dbGorm, validator, uuidGen, paginate, encryptor)
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

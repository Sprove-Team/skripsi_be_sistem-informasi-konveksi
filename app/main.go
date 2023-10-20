package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"

	"github.com/be-sistem-informasi-konveksi/app/config"
	"github.com/be-sistem-informasi-konveksi/common/handler_init/direktur"
	helper "github.com/be-sistem-informasi-konveksi/helper"
	midGlobal "github.com/be-sistem-informasi-konveksi/middleware/global"
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

	// handler init
	direkturHandler := direktur.NewDirekturHandlerInit(dbGorm, validator, uuidGen)

	// route
	app.Get("/direktur/produk/:id", midGlobal.TimeoutMid(direkturHandler.ProdukHandler().GetById, nil))
	app.Post("/direktur/produk", direkturHandler.ProdukHandler().Create)

	app.Post("/direktur/produk/kategori", direkturHandler.KategoriProdukHandler().Create)
	app.Put("/direktur/produk/:id", direkturHandler.KategoriProdukHandler().Update)
	app.Delete("/direktur/produk/kategori/:id", direkturHandler.KategoriProdukHandler().Delete)

	app.Post("/direktur/produk/harga_detail", direkturHandler.HargaDetailProdukHandler().Create)

	app.Listen(":8000")
}

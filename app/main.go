package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"

	"github.com/be-sistem-informasi-konveksi/app/config"
	"github.com/be-sistem-informasi-konveksi/common/handler_init"
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
	paginate := helper.NewPaginate()

	// handler init
	produkHandler := handler_init.NewProdukHandlerInit(dbGorm, validator, uuidGen, paginate)

	// route
	api := app.Group("/api")
	v1 := api.Group("/v1")
	direktur := v1.Group("/direktur", midGlobal.TimeoutMid(nil))
	{
		produkData := direktur.Group("/produk")
		kategoriProduk := produkData.Group("/kategori")
		hargaDetailProduk := produkData.Group("/harga_detail")
		{
			produkData.Get("", produkHandler.ProdukHandler().GetAll)
			kategoriProduk.Get("", produkHandler.KategoriProdukHandler().GetAll)
			hargaDetailProduk.Get("", produkHandler.HargaDetailProdukHandler().GetAll) // tidak perlu isi ini (lihat frontend dulu)

			produkData.Get("/:id", produkHandler.ProdukHandler().GetById)
			kategoriProduk.Get("/:id", produkHandler.KategoriProdukHandler().GetById)
			hargaDetailProduk.Get("/:produk_id", produkHandler.HargaDetailProdukHandler().GetByProdukId)

			produkData.Post("", produkHandler.ProdukHandler().Create)
			kategoriProduk.Post("", produkHandler.KategoriProdukHandler().Create)
			hargaDetailProduk.Post("", produkHandler.HargaDetailProdukHandler().Create)

			produkData.Put("/:id", produkHandler.ProdukHandler().Update)
			kategoriProduk.Put("/:id", produkHandler.KategoriProdukHandler().Update)
			hargaDetailProduk.Put("/:id", produkHandler.HargaDetailProdukHandler().Update)

			produkData.Delete("/:id", produkHandler.ProdukHandler().Delete)
			kategoriProduk.Delete("/:id", produkHandler.KategoriProdukHandler().Delete)
			hargaDetailProduk.Delete("/:id", produkHandler.HargaDetailProdukHandler().Delete)
			hargaDetailProduk.Delete("/:produk_id", produkHandler.HargaDetailProdukHandler().DeleteByProdukId)
		}
	}

	app.Listen(":8000")
}

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis/v3"
	_ "github.com/joho/godotenv/autoload"

	handler_init "github.com/be-sistem-informasi-konveksi/api/handler/_init"
	midGlobal "github.com/be-sistem-informasi-konveksi/api/middleware/global"
	"github.com/be-sistem-informasi-konveksi/app/config"
	helper "github.com/be-sistem-informasi-konveksi/helper"
)

func main() {
	cache := redis.New(redis.Config{
		URL: fmt.Sprintf("redis://default:%s@%s:%s/0",
			os.Getenv("REDIS_PASSWORD"),
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT")),
		Reset: false,
	})

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

	// handler init
	produkHandler := handler_init.NewProdukHandlerInit(dbGorm, validator, uuidGen, paginate)
	bordirHandler := handler_init.NewBordirHandlerInit(dbGorm, validator, uuidGen, paginate)
	sablonHandler := handler_init.NewSablonHandlerInit(dbGorm, validator, uuidGen, paginate)
	userHandler := handler_init.NewUserHandlerInit(dbGorm, validator, uuidGen, paginate, encryptor)
	// middleware
	// authMid := auth.NewAuthMiddleware()

	// route
	api := app.Group("/api")
	v1 := api.Group("/v1")
	direktur := v1.Group("/direktur", midGlobal.TimeoutMid(nil))
	{
		// produk
		produkData := direktur.Group("/produk")
		kategoriProduk := produkData.Group("/kategori")
		hargaDetailProduk := produkData.Group("/harga_detail")
		{
			produkData.Get("", produkHandler.ProdukHandler().GetAll)
			kategoriProduk.Get("", produkHandler.KategoriProdukHandler().GetAll)

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
		// bordir
		bordir := direktur.Group("/bordir")
		{
			bordir.Get("", bordirHandler.BordirHandler().GetAll)

			bordir.Get("/:id", bordirHandler.BordirHandler().GetById)

			bordir.Post("", bordirHandler.BordirHandler().Create)

			bordir.Put("/:id", bordirHandler.BordirHandler().Update)

			bordir.Delete("/:id", bordirHandler.BordirHandler().Delete)
		}
		// sablon
		sablon := direktur.Group("/sablon")
		{
			sablon.Get("", sablonHandler.SablonHandler().GetAll)

			sablon.Get("/:id", sablonHandler.SablonHandler().GetById)

			sablon.Post("", sablonHandler.SablonHandler().Create)

			sablon.Put("/:id", sablonHandler.SablonHandler().Update)

			sablon.Delete("/:id", sablonHandler.SablonHandler().Delete)
		}
		// user
		// user := direktur.Group("/user", authMid.Auth(os.Getenv("token_direktur"), "direktur"))
		user := direktur.Group("/user")
		{

			user.Get("", userHandler.UserHandler().GetAll)
			user.Post("", userHandler.UserHandler().Create)
			user.Put("/:id", userHandler.UserHandler().Update)
			user.Delete("/:id", userHandler.UserHandler().Delete)
		}
	}
	type redisReq struct {
		Key  string `json:"key"`
		Data string `json:"data"`
	}
	api.Post("/redis", func(c *fiber.Ctx) error {
		req := new(redisReq)
		c.BodyParser(req)
		log.Println(req)
		err := cache.Set(req.Key, []byte(req.Data), 20*time.Second)
		if err != nil {
			log.Println(err)
			return c.SendStatus(500)
		}
		return c.SendStatus(200)
	})
	api.Get("/redis/:key", func(c *fiber.Ctx) error {
		key := c.Params("key", "")
		log.Println(key)
		dat, err := cache.Get(key)
		if err != nil {
			return c.SendStatus(404)
		}
		log.Println(dat)
		return c.Status(200).SendString(string(dat))
	})

	app.Listen(":8000")
}

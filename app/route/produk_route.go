package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type ProdukRoute interface {
	Produk(router fiber.Router)
	KategoriProduk(router fiber.Router)
	HargaDetailProduk(router fiber.Router)
}

type produkRoute struct {
	h    handler_init.ProdukHandlerInit
	auth auth.AuthMidleware
}

func NewProdukRoute(h handler_init.ProdukHandlerInit, auth auth.AuthMidleware) ProdukRoute {
	return &produkRoute{h, auth}
}

func (ro *produkRoute) Produk(router fiber.Router) {
	router.Get("", ro.h.ProdukHandler().GetAll)
	router.Get("/:id", ro.h.ProdukHandler().GetById)
	router.Post("", ro.h.ProdukHandler().Create)
	router.Put("/:id", ro.h.ProdukHandler().Update)
	router.Delete("/:id", ro.h.ProdukHandler().Delete)
}

func (ro *produkRoute) KategoriProduk(router fiber.Router) {
	router.Get("", ro.h.KategoriProdukHandler().GetAll)
	router.Get("/:id", ro.h.KategoriProdukHandler().GetById)
	router.Post("", ro.h.KategoriProdukHandler().Create)
	router.Put("/:id", ro.h.KategoriProdukHandler().Update)
	router.Delete("/:id", ro.h.KategoriProdukHandler().Delete)
}

func (ro *produkRoute) HargaDetailProduk(router fiber.Router) {
	router.Get("/:produk_id", ro.h.HargaDetailProdukHandler().GetByProdukId)
	router.Post("", ro.h.HargaDetailProdukHandler().CreateByProdukId)
	router.Put("/:produk_id", ro.h.HargaDetailProdukHandler().UpdateByProdukId)
	router.Delete("/:id", ro.h.HargaDetailProdukHandler().Delete)
	router.Delete("/:produk_id", ro.h.HargaDetailProdukHandler().DeleteByProdukId)
}

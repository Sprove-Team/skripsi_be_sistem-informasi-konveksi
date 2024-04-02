package route

import (
	"github.com/gofiber/fiber/v2"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type ProdukRoute interface {
	Produk(router fiber.Router)
	KategoriProduk(router fiber.Router)
	HargaDetailProduk(router fiber.Router)
}

type produkRoute struct {
	h    handler_init.ProdukHandlerInit
	auth middleware_auth.AuthMidleware
}

func NewProdukRoute(h handler_init.ProdukHandlerInit, auth middleware_auth.AuthMidleware) ProdukRoute {
	return &produkRoute{h, auth}
}

func (ro *produkRoute) Produk(router fiber.Router) {
	auth := ro.auth.Authorization([]string{"DIREKTUR"})
	auth2 := ro.auth.Authorization([]string{"DIREKTUR", "ADMIN"})
	router.Get("", auth2, ro.h.ProdukHandler().GetAll)
	router.Get("/:id", auth, ro.h.ProdukHandler().GetById)
	router.Post("", auth, ro.h.ProdukHandler().Create)
	router.Put("/:id", auth, ro.h.ProdukHandler().Update)
	router.Delete("/:id", auth, ro.h.ProdukHandler().Delete)
}

func (ro *produkRoute) KategoriProduk(router fiber.Router) {
	auth := ro.auth.Authorization([]string{"DIREKTUR"})
	auth2 := ro.auth.Authorization([]string{"DIREKTUR", "ADMIN"})
	router.Get("", auth2, ro.h.KategoriProdukHandler().GetAll)
	router.Get("/:id", auth, ro.h.KategoriProdukHandler().GetById)
	router.Post("", auth, ro.h.KategoriProdukHandler().Create)
	router.Put("/:id", auth, ro.h.KategoriProdukHandler().Update)
	router.Delete("/:id", auth, ro.h.KategoriProdukHandler().Delete)
}

func (ro *produkRoute) HargaDetailProduk(router fiber.Router) {
	router.Use(ro.auth.Authorization([]string{"DIREKTUR"}))
	router.Get("/:produk_id", ro.h.HargaDetailProdukHandler().GetByProdukId)
	router.Post("", ro.h.HargaDetailProdukHandler().Create)
	router.Put("/:id", ro.h.HargaDetailProdukHandler().Update)
	router.Delete("/:id", ro.h.HargaDetailProdukHandler().Delete)
}

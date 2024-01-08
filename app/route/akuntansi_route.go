package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type AkuntansiRoute interface {
	Akun(router fiber.Router)
	KelompokAkun(router fiber.Router)
	Transaksi(router fiber.Router)
	Akuntansi(router fiber.Router)
}

type akuntansiRoute struct {
	h handler_init.AkuntansiHandlerInit
}

func NewAkuntansiRoute(h handler_init.AkuntansiHandlerInit) AkuntansiRoute {
	return &akuntansiRoute{h}
}

func (ro *akuntansiRoute) Akun(router fiber.Router) {
	router.Get("", ro.h.AkunHandler().GetAll)
	router.Get("/:id", ro.h.AkunHandler().GetById)
	router.Post("", ro.h.AkunHandler().Create)
	router.Put("/:id", ro.h.AkunHandler().Update)
	router.Delete("/:id", ro.h.AkunHandler().Delete)
}

func (ro *akuntansiRoute) KelompokAkun(router fiber.Router) {
	router.Post("", ro.h.KelompokAkunHandler().Create)
	router.Get("", ro.h.KelompokAkunHandler().GetAll)
	router.Put("/:id", ro.h.KelompokAkunHandler().Update)
	router.Delete("/:id", ro.h.KelompokAkunHandler().Delete)
}

func (ro *akuntansiRoute) Transaksi(router fiber.Router) {
	router.Get("/history", ro.h.Transaksi().GetHistory)
	router.Get("", ro.h.Transaksi().GetAll)
	router.Get("/:id", ro.h.Transaksi().GetById)
	router.Post("", ro.h.Transaksi().Create)
	router.Put("/:id", ro.h.Transaksi().Update)
	router.Delete("/:id", ro.h.Transaksi().Delete)
}

func (ro *akuntansiRoute) Akuntansi(router fiber.Router) {
	router.Get("/jurnal_umum", ro.h.Akuntansi().GetAllJU)
	router.Get("/buku_besar", ro.h.Akuntansi().GetAllBB)
	router.Get("/neraca_saldo", ro.h.Akuntansi().GetAllNC)
	router.Get("/laba_rugi", ro.h.Akuntansi().GetAllLBR)
}

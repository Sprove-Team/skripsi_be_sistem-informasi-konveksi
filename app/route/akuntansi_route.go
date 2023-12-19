package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type AkuntansiRoute interface {
	Akun(router fiber.Router)
	GolonganAkun(router fiber.Router)
	KelompokAkun(router fiber.Router)
	Transaksi(router fiber.Router)
}

type akuntansiRoute struct {
	h handler_init.AkuntansiHandlerInit
}

func NewAkuntansiRoute(h handler_init.AkuntansiHandlerInit) AkuntansiRoute {
	return &akuntansiRoute{h}
}

func (ro *akuntansiRoute) Akun(router fiber.Router) {
	router.Get("", ro.h.AkunHandler().GetAll)
	router.Post("", ro.h.AkunHandler().Create)
}

func (ro *akuntansiRoute) GolonganAkun(router fiber.Router) {
	router.Post("", ro.h.GolonganAkunHandler().Create)
	// Add other routes specific to Golongan Akun as needed
}

func (ro *akuntansiRoute) KelompokAkun(router fiber.Router) {
	router.Post("", ro.h.KelompokAkunHandler().Create)
	// Add other routes specific to Kelompok Akun as needed
}

func (ro *akuntansiRoute) Transaksi(router fiber.Router) {
	router.Post("", ro.h.Transaksi().Create)
}

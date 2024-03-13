package route

import (
	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/gofiber/fiber/v2"
)

type TugasRoute interface {
	Tugas(router fiber.Router)
	// User(router fiber.Router)
}

type tugasRoute struct {
	h    handler_init.TugasHandlerInit
	auth middleware_auth.AuthMidleware
}

func NewTugasRoute(h handler_init.TugasHandlerInit, auth middleware_auth.AuthMidleware) TugasRoute {
	return &tugasRoute{h, auth}
}

func (ro *tugasRoute) Tugas(router fiber.Router) {
	router.Use(ro.auth.Authorization([]string{"DIREKTUR", "MANAJER_PRODUKSI"}))
	router.Post("", ro.h.TugasHandler().Create)
	router.Get("/:invoice_id", ro.h.TugasHandler().GetByInvoiceId)
}

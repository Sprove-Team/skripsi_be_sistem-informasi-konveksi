package route

import (
	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/gofiber/fiber/v2"
)

type TugasRoute interface {
	Tugas(router fiber.Router)
	SubTugas(router fiber.Router)
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
	router.Put("/:id", ro.h.TugasHandler().Update)
	router.Delete("/:id", ro.h.TugasHandler().Delete)
	router.Get("", ro.h.TugasHandler().GetAll)
	router.Get("/:id", ro.h.TugasHandler().GetById)
	router.Get("/:invoice_id", ro.h.TugasHandler().GetByInvoiceId)
}

func (ro *tugasRoute) SubTugas(router fiber.Router) {
	router.Post("/:tugas_id", ro.h.SubTugasHandler().CreateByTugasId)
	router.Put("/:id", ro.h.SubTugasHandler().Update)
	router.Delete("/:id", ro.h.SubTugasHandler().Delete)
}

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
	auth := ro.auth.Authorization([]string{"DIREKTUR", "MANAJER_PRODUKSI"})
	auth2 := ro.auth.Authorization([]string{"DIREKTUR", "MANAJER_PRODUKSI", "SUPERVISOR"})

	router.Post("", auth, ro.h.TugasHandler().Create)
	router.Put("/:id", auth, ro.h.TugasHandler().Update)
	router.Delete("/:id", auth, ro.h.TugasHandler().Delete)
	router.Get("", auth2, ro.h.TugasHandler().GetAll)
	router.Get("/:id", auth2, ro.h.TugasHandler().GetById)
	router.Get("/invoice/:invoice_id", auth2, ro.h.TugasHandler().GetByInvoiceId)
}

func (ro *tugasRoute) SubTugas(router fiber.Router) {
	auth := ro.auth.Authorization([]string{"DIREKTUR", "MANAJER_PRODUKSI"})
	auth2 := ro.auth.Authorization([]string{"DIREKTUR", "MANAJER_PRODUKSI", "SUPERVISOR"})

	router.Post("/:tugas_id", auth, ro.h.SubTugasHandler().CreateByTugasId)
	router.Put("/:id", auth2, ro.h.SubTugasHandler().Update)
	router.Delete("/:id", auth, ro.h.SubTugasHandler().Delete)
}

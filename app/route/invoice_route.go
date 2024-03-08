package route

import (
	"github.com/gofiber/fiber/v2"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type InvoiceRoute interface {
	Invoice(router fiber.Router)
	DataBayarInvoice(router fiber.Router)
}

type invoiceRoute struct {
	h    handler_init.InvoiceHandlerInit
	auth middleware_auth.AuthMidleware
}

func NewInvoiceRoute(h handler_init.InvoiceHandlerInit, auth middleware_auth.AuthMidleware) InvoiceRoute {
	return &invoiceRoute{h, auth}
}

func (ro *invoiceRoute) Invoice(router fiber.Router) {
	router.Use(ro.auth.Authorization([]string{"DIREKTUR", "ADMIN", "MANAJER_PRODUKSI"}))
	router.Get("", ro.h.InvoiceHandler().GetAll)
	router.Get("/:id", ro.h.InvoiceHandler().GetById)
	router.Post("", ro.h.InvoiceHandler().Create)
	router.Put("/:id", ro.h.InvoiceHandler().Update)
	router.Delete("/:id", ro.h.InvoiceHandler().Delete)
}

func (ro *invoiceRoute) DataBayarInvoice(router fiber.Router) {
	router.Get("/:invoice_id", ro.h.DataBayarInvoiceHandler().GetByInvoiceId)
}

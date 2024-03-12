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
	auth1 := ro.auth.Authorization([]string{"DIREKTUR", "ADMIN"})
	auth2 := ro.auth.Authorization([]string{"DIREKTUR", "ADMIN", "MANAJER_PRODUKSI"})
	router.Get("", auth2, ro.h.InvoiceHandler().GetAll)
	router.Get("/:id", auth2, ro.h.InvoiceHandler().GetById)
	router.Post("", auth1, ro.h.InvoiceHandler().Create)
	router.Put("/:id", auth2, ro.h.InvoiceHandler().Update)
	router.Delete("/:id", auth2, ro.h.InvoiceHandler().Delete)
}

func (ro *invoiceRoute) DataBayarInvoice(router fiber.Router) {
	auth := ro.auth.Authorization([]string{"DIREKTUR", "ADMIN", "BENDAHARA"})
	router.Get("/:invoice_id", auth, ro.h.DataBayarInvoiceHandler().GetByInvoiceId)
	router.Put("/:invoice_id", auth, ro.h.DataBayarInvoiceHandler().Create)
	router.Put("/:id", auth, ro.h.DataBayarInvoiceHandler().Update)
	router.Put("/:id", auth, ro.h.DataBayarInvoiceHandler().Delete)
}

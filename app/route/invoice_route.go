package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type InvoiceRoute interface {
	Invoice(router fiber.Router)
}

type invoiceRoute struct {
	h    handler_init.InvoiceHandlerInit
	auth auth.AuthMidleware
}

func NewInvoiceRoute(h handler_init.InvoiceHandlerInit, auth auth.AuthMidleware) InvoiceRoute {
	return &invoiceRoute{h, auth}
}

func (ro *invoiceRoute) Invoice(router fiber.Router) {
	router.Post("", ro.h.InvoiceHandler().Create)
}

// func (ro *invoiceRoute) StatusProduksi(router fiber.Router) {

// }

package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type InvoiceRoute interface {
	Invoice(router fiber.Router)
}

type invoiceRoute struct {
	h handler_init.InvoiceHandlerInit
}

func NewInvoiceRoute(h handler_init.InvoiceHandlerInit) InvoiceRoute {
	return &invoiceRoute{h}
}

func (ro *invoiceRoute) Invoice(router fiber.Router) {
	router.Post("", ro.h.InvoiceHandler().Create)
}

package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type SablonRoute interface {
	Sablon(router fiber.Router)
}

type sablonRoute struct {
	h    handler_init.SablonHandlerInit
	auth auth.AuthMidleware
}

func NewSablonRoute(h handler_init.SablonHandlerInit, auth auth.AuthMidleware) SablonRoute {
	return &sablonRoute{h, auth}
}

func (ro *sablonRoute) Sablon(router fiber.Router) {
	router.Get("", ro.h.SablonHandler().GetAll)
	router.Get("/:id", ro.h.SablonHandler().GetById)
	router.Post("", ro.h.SablonHandler().Create)
	router.Put("/:id", ro.h.SablonHandler().Update)
	router.Delete("/:id", ro.h.SablonHandler().Delete)
}

package route

import (
	"github.com/gofiber/fiber/v2"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type BordirRoute interface {
	Bordir(router fiber.Router)
}

type bordirRoute struct {
	h    handler_init.BordirHandlerInit
	auth middleware_auth.AuthMidleware
}

func NewBordirRoute(h handler_init.BordirHandlerInit, auth middleware_auth.AuthMidleware) BordirRoute {
	return &bordirRoute{h, auth}
}

func (ro *bordirRoute) Bordir(router fiber.Router) {
	router.Use(ro.auth.Authorization([]string{"DIREKTUR"}))
	router.Get("", ro.h.BordirHandler().GetAll)
	router.Get("/:id", ro.h.BordirHandler().GetById)
	router.Post("", ro.h.BordirHandler().Create)
	router.Put("/:id", ro.h.BordirHandler().Update)
	router.Delete("/:id", ro.h.BordirHandler().Delete)
}

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
	auth := ro.auth.Authorization([]string{"DIREKTUR"})
	auth2 := ro.auth.Authorization([]string{"DIREKTUR", "ADMIN"})
	router.Get("", auth2, ro.h.BordirHandler().GetAll)
	router.Get("/:id", auth, ro.h.BordirHandler().GetById)
	router.Post("", auth, ro.h.BordirHandler().Create)
	router.Put("/:id", auth, ro.h.BordirHandler().Update)
	router.Delete("/:id", auth, ro.h.BordirHandler().Delete)
}

package route

import (
	"github.com/gofiber/fiber/v2"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type SablonRoute interface {
	Sablon(router fiber.Router)
}

type sablonRoute struct {
	h    handler_init.SablonHandlerInit
	auth middleware_auth.AuthMidleware
}

func NewSablonRoute(h handler_init.SablonHandlerInit, auth middleware_auth.AuthMidleware) SablonRoute {
	return &sablonRoute{h, auth}
}

func (ro *sablonRoute) Sablon(router fiber.Router) {
	auth := ro.auth.Authorization([]string{"DIREKTUR"})
	auth2 := ro.auth.Authorization([]string{"DIREKTUR", "ADMIN"})
	router.Get("", auth2, ro.h.SablonHandler().GetAll)
	router.Post("", auth, ro.h.SablonHandler().Create)
	router.Put("/:id", auth, ro.h.SablonHandler().Update)
	router.Delete("/:id", auth, ro.h.SablonHandler().Delete)
}

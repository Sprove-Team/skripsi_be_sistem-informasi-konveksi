package route

import (
	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/gofiber/fiber/v2"
)

type ProfileRoute interface {
	Profile(router fiber.Router)
}

type profileRoute struct {
	h    handler_init.ProfileHandlerInit
	auth middleware_auth.AuthMidleware
}

func NewProfileRoute(h handler_init.ProfileHandlerInit, auth middleware_auth.AuthMidleware) ProfileRoute {
	return &profileRoute{h, auth}
}

func (ro *profileRoute) Profile(router fiber.Router) {
	router.Use(ro.auth.Authorization([]string{}))
	router.Get("", ro.h.ProfileHandler().GetProfile)
	router.Put("", ro.h.ProfileHandler().Update)
}

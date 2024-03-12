package route

import (
	"github.com/gofiber/fiber/v2"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type AuthRoute interface {
	Auth(router fiber.Router)
}

type authRoute struct {
	h    handler_init.AuthHandlerInit
	auth middleware_auth.AuthMidleware
}

func NewAuthRoute(h handler_init.AuthHandlerInit, auth middleware_auth.AuthMidleware) AuthRoute {
	return &authRoute{h, auth}
}

func (ro *authRoute) Auth(router fiber.Router) {
	router.Get("/whoami", ro.auth.Authorization([]string{}), ro.h.Auth().WhoAmI)
	router.Post("/login", ro.h.Auth().Login)
	router.Post("/refresh_token", ro.h.Auth().RefreshToken)
}

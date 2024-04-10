package route

import (
	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/gofiber/fiber/v2"
)

type UserRoute interface {
	JenisSpv(router fiber.Router)
	User(router fiber.Router)
}

type userRoute struct {
	h    handler_init.UserHandlerInit
	auth middleware_auth.AuthMidleware
}

func NewUserRoute(h handler_init.UserHandlerInit, auth middleware_auth.AuthMidleware) UserRoute {
	return &userRoute{h, auth}
}

func (ro *userRoute) User(router fiber.Router) {
	auth := ro.auth.Authorization([]string{entity.RolesById[1]})
	auth2 := ro.auth.Authorization([]string{entity.RolesById[1], entity.RolesById[4]})
	router.Get("", auth2, ro.h.UserHandler().GetAll)
	router.Get("/:id", auth, ro.h.UserHandler().GetById)
	router.Post("", auth, ro.h.UserHandler().Create)
	router.Put("/:id", auth, ro.h.UserHandler().Update)
	router.Delete("/:id", auth, ro.h.UserHandler().Delete)
}

func (ro *userRoute) JenisSpv(router fiber.Router) {
	auth := ro.auth.Authorization([]string{entity.RolesById[1]})
	auth2 := ro.auth.Authorization([]string{entity.RolesById[1], entity.RolesById[4]})
	router.Get("", auth2, ro.h.JenisSpvHandler().GetAll)
	router.Post("", auth, ro.h.JenisSpvHandler().Create)
	router.Put("/:id", auth, ro.h.JenisSpvHandler().Update)
	router.Delete("/:id", auth, ro.h.JenisSpvHandler().Delete)
}

package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/be-sistem-informasi-konveksi/app/handler_init"
)

type UserRoute interface {
	JenisSpv(router fiber.Router)
	User(router fiber.Router)
}

type userRoute struct {
	h handler_init.UserHandlerInit
}

func NewUserRoute(h handler_init.UserHandlerInit) UserRoute {
	return &userRoute{h}
}

func (ro *userRoute) User(router fiber.Router) {
	router.Get("", ro.h.UserHandler().GetAll)
	router.Post("", ro.h.UserHandler().Create)
	router.Put("/:id", ro.h.UserHandler().Update)
	router.Delete("/:id", ro.h.UserHandler().Delete)
}

func (ro *userRoute) JenisSpv(router fiber.Router) {
	router.Get("", ro.h.JenisSpvHandler().GetAll)
	router.Post("", ro.h.JenisSpvHandler().Create)
	router.Put("/:id", ro.h.JenisSpvHandler().Update)
	router.Delete("/:id", ro.h.JenisSpvHandler().Delete)
}

package route

import (
	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/gofiber/fiber/v2"
)

type AkuntansiRoute interface {
	Akun(router fiber.Router)
	KelompokAkun(router fiber.Router)
	Transaksi(router fiber.Router)
	HutangPiutang(router fiber.Router)
	Kontak(router fiber.Router)
	Akuntansi(router fiber.Router)
}

type akuntansiRoute struct {
	h    handler_init.AkuntansiHandlerInit
	auth middleware_auth.AuthMidleware
}

func NewAkuntansiRoute(h handler_init.AkuntansiHandlerInit, auth middleware_auth.AuthMidleware) AkuntansiRoute {
	return &akuntansiRoute{h, auth}
}

func (ro *akuntansiRoute) Akun(router fiber.Router) {
	// this is how to use the authorization function
	// ~ router.Get("", ro.auth.Authorization([]string{"direktur", "bendahara"}), ro.h.AkunHandler().GetAll)
	router.Get("", ro.h.Akun().GetAll)
	router.Get("/:id", ro.h.Akun().GetById)
	router.Post("", ro.h.Akun().Create)
	router.Put("/:id", ro.h.Akun().Update)
	router.Delete("/:id", ro.h.Akun().Delete)
}

func (ro *akuntansiRoute) KelompokAkun(router fiber.Router) {
	router.Post("", ro.h.KelompokAkun().Create)
	router.Get("", ro.h.KelompokAkun().GetAll)
	router.Get("/:id", ro.h.KelompokAkun().GetById)
	router.Put("/:id", ro.h.KelompokAkun().Update)
	router.Delete("/:id", ro.h.KelompokAkun().Delete)
}

func (ro *akuntansiRoute) Transaksi(router fiber.Router) {
	router.Get("/history", ro.h.Transaksi().GetHistory)
	router.Get("", ro.h.Transaksi().GetAll)
	router.Get("/:id", ro.h.Transaksi().GetById)
	router.Post("", ro.h.Transaksi().Create)
	router.Put("/:id", ro.h.Transaksi().Update)
	router.Delete("/:id", ro.h.Transaksi().Delete)
}

func (ro *akuntansiRoute) HutangPiutang(router fiber.Router) {
	router.Post("", ro.h.HutangPiutang().Create)
	router.Post("/bayar/:hutang_piutang_id", ro.h.HutangPiutang().CreateBayar)
	router.Get("", ro.h.HutangPiutang().GetAll)
}

func (ro *akuntansiRoute) Kontak(router fiber.Router) {
	router.Post("", ro.h.Kontak().Create)
	router.Put("/:id", ro.h.Kontak().Update)
	router.Get("", ro.h.Kontak().GetAll)
	router.Get("/:id", ro.h.Kontak().GetById)
	router.Delete("/:id", ro.h.Kontak().Delete)
}

func (ro *akuntansiRoute) Akuntansi(router fiber.Router) {
	router.Use(ro.auth.Authorization([]string{"DIREKTUR", "BENDAHARA"}))
	router.Get("/jurnal_umum", ro.h.Akuntansi().GetAllJU)
	router.Get("/buku_besar", ro.h.Akuntansi().GetAllBB)
	router.Get("/neraca_saldo", ro.h.Akuntansi().GetAllNC)
	router.Get("/laba_rugi", ro.h.Akuntansi().GetAllLBR)
}

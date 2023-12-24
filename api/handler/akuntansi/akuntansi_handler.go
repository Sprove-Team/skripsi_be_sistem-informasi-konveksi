package akuntansi

import (
	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi"
	resGlobal "github.com/be-sistem-informasi-konveksi/common/response/global"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type AkuntansiHandler interface {
	GetAllJU(c *fiber.Ctx) error
	GetAllBB(c *fiber.Ctx) error
	GetAllNC(c *fiber.Ctx) error
}

type akuntansiHandler struct {
	uc        usecase.AkuntansiUsecase
	validator pkg.Validator
}

func NewAkuntansiHandler(uc usecase.AkuntansiUsecase, validator pkg.Validator) AkuntansiHandler {
	return &akuntansiHandler{uc, validator}
}

func (h *akuntansiHandler) GetAllJU(c *fiber.Ctx) error {
	reqU := new(req.GetAllJU)
	c.QueryParser(reqU)

	errValidate := h.validator.Validate(reqU)
	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}

	c.Accepts("application/json")
	ctx := c.UserContext()

	dataJU, err := h.uc.GetAllJU(ctx, *reqU)
	if err != nil {
		helper.LogsError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}
	data := fiber.Map{
		"jurnal_umum": dataJU,
	}
	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithData(data, "R"))
}

func (h *akuntansiHandler) GetAllBB(c *fiber.Ctx) error {
	reqU := new(req.GetAllBB)
	c.QueryParser(reqU)

	errValidate := h.validator.Validate(reqU)
	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}

	c.Accepts("application/json")
	ctx := c.UserContext()

	dataJU, err := h.uc.GetAllBB(ctx, *reqU)
	if err != nil {
		helper.LogsError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}
	data := fiber.Map{
		"buku_besar": dataJU,
	}
	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithData(data, "R"))
}

func (h *akuntansiHandler) GetAllNC(c *fiber.Ctx) error {
	reqU := new(req.GetAllNC)
	c.QueryParser(reqU)

	errValidate := h.validator.Validate(reqU)
	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}

	c.Accepts("application/json")
	ctx := c.UserContext()

	dataJU, err := h.uc.GetAllNC(ctx, *reqU)
	if err != nil {
		helper.LogsError(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}
	data := fiber.Map{
		"neraca_saldo": dataJU,
	}
	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithData(data, "R"))
}

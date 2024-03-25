package handler_akuntansi

import (
	"context"
	"strings"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi"
	res_global "github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type AkuntansiHandler interface {
	GetAllJU(c *fiber.Ctx) error
	GetAllBB(c *fiber.Ctx) error
	GetAllNC(c *fiber.Ctx) error
	GetAllLBR(c *fiber.Ctx) error
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
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	datasJU, err := h.uc.GetAllJU(ctx, *reqU)

	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(res_global.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datasJU))
}

func (h *akuntansiHandler) GetAllBB(c *fiber.Ctx) error {
	reqU := new(req.GetAllBB)
	c.QueryParser(reqU)
	akunIds := c.Query("akun_id")
	if akunIds != "" {
		reqU.AkunID = strings.Split(strings.Trim(akunIds, " "), ",")
	}
	errValidate := h.validator.Validate(reqU)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	datasBB, err := h.uc.GetAllBB(ctx, *reqU)

	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(res_global.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err != nil {
		return c.Status(fiber.StatusRequestTimeout).JSON(res_global.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datasBB))
}

func (h *akuntansiHandler) GetAllNC(c *fiber.Ctx) error {
	reqU := new(req.GetAllNC)
	c.QueryParser(reqU)

	errValidate := h.validator.Validate(reqU)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	datasNC, err := h.uc.GetAllNC(ctx, *reqU)

	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(res_global.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datasNC))
}

func (h *akuntansiHandler) GetAllLBR(c *fiber.Ctx) error {
	reqU := new(req.GetAllLBR)
	c.QueryParser(reqU)

	errValidate := h.validator.Validate(reqU)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	datasJU, err := h.uc.GetAllLBR(ctx, *reqU)

	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(res_global.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datasJU))
}

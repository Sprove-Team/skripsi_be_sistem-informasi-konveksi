package produk

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	resGlobal "github.com/be-sistem-informasi-konveksi/common/response/global"
	helper "github.com/be-sistem-informasi-konveksi/helper"
	usecase "github.com/be-sistem-informasi-konveksi/usecase/produk"
)

type ProdukHandler interface {
	Create(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
}

type produkHandler struct {
	uc        usecase.ProdukUsecase
	validator helper.Validator
}

func NewProdukHandler(uc usecase.ProdukUsecase, validator helper.Validator) ProdukHandler {
	return &produkHandler{uc, validator}
}

func (h *produkHandler) Create(c *fiber.Ctx) error {
	c.Accepts("application/json")
	req := new(req.CreateProduk)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}
	errValidate := h.validator.Validate(req)

	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}
	err := h.uc.Create(c.UserContext(), *req)
	if err != nil {
		if err.Error() == "duplicated key not allowed" {
			return c.Status(fiber.StatusConflict).JSON(resGlobal.ErrorResWithoutData(fiber.StatusConflict))
		}
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusBadRequest).JSON(resGlobal.CustomRes(fiber.StatusBadRequest, message.KategoriNotFound, nil))
		}

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusCreated).JSON(resGlobal.SuccessResWithoutData("C"))
}

func (h *produkHandler) GetById(c *fiber.Ctx) error {
	id := c.Params("id", "")
	ctx := c.UserContext()
	data, err := h.uc.GetById(ctx, id)

	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(resGlobal.ErrorResWithoutData(fiber.StatusRequestTimeout))
	}

	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(resGlobal.ErrorResWithoutData(fiber.StatusNotFound))
		}
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithData(data, "R"))
}

package produk

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	"github.com/be-sistem-informasi-konveksi/common/response"
	resGlobal "github.com/be-sistem-informasi-konveksi/common/response/global"
	helper "github.com/be-sistem-informasi-konveksi/helper"
	usecase "github.com/be-sistem-informasi-konveksi/usecase/produk"
)

type ProdukHandler interface {
	Delete(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
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
	ctx := c.UserContext()
	err := h.uc.Create(ctx, *req)
	

	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(resGlobal.ErrorResWithoutData(fiber.StatusRequestTimeout))
	}

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

func (h *produkHandler) Update(c *fiber.Ctx) error {
	req := new(req.UpdateProduk)
	if err := c.ParamsParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}
	c.Accepts("application/json")
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}

	errValidate := h.validator.Validate(req)

	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}
	ctx := c.UserContext()
	err := h.uc.Update(ctx, *req)
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(resGlobal.ErrorResWithoutData(fiber.StatusRequestTimeout))
	}
	if err != nil {
		if err.Error() == message.KategoriNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData([]response.BaseFormatError{
				{
					ValueInput:   req.KategoriID,
					ErrorMessage: message.KategoriNotFound,
				},
			}, fiber.StatusBadRequest))
		}

		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(resGlobal.ErrorResWithoutData(fiber.StatusNotFound))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithoutData("U"))
}

func (h *produkHandler) Delete(c *fiber.Ctx) error {
	req := new(req.ParamByID)
	if err := c.ParamsParser(req); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}
	errValidate := h.validator.Validate(req)
	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}
	ctx := c.UserContext()
	err := h.uc.Delete(ctx, req.ID)
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

	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithoutData("D"))
}

func (h *produkHandler) GetById(c *fiber.Ctx) error {
	req := new(req.ParamByID)
	if err := c.ParamsParser(req); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}
	errValidate := h.validator.Validate(req)
	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}
	ctx := c.UserContext()
	data, err := h.uc.GetById(ctx, req.ID)
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

func (h *produkHandler) GetAll(c *fiber.Ctx) error {
	reqU := new(req.GetAll)
	c.BodyParser(reqU)
	c.QueryParser(reqU)

	errValidate := h.validator.Validate(reqU)
	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}

	ctx := c.UserContext()
	data, currentPage, totalPage, err := h.uc.GetAll(ctx, *reqU)
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
	dataRes := fiber.Map{
		"produks":      data,
		"current_page": currentPage,
		"total_page":   totalPage,
	}
	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithData(dataRes, "R"))
}

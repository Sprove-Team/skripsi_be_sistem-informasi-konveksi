package handler_produk

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/produk"
	"github.com/be-sistem-informasi-konveksi/common/message"
	reqGlobal "github.com/be-sistem-informasi-konveksi/common/request/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	res_global "github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type ProdukHandler interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type produkHandler struct {
	uc        usecase.ProdukUsecase
	validator pkg.Validator
}

func NewProdukHandler(uc usecase.ProdukUsecase, validator pkg.Validator) ProdukHandler {
	return &produkHandler{uc, validator}
}

func errResponse(c *fiber.Ctx, err error) error {
	if err == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(res_global.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err.Error() == "record not found" {
		return c.Status(fiber.StatusNotFound).JSON(res_global.ErrorRes(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, nil))
	}

	if err.Error() == "duplicated key not allowed" {
		return c.Status(fiber.StatusConflict).JSON(res_global.ErrorRes(fiber.ErrConflict.Code, fiber.ErrConflict.Message, nil))
	}

	badRequest := make([]string, 0, 1)

	switch err.Error() {
	case message.KategoriProdukNotFound:
		badRequest = append(badRequest, err.Error())
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(res_global.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func (h *produkHandler) Create(c *fiber.Ctx) error {
	req := new(req.Create)
	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	err := h.uc.Create(ctx, *req)

	if err != nil {
		return errResponse(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(res_global.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *produkHandler) Update(c *fiber.Ctx) error {
	req := new(req.Update)

	c.ParamsParser(req)
	c.BodyParser(req)

	errValidate := h.validator.Validate(req)

	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	err := h.uc.Update(ctx, *req)

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *produkHandler) Delete(c *fiber.Ctx) error {
	req := new(reqGlobal.ParamByID)
	c.ParamsParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	err := h.uc.Delete(ctx, req.ID)

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *produkHandler) GetById(c *fiber.Ctx) error {
	req := new(reqGlobal.ParamByID)
	c.ParamsParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	data, err := h.uc.GetById(ctx, req.ID)

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, data))
}

func (h *produkHandler) GetAll(c *fiber.Ctx) error {
	reqU := new(req.GetAll)
	c.BodyParser(reqU)
	c.QueryParser(reqU)

	reqU.HargaDetail = strings.ToUpper(reqU.HargaDetail)

	errValidate := h.validator.Validate(reqU)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	datas, err := h.uc.GetAll(ctx, *reqU)

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datas))
}

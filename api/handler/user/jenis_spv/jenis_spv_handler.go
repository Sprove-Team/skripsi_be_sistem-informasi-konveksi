package user

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/user/jenis_spv"
	reqGlobal "github.com/be-sistem-informasi-konveksi/common/request/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/user/jenis_spv"
	resGlobal "github.com/be-sistem-informasi-konveksi/common/response/global"
	"github.com/be-sistem-informasi-konveksi/helper"
)

type JenisSpvHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type jenisSpvHandler struct {
	uc        usecase.JenisSpvUsecase
	validator helper.Validator
}

func NewJenisSpvHandler(uc usecase.JenisSpvUsecase, validator helper.Validator) JenisSpvHandler {
	return &jenisSpvHandler{uc, validator}
}

func (h *jenisSpvHandler) Create(c *fiber.Ctx) error {
	c.Accepts("application/json")
	req := new(req.Create)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}
	req.Nama = strings.ToLower(req.Nama)
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
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusCreated).JSON(resGlobal.SuccessResWithoutData("C"))
}

func (h *jenisSpvHandler) Update(c *fiber.Ctx) error {
	req := new(req.Update)
	if err := c.ParamsParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}
	c.Accepts("application/json")
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}
	req.Nama = strings.ToLower(req.Nama)
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
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(resGlobal.ErrorResWithoutData(fiber.StatusNotFound))
		}
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithoutData("U"))
}

func (h *jenisSpvHandler) Delete(c *fiber.Ctx) error {
	req := new(reqGlobal.ParamByID)
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

func (h *jenisSpvHandler) GetAll(c *fiber.Ctx) error {
	ctx := c.UserContext()
	datas, err := h.uc.GetAll(ctx)
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(resGlobal.ErrorResWithoutData(fiber.StatusRequestTimeout))
	}
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}
  dataRes := fiber.Map{
    "jenis_spv": datas,
  }
	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithData(dataRes, "R"))
}

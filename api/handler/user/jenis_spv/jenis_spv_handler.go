package user

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/user/jenis_spv"
	"github.com/be-sistem-informasi-konveksi/common/message"
	reqGlobal "github.com/be-sistem-informasi-konveksi/common/request/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/user/jenis_spv"
	"github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type JenisSpvHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type jenisSpvHandler struct {
	uc        usecase.JenisSpvUsecase
	validator pkg.Validator
}

func NewJenisSpvHandler(uc usecase.JenisSpvUsecase, validator pkg.Validator) JenisSpvHandler {
	return &jenisSpvHandler{uc, validator}
}

func (h *jenisSpvHandler) Create(c *fiber.Ctx) error {
	c.Accepts("application/json")
	req := new(req.Create)
	c.BodyParser(req)
	// if err := ; err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	// }
	req.Nama = strings.ToLower(req.Nama)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	err := h.uc.Create(ctx, *req)

	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err != nil {
		if err.Error() == "duplicated key not allowed" {
			return c.Status(fiber.StatusConflict).JSON(response.ErrorRes(fiber.ErrConflict.Code, fiber.ErrConflict.Message, nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *jenisSpvHandler) Update(c *fiber.Ctx) error {
	req := new(req.Update)
	c.ParamsParser(req)
	// if err := ; err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	// }
	c.BodyParser(req)
	// if err := c.BodyParser(req); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	// }
	req.Nama = strings.ToLower(req.Nama)
	errValidate := h.validator.Validate(req)

	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	err := h.uc.Update(ctx, *req)
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorRes(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *jenisSpvHandler) Delete(c *fiber.Ctx) error {
	req := new(reqGlobal.ParamByID)
	c.ParamsParser(req)
	// if err := ; err != nil {
	// 	log.Println(err)
	// 	return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	// }
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	err := h.uc.Delete(ctx, req.ID)
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorRes(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, nil))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *jenisSpvHandler) GetAll(c *fiber.Ctx) error {
	ctx := c.UserContext()
	datas, err := h.uc.GetAll(ctx)
	if ctx.Err() == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, datas))
}

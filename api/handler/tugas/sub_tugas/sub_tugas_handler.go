package handler_sub_tugas

import (
	"context"

	uc_sub_tugas "github.com/be-sistem-informasi-konveksi/api/usecase/tugas/sub_tugas"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_global "github.com/be-sistem-informasi-konveksi/common/request/global"
	req_sub_tugas "github.com/be-sistem-informasi-konveksi/common/request/tugas/sub_tugas"
	res_global "github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type SubTugasHandler interface {
	CreateByTugasId(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type subTugasHandler struct {
	uc        uc_sub_tugas.SubTugasUsecase
	validator pkg.Validator
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
	case message.TugasNotFound:
		badRequest = append(badRequest, err.Error())
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(res_global.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}

	if err.Error() == message.AkunCannotDeleted {
		return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorInterWithMessageRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, err.Error()))
	}

	return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func NewSubTugasHandler(uc uc_sub_tugas.SubTugasUsecase, validator pkg.Validator) SubTugasHandler {
	return &subTugasHandler{uc, validator}
}

func (h *subTugasHandler) CreateByTugasId(c *fiber.Ctx) error {
	req := new(req_sub_tugas.CreateByTugasId)
	c.ParamsParser(req)
	c.BodyParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	err := h.uc.CreateByTugasId(uc_sub_tugas.ParamCreateByTugasId{
		Ctx: ctx,
		Req: *req,
	})
	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(res_global.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *subTugasHandler) Update(c *fiber.Ctx) error {
	req := new(req_sub_tugas.Update)
	c.ParamsParser(req)
	c.BodyParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	claims, ok := c.Locals("user").(*pkg.Claims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	ctx := c.UserContext()
	err := h.uc.Update(uc_sub_tugas.ParamUpdate{
		Ctx:    ctx,
		Claims: claims,
		Req:    *req,
	})

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *subTugasHandler) Delete(c *fiber.Ctx) error {
	req := new(req_global.ParamByID)
	c.ParamsParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	err := h.uc.Delete(uc_sub_tugas.ParamDelete{
		Ctx: ctx,
		Req: *req,
	})

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}

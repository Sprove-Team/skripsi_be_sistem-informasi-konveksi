package handler_tugas

import (
	"context"

	uc_tugas "github.com/be-sistem-informasi-konveksi/api/usecase/tugas"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_tugas "github.com/be-sistem-informasi-konveksi/common/request/tugas"
	res_global "github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type TugasHandler interface {
	Create(c *fiber.Ctx) error
}

type tugasHandler struct {
	uc        uc_tugas.TugasUsecase
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

	// badRequest := map[string][]string{}

	// if err.Error() == message.KelompokAkunIdNotFound {
	// 	badRequest[""] = []string{err.Error()}
	// }

	// if len(badRequest) > 0 {
	// 	return c.Status(fiber.StatusBadRequest).JSON(response.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	// }

	if err.Error() == message.AkunCannotDeleted {
		return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorInterWithMessageRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, err.Error()))
	}

	return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func NewTugasHandler(uc uc_tugas.TugasUsecase, validator pkg.Validator) TugasHandler {
	return &tugasHandler{}
}

func (h *tugasHandler) Create(c *fiber.Ctx) error {
	req := new(req_tugas.Create)
	c.BodyParser(req)
	errValidate := h.validator.Validate(req)

	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	err := h.uc.Create(uc_tugas.ParamCreate{
		Ctx: ctx,
		Req: *req,
	})

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(res_global.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

package handler_akuntansi_akun

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/akun"
	req_global "github.com/be-sistem-informasi-konveksi/common/request/global"
	res_global "github.com/be-sistem-informasi-konveksi/common/response"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/akun"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type AkunHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error

	// Add other methods as needed
}

type akunHandler struct {
	uc        usecase.AkunUsecase
	validator pkg.Validator
}

func NewAkunHandler(uc usecase.AkunUsecase, validator pkg.Validator) AkunHandler {
	return &akunHandler{uc, validator}
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
	case message.CantModifiedDefaultData, message.KelompokAkunNotFound:
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

func (h *akunHandler) Create(c *fiber.Ctx) error {
	// Parse request body
	req := new(req.Create)
	c.BodyParser(req)
	// Validate request
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	// Create context
	ctx := c.UserContext()

	// Call usecase to create Akun
	err := h.uc.Create(ctx, *req)
	// Handle errors
	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(res_global.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *akunHandler) Update(c *fiber.Ctx) error {
	// Parse request body
	req := new(req.Update)
	c.ParamsParser(req)
	c.BodyParser(req)

	// Validate request
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	// Create context
	ctx := c.UserContext()

	// Call usecase to create Akun
	err := h.uc.Update(ctx, *req)
	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *akunHandler) Delete(c *fiber.Ctx) error {
	// Parse request body
	req := new(req_global.ParamByID)
	c.ParamsParser(req)

	// Validate request
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	// Create context
	ctx := c.UserContext()

	// Call usecase to create Akun
	err := h.uc.Delete(ctx, req.ID)
	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *akunHandler) GetAll(c *fiber.Ctx) error {
	req := new(req.GetAll)
	c.QueryParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()

	data, err := h.uc.GetAll(ctx, *req)
	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, data))
}

func (h *akunHandler) GetById(c *fiber.Ctx) error {
	// Parse request body
	req := new(req_global.ParamByID)
	c.ParamsParser(req)

	// Validate request
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	// Create context
	ctx := c.UserContext()

	// Call usecase to create Akun
	data, err := h.uc.GetById(ctx, req.ID)
	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, data))
}

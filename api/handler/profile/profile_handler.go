package handler_profile

import (
	"context"

	uc_profile "github.com/be-sistem-informasi-konveksi/api/usecase/profile"
	uc_user "github.com/be-sistem-informasi-konveksi/api/usecase/user"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_profile "github.com/be-sistem-informasi-konveksi/common/request/profile"
	res_global "github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type ProfileHandler interface {
	GetProfile(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type profileHandler struct {
	uc        uc_profile.ProfileUsecase
	ucUser    uc_user.UserUsecase
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
	case message.UsernameConflict,
		message.NotFitOldPassword:
		badRequest = append(badRequest, err.Error())
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(res_global.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}

	return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func NewProfileHandler(uc uc_profile.ProfileUsecase, ucUser uc_user.UserUsecase, validator pkg.Validator) ProfileHandler {
	return &profileHandler{uc, ucUser, validator}
}

func (h *profileHandler) GetProfile(c *fiber.Ctx) error {

	ctx := c.UserContext()
	claims, ok := c.Locals("user").(*pkg.Claims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	data, err := h.ucUser.GetById(uc_user.ParamGetById{
		Ctx: ctx,
		ID:  claims.ID,
	})

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, data))
}

func (h *profileHandler) Update(c *fiber.Ctx) error {
	req := new(req_profile.Update)

	c.BodyParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	claims, ok := c.Locals("user").(*pkg.Claims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
	}

	err := h.uc.Update(uc_profile.ParamUpdate{
		Ctx:    ctx,
		Claims: claims,
		Req:    *req,
	})

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}

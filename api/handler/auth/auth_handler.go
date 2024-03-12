package handler_auth

import (
	"context"

	uc_auth "github.com/be-sistem-informasi-konveksi/api/usecase/auth"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req_auth "github.com/be-sistem-informasi-konveksi/common/request/auth"
	res_global "github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	WhoAmI(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
}

type authHandler struct {
	uc        uc_auth.AuthUsecase
	validator pkg.Validator
}

func NewAuthHandler(uc uc_auth.AuthUsecase, validator pkg.Validator) AuthHandler {
	return &authHandler{uc, validator}
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
	case
		message.UserNotFoundOrDeleted,
		message.RefreshTokenExpired,
		message.InvalidRefreshToken,
		message.InvalidUsernameOrPassword:
		badRequest = append(badRequest, err.Error())
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(res_global.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func (h *authHandler) WhoAmI(c *fiber.Ctx) error {
	claims := c.Locals("user").(*pkg.Claims)
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, claims))
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	req := new(req_auth.Login)

	c.BodyParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	token, refToken, err := h.uc.Login(uc_auth.ParamLogin{
		Ctx: ctx,
		Req: *req,
	})
	// Handle errors
	if err != nil {
		return errResponse(c, err)
	}
	data := fiber.Map{
		"token":         token,
		"refresh_token": refToken,
	}
	// Respond with success status
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, data))
}

func (h *authHandler) RefreshToken(c *fiber.Ctx) error {
	req := new(req_auth.GetNewToken)

	c.BodyParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	newToken, err := h.uc.RefreshToken(uc_auth.ParamRefreshToken{
		Ctx: ctx,
		Req: *req,
	})

	// Handle errors
	if err != nil {
		return errResponse(c, err)
	}

	data := fiber.Map{
		"token": newToken,
	}
	// Respond with success status
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, data))
}

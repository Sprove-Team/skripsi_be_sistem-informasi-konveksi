package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/user"
	"github.com/be-sistem-informasi-konveksi/common/message"
	reqGlobal "github.com/be-sistem-informasi-konveksi/common/request/global"
	req "github.com/be-sistem-informasi-konveksi/common/request/user"
	"github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type UserHandler interface {
	// GetById(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type userHandler struct {
	uc        usecase.UserUsecase
	validator pkg.Validator
}

func NewUserHandler(uc usecase.UserUsecase, validator pkg.Validator) UserHandler {
	return &userHandler{uc, validator}
}

func errResponse(c *fiber.Ctx, err error) error {
	if err == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err.Error() == "record not found" {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorRes(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, nil))
	}

	if err.Error() == "duplicated key not allowed" {
		return c.Status(fiber.StatusConflict).JSON(response.ErrorRes(fiber.ErrConflict.Code, fiber.ErrConflict.Message, nil))
	}

	return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func (h *userHandler) Create(c *fiber.Ctx) error {
	req := new(req.Create)
	c.BodyParser(req)

	req.Role = strings.ToUpper(req.Role)
	errValidate := h.validator.Validate(req)

	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	dataUser, err := h.uc.CreateUserData(usecase.ParamCreateUserData{
		Ctx: ctx,
		Req: *req,
	})

	if err != nil {
		return errResponse(c, err)
	}

	err = h.uc.CreateCommitDB(usecase.ParamCreateCommitDB{
		Ctx:  ctx,
		User: dataUser,
	})

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *userHandler) GetAll(c *fiber.Ctx) error {
	req := new(req.GetAll)
	c.QueryParser(req)

	fmt.Println(req.Search)

	req.Search.Role = strings.ToUpper(req.Search.Role)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	data, err := h.uc.GetAll(usecase.ParamGetAll{
		Ctx: ctx,
		Req: *req,
	})
	if err != nil {
		return errResponse(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, data))
}

func (h *userHandler) Update(c *fiber.Ctx) error {
	req := new(req.Update)
	c.ParamsParser(req)

	c.BodyParser(req)
	req.Role = strings.ToUpper(req.Role)
	errValidate := h.validator.Validate(req)

	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	dataUser, err := h.uc.UpdateUserData(usecase.ParamUpdate{
		Ctx: ctx,
		Req: *req,
	})
	if err != nil {
		return errResponse(c, err)
	}

	err = h.uc.UpdateCommitDB(usecase.ParamUpdateCommitDB{
		Ctx:  ctx,
		User: dataUser,
	})

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *userHandler) Delete(c *fiber.Ctx) error {
	req := new(reqGlobal.ParamByID)
	c.ParamsParser(req)
	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, nil))
	}
	ctx := c.UserContext()
	err := h.uc.Delete(usecase.ParamDelete{
		Ctx: ctx,
		ID:  req.ID,
	})

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, nil))
}

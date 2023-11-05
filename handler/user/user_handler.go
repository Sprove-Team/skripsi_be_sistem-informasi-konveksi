package user

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	req "github.com/be-sistem-informasi-konveksi/common/request/user"
	resGlobal "github.com/be-sistem-informasi-konveksi/common/response/global"
	"github.com/be-sistem-informasi-konveksi/helper"
	usecase "github.com/be-sistem-informasi-konveksi/usecase/user"
)

type UserHandler interface {
	Create(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
}

type userHandler struct {
	uc        usecase.UserUsecase
	validator helper.Validator
}

func NewUserHandler(uc usecase.UserUsecase, validator helper.Validator) UserHandler {
	return &userHandler{uc, validator}
}

func (h *userHandler) Create(c *fiber.Ctx) error {
	req := new(req.CreateUser)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithoutData(fiber.StatusBadRequest))
	}
	req.Role = strings.ToUpper(req.Role)
	errValidate := h.validator.Validate(req)

	if len(errValidate) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
	}

	ctx := c.UserContext()
	err := h.uc.Create(ctx, *req)
	if err != nil {
		if err.Error() == "duplicated key not allowed" {
			return c.Status(fiber.StatusConflict).JSON(resGlobal.ErrorResWithoutData(fiber.StatusConflict))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(resGlobal.ErrorResWithoutData(fiber.StatusInternalServerError))
	}

	return c.Status(fiber.StatusCreated).JSON(resGlobal.SuccessResWithoutData("C"))
}

func (h *userHandler) GetAll(c *fiber.Ctx) error {
  req := new(req.GetAllUser)
  c.BodyParser(req)
	c.QueryParser(req)

  req.Search.Role = strings.ToUpper(req.Search.Role)
	errValidate := h.validator.Validate(req)
  if len(errValidate) > 0 {
    return c.Status(fiber.StatusBadRequest).JSON(resGlobal.ErrorResWithData(errValidate, fiber.StatusBadRequest))
  }
  ctx := c.UserContext()
	data, currentPage, totalPage, err := h.uc.GetAll(ctx, *req)

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
		"user":       data,
		"current_page": currentPage,
		"total_page":   totalPage,
	}
	return c.Status(fiber.StatusOK).JSON(resGlobal.SuccessResWithData(dataRes, "R"))
}

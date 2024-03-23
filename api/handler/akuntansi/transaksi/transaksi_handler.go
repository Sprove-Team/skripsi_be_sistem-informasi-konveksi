package handler_akuntansi_transaksi

import (
	"context"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/transaksi"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/transaksi"
	reqGlobal "github.com/be-sistem-informasi-konveksi/common/request/global"
	res_global "github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type TransaksiHandler interface {
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetById(c *fiber.Ctx) error
	GetHistory(c *fiber.Ctx) error
}

type transaksiHandler struct {
	uc        usecase.TransaksiUsecase
	validator pkg.Validator
}

func NewTransaksiHandler(uc usecase.TransaksiUsecase, validator pkg.Validator) TransaksiHandler {
	return &transaksiHandler{uc, validator}
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
	case message.AkunCannotBeSame,
		message.CreditDebitNotSame,
		message.AkunNotFound,
		message.AkunHutangPiutangNotEq2,
		message.BayarMustLessThanSisaTagihan,
		message.TotalHPMustGeOrEqToTotalByr,
		message.InvalidAkunHutangPiutang,
		message.AkunNotMatchWithJenisHPTr:
		badRequest = append(badRequest, err.Error())
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(res_global.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func (h *transaksiHandler) Create(c *fiber.Ctx) error {
	req := new(req.Create)

	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()
	// Call usecase to create KelompokAkun
	err := h.uc.Create(ctx, *req)
	// Handle errors
	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(res_global.SuccessRes(fiber.StatusCreated, message.Created, nil))
}

func (h *transaksiHandler) Update(c *fiber.Ctx) error {
	req := new(req.Update)

	c.ParamsParser(req)
	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	// Call usecase to create KelompokAkun
	err := h.uc.Update(ctx, *req)
	// Handle errors
	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *transaksiHandler) Delete(c *fiber.Ctx) error {
	reqU := new(reqGlobal.ParamByID)
	c.ParamsParser(reqU)

	errValidate := h.validator.Validate(reqU)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	err := h.uc.Delete(ctx, reqU.ID)
	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, nil))
}

func (h *transaksiHandler) GetById(c *fiber.Ctx) error {
	reqU := new(reqGlobal.ParamByID)
	c.ParamsParser(reqU)

	errValidate := h.validator.Validate(reqU)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	data, err := h.uc.GetById(ctx, reqU.ID)
	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, data))
}

func (h *transaksiHandler) GetAll(c *fiber.Ctx) error {
	reqU := new(req.GetAll)
	c.QueryParser(reqU)

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

func (h *transaksiHandler) GetHistory(c *fiber.Ctx) error {
	reqU := new(req.GetHistory)
	c.QueryParser(reqU)

	errValidate := h.validator.Validate(reqU)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	datas, err := h.uc.GetHistory(ctx, *reqU)

	if err != nil {
		return errResponse(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datas))
}

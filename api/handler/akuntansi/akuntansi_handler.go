package handler_akuntansi

import (
	"context"
	"fmt"
	"strings"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi"
	res_global "github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type AkuntansiHandler interface {
	GetAllJU(c *fiber.Ctx) error
	GetAllBB(c *fiber.Ctx) error
	GetAllNC(c *fiber.Ctx) error
	GetAllLBR(c *fiber.Ctx) error
}

type akuntansiHandler struct {
	uc        usecase.AkuntansiUsecase
	validator pkg.Validator
}

func NewAkuntansiHandler(uc usecase.AkuntansiUsecase, validator pkg.Validator) AkuntansiHandler {
	return &akuntansiHandler{uc, validator}
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
	case message.Timezoneunknown:
		badRequest = append(badRequest, err.Error())
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(res_global.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(res_global.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func (h *akuntansiHandler) GetAllJU(c *fiber.Ctx) error {
	reqU := new(req.GetAllJU)
	c.QueryParser(reqU)
	errValidate := h.validator.Validate(reqU)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	datasJU, err := h.uc.GetAllJU(ctx, *reqU)

	if err != nil {
		return errResponse(c, err)
	}

	if reqU.Download == "1" {
		buf, err := h.uc.DownloadJU(*reqU, datasJU)
		if err != nil {
			return errResponse(c, err)
		}
		name := fmt.Sprintf("Jurnal Umum (%s sampai %s).xlsx", reqU.StartDate, reqU.EndDate)
		c.Set("Content-Disposition", "attachment; filename="+name)
		return c.Status(fiber.StatusOK).Send(buf.Bytes())
	}
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datasJU))
}

func (h *akuntansiHandler) GetAllBB(c *fiber.Ctx) error {
	reqU := new(req.GetAllBB)
	c.QueryParser(reqU)
	reqU.AkunID = nil
	akunIds := c.Query("akun_id")
	if akunIds != "" {
		reqU.AkunID = strings.Split(strings.Trim(akunIds, " "), ",")
	}
	errValidate := h.validator.Validate(reqU)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	datasBB, err := h.uc.GetAllBB(ctx, *reqU)

	if err != nil {
		return errResponse(c, err)
	}

	if reqU.Download == "1" {
		buf, err := h.uc.DownloadBB(*reqU, datasBB)
		if err != nil {
			return errResponse(c, err)
		}
		name := fmt.Sprintf("Buku Besar (%s sampai %s).xlsx", reqU.StartDate, reqU.EndDate)
		c.Set("Content-Disposition", "attachment; filename="+name)
		return c.Status(fiber.StatusOK).Send(buf.Bytes())
	}
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datasBB))
}

func (h *akuntansiHandler) GetAllNC(c *fiber.Ctx) error {
	reqU := new(req.GetAllNC)
	c.QueryParser(reqU)

	errValidate := h.validator.Validate(reqU)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	datasNC, err := h.uc.GetAllNC(ctx, *reqU)

	if err != nil {
		return errResponse(c, err)
	}
	if reqU.Download == "1" {
		buf, err := h.uc.DownloadNC(*reqU, datasNC)
		if err != nil {
			return errResponse(c, err)
		}
		name := fmt.Sprintf("Neraca Saldo (%s).xlsx", reqU.Date)
		c.Set("Content-Disposition", "attachment; filename="+name)
		return c.Status(fiber.StatusOK).Send(buf.Bytes())
	}
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datasNC))
}

func (h *akuntansiHandler) GetAllLBR(c *fiber.Ctx) error {
	reqU := new(req.GetAllLBR)
	c.QueryParser(reqU)

	errValidate := h.validator.Validate(reqU)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}

	ctx := c.UserContext()

	datasLbr, err := h.uc.GetAllLBR(ctx, *reqU)

	if err != nil {
		return errResponse(c, err)
	}
	if reqU.Download == "1" {
		buf, err := h.uc.DownloadLBR(*reqU, datasLbr)
		if err != nil {
			return errResponse(c, err)
		}
		name := fmt.Sprintf("Laba Rugi (%s sampai %s).xlsx", reqU.StartDate, reqU.EndDate)
		c.Set("Content-Disposition", "attachment; filename="+name)
		return c.Status(fiber.StatusOK).Send(buf.Bytes())
	}
	return c.Status(fiber.StatusOK).JSON(res_global.SuccessRes(fiber.StatusOK, message.OK, datasLbr))
}

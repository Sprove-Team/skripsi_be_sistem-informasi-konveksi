package akuntansi

import (
	"context"
	"fmt"

	usecase "github.com/be-sistem-informasi-konveksi/api/usecase/akuntansi/hutang_piutang"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	"github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
)

type HutangPiutangHandler interface {
	Create(c *fiber.Ctx) error
	// Delete(c *fiber.Ctx) error
	// Update(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	// GetById(c *fiber.Ctx) error
	// GetHistory(c *fiber.Ctx) error
}

type hutangPiutangHandler struct {
	uc        usecase.HutangPiutangUsecase
	validator pkg.Validator
}

func NewHutangPiutangHandler(uc usecase.HutangPiutangUsecase, validator pkg.Validator) HutangPiutangHandler {
	return &hutangPiutangHandler{uc, validator}
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

	badRequest := map[string][]string{}

	if err.Error() == message.AkunCannotBeSame {
		badRequest["ayat_jurnal"] = []string{err.Error()}
	}

	if err.Error() == message.CreditDebitNotSame {
		badRequest["debit dan kredit"] = []string{err.Error()}
	}

	if err.Error() == message.AkunNotFound {
		badRequest["transaksi.ayat_jurnal.akun_id"] = []string{err.Error()}
	}

	if err.Error() == message.KontakNotFound {
		badRequest["kontak_id"] = []string{err.Error()}
	}

	if err.Error() == message.InvalidAkunHutangPiutang {
		badRequest["transaksi.ayat_jurnal.akun_id"] = []string{err.Error()}
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func (h *hutangPiutangHandler) Create(c *fiber.Ctx) error {
	req := new(req.Create)

	c.BodyParser(req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	err := h.uc.Create(ctx, *req)
	// Handle errors
	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusCreated).JSON(response.SuccessRes(fiber.StatusCreated, message.Created, nil))
}
func (h *hutangPiutangHandler) GetAll(c *fiber.Ctx) error {
	req := new(req.GetAll)

	c.QueryParser(req)
	fmt.Println(*req)

	errValidate := h.validator.Validate(req)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errValidate)
	}
	ctx := c.UserContext()
	data, err := h.uc.GetAll(ctx, *req)
	// Handle errors
	if err != nil {
		return errResponse(c, err)
	}

	// Respond with success status
	return c.Status(fiber.StatusOK).JSON(response.SuccessRes(fiber.StatusOK, message.OK, data))
}

package akuntansi

import (
	"context"
	"fmt"

	"github.com/be-sistem-informasi-konveksi/common/message"
	"github.com/be-sistem-informasi-konveksi/common/response"
	"github.com/gofiber/fiber/v2"
)

type errResponse struct{}

func (e *errResponse) errBase(c *fiber.Ctx, err error) error {
	fmt.Println(err)
	if err == context.DeadlineExceeded {
		return c.Status(fiber.StatusRequestTimeout).JSON(response.ErrorRes(fiber.ErrRequestTimeout.Code, fiber.ErrRequestTimeout.Message, nil))
	}

	if err.Error() == "record not found" {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorRes(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, nil))
	}

	if err.Error() == "duplicated key not allowed" {
		return c.Status(fiber.StatusConflict).JSON(response.ErrorRes(fiber.ErrConflict.Code, fiber.ErrConflict.Message, nil))
	}
	return nil
}

func (e *errResponse) errHP(c *fiber.Ctx, err error) error {
	if err := e.errBase(c, err); err != nil {
		return err
	}

	badRequest := make([]string, 0, 1)

	switch err.Error() {
	case message.AkunCannotBeSame,
		message.CreditDebitNotSame,
		message.AkunNotFound,
		message.KontakNotFound,
		message.InvalidAkunHutangPiutang,
		message.AkunNotMatchWithJenisHP,
		message.InvalidAkunBayar,
		message.BayarMustLessThanSisaTagihan,
		message.IncorrectPlacementOfCreditAndDebit:
		badRequest = append(badRequest, err.Error())
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

func (e *errResponse) errBayar(c *fiber.Ctx, err error) error {
	if err := e.errBase(c, err); err != nil {
		return err
	}

	badRequest := make([]string, 0, 1)

	switch err.Error() {
	case message.AkunNotFound,
		message.InvalidAkunBayar,
		message.KontakNotFound,
		message.HutangPiutangNotFound,
		message.BayarMustLessThanSisaTagihan:
		badRequest = append(badRequest, err.Error())
	}
	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

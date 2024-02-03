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

	badRequest := map[string][]string{}

	if err.Error() == message.AkunCannotBeSame {
		badRequest["ayat_jurnal"] = []string{err.Error()}
	}

	if err.Error() == message.CreditDebitNotSame {
		badRequest["transaksi.ayat_jurnal"] = []string{err.Error()}
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

	if err.Error() == message.AkunNotMatchWithJenisHP {
		badRequest["transaksi.ayat_jurnal.akun_id"] = []string{err.Error()}
	}

	if err.Error() == message.InvalidAkunBayar {
		badRequest["transaksi.ayat_jurnal.akun_id"] = []string{err.Error()}
	}

	if err.Error() == message.BayarMustLessThanSisaTagihan {
		badRequest["transaksi.ayat_jurnal.akun_id"] = []string{err.Error()}
	}

	if err.Error() == message.IncorrectPlacementOfCreditAndDebit {
		badRequest["transaksi.ayat_jurnal.akun_id"] = []string{err.Error()}
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

	badRequest := map[string][]string{}

	if err.Error() == message.AkunNotFound {
		badRequest["akun_bayar_id"] = []string{err.Error()}
	}

	if err.Error() == message.InvalidAkunBayar {
		badRequest["akun_bayar_id"] = []string{err.Error()}
	}

	if err.Error() == message.KontakNotFound {
		badRequest["kontak_id"] = []string{err.Error()}
	}

	if err.Error() == message.HutangPiutangNotFound {
		badRequest["hutang_piutang_id"] = []string{err.Error()}
	}

	if err.Error() == message.BayarMustLessThanSisaTagihan {
		badRequest["total"] = []string{err.Error()}
	}

	if len(badRequest) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorRes(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, badRequest))
	}
	return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorRes(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, nil))
}

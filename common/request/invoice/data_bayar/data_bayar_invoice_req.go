package invoice

import (
	"github.com/be-sistem-informasi-konveksi/entity"
)

type Create struct {
	InvoiceID       string                 `params:"invoice_id" validate:"required,ulid"`
	BuktiPembayaran entity.BuktiPembayaran `json:"bukti_pembayaran" validate:"required,dive,url"`
	Keterangan      string                 `json:"keterangan" validate:"required"`
	AkunID          string                 `json:"akun_id" validate:"required,ulid"`
	Total           float64                `json:"total" validate:"required,number,gt=0"`
}

type Update struct {
	ID              string                 `params:"id" validate:"required,ulid"`
	BuktiPembayaran entity.BuktiPembayaran `json:"bukti_pembayaran" validate:"omitempty,dive,url"`
	Keterangan      string                 `json:"keterangan" validate:"omitempty"`
	AkunID          string                 `json:"akun_id" validate:"omitempty,ulid"`
	Total           float64                `json:"total" validate:"omitempty,number,gt=0"`
}

type GetByInvoiceID struct {
	InvoiceID string `params:"invoice_id" validate:"required,ulid"`
	AkunID    string `query:"akun_id" validate:"omitempty,ulid"`
	Status    string `query:"statu" validate:"oneof=TERKONFIRMASI BELUM_TERKONFIRMASI"`
}

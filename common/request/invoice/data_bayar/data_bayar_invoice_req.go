package req_invoice_data_bayar

import (
	"mime/multipart"
)

type CreateByInvoiceId struct {
	InvoiceID       string                  `params:"id" validate:"required,ulid"`
	BuktiPembayaran []*multipart.FileHeader `form:"bukti_pembayaran" validate:"required"`
	Data            string                  `form:"data"`
	Keterangan      string                  `json:"keterangan" validate:"required"`
	AkunID          string                  `json:"akun_id" validate:"required,ulid"`
	Total           float64                 `json:"total" validate:"required,number,gt=0"`
}

type Update struct {
	ID              string                  `params:"id" validate:"required,ulid"`
	BuktiPembayaran []*multipart.FileHeader `form:"bukti_pembayaran" validate:"omitempty"`
	Data            string                  `form:"data"`
	Status          string                  `json:"status" validate:"omitempty,oneof=TERKONFIRMASI BELUM_TERKONFIRMASI"`
	Keterangan      string                  `json:"keterangan" validate:"omitempty"`
	AkunID          string                  `json:"akun_id" validate:"omitempty,ulid"`
	Total           float64                 `json:"total" validate:"omitempty,number,gt=0"`
}

type GetByInvoiceID struct {
	InvoiceID string `params:"id" validate:"required,ulid"`
	AkunID    string `query:"akun_id" validate:"omitempty,ulid"`
	Status    string `query:"status" validate:"omitempty,oneof=TERKONFIRMASI BELUM_TERKONFIRMASI"`
}

type GetAll struct {
	KontakID string `query:"kontak_id" validate:"omitempty,ulid"`
	AkunID   string `query:"akun_id" validate:"omitempty,ulid"`
	Status   string `query:"status" validate:"omitempty,oneof=TERKONFIRMASI BELUM_TERKONFIRMASI"`
	Sort     string `query:"sort" validate:"omitempty,oneof=ASC DESC"`
	Next     string `query:"next" validate:"omitempty,ulid"`
	Limit    int    `query:"limit" validate:"omitempty,number"`
}

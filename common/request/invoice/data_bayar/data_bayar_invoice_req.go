package invoice

import (
	"github.com/be-sistem-informasi-konveksi/common/request/invoice"
)

type Create struct {
	InvoiceID string `params:"invoice_id" validate:"required,ulid"`
	invoice.ReqBayar
}

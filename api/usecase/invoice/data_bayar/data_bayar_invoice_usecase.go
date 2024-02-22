package invoice

import (
	"context"

	req "github.com/be-sistem-informasi-konveksi/common/request/invoice/data_bayar"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type (
	ParamCreateDataBayar struct {
		Ctx context.Context
		Req req.Create
	}
	ParamCommitDB struct {
		Ctx       context.Context
		DataBayar *entity.DataBayarInvoice
	}
)

type DataBayarInvoice interface {
	CreateDataBayar(param ParamCreateDataBayar) (*entity.DataBayarInvoice, error)
}

type dataBayarInvoice struct {
	ulid pkg.UlidPkg
}

func NewDataInvoice(ulid pkg.UlidPkg) DataBayarInvoice {
	return &dataBayarInvoice{ulid}
}

func (u *dataBayarInvoice) CreateDataBayar(param ParamCreateDataBayar) (*entity.DataBayarInvoice, error) {
	dataBayar := entity.DataBayarInvoice{
		Base: entity.Base{
			ID: u.ulid.MakeUlid().String(),
		},
		InvoiceID:       param.Req.InvoiceID,
		AkunID:          param.Req.AkunID,
		Keterangan:      param.Req.Keterangan,
		BuktiPembayaran: param.Req.BuktiPembayaran,
		Total:           param.Req.Total,
	}
	return &dataBayar, nil
}

// func (u *dataBayarInvoice) CreateCommitDB(param ParamCommitDB) error {

// }

package akuntansi

import (
	"errors"
	"math"
	"time"

	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/akuntansi/hutang_piutang"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

func autoCreateDataBayarAwal(trByr req.ReqBayar, ayTagihan entity.AyatJurnal, kontakId, keterangan string, ulid pkg.UlidPkg) (*entity.DataBayarHutangPiutang, error) {
	tanggal, err := time.Parse(time.RFC3339, trByr.Tanggal)
	if err != nil {
		helper.LogsError(err)
		return nil, err
	}
	trByr.Total = math.Abs(trByr.Total)
	// cek jika total bayarnya lebih besar dari sisa tagihan
	if trByr.Total > ayTagihan.Debit {
		return nil, errors.New(message.BayarMustLessThanSisaTagihan)
	}
	byrHP := entity.DataBayarHutangPiutang{
		Base: entity.Base{
			ID: ulid.MakeUlid().String(),
		},
		Total: trByr.Total,
		Transaksi: entity.Transaksi{
			Base: entity.Base{
				ID: ulid.MakeUlid().String(),
			},
			Keterangan:      keterangan,
			BuktiPembayaran: trByr.BuktiPembayaran,
			Total:           trByr.Total,
			Tanggal:         tanggal,
			KontakID:        kontakId,
			AyatJurnals: []entity.AyatJurnal{
				{
					Base: entity.Base{
						ID: ulid.MakeUlid().String(),
					},
					AkunID: trByr.AkunBayarID,
					Debit:  trByr.Total,
					Saldo:  trByr.Total,
				},
				// ay kredit, bayar mengurangi akun utang/piutangnya
				{
					Base: entity.Base{
						ID: ulid.MakeUlid().String(),
					},
					AkunID: ayTagihan.AkunID,
					Kredit: trByr.Total,
					Saldo:  -trByr.Total,
				},
			},
		},
	}

	return &byrHP, nil
}

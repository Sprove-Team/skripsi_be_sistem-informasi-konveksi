package akuntansi

import "github.com/be-sistem-informasi-konveksi/entity"

type ResDataHutangPiutang struct {
	ID          string  `json:"id"`
	InvoiceID   string  `json:"invoice_id,omitempty"`
	Jenis       string  `json:"jenis"`
	TransaksiID string  `json:"transaksi_id"`
	Status      string  `json:"status"`
	Total       float64 `json:"total"`
	Sisa        float64 `json:"sisa"`
}

type GetAll struct {
	Nama          string                 `json:"nama"`
	TotalPiutang  float64                `json:"total_piutang"`
	SisaPiutang   float64                `json:"sisa_piutang"`
	TotalHutang   float64                `json:"total_hutang"`
	SisaHutang    float64                `json:"sisa_hutang"`
	HutangPiutang []ResDataHutangPiutang `json:"hutang_piutang"`
}

type GetById struct {
	ID          string                          `json:"id"`
	Nama        string                          `json:"nama"`
	InvoiceSlug string                          `json:"invoice_slug,omitempty"`
	Jenis       string                          `json:"jenis"`
	TransaksiID string                          `json:"transaksi_id"`
	Status      string                          `json:"status"`
	Total       float64                         `json:"total"`
	Sisa        float64                         `json:"sisa"`
	DataBayar   []entity.DataBayarHutangPiutang `json:"data_bayar"`
}

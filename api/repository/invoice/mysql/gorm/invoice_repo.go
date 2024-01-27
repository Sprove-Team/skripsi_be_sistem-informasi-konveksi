package invoice

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type CreateParam struct {
	Invoice       *entity.Invoice
	Kontak        *entity.Kontak
	Transaksi     *entity.Transaksi
	DetailInvoice []entity.DetailInvoice
}

type SearchFilter struct {
	StatusProduksi  string
	Kepada          string
	TanggalDeadline time.Time
	TanggalKirim    time.Time
	Order           string
	Next            string
	Limit           int
}

type InvoiceRepo interface {
	Create(ctx context.Context, param CreateParam) error
	GetLastInvoiceCrrYear(ctx context.Context) (entity.Invoice, error)
	GetAll(ctx context.Context, filter SearchFilter) ([]entity.Invoice, error)
}

type invoiceRepo struct {
	DB *gorm.DB
}

func NewInvoiceRepo(DB *gorm.DB) InvoiceRepo {
	return &invoiceRepo{DB}
}

func (r *invoiceRepo) Create(ctx context.Context, param CreateParam) error {
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(param.Kontak).Error; err != nil {
			return err
		}
		if err := tx.Create(param.Transaksi).Error; err != nil {
			return err
		}
		if err := tx.Create(param.Invoice).Error; err != nil {
			return err
		}
		if err := tx.Create(&param.DetailInvoice).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

func (r *invoiceRepo) GetLastInvoiceCrrYear(ctx context.Context) (entity.Invoice, error) {
	currentYear := time.Now().Year()
	invoice := new(entity.Invoice)
	err := r.DB.WithContext(ctx).Model(invoice).Select("nomor_referensi").Order("nomor_referensi DESC").First(invoice, "YEAR(created_at) = ?", currentYear).Error
	if err != nil {
		helper.LogsError(err)
		return *invoice, err
	}
	return *invoice, err
}

func (r *invoiceRepo) GetAll(ctx context.Context, filter SearchFilter) ([]entity.Invoice, error) {
	tx := r.DB.WithContext(ctx).Model(&entity.Invoice{})

	if filter.Order == "" {
		tx = tx.Order("id ASC")
	}

	s := reflect.ValueOf(filter)
	orderBy := strings.Split(filter.Order, " ")[1]
	for i := 0; i < s.NumField(); i++ {
		// kind := s.Field(i).Kind()
		key := s.Type().Field(i).Name
		value := s.Field(i).Interface()

		switch value.(type) {
		case string:
			if value != "" {
				switch key {
				case "Next":
					if orderBy == "DESC" {
						tx = tx.Where("id < ?", value)
					} else {
						tx = tx.Where("id > ?", value)
					}
				case "Order":
					tx = tx.Order(value)
				case "Kepada":
					tx = tx.Where("kepada LIKE ?", "%"+value.(string)+"%")
				case "StatusProduksi":
					tx = tx.Where("status_produksi = ?", value)
				}
			}
		case time.Time:
			if !value.(time.Time).IsZero() {
				switch key {
				case "TanggalDeadline":
					tx = tx.Where("DATE(tanggal_deadline) = ?", value)
				case "TanggalKirim":
					tx = tx.Where("DATE(tanggal_kirim) = ?", value)
				}
			}
		}
	}

	invoices := make([]entity.Invoice, filter.Limit)

	if err := tx.Limit(filter.Limit).Find(&invoices).Error; err != nil {
		return nil, err
	}

	return invoices, nil
}

package repo_akuntansi_hp

import (
	"context"

	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"gorm.io/gorm"
)

type (
	ParamCreate struct {
		Ctx           context.Context
		HutangPiutang *entity.HutangPiutang
	}
	ParamGetAll struct {
		Ctx      context.Context
		KontakID string
		Jenis    []string
		Status   []string
		Limit    int
		Next     string
	}
	ParamGetById struct {
		Ctx context.Context
		ID  string
	}
	ParamGetByInvoiceId struct {
		Ctx context.Context
		ID  string
	}
	ParamGetByTrId struct {
		Ctx context.Context
		ID  string
	}
	ParamGetHPForBayar struct {
		Ctx context.Context
		ID  string
	}
)

type HutangPiutangRepo interface {
	Create(param ParamCreate) error
	GetAll(param ParamGetAll) ([]entity.Kontak, error)
	GetById(param ParamGetById) (*entity.HutangPiutang, error)
	GetByInvoiceId(param ParamGetByInvoiceId) (*entity.HutangPiutang, error)
	GetByTrId(param ParamGetByTrId) (*entity.HutangPiutang, error)
	GetHPForBayar(param ParamGetHPForBayar) (*entity.HutangPiutang, error)
}

type hutangPiutangRepo struct {
	DB *gorm.DB
}

func NewHutangPiutangRepo(DB *gorm.DB) HutangPiutangRepo {
	return &hutangPiutangRepo{DB}
}

func (r *hutangPiutangRepo) Create(param ParamCreate) error {
	if err := r.DB.WithContext(param.Ctx).Create(param.HutangPiutang).Error; err != nil {
		return err
	}
	return nil
}

func (r *hutangPiutangRepo) GetAll(param ParamGetAll) ([]entity.Kontak, error) {
	datas := make([]entity.Kontak, 0, param.Limit)
	tx := r.DB.WithContext(param.Ctx).Model(datas).Order("id ASC")

	if param.KontakID != "" {
		tx = tx.Where("id = ?", param.KontakID)
	} else if param.Next != "" {
		tx = tx.Where("id > ?", param.Next)
	}

	tx = tx.Preload("Transaksi", func(db *gorm.DB) *gorm.DB {
		tx2 := db.Joins("INNER JOIN hutang_piutang ON transaksi.id = hutang_piutang.transaksi_id")
		if param.Jenis != nil && param.Status != nil {
			tx2 = tx2.Where("hutang_piutang.jenis IN (?) AND hutang_piutang.status IN (?)", param.Jenis, param.Status)
		} else if param.Jenis != nil {
			tx2 = tx2.Where("hutang_piutang.jenis IN (?)", param.Jenis)
		} else if param.Status != nil {
			tx2 = tx2.Where("hutang_piutang.status IN (?)", param.Status)
		}
		return tx2.Preload("HutangPiutang")
	})
	if err := tx.Limit(param.Limit).Find(&datas).Error; err != nil {
		helper.LogsError(err)
		return nil, err
	}
	return datas, nil
}

func (r *hutangPiutangRepo) GetById(param ParamGetById) (*entity.HutangPiutang, error) {
	data := new(entity.HutangPiutang)
	tx := r.DB.WithContext(param.Ctx).Model(data).Where("id = ?", param.ID).Order("id ASC")
	tx = tx.
		Preload("Transaksi").
		Preload("Transaksi.Kontak").
		Preload("DataBayarHutangPiutang").
		Preload("DataBayarHutangPiutang.Transaksi")

	if err := tx.First(data).Error; err != nil {
		helper.LogsError(err)
		return data, err
	}

	return data, nil
}

func (r *hutangPiutangRepo) GetByInvoiceId(param ParamGetByInvoiceId) (*entity.HutangPiutang, error) {
	data := new(entity.HutangPiutang)
	tx := r.DB.WithContext(param.Ctx).Model(data).Where("invoice_id = ?", param.ID).Order("id ASC")

	if err := tx.First(data).Error; err != nil {
		helper.LogsError(err)
		return data, err
	}
	return data, nil
}

func (r *hutangPiutangRepo) GetByTrId(param ParamGetByTrId) (*entity.HutangPiutang, error) {
	data := new(entity.HutangPiutang)
	tx := r.DB.WithContext(param.Ctx).Model(data).Where("transaksi_id = ?", param.ID)

	if err := tx.First(data).Error; err != nil {
		helper.LogsError(err)
		return data, err
	}

	return data, nil
}

func (r *hutangPiutangRepo) GetHPForBayar(param ParamGetHPForBayar) (*entity.HutangPiutang, error) {
	data := new(entity.HutangPiutang)
	tx := r.DB.WithContext(param.Ctx).Model(data).Where("id = ?", param.ID)
	tx = tx.
		Preload("Transaksi").
		Preload("Transaksi.AyatJurnals").
		Preload("Transaksi.AyatJurnals.Akun")

	if err := tx.First(data).Error; err != nil {
		return data, err
	}
	return data, nil
}

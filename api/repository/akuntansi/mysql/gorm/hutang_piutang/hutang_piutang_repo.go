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
		Ctx    context.Context
		Search SearchParam
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

type SearchParam struct {
	KontakID string
	Jenis    []string
	Status   []string
	Next     string
	Limit    int
}

type HutangPiutangRepo interface {
	Create(param ParamCreate) error
	GetAll(param ParamGetAll) ([]entity.HutangPiutang, error)
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

func (r *hutangPiutangRepo) GetAll(param ParamGetAll) ([]entity.HutangPiutang, error) {
	datas := []entity.HutangPiutang{}
	tx := r.DB.WithContext(param.Ctx).Model(&entity.HutangPiutang{}).Order("id ASC")

	if param.Search.Next != "" {
		tx = tx.Where("id > ?", param.Search.Next)
	}

	if param.Search.Jenis != nil {
		tx = tx.Where("jenis IN (?)", param.Search.Jenis)
	}

	if param.Search.Status != nil {
		tx = tx.Where("status IN (?)", param.Search.Status)
	}

	if param.Search.KontakID != "" {
		tx = tx.Where("kontak_id = ?", param.Search.KontakID)
	}

	tx = tx.
		Preload("Transaksi").
		Preload("Transaksi.Kontak").
		Preload("DataBayarHutangPiutang").
		Preload("DataBayarHutangPiutang.Transaksi").
		Limit(param.Search.Limit)

	if err := tx.Find(&datas).Error; err != nil {
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
		Preload("Transaksi.AyatJurnals.Akun")

	if err := tx.First(data).Error; err != nil {
		return data, err
	}
	return data, nil
}

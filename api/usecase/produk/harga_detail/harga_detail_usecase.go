package produk

import (
	"context"
	"errors"

	produkRepo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm"
	repo "github.com/be-sistem-informasi-konveksi/api/repository/produk/mysql/gorm/harga_detail"
	req "github.com/be-sistem-informasi-konveksi/common/request/produk/harga_detail"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type HargaDetailProdukUsecase interface {
	CreateByProdukId(ctx context.Context, reqHargaDetailProduk req.CreateByProdukId) error
	UpdateByProdukId(ctx context.Context, reqHargaDetailProduk req.UpdateByProdukId) error
	DeleteByProdukId(ctx context.Context, produkId string) error
	Delete(ctx context.Context, id string) error
	GetByProdukId(ctx context.Context, reqHargaDetailProduk req.GetByProdukId) ([]entity.HargaDetailProduk, error)
}

type hargaDetailProdukUsecase struct {
	repo    repo.HargaDetailProdukRepo
	produkR produkRepo.ProdukRepo
	ulid    pkg.UlidPkg
}

func NewHargaDetailProdukUsecase(repo repo.HargaDetailProdukRepo, produkR produkRepo.ProdukRepo, ulid pkg.UlidPkg) HargaDetailProdukUsecase {
	return &hargaDetailProdukUsecase{repo, produkR, ulid}
}

func (u *hargaDetailProdukUsecase) CreateByProdukId(ctx context.Context, reqHargaDetailProduk req.CreateByProdukId) error {
	_, err := u.produkR.GetById(ctx, reqHargaDetailProduk.ProdukId)
	if err != nil {
		return err
	}

	newData := make([]*entity.HargaDetailProduk, len(reqHargaDetailProduk.HargaDetail))
	qty := make([]uint, len(reqHargaDetailProduk.HargaDetail))

	for i, hargaDetail := range reqHargaDetailProduk.HargaDetail {
		newData[i] = &entity.HargaDetailProduk{
			Base: entity.Base{
				ID: u.ulid.MakeUlid().String(),
			},
			ProdukID: reqHargaDetailProduk.ProdukId,
			QTY:      hargaDetail.QTY,
			Harga:    hargaDetail.Harga,
		}
		qty[i] = hargaDetail.QTY
	}

	datas, err := u.repo.GetByInQtyProdukId(ctx, qty, reqHargaDetailProduk.ProdukId)

	if err != nil {
		return err
	}

	if len(datas) > 0 {
		return errors.New("duplicated key not allowed")
	}

	if err := u.repo.Create(ctx, newData); err != nil {
		return err
	}
	return nil
}

func (u *hargaDetailProdukUsecase) UpdateByProdukId(ctx context.Context, reqHargaDetailProduk req.UpdateByProdukId) error {
	_, err := u.produkR.GetById(ctx, reqHargaDetailProduk.ProdukId)
	if err != nil {
		return err
	}

	hargaDetails := make([]entity.HargaDetailProduk, len(reqHargaDetailProduk.HargaDetail))
	for i, hargaDetail := range reqHargaDetailProduk.HargaDetail {
		hargaDetails[i] = entity.HargaDetailProduk{
			Base: entity.Base{
				ID: hargaDetail.ID,
			},
			ProdukID: reqHargaDetailProduk.ProdukId,
			QTY:      hargaDetail.QTY,
			Harga:    hargaDetail.Harga,
		}
	}

	err = u.repo.UpdateByProdukId(ctx, reqHargaDetailProduk.ProdukId, hargaDetails)

	if err != nil {
		return err
	}
	return nil
}

func (u *hargaDetailProdukUsecase) DeleteByProdukId(ctx context.Context, produkId string) error {
	_, err := u.produkR.GetById(ctx, produkId)
	if err != nil {
		return err
	}
	err = u.repo.DeleteByProdukId(ctx, produkId)
	if err != nil {
		return err
	}
	return nil
}
func (u *hargaDetailProdukUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *hargaDetailProdukUsecase) GetByProdukId(ctx context.Context, reqHargaDetailProduk req.GetByProdukId) ([]entity.HargaDetailProduk, error) {
	datas, err := u.repo.GetByProdukId(ctx, reqHargaDetailProduk.ProdukId)
	if len(datas) == 0 {
		return nil, err
	}
	return datas, nil
}

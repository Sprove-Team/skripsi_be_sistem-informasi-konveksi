package produk

import (
	"context"

	req "github.com/be-sistem-informasi-konveksi/common/request/produk"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/produk/mysql/gorm"
)

type KategoriProdukUsecase interface {
	Create(ctx context.Context, kategoriProduk req.CreateKategoriProduk) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, kategoriProduk req.UpdateKategoriProduk) error
	GetAll(ctx context.Context) ([]entity.KategoriProduk, error)
	GetById(ctx context.Context, id string) (entity.KategoriProduk, error)
}

type kategoriProdukUsecase struct {
	repo repo.KategoriProdukRepo
	uuidGen         helper.UuidGenerator
}

func NewKategoriProdukUsecase(repo repo.KategoriProdukRepo, uuidGen helper.UuidGenerator) KategoriProdukUsecase {
	return &kategoriProdukUsecase{repo, uuidGen}
}

func (u *kategoriProdukUsecase) Create(ctx context.Context, kategoriProduk req.CreateKategoriProduk) error {
	id, _ := u.uuidGen.GenerateUUID()
	kategoriProdukR := entity.KategoriProduk{
		ID: id,
		Nama: kategoriProduk.Nama,
	}
	return u.repo.Create(ctx, &kategoriProdukR)
}

func (u *kategoriProdukUsecase) Update(ctx context.Context, kategoriProduk req.UpdateKategoriProduk) error {
	_, err := u.repo.GetById(ctx, kategoriProduk.ID)
	if err != nil {
		return err
	}
	ketegoriProdukR := entity.KategoriProduk{
		ID:   kategoriProduk.ID,
		Nama: kategoriProduk.Nama,
	}

	return u.repo.Update(ctx, &ketegoriProdukR)
}

func (u *kategoriProdukUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.repo.Delete(ctx, id)
	return err
}

func (u *kategoriProdukUsecase) GetById(ctx context.Context, id string) (entity.KategoriProduk, error){
	kategoriProduk, err := u.repo.GetById(ctx, id)
	return kategoriProduk, err
}

func (u *kategoriProdukUsecase) GetAll(ctx context.Context) ([]entity.KategoriProduk, error) {
	kategoriProduks, err := u.repo.GetAll(ctx)
	return kategoriProduks, err
}

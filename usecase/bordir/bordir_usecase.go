package bordir

import (
	"context"

	req "github.com/be-sistem-informasi-konveksi/common/request/bordir"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/bordir/mysql/gorm"
)

type BordirUsecase interface {
	Create(ctx context.Context, kategoriProduk req.CreateBordir) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, kategoriProduk req.UpdateBordir) error
	GetAll(ctx context.Context, get req.GetAllBordir) ([]entity.Bordir, int, int,error)
	GetById(ctx context.Context, id string) (entity.Bordir, error)
}

type bordirUsecase struct {
	repo    repo.BordirRepo
	uuidGen helper.UuidGenerator
	paginate helper.Paginate
}

func NewKategoriProdukUsecase(repo repo.BordirRepo, uuidGen helper.UuidGenerator, paginate helper.Paginate) BordirUsecase {
	return &bordirUsecase{repo, uuidGen, paginate}
}

func (u *bordirUsecase) Create(ctx context.Context, kategoriProduk req.CreateBordir) error {
	id, _ := u.uuidGen.GenerateUUID()
	kategoriProdukR := entity.Bordir{
		ID:   id,
		Nama: kategoriProduk.Nama,
	}
	return u.repo.Create(ctx, &kategoriProdukR)
}

func (u *bordirUsecase) Update(ctx context.Context, bordir req.UpdateBordir) error {
	_, err := u.repo.GetById(ctx, bordir.ID)
	if err != nil {
		return err
	}
	bordirR := entity.Bordir{
		ID:   bordir.ID,
		Nama: bordir.Nama,
		Harga: bordir.Harga,
	}

	return u.repo.Update(ctx, &bordirR)
}

func (u *bordirUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.repo.Delete(ctx, id)
	return err
}

func (u *bordirUsecase) GetById(ctx context.Context, id string) (entity.Bordir, error) {
	kategoriProduk, err := u.repo.GetById(ctx, id)
	return kategoriProduk, err
}

func (u *bordirUsecase) GetAll(ctx context.Context, get req.GetAllBordir) ([]entity.Bordir, int, int,error) {
	
	currentPage, offset, limit := u.paginate.GetPaginateData(get.Page, get.Limit)

	bordirs, totalData, err := u.repo.GetAll(ctx, repo.SearchParams{
		Nama:             get.Search.Nama,
		Limit:            limit,
		Offset:           offset,
	})

	if err != nil {
		return nil, currentPage, 0, err
	}

	totalPage := u.paginate.GetTotalPages(int(totalData), limit)

	return bordirs, currentPage, totalPage, err
}

package uc_bordir

import (
	"context"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/bordir/mysql/gorm"
	req "github.com/be-sistem-informasi-konveksi/common/request/bordir"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type BordirUsecase interface {
	Create(ctx context.Context, reqBordir req.Create) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, reqBordir req.Update) error
	GetAll(ctx context.Context, reqBordir req.GetAll) ([]entity.Bordir, error)
	GetById(ctx context.Context, id string) (entity.Bordir, error)
}

type bordirUsecase struct {
	repo repo.BordirRepo
	// uuidGen  pkg.UuidGenerator
	ulid pkg.UlidPkg
}

func NewBordirUsecase(repo repo.BordirRepo, ulid pkg.UlidPkg) BordirUsecase {
	return &bordirUsecase{repo, ulid}
}

func (u *bordirUsecase) Create(ctx context.Context, reqBordir req.Create) error {
	id := u.ulid.MakeUlid().String()
	data := entity.Bordir{
		Base: entity.Base{
			ID: id,
		},
		Nama:  reqBordir.Nama,
		Harga: reqBordir.Harga,
	}
	return u.repo.Create(ctx, &data)
}

func (u *bordirUsecase) Update(ctx context.Context, reqBordir req.Update) error {
	_, err := u.repo.GetById(ctx, reqBordir.ID)
	if err != nil {
		return err
	}
	data := entity.Bordir{
		Base: entity.Base{
			ID: reqBordir.ID,
		},
		Nama:  reqBordir.Nama,
		Harga: reqBordir.Harga,
	}

	return u.repo.Update(ctx, &data)
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
	data, err := u.repo.GetById(ctx, id)
	return data, err
}

func (u *bordirUsecase) GetAll(ctx context.Context, reqBordir req.GetAll) ([]entity.Bordir, error) {
	// currentPage, offset, limit := u.paginate.GetPaginateData(reqBordir.Page, reqBordir.Limit)

	if reqBordir.Limit <= 0 {
		reqBordir.Limit = 10
	}

	datas, err := u.repo.GetAll(ctx, repo.SearchBordir{
		Nama:  reqBordir.Nama,
		Limit: reqBordir.Limit,

		// Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	// totalPage := u.paginate.GetTotalPages(int(totalData), limit)

	return datas, err
}

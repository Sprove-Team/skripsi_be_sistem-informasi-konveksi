package bordir

import (
	"context"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/bordir/mysql/gorm"
	req "github.com/be-sistem-informasi-konveksi/common/request/bordir"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type BordirUsecase interface {
	Create(ctx context.Context, reqBordir req.CreateBordir) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, reqBordir req.UpdateBordir) error
	GetAll(ctx context.Context, reqBordir req.GetAllBordir) ([]entity.Bordir, int, int, error)
	GetById(ctx context.Context, id string) (entity.Bordir, error)
}

type bordirUsecase struct {
	repo     repo.BordirRepo
	uuidGen  pkg.UuidGenerator
	paginate helper.Paginate
}

func NewBordirUsecase(repo repo.BordirRepo, uuidGen pkg.UuidGenerator, paginate helper.Paginate) BordirUsecase {
	return &bordirUsecase{repo, uuidGen, paginate}
}

func (u *bordirUsecase) Create(ctx context.Context, reqBordir req.CreateBordir) error {
	id, _ := u.uuidGen.GenerateUUID()
	data := entity.Bordir{
		ID:    id,
		Nama:  reqBordir.Nama,
		Harga: reqBordir.Harga,
	}
	return u.repo.Create(ctx, &data)
}

func (u *bordirUsecase) Update(ctx context.Context, reqBordir req.UpdateBordir) error {
	_, err := u.repo.GetById(ctx, reqBordir.ID)
	if err != nil {
		return err
	}
	data := entity.Bordir{
		ID:    reqBordir.ID,
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

func (u *bordirUsecase) GetAll(ctx context.Context, reqBordir req.GetAllBordir) ([]entity.Bordir, int, int, error) {
	currentPage, offset, limit := u.paginate.GetPaginateData(reqBordir.Page, reqBordir.Limit)

	datas, totalData, err := u.repo.GetAll(ctx, repo.SearchBordir{
		Nama:   reqBordir.Search.Nama,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, currentPage, 0, err
	}

	totalPage := u.paginate.GetTotalPages(int(totalData), limit)

	return datas, currentPage, totalPage, err
}

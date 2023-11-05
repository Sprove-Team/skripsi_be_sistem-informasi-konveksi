package sablon

import (
	"context"

	req "github.com/be-sistem-informasi-konveksi/common/request/sablon"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/sablon/mysql/gorm"
)

type SablonUsecase interface {
	Create(ctx context.Context, reqSablon req.CreateSablon) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, reqSablon req.UpdateSablon) error
	GetAll(ctx context.Context, reqSablon req.GetAllSablon) ([]entity.Sablon, int, int, error)
	GetById(ctx context.Context, id string) (entity.Sablon, error)
}

type sablonUsecase struct {
	repo     repo.SablonRepo
	uuidGen  helper.UuidGenerator
	paginate helper.Paginate
}

func NewSablonUsecase(repo repo.SablonRepo, uuidGen helper.UuidGenerator, paginate helper.Paginate) SablonUsecase {
	return &sablonUsecase{repo, uuidGen, paginate}
}

func (u *sablonUsecase) Create(ctx context.Context, sablon req.CreateSablon) error {
	id, _ := u.uuidGen.GenerateUUID()
	data := entity.Sablon{
		ID:    id,
		Nama:  sablon.Nama,
		Harga: sablon.Harga,
	}
	return u.repo.Create(ctx, &data)
}

func (u *sablonUsecase) Update(ctx context.Context, reqSablon req.UpdateSablon) error {
	_, err := u.repo.GetById(ctx, reqSablon.ID)
	if err != nil {
		return err
	}
	data := entity.Sablon{
		ID:    reqSablon.ID,
		Nama:  reqSablon.Nama,
		Harga: reqSablon.Harga,
	}

	return u.repo.Update(ctx, &data)
}

func (u *sablonUsecase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	err = u.repo.Delete(ctx, id)
	return err
}

func (u *sablonUsecase) GetById(ctx context.Context, id string) (entity.Sablon, error) {
	data, err := u.repo.GetById(ctx, id)
	return data, err
}

func (u *sablonUsecase) GetAll(ctx context.Context, reqSablon req.GetAllSablon) ([]entity.Sablon, int, int, error) {
	currentPage, offset, limit := u.paginate.GetPaginateData(reqSablon.Page, reqSablon.Limit)

	datas, totalData, err := u.repo.GetAll(ctx, repo.SearchSablon{
		Nama:   reqSablon.Search.Nama,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, currentPage, 0, err
	}

	totalPage := u.paginate.GetTotalPages(int(totalData), limit)

	return datas, currentPage, totalPage, err
}

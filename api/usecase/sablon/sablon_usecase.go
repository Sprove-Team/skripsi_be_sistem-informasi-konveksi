package uc_sablon

import (
	"context"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/sablon/mysql/gorm"
	req "github.com/be-sistem-informasi-konveksi/common/request/sablon"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type SablonUsecase interface {
	Create(ctx context.Context, reqSablon req.Create) error
	Update(ctx context.Context, reqSablon req.Update) error
	Delete(ctx context.Context, id string) error
	GetById(ctx context.Context, id string) (entity.Sablon, error)
	GetAll(ctx context.Context, reqSablon req.GetAll) ([]entity.Sablon, error)
}

type sablonUsecase struct {
	repo repo.SablonRepo
	ulid pkg.UlidPkg
}

func NewSablonUsecase(repo repo.SablonRepo, ulid pkg.UlidPkg) SablonUsecase {
	return &sablonUsecase{repo, ulid}
}

func (u *sablonUsecase) Create(ctx context.Context, sablon req.Create) error {
	id := u.ulid.MakeUlid().String()
	data := entity.Sablon{
		Base: entity.Base{
			ID: id,
		},
		Nama:  sablon.Nama,
		Harga: sablon.Harga,
	}
	return u.repo.Create(ctx, &data)
}

func (u *sablonUsecase) Update(ctx context.Context, reqSablon req.Update) error {
	_, err := u.repo.GetById(ctx, reqSablon.ID)
	if err != nil {
		return err
	}
	data := entity.Sablon{
		Base: entity.Base{
			ID: reqSablon.ID,
		},
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

func (u *sablonUsecase) GetAll(ctx context.Context, reqSablon req.GetAll) ([]entity.Sablon, error) {

	if reqSablon.Limit <= 0 {
		reqSablon.Limit = 10
	}

	datas, err := u.repo.GetAll(ctx, repo.SearchSablon{
		Nama:  reqSablon.Nama,
		Limit: reqSablon.Limit,
		Next:  reqSablon.Next,
	})
	if err != nil {
		return nil, err
	}

	// totalPage := u.paginate.GetTotalPages(int(totalData), limit)

	return datas, err
}

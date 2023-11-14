package user

import (
	"context"

	"golang.org/x/sync/errgroup"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	req "github.com/be-sistem-informasi-konveksi/common/request/user"
	res "github.com/be-sistem-informasi-konveksi/common/response/user"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
)

type UserUsecase interface {
	Create(ctx context.Context, reqUser req.CreateUser) error
	Update(ctx context.Context, reqUser req.UpdateUser) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, reqUser req.GetAllUser) ([]res.DataGetAllUserRes, int, int, error)
}

type userUsecase struct {
	repo         repo.UserRepo
	jenisSpvRepo repo.JenisSpvRepo
	uuidGen      helper.UuidGenerator
	paginate     helper.Paginate
	encryptor    helper.Encryptor
}

func NewUserUsecase(repo repo.UserRepo, jenisSpvRepo repo.JenisSpvRepo, uuidGen helper.UuidGenerator, paginate helper.Paginate, encryptor helper.Encryptor) UserUsecase {
	return &userUsecase{repo, jenisSpvRepo, uuidGen, paginate, encryptor}
}

func (u *userUsecase) Create(ctx context.Context, reqUser req.CreateUser) error {
	id, _ := u.uuidGen.GenerateUUID()
	pass, err := u.encryptor.HashPassword(reqUser.Password)
	if err != nil {
		return err
	}
	user := entity.User{
		ID:       id,
		Nama:     reqUser.Nama,
		Role:     reqUser.Role,
		NoTelp:   reqUser.NoTelp,
		Alamat:   reqUser.Alamat,
		Username: reqUser.Username,
		Password: pass,
	}
	if reqUser.JenisSpvID != "" {
		user.JenisSpvID = reqUser.JenisSpvID
	}
	err = u.repo.Create(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) Update(ctx context.Context, reqUser req.UpdateUser) error {
	user := entity.User{
		ID:       reqUser.ID,
		Nama:     reqUser.Nama,
		Role:     reqUser.Role,
		NoTelp:   reqUser.NoTelp,
		Alamat:   reqUser.Alamat,
		Username: reqUser.Username,
	}

	if reqUser.Password != "" {
		pass, err := u.encryptor.HashPassword(reqUser.Password)
		if err != nil {
			return err
		}
		user.Password = pass
	}

	if reqUser.JenisSpvID != "" {
		user.JenisSpvID = reqUser.JenisSpvID
	}

	err := u.repo.Update(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *userUsecase) GetAll(ctx context.Context, reqUser req.GetAllUser) ([]res.DataGetAllUserRes, int, int, error) {
	currentPage, offset, limit := u.paginate.GetPaginateData(reqUser.Page, reqUser.Limit)
	datas, totalData, err := u.repo.GetAll(ctx, repo.SearchUser{
		Nama:       reqUser.Search.Nama,
		Alamat:     reqUser.Search.Alamat,
		Username:   reqUser.Search.Username,
		NoTelp:     reqUser.Search.NoTelp,
		Role:       reqUser.Search.Role,
		JenisSpvId: reqUser.Search.JenisSpvID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, currentPage, 0, err
	}
	datasRes := make([]res.DataGetAllUserRes, len(datas))
	g := &errgroup.Group{}
	for i, d := range datas {
		i := i
		d := d
		g.Go(func() error {
			jenisSpv, err := u.jenisSpvRepo.GetById(ctx, d.JenisSpvID)
			if err != nil {
				return err
			}

			datasRes[i] = res.DataGetAllUserRes{
				ID:       d.ID,
				Nama:     d.Nama,
				Alamat:   d.Alamat,
				Username: d.Username,
				NoTelp:   d.NoTelp,
				Role:     d.Role,
			}

			if jenisSpv.Nama != "" {
				datasRes[i].JenisSpv = jenisSpv
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, currentPage, 0, err
	}

	totalPage := u.paginate.GetTotalPages(int(totalData), limit)

	return datasRes, currentPage, totalPage, err
}

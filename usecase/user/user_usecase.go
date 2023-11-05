package user

import (
	"context"

	req "github.com/be-sistem-informasi-konveksi/common/request/user"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	repo "github.com/be-sistem-informasi-konveksi/repository/user/mysql/gorm"
)

type UserUsecase interface {
	Create(ctx context.Context, reqUser req.CreateUser) error
	GetAll(ctx context.Context, reqUser req.GetAllUser) ([]entity.User, int, int, error)
}

type userUsecase struct {
	repo      repo.UserRepo
	uuidGen   helper.UuidGenerator
	paginate  helper.Paginate
	encryptor helper.Encryptor
}

func NewUserUsecase(repo repo.UserRepo, uuidGen helper.UuidGenerator, paginate helper.Paginate, encryptor helper.Encryptor) UserUsecase {
	return &userUsecase{repo, uuidGen, paginate, encryptor}
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

func (u *userUsecase) GetAll(ctx context.Context, reqUser req.GetAllUser) ([]entity.User, int, int ,error) {
	currentPage, offset, limit := u.paginate.GetPaginateData(reqUser.Page, reqUser.Limit)
  datas, totalData, err := u.repo.GetAll(ctx, repo.SearchUser{
		Nama:   reqUser.Search.Nama,
    Alamat: reqUser.Search.Alamat,
    Username: reqUser.Search.Username,
    NoTelp: reqUser.Search.NoTelp,
    Role: reqUser.Search.Role,
    JenisSpvId: reqUser.Search.JenisSpvID,
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, currentPage, 0, err
	} 

  totalPage := u.paginate.GetTotalPages(int(totalData), limit)

	return datas, currentPage, totalPage, err
}

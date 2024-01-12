package user

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"

	repo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	jenisSpvRepo "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm/jenis_spv"
	"github.com/be-sistem-informasi-konveksi/common/message"
	req "github.com/be-sistem-informasi-konveksi/common/request/user"
	res "github.com/be-sistem-informasi-konveksi/common/response/user"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type UserUsecase interface {
	Create(ctx context.Context, reqUser req.Create) error
	Update(ctx context.Context, reqUser req.Update) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, reqUser req.GetAll) ([]res.DataGetAllUserRes, error)
}

type userUsecase struct {
	repo         repo.UserRepo
	jenisSpvRepo jenisSpvRepo.JenisSpvRepo
	ulid         pkg.UlidPkg
	paginate     helper.Paginate
	encryptor    helper.Encryptor
}

func NewUserUsecase(
	repo repo.UserRepo,
	jenisSpvRepo jenisSpvRepo.JenisSpvRepo,
	ulid pkg.UlidPkg,
	paginate helper.Paginate,
	encryptor helper.Encryptor,
) UserUsecase {
	return &userUsecase{repo, jenisSpvRepo, ulid, paginate, encryptor}
}

func (u *userUsecase) Create(ctx context.Context, reqUser req.Create) error {
	id := u.ulid.MakeUlid().String()
	pass, err := u.encryptor.HashPassword(reqUser.Password)
	if err != nil {
		return err
	}
	user := entity.User{
		Base: entity.Base{
			ID: id,
		},
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

func (u *userUsecase) Update(ctx context.Context, reqUser req.Update) error {
	_, err := u.repo.GetById(ctx, reqUser.ID)
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New(message.UserNotFound)
		}
		return err
	}

	user := entity.User{
		Base: entity.Base{
			ID: reqUser.ID,
		},
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

	err = u.repo.Update(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) Delete(ctx context.Context, id string) error {
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

func (u *userUsecase) GetAll(ctx context.Context, reqUser req.GetAll) ([]res.DataGetAllUserRes, error) {
	// currentPage, offset, limit := u.paginate.GetPaginateData(reqUser.Page, reqUser.Limit)
	datas, err := u.repo.GetAll(ctx, repo.SearchUser{
		Nama:       reqUser.Search.Nama,
		Alamat:     reqUser.Search.Alamat,
		Username:   reqUser.Search.Username,
		NoTelp:     reqUser.Search.NoTelp,
		Role:       reqUser.Search.Role,
		JenisSpvId: reqUser.Search.JenisSpvID,
		Limit:      reqUser.Limit,
		Next:       reqUser.Next,
		// Offset:     offset,
	})
	if err != nil {
		return nil, err
	}
	datasRes := make([]res.DataGetAllUserRes, len(datas))
	g := &errgroup.Group{}
	for i, d := range datas {
		i := i
		d := d
		g.Go(func() error {
			jenisSpv, err := u.jenisSpvRepo.GetById(ctx, d.JenisSpvID)
			if err != nil && err.Error() != "record not found" {
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
		return nil, err
	}

	// totalPage := u.paginate.GetTotalPages(int(totalData), limit)

	return datasRes, err
}

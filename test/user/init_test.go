package test_user

import (
	"os"
	"testing"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

var dbt *gorm.DB
var tokens map[string]string
var app = fiber.New()
var idsDefaultUser []string
var idsDefaultSpv []string

func cleanUp() {
	// clean spv with tugas
	dbt.Unscoped().Where("id = ?", idKontak).Delete(&entity.Kontak{})
	dbt.Unscoped().Where("id = ?", idInvoice).Delete(&entity.Invoice{})
	dbt.Unscoped().Where("id = ?", idTugas).Delete(&entity.Tugas{})
	dbt.Unscoped().Where("id IN (?)", idSubTugas).Delete(&entity.SubTugas{})

	idsSpv := make([]string, len(static_data.DefaultSupervisor))
	for i, v := range static_data.DefaultSupervisor {
		idsSpv[i] = v.ID
	}
	dbt.Unscoped().Where("id NOT IN (?)", idsSpv).Delete(&entity.JenisSpv{})

	ids := make([]string, len(static_data.DefaultUsers))
	for i, v := range static_data.DefaultUsers {
		ids[i] = v.ID
	}
	dbt.Unscoped().Where("id NOT IN (?)", ids).Delete(&entity.User{})
}

func TestMain(m *testing.M) {
	test.GetDB()
	dbt = test.DBT

	userH := handler_init.NewUserHandlerInit(dbt, test.Validator, test.UlidPkg, test.Encryptor)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	userRoute := route.NewUserRoute(userH, authMid)

	test.GetTokens(dbt, authMid)
	tokens = test.Tokens
	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	userGroup := v1.Group("/user")
	userGroup.Route("/jenis_spv", userRoute.JenisSpv)
	userGroup.Route("/", userRoute.User)

	idsDefaultUser = make([]string, len(static_data.DefaultUsers))
	for i, v := range static_data.DefaultUsers {
		idsDefaultUser[i] = v.ID
	}
	idsDefaultSpv = make([]string, len(static_data.DefaultSupervisor))
	for i, v := range static_data.DefaultSupervisor {
		idsDefaultSpv[i] = v.ID
	}
	// Run tests
	exitVal := m.Run()
	cleanUp()
	os.Exit(exitVal)
}

func TestEndPointUser(t *testing.T) {
	//? spv
	UserCreateSpv(t)
	UserUpdateSpv(t)
	UserGetAllSpv(t)
	//? user
	UserCreate(t)
	UserUpdate(t)
	UserGetAll(t)
	UserGet(t)
	//? delete
	UserDelete(t)
	UserDeleteSpv(t)
}

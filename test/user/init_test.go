package test_user

import (
	"os"
	"testing"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

var dbt *gorm.DB
var token string
var app = fiber.New()

func TestMain(m *testing.M) {
	dbt = test.GetDB()
	userH := handler_init.NewUserHandlerInit(dbt, test.Validator, test.UlidPkg, test.Encryptor)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	userRoute := route.NewUserRoute(userH, authMid)

	token = test.GetToken(dbt, authMid)

	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	userGroup := v1.Group("/user")
	userGroup.Route("/jenis_spv", userRoute.JenisSpv)
	userGroup.Route("/", userRoute.User)
	// Run tests
	exitVal := m.Run()
	dbt.Unscoped().Where("1 = 1").Delete(&entity.JenisSpv{})
	dbt.Unscoped().Where("1 = 1").Delete(&entity.User{})
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
	UserDeleteSpv(t)
	UserDelete(t)
}

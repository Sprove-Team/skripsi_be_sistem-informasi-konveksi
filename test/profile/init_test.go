package test_profile

import (
	"os"
	"testing"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

var dbt *gorm.DB
var tokens map[string]string
var app = fiber.New()

func rollBackUpdate() {
	if err := dbt.Updates(&static_data.DefaultUsers[0]).Error; err != nil {
		panic(helper.LogsError(err))
		return
	}
	if err := dbt.Updates(&static_data.DefaultUsers[1]).Error; err != nil {
		panic(helper.LogsError(err))
		return
	}
}

func TestMain(m *testing.M) {
	test.GetDB()
	dbt = test.DBT
	profileH := handler_init.NewProfileHandlerInit(dbt, test.Validator, test.UlidPkg, test.Encryptor)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	profilRoute := route.NewProfileRoute(profileH, authMid)

	test.GetTokens(dbt, authMid)
	tokens = test.Tokens
	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	tugasGroup := v1.Group("/profile")
	tugasGroup.Route("/", profilRoute.Profile)
	// Run tests
	exitVal := m.Run()
	rollBackUpdate()
	os.Exit(exitVal)
}

func TestEndPointProfile(t *testing.T) {
	//? profile
	ProfileGet(t)
	ProfileUpdate(t)
}

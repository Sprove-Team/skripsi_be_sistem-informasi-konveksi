//go:build test_exclude

package test_sablon

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
var tokens map[string]string
var app = fiber.New()

func TestMain(m *testing.M) {
	test.GetDB()
	dbt = test.DBT
	sablonH := handler_init.NewSablonHandlerInit(dbt, test.Validator, test.UlidPkg)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	sablonRoute := route.NewSablonRoute(sablonH, authMid)

	test.GetTokens(dbt, authMid)
	tokens = test.Tokens

	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	sablonGroup := v1.Group("/sablon")
	sablonGroup.Route("/", sablonRoute.Sablon)
	// Run tests
	exitVal := m.Run()
	dbt.Unscoped().Where("1 = 1").Delete(&entity.Sablon{})
	os.Exit(exitVal)
}

func TestEndPointSablon(t *testing.T) {
	SablonCreate(t)
	SablonUpdate(t)
	SablonGetAll(t)
	SablonDelete(t)
}

package test_bordir

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
	bordirH := handler_init.NewBordirHandlerInit(dbt, test.Validator, test.UlidPkg)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	bordirRoute := route.NewBordirRoute(bordirH, authMid)

	token = test.GetToken(dbt, authMid)

	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	bordirGroup := v1.Group("/bordir")
	bordirGroup.Route("/", bordirRoute.Bordir)
	// Run tests
	exitVal := m.Run()
	dbt.Unscoped().Where("1 = 1").Delete(&entity.Bordir{})
	os.Exit(exitVal)
}

func TestEndPointBordir(t *testing.T) {
	BordirCreate(t)
	BordirUpdate(t)
	BordirGetAll(t)
	BordirDelete(t)
}

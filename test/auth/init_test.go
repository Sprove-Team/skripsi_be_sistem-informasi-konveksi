package test_auth

import (
	"os"
	"testing"
	"time"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var dbt *gorm.DB
var app = fiber.New()
var refreshTokenUserDoesntExits string
var tokens map[string]string
var jwtPkg pkg.JwtC

func TestMain(m *testing.M) {
	test.GetDB()
	dbt = test.DBT
	jwtPkg = pkg.NewJwt(os.Getenv("JWT_TOKEN"), os.Getenv("JWT_REFTOKEN"))
	refreshTokenUserDoesntExits, _ = jwtPkg.CreateToken(true, &pkg.Claims{
		ID: test.UlidPkg.MakeUlid().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "refresh_token",
		},
	}, time.Now().Add(time.Second*30))

	authH := handler_init.NewAuthHandlerInit(dbt, jwtPkg, test.Validator, test.UlidPkg, test.Encryptor)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	authRoute := route.NewAuthRoute(authH, authMid)

	test.GetTokens(dbt, authMid)
	tokens = test.Tokens
	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	akuntansiGroup := v1.Group("/auth")
	akuntansiGroup.Route("", authRoute.Auth)
	// Run tests
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestEndpointAuth(t *testing.T) {
	AuthLogin(t)
	AuthRefreshToken(t)
	AuthWhoAmI(t)
}

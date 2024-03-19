package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"regexp"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/config"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/app/static_data"
	req_auth "github.com/be-sistem-informasi-konveksi/common/request/auth"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type Response struct {
	Status         string   `json:"status"`
	Code           int      `json:"code"`
	Data           any      `json:"data,omitempty"`
	ErrorsMessages []string `json:"errors_messages,omitempty"`
}

var Validator = pkg.NewValidator()
var UlidPkg = pkg.NewUlidPkg()

func GetToken(dbt *gorm.DB, authMid middleware_auth.AuthMidleware) (token string) {
	app := fiber.New()
	api := app.Group("/api/v1")
	app.Use(recover.New())

	pkgJwt := pkg.NewJwt(os.Getenv("JWT_TOKEN"), os.Getenv("JWT_REFTOKEN"))
	authH := handler_init.NewAuthHandlerInit(dbt, pkgJwt, Validator, UlidPkg, helper.NewEncryptor())
	authRoute := route.NewAuthRoute(authH, authMid)
	authGroup := api.Group("/auth")
	{
		authGroup.Route("", authRoute.Auth)
	}
	_, body, err := GetJsonTestRequestResponse(app, "POST", "/api/v1/auth/login", req_auth.Login{
		Username: static_data.DefaultUserDirektur.Username,
		Password: static_data.PlainPassword,
	}, nil)
	if err != nil {
		helper.LogsError(err)
		return
	}

	var dataAuth map[string]interface{}
	if err := mapstructure.Decode(body.Data, &dataAuth); err != nil {
		helper.LogsError(err)
		return
	}
	token = dataAuth["token"].(string)
	return
}

func GetJsonTestRequestResponse(app *fiber.App, method string, url string, reqBody any, token *string) (code int, respBody Response, err error) {
	bodyJson := []byte("")
	if reqBody != nil {
		bodyJson, _ = json.Marshal(reqBody)
	}

	req := httptest.NewRequest(method, url, bytes.NewReader(bodyJson))
	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != nil {
		req.Header.Set("Authorization", "Bearer "+*token)
	}

	resp, err := app.Test(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	code = resp.StatusCode
	// If error we're done
	if err != nil {
		return
	}
	// If no body content, we're done
	if resp.ContentLength == 0 {
		return
	}
	bodyData := make([]byte, resp.ContentLength)
	_, _ = resp.Body.Read(bodyData)
	err = json.Unmarshal(bodyData, &respBody)
	return
}

func loadEnv() {
	const projectDirName = "be-sistem-informasi-konveksi" // change to relevant project name
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GetDB() *gorm.DB {
	loadEnv()
	dbGormConf := config.DBGorm{
		DB_Username: "root",
		DB_Password: os.Getenv("DB_PASSWORD_ROOT"),
		DB_Name:     os.Getenv("DB_NAME_TEST"),
		DB_HOST:     os.Getenv("DB_HOST_TEST"),
		DB_Port:     os.Getenv("DB_PORT_TEST"),
	}

	ulidPkg := pkg.NewUlidPkg()
	dbGorm := dbGormConf.InitDBGorm(ulidPkg)
	return dbGorm
}

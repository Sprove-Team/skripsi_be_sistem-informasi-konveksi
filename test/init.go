package test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"regexp"
	"time"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	"github.com/be-sistem-informasi-konveksi/app/config"
	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Response struct {
	Status         string   `json:"status"`
	Code           int      `json:"code"`
	Data           any      `json:"data,omitempty"`
	ErrorsMessages []string `json:"errors_messages,omitempty"`
}

var Validator = pkg.NewValidator()
var UlidPkg = pkg.NewUlidPkg()
var Excelize = pkg.NewExcelizePkg()
var Encryptor = helper.NewEncryptor()
var Tokens = make(map[string]string)
var DBT *gorm.DB

func GetTokens(dbt *gorm.DB, authMid middleware_auth.AuthMidleware) {
	if len(Tokens) == len(static_data.DefaultUsers) {
		return
	}
	pkgJwt := pkg.NewJwt(os.Getenv("JWT_TOKEN"), os.Getenv("JWT_REFTOKEN"))

	for _, v := range static_data.DefaultUsers {
		claims := &pkg.Claims{
			Nama:     v.Nama,
			ID:       v.ID,
			Username: v.Username,
			Role:     v.Role,
		}
		claims.Subject = "access_token"
		token, err := pkgJwt.CreateToken(false, claims, time.Now().Add(time.Hour*8))

		if err != nil {
			helper.LogsError(err)
			os.Exit(1)
		} else {
			Tokens[v.Role] = token
		}
	}
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

func GetAttachTestRequestResponse(app *fiber.App, method string, url string, reqBody any, fileName string, token *string) (code int, err error) {
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

	resBuf := bytes.Buffer{}
	_, err = io.Copy(&resBuf, resp.Body)
	if err != nil {
		return
	}
	// save
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = resBuf.WriteTo(file)
	if err != nil {
		return
	}
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

func GetDB() {
	if DBT != nil {
		return
	}
	loadEnv()
	dbGormConf := config.DBGorm{
		DB_Username: "root",
		DB_Password: os.Getenv("DB_PASSWORD_ROOT"),
		DB_Name:     os.Getenv("DB_NAME_TEST"),
		DB_HOST:     os.Getenv("DB_HOST_TEST"),
		DB_Port:     os.Getenv("DB_PORT_TEST"),
		LogLevel:    logger.Silent,
	}
	if os.Getenv("ENVIRONMENT") == "DEVELOPMENT" {
		dbGormConf.DB_HOST = "localhost"
	}

	ulidPkg := pkg.NewUlidPkg()
	DBT = dbGormConf.InitDBGorm(ulidPkg)
}

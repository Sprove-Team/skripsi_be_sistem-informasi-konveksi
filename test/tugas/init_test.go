package test_tugas

import (
	"os"
	"testing"
	"time"

	middleware_auth "github.com/be-sistem-informasi-konveksi/api/middleware/auth"
	repo_user "github.com/be-sistem-informasi-konveksi/api/repository/user/mysql/gorm"
	"github.com/be-sistem-informasi-konveksi/app/handler_init"
	"github.com/be-sistem-informasi-konveksi/app/route"
	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/test"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

var dbt *gorm.DB
var tokens map[string]string
var app = fiber.New()

var invoiceId = []string{"01HTM8HJJA7KF5JZE0FDRDG0Y0", "01HTM8HJJA7KF5JZE0FFHXC27T"}
var kontakId = "01HTM8KDD5YAEEJVARYWNK21NK"
var spvId = "01HTM8HJD09FDNV2TTXB9TGHP6"
var userId = []string{"01HTM8HJD09FDNV2TTXC0GYJQE", "01HTM8HJEP97T9WSVH7SCZWK57"}

func setUpData() {
	kontak := entity.Kontak{
		Base: entity.Base{
			ID: kontakId,
		},
		Nama:       "John Doe",
		NoTelp:     "+6281234567890",
		Alamat:     "123 Main Street",
		Email:      "john.doe@example.com",
		Keterangan: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	}

	if err := dbt.Create(&kontak).Error; err != nil {
		helper.LogsError(err)
		return
	}
	userSpv := entity.JenisSpv{
		Base: entity.Base{
			ID: spvId,
		},
		Nama: "bordir",
		Users: []entity.User{
			{
				Base: entity.Base{
					ID: userId[0],
				},
				Nama:     "user spv bordir",
				Role:     entity.RolesById[5],
				Username: "spv_bordir1",
				Password: func() string {
					pass, _ := test.Encryptor.HashPassword("spvbordir1")
					return pass
				}(),
				NoTelp: "+6289123123124",
				Alamat: "test",
			},
			{
				Base: entity.Base{
					ID: userId[1],
				},
				Nama:     "user spv bordir2",
				Role:     entity.RolesById[5],
				Username: "spv_bordir2",
				Password: func() string {
					pass, _ := test.Encryptor.HashPassword("spvbordir2")
					return pass
				}(),
				NoTelp: "+6289123123125",
				Alamat: "test2",
			},
		},
	}

	if err := dbt.Create(&userSpv).Error; err != nil {
		helper.LogsError(err)
		return
	}

	tt := time.Now()
	invoices := []entity.Invoice{
		{
			Base: entity.Base{
				ID: invoiceId[0],
			},
			UserID:         static_data.DefaultUsers[1].ID,
			KontakID:       kontak.ID,
			NomorReferensi: "001",
			TotalQty:       10,
			TotalHarga:     10000,
			Keterangan:     "test tugas",
			TanggalKirim:   &tt,
		},
		{
			Base: entity.Base{
				ID: invoiceId[1],
			},
			UserID:         static_data.DefaultUsers[1].ID,
			KontakID:       kontak.ID,
			NomorReferensi: "002",
			TotalQty:       15,
			TotalHarga:     15000,
			Keterangan:     "test tugas2",
			TanggalKirim:   &tt,
		},
	}
	if err := dbt.Create(&invoices).Error; err != nil {
		helper.LogsError(err)
		return
	}
}

func cleanUp() {
	if err := dbt.Unscoped().Where("1 = 1").Delete(&entity.Tugas{}).Error; err != nil {
		helper.LogsError(err)
		return
	}
	if err := dbt.Unscoped().Where("id = ?", spvId).Delete(&entity.JenisSpv{}).Error; err != nil {
		helper.LogsError(err)
		return
	}
	if err := dbt.Unscoped().Where("id IN (?)", userId).Delete(&entity.User{}).Error; err != nil {
		helper.LogsError(err)
		return
	}
	if err := dbt.Unscoped().Where("id = ?", kontakId).Delete(&entity.Kontak{}).Error; err != nil {
		helper.LogsError(err)
		return
	}
	if err := dbt.Unscoped().Where("id IN (?)", invoiceId).Delete(&entity.Invoice{}).Error; err != nil {
		helper.LogsError(err)
		return
	}
}
func TestMain(m *testing.M) {
	test.GetDB()
	dbt = test.DBT
	setUpData()
	tugasH := handler_init.NewTugasHandlerInit(dbt, test.Validator, test.UlidPkg)

	userRepo := repo_user.NewUserRepo(dbt)
	authMid := middleware_auth.NewAuthMiddleware(userRepo)
	tugasRoute := route.NewTugasRoute(tugasH, authMid)

	test.GetTokens(dbt, authMid)
	tokens = test.Tokens
	// app
	app.Use(recover.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	tugasGroup := v1.Group("/tugas")
	tugasGroup.Route("/sub_tugas", tugasRoute.SubTugas)
	tugasGroup.Route("/", tugasRoute.Tugas)
	// Run tests
	exitVal := m.Run()
	cleanUp()
	os.Exit(exitVal)
}

func TestEndPointTugas(t *testing.T) {
	//? tugas
	TugasCreate(t)
	TugasUpdate(t)
	TugasGetAll(t)
	TugasGetByInvoiceId(t)
	TugasGet(t)
	TugasDelete(t)
}

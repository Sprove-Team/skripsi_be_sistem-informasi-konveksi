package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/be-sistem-informasi-konveksi/app/static_data"
	"github.com/be-sistem-informasi-konveksi/entity"
	"github.com/be-sistem-informasi-konveksi/helper"
	"github.com/be-sistem-informasi-konveksi/pkg"
)

type DBGorm struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Host     string
	DB_Name     string
}

func (dbgc *DBGorm) InitDBGorm(ulid pkg.UlidPkg) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbgc.DB_Username,
		dbgc.DB_Password,
		dbgc.DB_Host,
		dbgc.DB_Port,
		dbgc.DB_Name)

	logLevel := logger.Silent

	if os.Getenv("PRODUCTION") == "" {
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:         newLogger,
		TranslateError: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}
	// produk
	autoMigrateEntities(db, &entity.KategoriProduk{}, &entity.Produk{}, &entity.HargaDetailProduk{})

	// bordir
	autoMigrateEntities(db, &entity.Bordir{}, &entity.Sablon{})

	// user & jenis spv
	autoMigrateEntities(db, &entity.JenisSpv{}, &entity.User{})

	// akuntansi
	autoMigrateEntities(db, &entity.KelompokAkun{}, &entity.Akun{}, &entity.Transaksi{}, &entity.AyatJurnal{})
	// default value for akuntansi
	klmpkData := static_data.KelompokAkuns(ulid)
	akun := static_data.Akuns(klmpkData, ulid)

	err = addDefultValues(db, klmpkData, akun)

	if err != nil {
		if err.Error() != "duplicated key not allowed" {
			helper.LogsError(err)
			os.Exit(1)
		}
	}

	// invoice
	autoMigrateEntities(db, &entity.StatusProduksi{}, &entity.Invoice{}, &entity.DetailInvoice{})

	return db
}

func autoMigrateEntities(db *gorm.DB, entities ...interface{}) {
	for _, entity := range entities {
		if err := db.AutoMigrate(entity); err != nil {
			panic(err)
		}
	}
}

func addDefultValues(db *gorm.DB, values ...interface{}) error {
	// fmt.Println(golAkun)
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, v := range values {
			if err := tx.Create(v).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		helper.LogsError(err)
		return err
	}
	return nil
}

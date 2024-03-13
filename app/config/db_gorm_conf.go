package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
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
	DB_HOST     string
	DB_Port     string
	DB_Name     string
}

func (dbgc *DBGorm) InitDBGorm(ulid pkg.UlidPkg) *gorm.DB {
	host := "localhost"
	logLevel := logger.Info
	if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		host = dbgc.DB_HOST
		logLevel = logger.Silent
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbgc.DB_Username,
		dbgc.DB_Password,
		host,
		dbgc.DB_Port,
		dbgc.DB_Name)

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

	g := &errgroup.Group{}
	g.SetLimit(10)
	// produk
	g.Go(func() error {
		autoMigrateEntities(db, &entity.KategoriProduk{}, &entity.Produk{}, &entity.HargaDetailProduk{})
		return nil
	})

	// bordir & sablon
	g.Go(func() error {
		autoMigrateEntities(db, &entity.Bordir{}, &entity.Sablon{})
		return nil
	})

	// user & jenis spv
	g.Go(func() error {
		autoMigrateEntities(db, &entity.JenisSpv{}, &entity.User{})
		return nil
	})

	//  akuntansi, invoice & tugas
	g.Go(func() error {
		// akuntansi
		autoMigrateEntities(db, &entity.KelompokAkun{}, &entity.Akun{}, &entity.Transaksi{}, &entity.AyatJurnal{})
		autoMigrateEntities(db, &entity.Kontak{}, &entity.HutangPiutang{}, &entity.DataBayarHutangPiutang{})
		// invoice
		autoMigrateEntities(db, &entity.Invoice{}, &entity.DetailInvoice{}, &entity.DataBayarInvoice{})
		// tugas
		autoMigrateEntities(db, &entity.Tugas{}, &entity.SubTugas{})
		return nil
	})

	if err := g.Wait(); err != nil {
		helper.LogsError(err)
		os.Exit(1)
	}

	// add default value
	// default value for akuntansi
	g.Go(func() error {
		klmpkData := static_data.KelompokAkuns(ulid)
		akun := static_data.Akuns(klmpkData, ulid)

		err = addDefultValues(db, klmpkData, akun)

		if err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		if err.Error() != "duplicated key not allowed" {
			helper.LogsError(err)
			os.Exit(1)
		}
	}

	return db
}

func autoMigrateEntities(db *gorm.DB, entities ...interface{}) {
	for _, entity := range entities {
		if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
			if db.Migrator().HasTable(entity) {
				continue
			}
		}
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

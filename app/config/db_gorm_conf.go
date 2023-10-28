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

	"github.com/be-sistem-informasi-konveksi/entity"
)

type DBGormConf struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Host     string
	DB_Name     string
}

func (dbgc *DBGormConf) InitDBGormConf() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbgc.DB_Username,
		dbgc.DB_Password,
		dbgc.DB_Host,
		dbgc.DB_Port,
		dbgc.DB_Name)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
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
		log.Fatal(err)
		os.Exit(1)
	}
	// produk
	autoMigrateEntities(db, &entity.Produk{}, &entity.HargaDetailProduk{}, &entity.KategoriProduk{})

	// user & jenis spv
	autoMigrateEntities(db,&entity.JenisSpv{}, &entity.User{})
	
	return db
}

func autoMigrateEntities(db *gorm.DB, entities ...interface{}) {
    for _, entity := range entities {
        if err := db.AutoMigrate(entity); err != nil {
            panic(err)
        }
    }
}


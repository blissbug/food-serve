package db

import (
	"fmt"
	"time"

	"food-serve.com/pkg/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/labstack/gommon/log"
)

func NewDBStorage(dsn string) (db *gorm.DB, err error) {
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Errorf("DB Couldnt connect??")
		return nil, err
	}

	sqlDb, err := db.DB()

	if err != nil {
		log.Errorf("DB Couldnt connect??")
		return nil, err
	}

	err = sqlDb.Ping()

	if err != nil {
		log.Errorf("DB Couldnt connect??")
		return nil, err
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)

	fmt.Println("DB Connected yayy!!")

	//TODO: make a separate function for this
	err = db.AutoMigrate(&types.User{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("tables created")
	return db, nil
}

package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/labstack/gommon/log"
)

func NewDBStorage(dsn string)(db *gorm.DB, err error){
  db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err != nil{
    log.Errorf("DB Couldnt connect??")
    return nil, err
  }
  fmt.Println("DB Connected yayy!!")
  return db, nil
}
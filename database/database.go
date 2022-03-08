package database

import (
	"os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func Connect() error {

	dsn := os.Getenv("DSN") //"root:1q1q1q!Q@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	db.Exec("SET sql_mode = ''")

	db.Logger = logger.Default.LogMode(logger.Error)

	Db = db

	return nil

}

func AutoMigrate(i ...interface{}) error {

	err := Db.AutoMigrate(i...)

	return err

}

package dao

import (
	"LarsWebV0/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	logger "github.com/sirupsen/logrus"
	"time"
)

var db *gorm.DB

func SetupDB() {
	var err error
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/larswebv0?charset=utf8mb4&parseTime=true&loc=Local", config.Username, config.Password, config.Hostname, config.Port)
	db, err = gorm.Open("mysql", dsn)

	if err != nil {
		logger.Errorf("models.Setup err: %v", err)
		return
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxIdleTime(time.Hour)
}

func CloseDB() {
	defer db.Close()
}

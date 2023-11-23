package utils

import (
	"fmt"
	"sync"
	"time"
	"user_system/config"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

// 连接数据库
func openDB() {
	mysqlConf := config.GetGlobalConf().DbConfig
	// dsn := "username:password@tcp(127.0.0.1:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local"
	connArgs := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConf.User, mysqlConf.Password, mysqlConf.Host, mysqlConf.Port, mysqlConf.Dbname)

	var err error
	db, err = gorm.Open(mysql.Open(connArgs), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database")
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Panic("fetch dbconnection err:" + err.Error())
		return
	}

	sqlDB.SetMaxIdleConns(mysqlConf.Max_idle_conn) //设置最大空闲连接
	sqlDB.SetMaxOpenConns(mysqlConf.Max_open_conn) //设置最大打开的连接
	sqlDB.SetConnMaxLifetime(time.Duration(mysqlConf.Max_idle_time * int64(time.Second)))
	log.Info("connect mysql success")
}

// 获得数据库连接
func GetDB() *gorm.DB {
	dbOnce.Do(openDB)
	return db
}

//关闭数据库连接

func CloseDB() {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error("[CloseDB] fetch db connection error:", err.Error())
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Error("[CloseDB] close db connection error", err.Error())
		return
	}

	log.Info("[CloseDB] Database connection cloded")
}

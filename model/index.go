package model

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

// DB 单例模式数据库引擎
var DB *gorm.DB

// Database 初始化数据库引擎
func Database(connString string) {
	db, err := gorm.Open("mysql", connString)
	db.LogMode(func(mode string) bool {
		if mode == "debug" {
			return true
		}
		return false
	}(os.Getenv("DB_MODE")))
	db.SingularTable(true)
	if err != nil {
		log.Info(err)
	}
	//设置连接池
	//空闲
	db.DB().SetMaxIdleConns(50)
	//打开
	db.DB().SetMaxOpenConns(100)
	//超时
	db.DB().SetConnMaxLifetime(time.Second * 30)
	DB = db

	migration()
}

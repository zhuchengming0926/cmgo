package mysql_pool

/**
 * @Author: chengming1
 * @Date: 2023/2/1 上午11:49
 * @Desc:
 */

import (
	"fmt"
	"log"
	"time"

	"cmgo/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go.uber.org/zap"
)

// init mysql
func InitJinzhuMysql() *gorm.DB {
	mysqlDB, err := mysqlConn()
	if err != nil {
		log.Println("mysql conn error: ", zap.Error(err))
	}
	log.Println("mysql connect success!")
	return mysqlDB
}

// mysql connect
func mysqlConn() (*gorm.DB, error) {
	mysqlConfig := conf.GetMysqlConfig()

	var mysqldb *gorm.DB
	if mysqlConfig.Driver == "sqlite3" {
		db, err := gorm.Open(mysqlConfig.Driver, mysqlConfig.Database)
		if err != nil {
			return nil, err
		}
		mysqldb = db
	} else {
		connArgs := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
			mysqlConfig.User, mysqlConfig.Passwd, mysqlConfig.Host,
			mysqlConfig.Port, mysqlConfig.Database, beijingTimeLocaltion)
		db, err := gorm.Open(mysqlConfig.Driver, connArgs)
		if err != nil {
			log.Fatalf("db open error: %s", err.Error())
			return nil, err
		}
		db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")

		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(1000)
		db.DB().SetConnMaxLifetime(600 * time.Second)
		mysqldb = db
	}
	return mysqldb, nil
}

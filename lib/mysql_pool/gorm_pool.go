package mysql_pool

/**
* @Author: chengming1
* @Date: 2023/1/30 下午5:09
* @Desc: 参考：https://www.cnblogs.com/davis12/p/16357728.html
* mysql-gorm 建立连接池
 */

//其实gorm的连接池设置，底层还是用的database/sql的设置连接池的方法，无非就是加一层gorm自身的一些设置。
//以下示例为gorm v2版本，v1版本通过github.com/jinzhu/gorm，如mysql的驱动导入需要加_。

import (
	"cmgo/conf"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	beijingTimeLocaltion = "Asia%2fShanghai"
)

var gorm_db *gorm.DB

func InitGormPool() {
	var dbURi string
	var dialector gorm.Dialector
	mysqlConfig := conf.GetMysqlConfig()
	if mysqlConfig.Driver == "mysql" {
		dbURi = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=%s",
			mysqlConfig.User,
			mysqlConfig.Passwd,
			mysqlConfig.Host,
			mysqlConfig.Port,
			mysqlConfig.Database,
			beijingTimeLocaltion)
		dialector = mysql.New(mysql.Config{
			DSN:                       dbURi, // data source name
			DefaultStringSize:         256,   // default size for string fields
			DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
		})
	} else if mysqlConfig.Driver == "sqlite3" {
		dbURi = fmt.Sprintf("test.db")
		dialector = sqlite.Open("test.db")
	}

	conn, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatal("connect db server failed.")
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(600 * time.Second)

	gorm_db = conn
}

// open api
func GetDB() *gorm.DB {
	sqlDB, err := gorm_db.DB()
	if err != nil {
		log.Println("connect db server failed.")
	}
	if err = sqlDB.Ping(); err != nil {
		sqlDB.Close()
	}

	return gorm_db
}

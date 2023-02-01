package mysql_pool

/**
 * @Author: chengming1
 * @Date: 2023/1/30 下午4:59
 * @Desc: 参考：https://www.cnblogs.com/davis12/p/16357728.html
 * mysql-通过sql建立连接池
 */

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var sql_DB *sql.DB

// init方法只要被import就会执行
func InitSqlPool() {
	sql_DB, _ = sql.Open("mysql", "root:root@tcp(10.41.11.211:3306)/ad_feature_platform?charset=utf8&parseTime=True&loc=Local") // 使用本地时间，即东八区，北京时间
	// set pool params
	sql_DB.SetMaxOpenConns(1000)
	sql_DB.SetMaxIdleConns(10)
	sql_DB.SetConnMaxLifetime(time.Minute * 60) // mysql default conn timeout=8h, should < mysql_timeout
	err := sql_DB.Ping()
	if err != nil {
		log.Fatalf("database init failed, err: ", err)
	}
	log.Println("mysql conn pool has initiated.")
}

func GetSqlDb() *sql.DB {
	return sql_DB
}

func GetFeature() {
	rows, err := sql_DB.Query("SELECT `name` FROM features LIMIT 10")
	fmt.Println(err)
	for rows.Next() {
		var feature_name string
		err = rows.Scan(&feature_name)
		fmt.Println(feature_name)
	}
}

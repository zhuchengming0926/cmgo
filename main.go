package main

import (
	"cmgo/lib/mysql_pool"
	"fmt"
)

func main() {
	fmt.Println("人兴财旺")
	mysql_pool.InitSqlPool()
	mysql_pool.GetFeature()
}

package gomock

import (
	"cmgo/conf"
	"cmgo/lib/mysql_pool"
	gomonkey "github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"
	"testing"
)

/**
 * @Author: chengming1
 * @Date: 2023/1/31 下午2:42
 * @Desc: 使用gomock执行单元测试
 */

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func TestTtt(t *testing.T) {
	// 函数打桩，改变conf.GetMysqlConfig这个函数返回值为使用打桩中设置的值
	convey.Convey("conf.GetMysqlConfig", t, func() {

		var patches = gomonkey.ApplyFunc(conf.GetMysqlConfig, func() *conf.MysqlConfig {
			return &conf.MysqlConfig{
				Driver: "sqlite3",
			}
		})
		defer patches.Reset()

		mysql_pool.InitJinzhuMysql()
		db := mysql_pool.GetDB()

		db.AutoMigrate(&Product{})
		db.Create(&Product{Code: "D42", Price: 100})

		var product Product
		db.First(&product, 1) // 根据整形主键查找
		t.Log(product)
		db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
		t.Log(product)

		db.Model(&product).Update("price", 200)
		t.Log(product)

		db.Model(&product).Updates(map[string]interface{}{"price": 200, "code": "F42"})
		t.Log(product)

		// 清理数据
		db.Exec("delete from product where 1;")
		db.Exec("delete from sqlite_sequence where name='product'")
	})
}

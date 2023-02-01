package conf

/**
 * @Author: chengming1
 * @Date: 2023/1/31 下午3:05
 * @Desc:
 */

type MysqlConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Passwd   string `json:"passwd"`
	Database string `json:"database"`
}

func GetMysqlConfig() *MysqlConfig {
	return &MysqlConfig{
		Driver:   "mysql",
		Host:     "10.41.11.211",
		Port:     "3306",
		User:     "root",
		Passwd:   "root",
		Database: "ad_feature_platform",
	}
}

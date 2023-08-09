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

type Consul struct {
	Token       string `json:"token"`
	Port        string `json:"port"`
	ServiceName string `json:"service_name"`
}

type ConsulKVConfig struct {
	ConsulAddress      string `json:"consul_address"`
	UseLocalAgent      bool   `json:"use_local_agent"`
	FeatureServicePath string `json:"feature_service_path"`
	Port               int    `json:"feature_service_port"`
	CheckConsulTimeout int    `json:"check_consul_timeout"`
}

func GetConsulKVConfig() *ConsulKVConfig {
	return &ConsulKVConfig{}
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

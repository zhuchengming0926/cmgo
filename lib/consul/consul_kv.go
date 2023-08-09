package consul

import (
	"fmt"
	"log"
	"time"

	"cmgo/conf"
	"cmgo/lib/logger"

	kv "github.com/lockp111/consul-kv"
	"go.uber.org/zap"
)

const (
	DefaultConsulValueThrift     = "thrift"
	DefaultConsulValueDc         = "default"
	DefaultConsulValueStatus     = "available"
	DefaultConsulValueIsEdgeNode = false
)

type ConsulKVClient struct {
	ConsulKVClient *kv.Config
}

func NewConsulKVClient() ConsulKVClient {
	consulAddress := conf.GetConsulKVConfig().ConsulAddress
	cli := kv.NewConfig(kv.WithAddress(consulAddress))
	err := cli.Init()
	if err != nil {
		log.Fatalln(err)
	}

	return ConsulKVClient{
		ConsulKVClient: cli,
	}
}

// ads-core/services/feature_service/nodes/thrift:10.133.136.146:7777
func MakeConsulKey(path, servicename, ip string, port int) string {
	return fmt.Sprintf("%s/%s/nodes/thrift:%s:%d", path, servicename, ip, port)
}

type ConsulValue struct {
	PartitionList              []int    `json:"PartitionList"`
	Idc                        string   `json:"Idc"`
	FollowerShardList          []int    `json:"FollowerShardList"`
	ServiceName                string   `json:"ServiceName"`
	UpdateTime                 int64    `json:"UpdateTime"`
	Host                       string   `json:"Host"`
	Protocol                   string   `json:"Protocol"`
	Port                       int      `json:"Port"`
	Weight                     int      `json:"Weight"`
	HostLong                   int      `json:"HostLong"`
	ShardList                  []int    `json:"ShardList"`
	AvailableShardList         []int    `json:"AvailableShardList"`
	Dc                         string   `json:"Dc"`
	IsEdgeNode                 bool     `json:"IsEdgeNode"`
	Status                     string   `json:"Status"`
	FollowerAvailableShardList []int    `json:"FollowerAvailableShardList"`
	OtherSettings              struct{} `json:"OtherSettings"`
}

func MakeConsulValue(serviceName, ip string, weight, port int) ConsulValue {
	return ConsulValue{
		PartitionList:              []int{},
		Idc:                        "",
		FollowerShardList:          []int{},
		ServiceName:                serviceName,
		UpdateTime:                 time.Now().UnixNano() / 1e6,
		Host:                       ip,
		Protocol:                   DefaultConsulValueThrift,
		Port:                       port,
		Weight:                     weight,
		HostLong:                   0,
		ShardList:                  []int{},
		AvailableShardList:         []int{},
		Dc:                         DefaultConsulValueDc,
		IsEdgeNode:                 DefaultConsulValueIsEdgeNode,
		Status:                     DefaultConsulValueStatus,
		FollowerAvailableShardList: []int{},
		OtherSettings:              struct{}{},
	}
}

// 从Consul KV根据key读取value
func (kv *ConsulKVClient) Get(key string) (string, error) {
	ret := kv.ConsulKVClient.Get(key)
	if ret.Err() != nil {
		logger.Error("consul kv client get err", zap.String("key", key), zap.Error(ret.Err()))
		return "", ret.Err()
	}

	return ret.String(), nil
}

// 创建/更新key；无则新建（路径也会新建），有则更新
func (kv *ConsulKVClient) Put(key string, val interface{}) error {
	err := kv.ConsulKVClient.Put(key, val)
	if err != nil {
		logger.Error("consul kv client put err",
			zap.String("key", key), zap.Any("val", val), zap.Error(err))
		return err
	}

	return nil
}

// 删除key或者路径，没有会报错
func (kv *ConsulKVClient) Delete(key string) error {
	err := kv.ConsulKVClient.Delete(key)
	if err != nil {
		if err.Error() == "key not found" {
			return nil
		}

		logger.Error("consul kv client delete err", zap.String("key", key), zap.Error(err))
		return err
	}

	return nil
}

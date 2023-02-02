package kafka

import (
	"errors"
	"log"
	"strings"

	"cmgo/lib/consistent"

	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
)

type KafkaProducer struct {
	Version    string   `toml:"version"`
	Brokers    []string `toml:"brokers"`
	Topic      string   `toml:"topic"`
	User       string   `toml:"user"`
	Password   string   `toml:"password"`
	EnableAuth bool     `toml:"enable_auth"`
}

type CustomizeProducer struct {
	sarama.Partitioner //组合，下边重载实现Partition方法

	Consistent *consistent.Consistent
}

type Message struct {
	HashID string `json:"hash_id"`
}

func (cp *CustomizeProducer) Partition(message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	// 一致性哈希
	value, err := message.Value.Encode()
	if err != nil {
		return 0, err
	}

	valueStr := strings.TrimSpace(string(value))
	if valueStr == "" {
		return 0, errors.New("msg is null")
	}

	msg := &Message{}
	err = jsoniter.UnmarshalFromString(valueStr, msg)
	if err != nil {
		return 0, err
	}
	partition, err := cp.Consistent.Get(msg.HashID, int(numPartitions))
	return int32(partition), err
}

func (cp *CustomizeProducer) RequiresConsistency() bool {
	return true
}

func NewKafkaProducer(producerCfg KafkaProducer) sarama.SyncProducer {
	//初始化配置
	kafkaVersion, err := sarama.ParseKafkaVersion(producerCfg.Version)
	if err != nil {
		log.Fatalf("[kafka init]sarama.ParseKafkaVersion failed, err:%v", err)
	}

	config := sarama.NewConfig()
	config.Version = kafkaVersion
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = func(topic string) sarama.Partitioner {
		return &CustomizeProducer{
			//Partitioner: sarama.NewManualPartitioner(topic), 感觉这里不用赋值
			Consistent: consistent.NewConsistent(10),
		}
	}

	//添加sasl认证配置
	config.Net.SASL.Enable = producerCfg.EnableAuth
	config.Net.SASL.User = producerCfg.User
	config.Net.SASL.Password = producerCfg.Password

	//生产者
	producer, err := sarama.NewSyncProducer(producerCfg.Brokers, config)
	if err != nil {
		log.Fatalf("Error creating kafka producer client: %v", err)
		return nil
	}
	return producer
}

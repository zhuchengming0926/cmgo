package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

type Kafka struct {
	Version    string   `toml:"version"`
	Brokers    []string `toml:"brokers"`
	Topics     []string `toml:"topics"`
	Group      string   `toml:"group"`
	User       string   `toml:"user"`
	Password   string   `toml:"password"`
	EnableAuth bool     `toml:"enable_auth"`
}

//1、这里初始化kafka消费者组
//2、提供一个consumerGroup全局变量，初始化这个变量
func NewKafkaConsumerGroup(kafkaCfg Kafka) sarama.ConsumerGroup {
	log.Println("Starting a new Sarama consumer group")

	kafkaVersion, err := sarama.ParseKafkaVersion(kafkaCfg.Version)
	if err != nil {
		log.Fatalf("[kafka init]sarama.ParseKafkaVersion failed, err:%v", err)
	}
	config := sarama.NewConfig()
	config.Version = kafkaVersion
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	//添加sasl认证配置
	config.Net.SASL.Enable = kafkaCfg.EnableAuth
	config.Net.SASL.User = kafkaCfg.User
	config.Net.SASL.Password = kafkaCfg.Password

	client, err := sarama.NewConsumerGroup(kafkaCfg.Brokers, kafkaCfg.Group, config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}

	return client
}

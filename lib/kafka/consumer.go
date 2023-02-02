package kafka

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

// kafka消费者，消费指定分区
// machineCnt是机器数量
// machineIndex机器编号
func NewKafkaPartionConsumerMap(kafkaCfg Kafka, machineCnt, machineIndex int32) (
	sarama.Consumer, map[int32]sarama.PartitionConsumer) {
	log.Println("Starting a new Sarama consumer")

	kafkaVersion, err := sarama.ParseKafkaVersion(kafkaCfg.Version)
	if err != nil {
		log.Fatalf("[kafka init]sarama.ParseKafkaVersion failed, err:%v", err)
	}
	config := sarama.NewConfig()
	config.Version = kafkaVersion
	config.Consumer.Return.Errors = true

	//添加sasl认证配置
	config.Net.SASL.Enable = kafkaCfg.EnableAuth
	config.Net.SASL.User = kafkaCfg.User
	config.Net.SASL.Password = kafkaCfg.Password

	consumer, err := sarama.NewConsumer(kafkaCfg.Brokers, config)
	if err != nil {
		log.Fatalf("Error creating consumer client: %v", err)
	}

	allPartitions, err := consumer.Partitions(kafkaCfg.Topics[0]) // 根据topic取到所有的分区
	if err != nil {
		log.Fatalf("Error creating consumer, get Partitions fail, err: %v", err)
	}
	log.Println(fmt.Sprintf("%s partition total is %d", kafkaCfg.Topics[0], len(allPartitions)))

	partitionConsumerMap := make(map[int32]sarama.PartitionConsumer)
	for _, partitionNum := range allPartitions {
		if partitionNum%machineCnt != machineIndex {
			continue
		}
		partitionConsumer, err := consumer.ConsumePartition(kafkaCfg.Topics[0], partitionNum, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error creating partitionConsumer, err: %v, partition:%d", err, partitionNum)
		}
		partitionConsumerMap[partitionNum] = partitionConsumer
	}

	return consumer, partitionConsumerMap
}

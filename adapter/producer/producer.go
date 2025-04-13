package producer

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/ether-echo/telegram-api-service/internal/model"
	"github.com/ether-echo/telegram-api-service/pkg/logger"
	"time"
)

var (
	log = logger.Logger().Named("producer").Sugar()
)

type KafkaProducer struct {
	AsyncProducer sarama.AsyncProducer
}

func NewKafkaProducer(brokerList []string) *KafkaProducer {

	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForLocal     // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}

	return &KafkaProducer{AsyncProducer: producer}
}

// Отправка сообщения в Kafka
func (kp *KafkaProducer) SendMessageToKafka(message model.MessageRequest) error {

	value, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	var topic string

	switch message.Message {
	case "/start":
		topic = "start"
	case "/support":
		topic = "support"
	default:
		topic = "message"
	}
	// Создаем сообщение Kafka
	kafkaMessage := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", message.ChatId)),
		Value: sarama.ByteEncoder(value),
	}

	// Отправляем сообщение в Kafka
	kp.AsyncProducer.Input() <- kafkaMessage

	// Ждем подтверждения или ошибки
	select {
	case <-kp.AsyncProducer.Successes():
		log.Infof("%v", message)
		return nil
	case err := <-kp.AsyncProducer.Errors():
		return fmt.Errorf("sarama producer failed: %v", err.Err)
	}
}

func (kp *KafkaProducer) Close() {
	err := kp.AsyncProducer.Close()
	if err != nil {
		log.Errorf("Failed to close Sarama producer: %v", err)
	}
}

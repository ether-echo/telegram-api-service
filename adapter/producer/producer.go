package producer

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/ether-echo/telegram-api-service/internal/model"
	"github.com/ether-echo/telegram-api-service/pkg/logger"
	"sync"
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
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}

	go func() {
		for {
			select {
			case successMsg := <-producer.Successes():
				log.Infof("Message delivered to topic %s, partition %d, offset %d", successMsg.Topic, successMsg.Partition, successMsg.Offset)
			case errMsg := <-producer.Errors():
				log.Errorf("Failed to deliver message: %v", errMsg.Err)
			}
		}
	}()

	return &KafkaProducer{AsyncProducer: producer}
}

// Отправка сообщения в Kafka
func (kp *KafkaProducer) SendMessageToKafka(message model.MessageRequest) error {

	var mu sync.Mutex

	value, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	isAdmin := message.ChatId == 480842950 || message.ChatId == 689105464

	messageStates := make(map[int64]struct{})

	var topic string

	switch message.Message {
	case "/start":
		topic = "start"
	case "/support":
		topic = "support"
	case "🔮 Расклад ТАРО":
		topic = "taro"
	case "💸 Нумерология":
		topic = "numerology"
	case "Админ-панель":
		if isAdmin {
			topic = "admin"
		} else {
			topic = "message"
		}
	case "Отправить сообщение всем":
		if isAdmin {
			mu.Lock()
			messageStates[message.ChatId] = struct{}{}
			mu.Unlock()
		} else {
			topic = "message"
		}
	case "Вывести всех пользователей":
		if isAdmin {
			topic = "get_all_users"
		} else {
			topic = "message"
		}
	default:
		topic = "message"
		if _, ok := messageStates[message.ChatId]; ok {
			topic = "send_notification"
			mu.Lock()
			delete(messageStates, message.ChatId)
			mu.Unlock()
		}
	}
	// Создаем сообщение Kafka
	kafkaMessage := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", message.ChatId)),
		Value: sarama.ByteEncoder(value),
	}

	// Отправляем сообщение в Kafka
	kp.AsyncProducer.Input() <- kafkaMessage

	return nil
}

func (kp *KafkaProducer) Close() {
	err := kp.AsyncProducer.Close()
	if err != nil {
		log.Errorf("Failed to close Sarama producer: %v", err)
	}
}

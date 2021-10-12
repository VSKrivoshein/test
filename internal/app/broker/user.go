package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VSKrivoshein/test/internal/app/service"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	topic     = "topic"
	partition = 0
)

type broker struct {
	conn *kafka.Conn
}

func New() service.Broker {
	host := fmt.Sprintf("%v:%v", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	log.Infof("kafka host and port: %v", host)
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		host,
		topic,
		partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}


	return &broker{
		conn: conn,
	}
}

type Message struct {
	Timestamp int64 `json:"timestamp"`
	Message   string    `json:"message"`
}

func (b *broker) Log(mail string) {
	msg, err := json.Marshal(Message{
		Timestamp: time.Now().Unix(),
		Message:   mail,
	})

	if err = b.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		log.Fatal("kafka SetWriteDeadline err:", err)
	}

	if err != nil {
		return
	}

	_, err = b.conn.WriteMessages(
		kafka.Message{Value: msg},
	)

	if err != nil {
		log.Errorf("kafka WriteMessages err: %v", err)
	}
}

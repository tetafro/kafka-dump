package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	kafka "github.com/segmentio/kafka-go"
	_ "github.com/segmentio/kafka-go/gzip"
	_ "github.com/segmentio/kafka-go/snappy"
	log "github.com/sirupsen/logrus"
)

// Message is a kafka message with timestamp.
type Message struct {
	time time.Time
	data map[string]interface{}
}

// Consumer describes source of messages.
type Consumer interface {
	Read(context.Context) (Message, error)
	Close()
}

// KafkaConsumer is a consumer that reads messages from kafka.
type KafkaConsumer struct {
	conf     kafka.ReaderConfig
	kafka    *kafka.Reader
	messages chan Message
}

// NewKafkaConsumer creates new kafka consumer.
func NewKafkaConsumer(conf KafkaConf) (*KafkaConsumer, error) {
	if len(conf.Brokers) == 0 {
		return nil, fmt.Errorf("brokers list is empty")
	}
	if conf.Topic == "" {
		return nil, fmt.Errorf("topic is empty")
	}
	if conf.GroupID == "" {
		return nil, fmt.Errorf("group id is empty")
	}
	if conf.Offset == 0 {
		conf.Offset = -1
	}

	c := &KafkaConsumer{
		conf: kafka.ReaderConfig{
			Brokers:        conf.Brokers,
			Topic:          conf.Topic,
			GroupID:        conf.GroupID,
			CommitInterval: time.Second,
		},
		messages: make(chan Message),
	}

	switch conf.Offset {
	case -1:
		c.conf.StartOffset = kafka.LastOffset
	case -2:
		c.conf.StartOffset = kafka.FirstOffset
	default:
		err := setOffset(context.TODO(), conf.Brokers, conf.Topic, conf.GroupID, conf.Offset)
		if err != nil {
			return nil, fmt.Errorf("set offsets: %v", err)
		}
	}

	c.kafka = kafka.NewReader(c.conf)

	return c, nil
}

// Read reads next message from kafka.
func (c *KafkaConsumer) Read(ctx context.Context) (Message, error) {
	msg, err := c.kafka.ReadMessage(ctx)
	if err != nil {
		return Message{}, fmt.Errorf("read message: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		return Message{}, fmt.Errorf("invalid json: %s", msg.Value)
	}

	return Message{time: msg.Time, data: data}, nil
}

// Close properly closes kafka connection.
func (c *KafkaConsumer) Close() {
	c.kafka.Close() // nolint: errcheck,gosec
}

func setOffset(ctx context.Context, brokers []string, topic, gid string, ts int64) error {
	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("create connection: %v", err)
	}
	defer conn.Close() // nolint: errcheck

	log.Debug("Read partitions list")
	parts, err := conn.ReadPartitions(topic)
	if err != nil {
		return fmt.Errorf("get partitions: %v", err)
	}

	log.Debugf("Get offsets by timestamp %s", time.Unix(ts, 0))
	offsets := make(map[int]int64, len(parts))
	for _, p := range parts {
		addr := fmt.Sprintf("%s:%d", p.Leader.Host, p.Leader.Port)
		c, err := kafka.DialLeader(ctx, "tcp", addr, topic, p.ID)
		if err != nil {
			return fmt.Errorf("create connection to partition %d: %v", p.ID, err)
		}
		defer c.Close() // nolint: errcheck
		offset, err := c.ReadOffset(time.Unix(ts, 0))
		if err != nil {
			return fmt.Errorf("read offset of partition %d: %v", p.ID, err)
		}
		offsets[p.ID] = offset
	}

	log.Debugf("Set offsets (%d partitions) for %s/%s", len(parts), topic, gid)
	group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		Brokers: brokers,
		Topics:  []string{topic},
		ID:      gid,
	})
	if err != nil {
		return fmt.Errorf("create consumer group: %v", err)
	}
	defer group.Close() // nolint: errcheck

	gen, err := group.Next(ctx)
	if err != nil {
		return fmt.Errorf("get next generation: %v", err)
	}
	err = gen.CommitOffsets(map[string]map[int]int64{topic: offsets})
	if err != nil {
		return fmt.Errorf("commit offsets: %v", err)
	}
	return nil
}

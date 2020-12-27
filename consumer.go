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

type message struct {
	time time.Time
	data map[string]interface{}
}

type consumer struct {
	conf     kafka.ReaderConfig
	kafka    *kafka.Reader
	messages chan message
}

func newConsumer(brokers []string, topic, gid string, start int64) (*consumer, error) {
	c := &consumer{
		conf: kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        gid,
			CommitInterval: time.Second,
		},
		messages: make(chan message),
	}

	switch start {
	case -1:
		c.conf.StartOffset = kafka.LastOffset
	case -2:
		c.conf.StartOffset = kafka.FirstOffset
	case 0:
		c.conf.StartOffset = kafka.FirstOffset
	default:
		if err := setOffset(context.TODO(), brokers, topic, gid, start); err != nil {
			return nil, fmt.Errorf("set offsets: %v", err)
		}
	}

	c.kafka = kafka.NewReader(c.conf)

	return c, nil
}

func (c *consumer) run(ctx context.Context) error {
	defer c.kafka.Close()
	defer close(c.messages)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		msg, err := c.kafka.ReadMessage(ctx)
		if err == context.Canceled {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read message: %v", err)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			log.Errorf("Invalid json: %s", string(msg.Value))
			continue
		}

		c.messages <- message{
			time: msg.Time,
			data: data,
		}
	}
}

func setOffset(ctx context.Context, brokers []string, topic, gid string, ts int64) error {
	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("create connection: %v", err)
	}
	defer conn.Close()

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
		defer c.Close()
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
	defer group.Close()

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

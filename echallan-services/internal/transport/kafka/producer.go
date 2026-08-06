package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	kafkago "github.com/segmentio/kafka-go"
)

// Producer wraps the Kafka writer for pushing messages to topics.
// This replicates the Java Producer.java from tracer library.
type Producer struct {
	writers map[string]*kafkago.Writer
	brokers []string
}

// NewProducer creates a new Kafka producer connected to the specified bootstrap servers.
func NewProducer(brokers string) *Producer {
	return &Producer{
		writers: make(map[string]*kafkago.Writer),
		brokers: []string{brokers},
	}
}

// getWriter returns an existing writer for a topic or creates a new one.
func (p *Producer) getWriter(topic string) *kafkago.Writer {
	if w, ok := p.writers[topic]; ok {
		return w
	}

	w := &kafkago.Writer{
		Addr:     kafkago.TCP(p.brokers...),
		Topic:    topic,
		Balancer: &kafkago.LeastBytes{},
	}
	p.writers[topic] = w
	return w
}

// Push serializes the payload to JSON and pushes it to the specified Kafka topic.
// This is the Go equivalent of the Java producer.push(topic, key, value) call.
func (p *Producer) Push(topic string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload for kafka: %w", err)
	}

	writer := p.getWriter(topic)
	err = writer.WriteMessages(context.Background(),
		kafkago.Message{
			Value: data,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to push to kafka topic %s: %w", topic, err)
	}

	log.Printf("📨 Pushed message to Kafka topic: %s (%d bytes)", topic, len(data))
	return nil
}

// Close gracefully shuts down all Kafka writers.
func (p *Producer) Close() {
	for topic, w := range p.writers {
		if err := w.Close(); err != nil {
			log.Printf("⚠️ Error closing kafka writer for topic %s: %v", topic, err)
		}
	}
}

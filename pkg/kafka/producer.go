package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

type Producer struct {
	Broker []string
	Topic  string
	Key    string
	Body   interface{}
}

func (p Producer) Create(ctx context.Context) {

	bodyAsBytes, err_ := json.Marshal(p.Body)
	if err_ != nil {
		_ = errs.ErrMarshalingJson.Throwf(applog.Log, errs.ErrFmt, err_.Error())
	}
	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: p.Broker,
		Topic:   p.Topic,
	})
	w.AllowAutoTopicCreation = true

	// each kafka message has a key and value. The key is used
	// to decide which partition (and consequently, which broker)
	// the message gets published on
	err := w.WriteMessages(ctx, kafka.Message{
		Key: []byte(p.Key),
		// create an arbitrary message payload for the value
		Value: bodyAsBytes,
	})
	if err != nil {
		_ = errs.ErrWriteKafkaMessage.Throwf(applog.Log, errs.ErrFmt, err.Error())
	}

	// log a confirmation once the message is written
	applog.Log.Infof("Sending message on Topic: %s with Value: %s", p.Topic, p.Body)

}

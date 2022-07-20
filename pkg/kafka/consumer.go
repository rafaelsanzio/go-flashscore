package kafka

import (
	"context"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
	"github.com/rafaelsanzio/go-flashscore/pkg/notification"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Broker      []string
	GroupTopics []string
	GroupID     string
}

func (c Consumer) Init(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.Broker,
		GroupTopics: c.GroupTopics,
		GroupID:     c.GroupID,
	})
	defer r.Close()

	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			_ = errs.ErrReadingKafkaMessage.Throwf(applog.Log, errs.ErrFmt, err.Error())
		}

		err = notification.Handler(ctx, string(msg.Value), string(msg.Key))
		if err != nil {
			_ = errs.ErrHandlingKafkaMessage.Throwf(applog.Log, errs.ErrFmt, err.Error())
		}

		// after receiving the message, log its value
		applog.Log.Infof("Receiving message with Value: %s", string(msg.Value))
	}

}

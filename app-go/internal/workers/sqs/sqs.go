package sqs

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
)

type sqs struct {
	client   *awssqs.Client
	handlers []IHandler
}

func New(awsConfig aws.Config) *sqs {
	sqsclient := awssqs.NewFromConfig(awsConfig)
	return &sqs{
		client:   sqsclient,
		handlers: []IHandler{},
	}
}

type IHandler interface {
	// Handle is to handle message consumed from sqs
	Handle(ctx context.Context, body string) error

	// GetQueueUrl retrieve queue which worker listen to
	GetQueueUrl() *string

	// injectClient is to inject sqs client onto the worker
	injectClient(client *awssqs.Client)
}

func (s *sqs) RegisterHandlers(handlers ...IHandler) {
	s.handlers = append(s.handlers, handlers...)
}

func (s *sqs) Start() {

	for _, handler := range s.handlers {
		go s.runHandler(handler)
	}

}

// runHandler define how to consume message per handler
func (s *sqs) runHandler(handler IHandler) {
	handler.injectClient(s.client)

	// consume message every 1 second
	for range time.Tick(time.Second) {
		log.Printf("fetching message for '%s'", *handler.GetQueueUrl())
		messageCh, _ := s.client.ReceiveMessage(context.TODO(), &awssqs.ReceiveMessageInput{
			QueueUrl:            handler.GetQueueUrl(),
			MaxNumberOfMessages: 1, // NOTE: need to increase this in future
		})

		if messageCh == nil {
			continue
		}

		for _, message := range messageCh.Messages {
			ctx := context.Background()

			log.Printf("consuming message. messageId=%s. message=%s", *message.MessageId, *message.Body)

			// trigger message handler
			err := handler.Handle(ctx, *message.Body)
			if err != nil {
				log.Printf("error on handling message. queueUrl=%s. messsage=%s", *handler.GetQueueUrl(), *message.Body)
				continue
			}

			// when no error found
			// then remove message from queue
			_, err = s.client.DeleteMessage(ctx, &awssqs.DeleteMessageInput{
				QueueUrl:      handler.GetQueueUrl(),
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Printf("error on deleting message from queue. queueUrl=%s. message=%s", *handler.GetQueueUrl(), *message.Body)
			}
		}

	}

}

func (s *sqs) Stop() {
	// sqs doesnt seems to provide this
}

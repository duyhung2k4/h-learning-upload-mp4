package videoservice

import (
	"app/internal/connection"
	constant "app/internal/constants"
	queuepayload "app/internal/dto/queue_payload"
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func (s *videoService) SendMessQueueQuantity(queue constant.QUEUE_QUANTITY, uuidVideo string) error {
	path := fmt.Sprintf("%s.mp4", uuidVideo)
	ipServer := fmt.Sprintf(
		"http://%s:%s/api/v1/video",
		connection.GetConnect().UploadMp4Service.Host,
		connection.GetConnect().UploadMp4Service.Port,
	)

	payload := queuepayload.QueueMp4QuantityPayload{
		Path:     path,
		Uuid:     uuidVideo,
		IpServer: ipServer,
	}

	payloadJsonString, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ch, err := s.connQueue.Channel()
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(context.Background(),
		"",
		string(queue),
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        payloadJsonString,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

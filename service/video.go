package service

import (
	"app/config"
	"app/constant"
	queuepayload "app/dto/queue_payload"
	"app/dto/request"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

type videoService struct {
	connQueue *amqp091.Connection
}

type VideoService interface {
	CreateVideo(payload request.InfoVideo, r *http.Request) error
	SendMessQueueQuantity(queue constant.QUEUE_QUANTITY, uuidVideo string) error
}

func (s *videoService) CreateVideo(payload request.InfoVideo, r *http.Request) error {
	file, _, err := r.FormFile("video")
	if err != nil {
		return err
	}
	defer file.Close()

	fileName := fmt.Sprintf("%s.mp4", payload.Uuid)

	outFile, err := os.Create("video/" + fileName)
	if err != nil {
		return err
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, file)

	return err
}

func (s *videoService) SendMessQueueQuantity(queue constant.QUEUE_QUANTITY, uuidVideo string) error {
	path := fmt.Sprintf("%s.mp4", uuidVideo)
	ipServer := fmt.Sprintf("http://%s:%s/api/v1/video", config.GetAppHost(), config.GetAppPort())

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

func NewVideoService() VideoService {
	return &videoService{
		connQueue: config.GetRabbitmq(),
	}
}

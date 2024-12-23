package videoservice

import (
	"app/internal/connection"
	constant "app/internal/constants"
	requestdata "app/internal/dto/client"
	queuepayload "app/internal/dto/queue_payload"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
)

type videoService struct {
	connQueue *amqp091.Connection
}

type VideoService interface {
	CreateVideo(ctx *gin.Context, payload requestdata.InfoVideo) error
	DeleteVideoMp4(payload queuepayload.QueueFileDeleteMp4) error
	SendMessQueueQuantity(queue constant.QUEUE_QUANTITY, uuidVideo string) error
}

func Register() VideoService {
	return &videoService{
		connQueue: connection.GetRabbitmq(),
	}
}

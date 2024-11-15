package queue

import (
	"app/config"
	"app/constant"
	queuepayload "app/dto/queue_payload"
	"app/service"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type queueFileDeleteMp4 struct {
	videoService service.VideoService
}

type QueueFileDeleteMp4 interface {
	Worker()
}

func (q *queueFileDeleteMp4) Worker() {
	conn := config.GetRabbitmq()
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error chanel: ", err)
		return
	}

	qu, err := ch.QueueDeclare(
		string(constant.QUEUE_DELETE_MP4),
		true,
		false,
		false,
		false,
		amqp091.Table{},
	)
	if err != nil {
		log.Println("error queue declare: ", err)
		return
	}
	log.Printf("start %s", string(constant.QUEUE_DELETE_MP4))

	msgs, err := ch.Consume(
		qu.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("error consumer: ", err)
		return
	}

	for d := range msgs {
		go func(mess amqp091.Delivery) {
			var payload queuepayload.QueueFileDeleteMp4
			err := json.Unmarshal(mess.Body, &payload)
			if err != nil {
				log.Println("error msg: ", err)
				mess.Reject(true)
				return
			}

			err = q.videoService.DeleteVideoMp4(payload)
			if err != nil {
				log.Println("error delete file: ", err)
				mess.Reject(true)
				return
			}

			mess.Ack(false)
		}(d)
	}
}

func NewQueueFileM3U8() QueueFileDeleteMp4 {
	return &queueFileDeleteMp4{
		videoService: service.NewVideoService(),
	}
}

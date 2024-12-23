package videoservice

import (
	queuepayload "app/internal/dto/queue_payload"
	"fmt"
	"os"
)

func (s *videoService) DeleteVideoMp4(payload queuepayload.QueueFileDeleteMp4) error {
	path := fmt.Sprintf("cmd/upload-mp4-service/data/video/%s.mp4", payload.Uuid)
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

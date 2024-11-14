package controller

import (
	"app/constant"
	"app/dto/request"
	"app/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/render"
)

type videoController struct {
	videoService service.VideoService
}

type VideoController interface {
	Upload(w http.ResponseWriter, r *http.Request)
}

func (c *videoController) Upload(w http.ResponseWriter, r *http.Request) {
	var videoPayload request.InfoVideo
	metadata := r.FormValue("metadata")
	err := json.Unmarshal([]byte(metadata), &videoPayload)
	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	// Đọc kích thước của file tải lên từ header, ngăn không cho tải lên vượt quá giới hạn
	const maxUploadSize = 5 << 30 // 5GB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	err = c.videoService.CreateVideo(videoPayload, r)

	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	listErr := []error{}
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := c.videoService.SendMessQueueQuantity(constant.QUEUE_MP4_360_P, videoPayload.Uuid)
		if err != nil {
			mutex.Lock()
			listErr = append(listErr, err)
			mutex.Unlock()
		}
	}()
	wg.Wait()

	if len(listErr) > 0 {
		log.Println(listErr)
		InternalServerError(w, r, errors.New("error send mess"))
		return
	}

	res := Response{
		Data:    nil,
		Message: "OK",
		Status:  200,
		Error:   nil,
	}

	render.JSON(w, r, res)
}

func NewVideoController() VideoController {
	return &videoController{
		videoService: service.NewVideoService(),
	}
}

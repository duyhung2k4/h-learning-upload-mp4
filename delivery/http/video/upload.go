package videohandle

import (
	constant "app/internal/constants"
	requestdata "app/internal/dto/client"
	"app/internal/entity"
	httpresponse "app/pkg/http_response"
	logapp "app/pkg/log"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func (h *videoHandle) Upload(ctx *gin.Context) {
	var videoPayload requestdata.InfoVideo
	metadata := ctx.Request.FormValue("metadata")
	err := json.Unmarshal([]byte(metadata), &videoPayload)

	if err != nil {
		httpresponse.BadRequest(ctx, err)
		logapp.Logger(constant.TITLE_GET_PAYLOAD, err.Error(), constant.ERROR_LOG)
		return
	}

	profileId := ctx.GetUint(string(constant.PROFILE_ID_KEY))

	videoLession, err := h.service.QueryVideoLession.First(requestdata.QueryReq[entity.VideoLession]{
		Joins: []string{
			"JOIN lessions AS l ON l.id = video_lessions.lession_id",
			"JOIN courses AS c ON c.id = l.course_id",
		},
		Condition: `
			video_lessions.code = ?
			AND l.id = ?
			AND c.create_id = ?
		`,
		Args: []interface{}{
			videoPayload.Uuid,
			videoPayload.LessionId,
			profileId,
		},
	})
	if err != nil {
		httpresponse.InternalServerError(ctx, err)
		logapp.Logger("get-video-lession", err.Error(), constant.ERROR_LOG)
		return
	}

	if videoLession.Url360p != nil {
		httpresponse.InternalServerError(ctx, err)
		logapp.Logger("upload-video", "video uploaded", constant.ERROR_LOG)
		return
	}

	const maxUploadSize = 5 << 30 // 5GB
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxUploadSize)
	err = h.service.VideoService.CreateVideo(ctx, videoPayload)

	if err != nil {
		httpresponse.InternalServerError(ctx, err)
		logapp.Logger("read-file", err.Error(), constant.ERROR_LOG)
		return
	}

	ctx.Set("video_lession_id", videoLession.ID)
	err = h.service.VideoService.ChangeStatus(ctx, entity.VIDEO_LESSION_PENDING)
	if err != nil {
		httpresponse.InternalServerError(ctx, err)
		logapp.Logger("change-status", err.Error(), constant.ERROR_LOG)
		return
	}

	listErr := []error{}
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// push mess to queue with quantity
	queues := []constant.QUEUE_QUANTITY{
		constant.QUEUE_MP4_360_P,
		// constant.QUEUE_MP4_480_P,
		// constant.QUEUE_MP4_720_P,
		// constant.QUEUE_MP4_1080_P,
	}

	for _, q := range queues {
		wg.Add(1)
		go func(q constant.QUEUE_QUANTITY) {
			defer wg.Done()
			err := h.service.VideoService.SendMessQueueQuantity(q, videoPayload.Uuid)
			if err != nil {
				mutex.Lock()
				listErr = append(listErr, err)
				mutex.Unlock()
			}
		}(q)
	}

	wg.Wait()

	if len(listErr) > 0 {
		httpresponse.InternalServerError(ctx, errors.New("error send message"))
		for _, e := range listErr {
			logapp.Logger("send-mess", e.Error(), constant.ERROR_LOG)
		}
		return
	}

	httpresponse.Success(ctx, nil)
}

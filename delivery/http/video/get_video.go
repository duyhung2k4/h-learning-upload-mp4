package videohandle

import (
	constant "app/internal/constants"
	httpresponse "app/pkg/http_response"
	logapp "app/pkg/log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (h *videoHandle) GetVideo(ctx *gin.Context) {
	filename := ctx.Param("filename")
	videoPath := filepath.Join("cmd/upload-mp4-service/data/video", filename) // Thay đổi đường dẫn này

	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		httpresponse.BadRequest(ctx, err)
		logapp.Logger("get-video-mp4", err.Error(), constant.ERROR_LOG)
		return
	}

	ctx.FileAttachment(videoPath, filename)
}

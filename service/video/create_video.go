package videoservice

import (
	requestdata "app/internal/dto/client"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func (s *videoService) CreateVideo(ctx *gin.Context, payload requestdata.InfoVideo) error {
	file, _, err := ctx.Request.FormFile("video")
	if err != nil {
		return err
	}
	defer file.Close()

	fileName := fmt.Sprintf("%s.mp4", payload.Uuid)

	outFile, err := os.Create("cmd/upload-mp4-service/data/video/" + fileName)
	if err != nil {
		return err
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, file)

	return err
}

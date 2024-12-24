package videoservice

import (
	"app/internal/entity"

	"github.com/gin-gonic/gin"
)

func (s *videoService) ChangeStatus(ctx *gin.Context, status entity.VIDEO_LESSION_STATUS) error {
	videoLessionId := ctx.GetUint("video_lession_id")
	err := s.psql.
		Model(&entity.VideoLession{}).
		Where("id = ?", videoLessionId).
		Updates(&entity.VideoLession{
			Status: &status,
		}).Error

	if err != nil {
		return err
	}

	return nil
}

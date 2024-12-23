package service

import (
	videoservice "app/cmd/upload-mp4-service/service/video"
	"app/internal/entity"
	query "app/pkg/query/basic"
)

type Service struct {
	VideoService      videoservice.VideoService
	QueryVideoLession query.QueryService[entity.VideoLession]
}

package service

import (
	videoservice "app/cmd/upload-mp4-service/service/video"
	"app/internal/entity"
	query "app/pkg/query/basic"
)

func Register() Service {
	return Service{
		VideoService:      videoservice.Register(),
		QueryVideoLession: query.Register[entity.VideoLession](),
	}
}

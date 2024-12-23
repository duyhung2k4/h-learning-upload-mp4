package videohandle

import (
	"app/cmd/upload-mp4-service/service"
	constant "app/internal/constants"
	middlewareapp "app/internal/middleware"
	routerconfig "app/internal/router_config"

	"github.com/gin-gonic/gin"
)

type videoHandle struct {
	service service.Service
}

type VideoHandle interface {
	Upload(ctx *gin.Context)
	GetVideo(ctx *gin.Context)
}

func NewHandle() VideoHandle {
	return &videoHandle{
		service: service.Register(),
	}
}

func Register(r *gin.Engine) {
	handle := NewHandle()

	routerconfig.AddRouter(r, routerconfig.RouterConfig{
		Method:   constant.POST_HTTP,
		Endpoint: "video/upload",
		Middleware: []gin.HandlerFunc{
			middlewareapp.ValidateToken,
			middlewareapp.GetProfileId,
		},
		Handle: handle.Upload,
	})

	routerconfig.AddRouter(r, routerconfig.RouterConfig{
		Method:     constant.GET_HTTP,
		Endpoint:   "video/:filename",
		Middleware: []gin.HandlerFunc{},
		Handle:     handle.GetVideo,
	})
}

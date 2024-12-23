package httpapp

import (
	videohandle "app/cmd/upload-mp4-service/delivery/http/video"
	"app/internal/connection"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Register() http.Handler {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	videohandle.Register(r)

	log.Printf(
		"Server h-learning-upload-mp4 starting success! URL: http://%s:%s",
		connection.GetConnect().UploadMp4Service.Host,
		connection.GetConnect().UploadMp4Service.Port,
	)

	return r
}

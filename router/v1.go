package router

import (
	"app/controller"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func apiV1(router chi.Router) {
	// middlewares := middlewares.NewMiddlewares()

	videoController := controller.NewVideoController()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]interface{}{
			"mess": "done",
		})
	})

	router.Route("/public", func(public chi.Router) {
	})

	router.Route("/protected", func(protected chi.Router) {
		// protected.Use(jwtauth.Verifier(config.GetJWT()))
		// protected.Use(middlewares.ValidateExpAccessToken())

		protected.Route("/video", func(video chi.Router) {
			video.Post("/upload", videoController.Upload)
		})
	})

	router.Get("/video/{filename}", func(w http.ResponseWriter, r *http.Request) {
		filename := chi.URLParam(r, "filename")
		imagePath := filepath.Join("video", filename) // Thay đổi đường dẫn này

		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "video/mp4") // Hoặc loại hình ảnh tương ứng
		http.ServeFile(w, r, imagePath)
	})
}

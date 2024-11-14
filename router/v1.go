package router

import (
	"app/controller"
	"net/http"

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
}

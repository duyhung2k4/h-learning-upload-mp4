package controller

import "net/http"

type uploadController struct{}

type UploadController interface {
	Upload(w http.ResponseWriter, r *http.Request)
}

func (c *uploadController) Upload(w http.ResponseWriter, r *http.Request) {}

func NewUploadController() UploadController {
	return &uploadController{}
}

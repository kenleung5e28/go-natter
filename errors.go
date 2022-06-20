package main

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	Error          error  `json:"-"`
	HttpStatusCode int    `json:"-"`
	Message        string `json:"message"`
	ErrorText      string `json:"error"`
}

func (e ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HttpStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Error:          err,
		HttpStatusCode: 400,
		Message:        "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrServer(err error) render.Renderer {
	return &ErrResponse{
		Error:          err,
		HttpStatusCode: 500,
		Message:        "Server error.",
		ErrorText:      err.Error(),
	}
}

package main

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	Error          error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	Message        string `json:"message"`
	ErrorText      string `json:"error"`
}

func (e ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Error:          err,
		HTTPStatusCode: 400,
		Message:        "Invalid request.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, Message: "Resource not found."}

func ErrServer(err error) render.Renderer {
	return &ErrResponse{
		Error:          err,
		HTTPStatusCode: 500,
		Message:        "Server error.",
		ErrorText:      err.Error(),
	}
}

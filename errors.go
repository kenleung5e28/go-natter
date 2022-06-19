package main

import (
	"github.com/go-chi/render"
	"net/http"
)

type ResponseError struct {
	Error          error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	Message        string `json:"message"`
	ErrorText      string `json:"error"`
}

func (e *ResponseError) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func InvalidRequestError(err error) render.Renderer {
	return &ResponseError{
		Error:          err,
		HTTPStatusCode: 400,
		Message:        "Invalid request.",
		ErrorText:      err.Error(),
	}
}

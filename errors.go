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

func renderError(w http.ResponseWriter, r *http.Request, e ErrResponse) {
	w.WriteHeader(e.HTTPStatusCode)
	render.JSON(w, r, e)
}

func renderServerError(w http.ResponseWriter, r *http.Request, err error) {
	renderError(w, r, ErrResponse{
		Error:          err,
		HTTPStatusCode: 500,
		Message:        "Server error.",
		ErrorText:      err.Error(),
	})
}

func renderInvalidRequest(w http.ResponseWriter, r *http.Request, err error) {
	renderError(w, r, ErrResponse{
		Error:          err,
		HTTPStatusCode: 400,
		Message:        "Invalid request.",
		ErrorText:      err.Error(),
	})
}

func renderNotFound(w http.ResponseWriter, r *http.Request) {
	renderError(w, r, ErrResponse{HTTPStatusCode: 404, Message: "Resource not found."})
}

package main

import (
	"database/sql"
	"github.com/go-chi/render"
	"net/http"
)

type DbContext struct {
	db *sql.DB
}

func (c DbContext) createSpace(w http.ResponseWriter, r *http.Request) {
	data := &CreateSpaceRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, InvalidRequestError(err))
		return
	}

}

type CreateSpaceRequest struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

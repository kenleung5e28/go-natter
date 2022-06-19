package main

import (
	"database/sql"
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type Env struct {
	db *sql.DB
}

func (e Env) CreateSpace(w http.ResponseWriter, r *http.Request) {
	data := &CreateSpaceRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, InvalidRequestError(err))
		return
	}
	ctx := r.Context()
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		render.Render(w, r, ServerError(err))
		return
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx,
		"INSERT INTO spaces(name, owner) VALUES (?, ?);",
		data.Name, data.Owner)
	if err != nil {
		render.Render(w, r, ServerError(err))
		return
	}
	spaceId, err := res.LastInsertId()
	if err != nil {
		render.Render(w, r, ServerError(err))
		return
	}
	if err = tx.Commit(); err != nil {
		render.Render(w, r, ServerError(err))
		return
	}
	render.Render(w, r, &CreateSpaceResponse{
		Name: data.Name,
		URI:  "/spaces/" + strconv.FormatInt(spaceId, 10),
	})
}

type CreateSpaceRequest struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

func (c CreateSpaceRequest) Bind(_ *http.Request) error {
	if c.Name == "" {
		return errors.New("name must be non-empty")
	}
	if c.Owner == "" {
		return errors.New("owner must be non-empty")
	}
	return nil
}

type CreateSpaceResponse struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

func (c *CreateSpaceResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, 201)
	return nil
}

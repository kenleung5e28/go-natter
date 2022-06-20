package main

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"regexp"
	"strconv"
)

func (e Env) CreateSpace(w http.ResponseWriter, r *http.Request) {
	data := &CreateSpaceRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	ctx := r.Context()
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx,
		"INSERT INTO spaces(name, owner) VALUES (?, ?);",
		data.Name, data.Owner)
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}
	spaceId, err := res.LastInsertId()
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}
	if err = tx.Commit(); err != nil {
		render.Render(w, r, ErrServer(err))
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
	if c.Name == "" || len(c.Name) > 255 {
		return errors.New("length of name must between 1 and 255")
	}
	pattern := regexp.MustCompile("[a-zA-Z][a-zA-Z\\d]{1,29}")
	if !pattern.MatchString(c.Owner) {
		return errors.New("invalid owner: " + c.Owner)
	}
	return nil
}

type CreateSpaceResponse struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

func (CreateSpaceResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, 201)
	return nil
}

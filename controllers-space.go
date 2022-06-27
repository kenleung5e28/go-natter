package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
	"regexp"
)

func (e Env) CreateSpace(w http.ResponseWriter, r *http.Request) {
	data := &CreateSpaceRequest{}
	if err := render.Bind(r, data); err != nil {
		renderInvalidRequest(w, r, err)
		return
	}
	var spaceId int64
	err := transact(e.db, r.Context(), func(tx *sql.Tx, ctx context.Context) error {
		res, err := tx.ExecContext(ctx,
			"INSERT INTO spaces(name, owner) VALUES (?, ?);",
			data.Name, data.Owner)
		if err != nil {
			return err
		}
		spaceId, err = res.LastInsertId()
		return err
	})
	if err != nil {
		renderServerError(w, r, err)
		return
	}
	uri := fmt.Sprintf("/spaces/%d", spaceId)
	w.WriteHeader(201)
	w.Header().Set("Location", uri)
	render.JSON(w, r, CreateSpaceResponse{
		Name: data.Name,
		URI:  uri,
	})
}

type CreateSpaceRequest struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

func (c CreateSpaceRequest) Bind(_ *http.Request) error {
	if c.Name == "" || len(c.Name) > 255 {
		return errors.New("length of name must be between 1 and 255")
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

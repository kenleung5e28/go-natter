package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"regexp"
	"time"
)

func (e Env) AddMessage(w http.ResponseWriter, r *http.Request) {
	spaceId := chi.URLParam(r, "spaceId")
	data := &AddMessageRequest{}
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
	var count int64
	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM spaces WHERE space_id = ?;", spaceId).Scan(&count)
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}
	if count == 0 {
		render.Render(w, r, ErrNotFound)
		return
	}
	res, err := tx.ExecContext(ctx,
		"INSERT INTO messages(space_id, author, msg_text) VALUES (?, ?, ?);",
		spaceId, data.Author, data.Text)
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}
	messageId, err := res.LastInsertId()
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}
	if err = tx.Commit(); err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}
	render.Render(w, r, &AddMessageResponse{
		URI: fmt.Sprintf("/spaces/%s/messages/%d", spaceId, messageId),
	})
}

func (e Env) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	spaceId := chi.URLParam("spaceId")
	sinceRaw := r.URL.Query().Get("since")
	_, err := time.Parse("YYYY-MM-DD hh:mm:ss", sinceRaw)
	sinceValid := err != nil
	ctx := r.Context()
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}
	defer tx.Rollback()
	query := "SELECT msg_id FROM messages WHERE spaceId = ?"
	if sinceValid {
		query += " AND msg_time >= ?"
	}
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		render.Render(w, r, ErrServer(err))
		return
	}
	// TODO
}

type AddMessageRequest struct {
	Author string `json:"author"`
	Text   string `json:"msg_text"`
}

func (a AddMessageRequest) Bind(_ *http.Request) error {
	pattern := regexp.MustCompile("[a-zA-Z][a-zA-Z\\d]{1,29}")
	if !pattern.MatchString(a.Author) {
		return errors.New("invalid author: " + a.Author)
	}
	if a.Text == "" || len(a.Text) > 1024 {
		return errors.New("length of msg_text must be between 1 and 1024")
	}
	return nil
}

type AddMessageResponse struct {
	URI string `json:"uri"`
}

func (AddMessageResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, 201)
	return nil
}

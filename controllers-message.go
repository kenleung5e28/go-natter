package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func (e Env) AddMessage(w http.ResponseWriter, r *http.Request) {
	spaceId := chi.URLParam(r, "spaceId")
	data := &AddMessageRequest{}
	if err := render.Bind(r, data); err != nil {
		renderInvalidRequest(w, r, err)
		return
	}
	notFound := false
	var messageId int64
	err := transact(e.db, r.Context(), func(tx *sql.Tx, ctx context.Context) error {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM spaces WHERE space_id = ?;", spaceId).Scan(&count)
		if err != nil {
			return err
		}
		if count == 0 {
			notFound = true
			return nil
		}
		res, err := tx.ExecContext(ctx,
			"INSERT INTO messages(space_id, author, msg_text) VALUES (?, ?, ?);",
			spaceId, data.Author, data.Text)
		if err != nil {
			return err
		}
		messageId, err = res.LastInsertId()
		return err
	})
	if err != nil {
		renderServerError(w, r, err)
		return
	}
	if notFound {
		renderNotFound(w, r)
		return
	}
	w.WriteHeader(201)
	render.JSON(w, r, AddMessageResponse{
		URI: fmt.Sprintf("/spaces/%s/messages/%d", spaceId, messageId),
	})
}

func (e Env) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	spaceId, err := strconv.ParseInt(chi.URLParam(r, "spaceId"), 10, 64)
	if err != nil {
		renderNotFound(w, r)
		return
	}
	since := r.URL.Query().Get("since")
	if since != "" {
		if _, err := time.Parse(time.RFC3339, since); err != nil {
			renderInvalidRequest(w, r, err)
			return
		}
	}
	var messageIds []string
	err = transact(e.db, r.Context(), func(tx *sql.Tx, ctx context.Context) error {
		var rows *sql.Rows
		if since != "" {
			query := "SELECT msg_id FROM messages WHERE space_id = ? AND msg_time >= ?"
			rows, err = tx.QueryContext(ctx, query, spaceId, since)
		} else {
			query := "SELECT msg_id FROM messages WHERE space_id = ?"
			rows, err = tx.QueryContext(ctx, query, spaceId)
		}
		if err != nil {
			return err
		}
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return err
			}
			messageIds = append(messageIds, id)
		}
		return nil
	})
	if err != nil {
		renderServerError(w, r, err)
		return
	}
	render.JSON(w, r, messageIds)
}

func (e Env) GetMessage(w http.ResponseWriter, r *http.Request) {
	spaceId, err := strconv.ParseInt(chi.URLParam(r, "spaceId"), 10, 64)
	if err != nil {
		renderNotFound(w, r)
		return
	}
	messageId, err := strconv.ParseInt(chi.URLParam(r, "messageId"), 10, 64)
	if err != nil {
		renderNotFound(w, r)
		return
	}
	notFound := false
	data := &GetMessageResponse{}
	err = transact(e.db, r.Context(), func(tx *sql.Tx, ctx context.Context) error {
		err := tx.QueryRowContext(ctx,
			"SELECT space_id, msg_id, author, msg_time, msg_text FROM messages WHERE space_id = ? AND msg_id = ?",
			spaceId, messageId).Scan(&data.SpaceId, &data.MessageId, &data.Author, &data.Time, &data.Message)
		if errors.Is(err, sql.ErrNoRows) {
			notFound = true
			return nil
		}
		return err
	})
	if err != nil {
		renderServerError(w, r, err)
		return
	}
	if notFound {
		renderNotFound(w, r)
		return
	}
	render.JSON(w, r, data)
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

type GetMessageResponse struct {
	SpaceId   int64     `json:"spaceId"`
	MessageId int64     `json:"msgId"`
	Author    string    `json:"author"`
	Time      time.Time `json:"time"`
	Message   string    `json:"message"`
}

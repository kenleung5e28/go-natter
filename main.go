package main

import (
	"database/sql"
	_ "embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

//go:embed schema.sql
var schema string

type Env struct {
	db *sql.DB
}

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Panic(err)
	}
	_, err = db.Exec(schema)
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db: db}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	r.Use(middleware.SetHeader("X-Frame-Options", "DENY"))
	r.Use(middleware.SetHeader("X-XSS-Protection", "0"))
	r.Use(middleware.SetHeader("Cache-Control", "no-store"))
	r.Use(middleware.SetHeader("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; sandbox"))
	r.Route("/spaces", func(r chi.Router) {
		r.Post("/", env.CreateSpace)
		r.Route("/{spaceId}/messages", func(r chi.Router) {
			r.Post("/", env.AddMessage)
			// r.Get("/", env.GetAllMessages)
			// r.Get("/{messageId}", env.GetMessage)
		})
	})
	http.ListenAndServe(":8000", r)
}

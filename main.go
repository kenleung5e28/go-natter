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
	r.Route("/spaces", func(r chi.Router) {
		r.Post("/", env.CreateSpace)
	})
	http.ListenAndServe(":8000", r)
}

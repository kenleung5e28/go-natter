package main

import (
	"database/sql"
	_ "embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	dbContext := &DbContext{db: db}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("你老闆"))
	})
	http.ListenAndServe(":8000", r)
}

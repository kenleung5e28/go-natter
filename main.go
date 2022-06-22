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

func filterInvalidContentTypeRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && render.GetRequestContentType(r) != render.ContentTypeJSON {
			render.Render(w, r, &ErrResponse{
				HTTPStatusCode: 415,
				Message:        http.StatusText(415),
				ErrorText:      "Only application/json supported",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func addSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; sandbox")
		next.ServeHTTP(w, r)
	})
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
	r.Use(addSecurityHeaders)
	r.Use(filterInvalidContentTypeRequests)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Route("/spaces", func(r chi.Router) {
		r.Post("/", env.CreateSpace)
		r.Route("/{spaceId}/messages", func(r chi.Router) {
			r.Post("/", env.AddMessage)
			r.Get("/", env.GetAllMessages)
			r.Get("/{messageId}", env.GetMessage)
		})
	})
	http.ListenAndServe(":8000", r)
}

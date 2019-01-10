package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type key string

// PortKey is a context key for port
const portKey key = "port"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(setCtx)
	r.Get("/", helloChi)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	port = fmt.Sprintf(":%s", port)
	http.ListenAndServe(port, r)
}

func setCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		port := r.Host
		ctx := context.WithValue(r.Context(), portKey, port)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func helloChi(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	val := ctx.Value(portKey).(string)
	log.Printf("hi, this server implements with chi %s", val)
	os.Stdout.Write([]byte("port: " + val + "\n"))
	w.Write([]byte("hi, this server implements with chi"))
}

package server

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	swaggerassets "github.com/Baymax6s/KOBE-Tech/api/swagger"
)

func registerSwaggerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /swagger/{$}", func(w http.ResponseWriter, r *http.Request) {
		serveSwaggerFile(w, r, "index.html", "text/html; charset=utf-8")
	})
	mux.HandleFunc("GET /swagger/index.html", func(w http.ResponseWriter, r *http.Request) {
		serveSwaggerFile(w, r, "index.html", "text/html; charset=utf-8")
	})
	mux.HandleFunc("GET /swagger/openapi.yml", func(w http.ResponseWriter, r *http.Request) {
		serveSwaggerFile(w, r, "openapi.yml", "application/yaml")
	})
	mux.HandleFunc("GET /swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		serveSwaggerFile(w, r, "swagger.json", "application/json; charset=utf-8")
	})
}

func serveSwaggerFile(w http.ResponseWriter, r *http.Request, name, contentType string) {
	body, err := fs.ReadFile(swaggerassets.Files, name)
	if err != nil {
		http.Error(w, fmt.Sprintf("swagger asset %q could not be loaded", name), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(body))
}

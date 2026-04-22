package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSwaggerRedirect(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/swagger", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMovedPermanently)
	}

	if location := rec.Header().Get("Location"); location != "/swagger/" {
		t.Fatalf("location = %q, want %q", location, "/swagger/")
	}
}

func TestSwaggerIndex(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/swagger/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if contentType := rec.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "text/html") {
		t.Fatalf("content-type = %q, want text/html", contentType)
	}

	if !strings.Contains(rec.Body.String(), `id="swagger-ui"`) {
		t.Fatalf("swagger ui container was not rendered: %q", rec.Body.String())
	}
}

func TestSwaggerOpenAPI(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/swagger/openapi.yml", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if contentType := rec.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/yaml") {
		t.Fatalf("content-type = %q, want application/yaml", contentType)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `swagger: "2.0"`) {
		t.Fatalf("swagger version missing: %q", body)
	}

	if !strings.Contains(body, "/api/v1/auth/login:") {
		t.Fatalf("login path missing from spec: %q", body)
	}
}

func TestSwaggerJSON(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/swagger/swagger.json", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if contentType := rec.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
		t.Fatalf("content-type = %q, want application/json", contentType)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `"swagger": "2.0"`) {
		t.Fatalf("swagger version missing: %q", body)
	}

	if !strings.Contains(body, `"/api/v1/articles"`) {
		t.Fatalf("article path missing from spec: %q", body)
	}
}

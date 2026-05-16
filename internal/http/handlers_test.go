package httpapp

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"hotel_site/internal/service"
)

// stubRepository is an in-memory test double implementing contentRepository.
// It allows tests to control success/error responses from HomePageData.
type stubRepository struct {
	// data is returned when homeErr is nil.
	data service.HomePageData
	// homeErr forces repository failures to test error handling paths.
	homeErr error
}

// HomePageData returns either configured error or preloaded data.
func (s stubRepository) HomePageData(context.Context) (service.HomePageData, error) {
	if s.homeErr != nil {
		return service.HomePageData{}, s.homeErr
	}
	return s.data, nil
}

// TestMain changes working directory to repository root so template relative
// paths used by New() resolve correctly during test execution.
func TestMain(m *testing.M) {
	if err := os.Chdir("../.."); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// TestHandleHomeRendersRoot verifies that "/" renders home page successfully.
func TestHandleHomeRendersRoot(t *testing.T) {
	app := newTestApp(t, stubRepository{
		data: service.HomePageData{
			Site: service.DefaultSiteContent(),
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	app.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "Уютный отель с двориком у моря в Феодосии") {
		t.Fatalf("response does not contain hero title: %s", rec.Body.String())
	}
}

// TestHandleHomeAcceptsIndexHTML verifies that "/index.html" is treated as home.
func TestHandleHomeAcceptsIndexHTML(t *testing.T) {
	app := newTestApp(t, stubRepository{
		data: service.HomePageData{
			Site: service.DefaultSiteContent(),
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	rec := httptest.NewRecorder()

	app.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}
}

// TestHandleHomeRenders404Page verifies that unknown paths render custom 404 template.
func TestHandleHomeRenders404Page(t *testing.T) {
	app := newTestApp(t, stubRepository{
		data: service.HomePageData{
			Site: service.DefaultSiteContent(),
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()

	app.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNotFound)
	}
	if !strings.Contains(rec.Body.String(), "Страница не найдена") {
		t.Fatalf("expected custom 404 page, got %s", rec.Body.String())
	}
}

// TestHandleHomeReturns500OnRepositoryError verifies repository failures map to HTTP 500.
func TestHandleHomeReturns500OnRepositoryError(t *testing.T) {
	app := newTestApp(t, stubRepository{homeErr: errors.New("boom")})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	app.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusInternalServerError)
	}
}

// TestHandleHomeAPIMethodNotAllowed verifies non-GET requests to API are rejected with 405.
func TestHandleHomeAPIMethodNotAllowed(t *testing.T) {
	app := newTestApp(t, stubRepository{})

	req := httptest.NewRequest(http.MethodPost, "/api/home", nil)
	rec := httptest.NewRecorder()

	app.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

// newTestApp constructs App for tests and fails test immediately on startup error.
func newTestApp(t *testing.T, repo stubRepository) *App {
	t.Helper()

	app, err := New(repo)
	if err != nil {
		t.Fatalf("new app: %v", err)
	}

	return app
}

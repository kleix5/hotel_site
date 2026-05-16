package httpapp

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"hotel_site/internal/repository"
	"hotel_site/internal/service"
)

// contentRepository is the minimal read contract needed by public HTTP handlers.
// Keeping this interface narrow simplifies tests and decouples handlers from storage details.
type contentRepository interface {
	HomePageData(ctx context.Context) (service.HomePageData, error)
}

// App wires HTTP handlers with templates and a repository abstraction.
type App struct {
	// repo provides homepage data used by HTML and JSON handlers.
	repo contentRepository
	// tmpl stores parsed templates shared across all requests.
	tmpl *template.Template
}

// New builds an App by parsing required template files and storing repository dependency.
func New(repo contentRepository) (*App, error) {
	// ParseFiles validates templates during startup so request handlers fail less often at runtime.
	tmpl, err := template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/index.html",
		"web/templates/404.html",
	)
	if err != nil {
		return nil, err
	}

	return &App{
		repo: repo,
		tmpl: tmpl,
	}, nil
}

// Routes registers all HTTP endpoints and wraps them with common middleware.
func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()
	// Public homepage in HTML format.
	mux.HandleFunc("/", a.handleHome)
	// JSON endpoint for the same homepage data.
	mux.HandleFunc("/api/home", a.handleHomeAPI)
	// Static files (CSS and uploaded media) served from web/static.
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Apply request logging around all routes.
	return loggingMiddleware(mux)
}

// handleHome renders the homepage template for "/" and "/index.html".
// Any other path is treated as an application-level 404.
func (a *App) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		a.renderNotFound(w, r)
		return
	}

	// Load all data needed for server-side template rendering.
	data, err := a.repo.HomePageData(r.Context())
	if err != nil {
		http.Error(w, "не удалось загрузить главную страницу", http.StatusInternalServerError)
		return
	}

	// Render layout template which includes content template blocks.
	if err := a.tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "не удалось отрисовать страницу", http.StatusInternalServerError)
	}
}

// Compile-time assertion that concrete repository implements handler dependency contract.
var _ contentRepository = (*repository.Repository)(nil)

// renderNotFound renders custom 404 page including requested unresolved path.
func (a *App) renderNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if err := a.tmpl.ExecuteTemplate(w, "not_found", map[string]any{
		"Path": r.URL.Path,
	}); err != nil {
		http.Error(w, "страница не найдена", http.StatusNotFound)
	}
}

// handleHomeAPI exposes homepage data as JSON for client-side or integration use.
// Only GET is accepted.
func (a *App) handleHomeAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	data, err := a.repo.HomePageData(r.Context())
	if err != nil {
		http.Error(w, "не удалось загрузить данные главной страницы", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

// writeJSON writes a JSON payload with explicit status and content type.
// Encoding errors are ignored because this helper is used for already-serializable structs.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// loggingMiddleware logs method, URL path, and request processing latency.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

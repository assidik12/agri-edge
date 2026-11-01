package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"agriedge/backend/internal/domain"
	"agriedge/backend/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Handler struct {
	metricUC   *usecase.QueryMetricsUseCase
	reportRepo domain.ReportRepository
}

func NewHandler(metricUC *usecase.QueryMetricsUseCase, reportRepo domain.ReportRepository) *Handler {
	return &Handler{
		metricUC:   metricUC,
		reportRepo: reportRepo,
	}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// CORS configuration for Frontend Web Dashboard
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/healthz", h.HealthCheck)

	r.Route("/api/v1/nodes/{id}", func(r chi.Router) {
		r.Get("/metrics", h.GetMetrics)
		r.Get("/metrics/latest", h.GetLatestMetric)
		r.Get("/reports/latest", h.GetLatestReport)
		r.Get("/reports", h.GetReportByDate)
	})

	return r
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// GetMetrics handles GET /api/v1/nodes/{id}/metrics?from=&to=
func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "id")
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid 'from' timestamp: use RFC3339")
			return
		}
	}
	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid 'to' timestamp: use RFC3339")
			return
		}
	}

	metrics, err := h.metricUC.GetMetricsByNode(r.Context(), nodeID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if metrics == nil {
		metrics = []domain.DailyMetric{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"node_id": nodeID,
		"count":   len(metrics),
		"data":    metrics,
	})
}

// GetLatestMetric handles GET /api/v1/nodes/{id}/metrics/latest
func (h *Handler) GetLatestMetric(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "id")
	metric, err := h.metricUC.GetLatestMetric(r.Context(), nodeID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, http.StatusNotFound, "metric not found for node")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, metric)
}

// GetLatestReport handles GET /api/v1/nodes/{id}/reports/latest
func (h *Handler) GetLatestReport(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "id")
	if h.reportRepo == nil {
		respondError(w, http.StatusServiceUnavailable, "report repository unavailable")
		return
	}

	report, err := h.reportRepo.GetLatestReport(r.Context(), nodeID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, http.StatusNotFound, "report not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, report)
}

// GetReportByDate handles GET /api/v1/nodes/{id}/reports?date=2026-10-24
func (h *Handler) GetReportByDate(w http.ResponseWriter, r *http.Request) {
	nodeID := chi.URLParam(r, "id")
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		h.GetLatestReport(w, r)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		return
	}

	if h.reportRepo == nil {
		respondError(w, http.StatusServiceUnavailable, "report repository unavailable")
		return
	}

	report, err := h.reportRepo.GetReportByDate(r.Context(), nodeID, date)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondError(w, http.StatusNotFound, "report not found for specified date")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, report)
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"error": message,
	})
}

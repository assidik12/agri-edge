package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"agriedge/backend/internal/domain"
	"agriedge/backend/internal/usecase"
)

type dummyMetricRepo struct {
	metrics []domain.DailyMetric
	latest  *domain.DailyMetric
}

func (d *dummyMetricRepo) SaveMetric(ctx context.Context, metric *domain.DailyMetric) error {
	return nil
}

func (d *dummyMetricRepo) GetMetricsByNode(ctx context.Context, nodeID string, from, to time.Time) ([]domain.DailyMetric, error) {
	return d.metrics, nil
}

func (d *dummyMetricRepo) GetLatestMetric(ctx context.Context, nodeID string) (*domain.DailyMetric, error) {
	if d.latest == nil {
		return nil, domain.ErrNotFound
	}
	return d.latest, nil
}

func TestHealthCheck(t *testing.T) {
	h := NewHandler(nil, nil)
	r := h.Routes()

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res map[string]string
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if res["status"] != "ok" {
		t.Errorf("expected status ok, got %s", res["status"])
	}
}

func TestGetMetrics(t *testing.T) {
	repo := &dummyMetricRepo{
		metrics: []domain.DailyMetric{
			{
				NodeID:            "node1",
				SoilMoisturePct:   40.0,
				PumpDurationSec:   120,
				RecordedAt:        time.Now(),
			},
		},
	}
	uc := usecase.NewQueryMetricsUseCase(repo)
	h := NewHandler(uc, nil)
	r := h.Routes()

	req := httptest.NewRequest("GET", "/api/v1/nodes/node1/metrics", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

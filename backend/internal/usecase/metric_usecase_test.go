package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"agriedge/backend/internal/domain"
)

type mockMetricRepo struct {
	saved  *domain.DailyMetric
	err    error
	list   []domain.DailyMetric
	latest *domain.DailyMetric
}

func (m *mockMetricRepo) SaveMetric(ctx context.Context, metric *domain.DailyMetric) error {
	if m.err != nil {
		return m.err
	}
	m.saved = metric
	return nil
}

func (m *mockMetricRepo) GetMetricsByNode(ctx context.Context, nodeID string, from, to time.Time) ([]domain.DailyMetric, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.list, nil
}

func (m *mockMetricRepo) GetLatestMetric(ctx context.Context, nodeID string) (*domain.DailyMetric, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.latest, nil
}

func TestIngestMetricUseCase_Execute(t *testing.T) {
	repo := &mockMetricRepo{}
	uc := NewIngestMetricUseCase(repo)

	payload := &MetricPayload{
		Timestamp: "2026-10-24T17:00:00Z",
	}
	payload.Metrics.SoilMoisturePct = 42.5
	payload.Metrics.TemperatureC = 33.2
	payload.Metrics.AirHumidityPct = 58.0
	payload.Metrics.LightAccum = 8950
	payload.AIDecision.PumpDurationSec = 150
	payload.AIDecision.InferenceTimeMs = 12

	ctx := context.Background()
	err := uc.Execute(ctx, "node1", payload)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if repo.saved == nil {
		t.Fatal("expected saved metric not nil")
	}
	if repo.saved.NodeID != "node1" {
		t.Errorf("expected node1, got %s", repo.saved.NodeID)
	}
	if repo.saved.PumpDurationSec != 150 {
		t.Errorf("expected pump duration 150, got %d", repo.saved.PumpDurationSec)
	}
}

func TestIngestMetricUseCase_ValidationFail(t *testing.T) {
	repo := &mockMetricRepo{}
	uc := NewIngestMetricUseCase(repo)

	// Invalid pump duration > 300
	payload := &MetricPayload{
		Timestamp: "2026-10-24T17:00:00Z",
	}
	payload.Metrics.SoilMoisturePct = 42.5
	payload.AIDecision.PumpDurationSec = 500

	ctx := context.Background()
	err := uc.Execute(ctx, "node1", payload)
	if err == nil {
		t.Fatal("expected error for pump duration > 300, got nil")
	}
	if !errors.Is(err, domain.ErrInvalidPayload) {
		t.Errorf("expected ErrInvalidPayload wrapped, got %v", err)
	}
}

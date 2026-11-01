package usecase

import (
	"context"
	"fmt"
	"time"

	"agriedge/backend/internal/domain"
)

// MetricPayload matches the MQTT JSON contract described in docs/MQTT_PAYLOAD.md
type MetricPayload struct {
	Timestamp string `json:"timestamp"`
	Metrics   struct {
		SoilMoisturePct  float32 `json:"soil_moisture_pct"`
		TemperatureC     float32 `json:"temperature_c"`
		AirHumidityPct   float32 `json:"air_humidity_pct"`
		LightAccum       int     `json:"light_accumulation"`
	} `json:"metrics"`
	AIDecision struct {
		PumpDurationSec int `json:"pump_duration_sec"`
		InferenceTimeMs int `json:"inference_time_ms"`
	} `json:"ai_decision"`
}

// IngestMetricUseCase orchestrates incoming telemetry and stores it to TimescaleDB
type IngestMetricUseCase struct {
	repo domain.MetricRepository
}

func NewIngestMetricUseCase(repo domain.MetricRepository) *IngestMetricUseCase {
	return &IngestMetricUseCase{repo: repo}
}

func (uc *IngestMetricUseCase) Execute(ctx context.Context, nodeID string, payload *MetricPayload) error {
	if payload == nil {
		return fmt.Errorf("%w: payload is nil", domain.ErrInvalidPayload)
	}

	// Parse ISO-8601 timestamp
	t, err := time.Parse(time.RFC3339, payload.Timestamp)
	if err != nil {
		// Fallback parse without timezone if naive
		t, err = time.Parse("2006-01-02T15:04:05", payload.Timestamp)
		if err != nil {
			return fmt.Errorf("%w: invalid timestamp format (%v)", domain.ErrInvalidPayload, err)
		}
	}

	metric := &domain.DailyMetric{
		RecordedAt:        t,
		NodeID:            nodeID,
		SoilMoisturePct:   payload.Metrics.SoilMoisturePct,
		TemperatureC:      payload.Metrics.TemperatureC,
		AirHumidityPct:    payload.Metrics.AirHumidityPct,
		LightAccumulation: payload.Metrics.LightAccum,
		PumpDurationSec:   payload.AIDecision.PumpDurationSec,
		InferenceTimeMs:   payload.AIDecision.InferenceTimeMs,
		CreatedAt:         time.Now(),
	}

	if err := metric.Validate(); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrInvalidPayload, err)
	}

	return uc.repo.SaveMetric(ctx, metric)
}

// QueryMetricsUseCase retrieves metrics history for frontend visualization
type QueryMetricsUseCase struct {
	repo domain.MetricRepository
}

func NewQueryMetricsUseCase(repo domain.MetricRepository) *QueryMetricsUseCase {
	return &QueryMetricsUseCase{repo: repo}
}

func (uc *QueryMetricsUseCase) GetMetricsByNode(ctx context.Context, nodeID string, from, to time.Time) ([]domain.DailyMetric, error) {
	if nodeID == "" {
		return nil, fmt.Errorf("nodeID is required")
	}
	if from.IsZero() {
		from = time.Now().AddDate(0, 0, -30) // default 30 days
	}
	if to.IsZero() {
		to = time.Now()
	}
	return uc.repo.GetMetricsByNode(ctx, nodeID, from, to)
}

func (uc *QueryMetricsUseCase) GetLatestMetric(ctx context.Context, nodeID string) (*domain.DailyMetric, error) {
	if nodeID == "" {
		return nil, fmt.Errorf("nodeID is required")
	}
	return uc.repo.GetLatestMetric(ctx, nodeID)
}

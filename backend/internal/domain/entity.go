package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidPayload = errors.New("invalid metric payload")
	ErrNotFound       = errors.New("resource not found")
)

// DailyMetric represents sensory inputs and edge AI pump decision from an IoT node.
type DailyMetric struct {
	RecordedAt        time.Time `json:"recorded_at"`
	NodeID            string    `json:"node_id"`
	SoilMoisturePct   float32   `json:"soil_moisture_pct"`
	TemperatureC      float32   `json:"temperature_c"`
	AirHumidityPct    float32   `json:"air_humidity_pct"`
	LightAccumulation int       `json:"light_accumulation"`
	PumpDurationSec   int       `json:"pump_duration_sec"`
	InferenceTimeMs   int       `json:"inference_time_ms"`
	CreatedAt         time.Time `json:"created_at"`
}

// Validate checks the metric fields against domain safety boundaries.
func (m *DailyMetric) Validate() error {
	if m.NodeID == "" {
		return errors.New("node_id is required")
	}
	if m.RecordedAt.IsZero() {
		return errors.New("recorded_at is required")
	}
	if m.SoilMoisturePct < 0 || m.SoilMoisturePct > 100 {
		return errors.New("soil_moisture_pct must be between 0 and 100")
	}
	if m.TemperatureC < -10 || m.TemperatureC > 60 {
		return errors.New("temperature_c must be between -10 and 60")
	}
	if m.AirHumidityPct < 0 || m.AirHumidityPct > 100 {
		return errors.New("air_humidity_pct must be between 0 and 100")
	}
	if m.LightAccumulation < 0 {
		return errors.New("light_accumulation must be non-negative")
	}
	// PUMP_MAX_SEC safety clamp is 300 seconds
	if m.PumpDurationSec < 0 || m.PumpDurationSec > 300 {
		return errors.New("pump_duration_sec must be between 0 and 300")
	}
	if m.InferenceTimeMs < 0 {
		return errors.New("inference_time_ms must be non-negative")
	}
	return nil
}

// DailyReport represents LLM-generated narrative summary for daily farm condition.
type DailyReport struct {
	ID         int64     `json:"id"`
	NodeID     string    `json:"node_id"`
	ReportDate time.Time `json:"report_date"`
	Narrative  string    `json:"narrative"`
	Model      string    `json:"model"`
	CreatedAt  time.Time `json:"created_at"`
}

// MetricRepository defines the persistence contract for metric data (TimescaleDB).
type MetricRepository interface {
	SaveMetric(ctx context.Context, metric *DailyMetric) error
	GetMetricsByNode(ctx context.Context, nodeID string, from, to time.Time) ([]DailyMetric, error)
	GetLatestMetric(ctx context.Context, nodeID string) (*DailyMetric, error)
}

// ReportRepository defines the persistence contract for LLM daily reports.
type ReportRepository interface {
	SaveReport(ctx context.Context, report *DailyReport) error
	GetLatestReport(ctx context.Context, nodeID string) (*DailyReport, error)
	GetReportByDate(ctx context.Context, nodeID string, date time.Time) (*DailyReport, error)
}

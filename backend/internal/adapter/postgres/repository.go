package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"agriedge/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// SaveMetric inserts a daily metric record into the TimescaleDB hypertable
func (r *Repository) SaveMetric(ctx context.Context, m *domain.DailyMetric) error {
	query := `
		INSERT INTO daily_metrics (
			recorded_at, node_id, soil_moisture_pct, temperature_c, air_humidity_pct,
			light_accumulation, pump_duration_sec, inference_time_ms, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (recorded_at, node_id) DO UPDATE SET
			soil_moisture_pct = EXCLUDED.soil_moisture_pct,
			temperature_c = EXCLUDED.temperature_c,
			air_humidity_pct = EXCLUDED.air_humidity_pct,
			light_accumulation = EXCLUDED.light_accumulation,
			pump_duration_sec = EXCLUDED.pump_duration_sec,
			inference_time_ms = EXCLUDED.inference_time_ms;
	`
	createdAt := m.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	_, err := r.pool.Exec(ctx, query,
		m.RecordedAt,
		m.NodeID,
		m.SoilMoisturePct,
		m.TemperatureC,
		m.AirHumidityPct,
		m.LightAccumulation,
		m.PumpDurationSec,
		m.InferenceTimeMs,
		createdAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save metric: %w", err)
	}
	return nil
}

// GetMetricsByNode retrieves metrics for a node within a time window
func (r *Repository) GetMetricsByNode(ctx context.Context, nodeID string, from, to time.Time) ([]domain.DailyMetric, error) {
	query := `
		SELECT recorded_at, node_id, soil_moisture_pct, temperature_c, air_humidity_pct,
		       light_accumulation, pump_duration_sec, inference_time_ms, created_at
		FROM daily_metrics
		WHERE node_id = $1 AND recorded_at >= $2 AND recorded_at <= $3
		ORDER BY recorded_at ASC;
	`
	rows, err := r.pool.Query(ctx, query, nodeID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}
	defer rows.Close()

	var results []domain.DailyMetric
	for rows.Next() {
		var m domain.DailyMetric
		if err := rows.Scan(
			&m.RecordedAt,
			&m.NodeID,
			&m.SoilMoisturePct,
			&m.TemperatureC,
			&m.AirHumidityPct,
			&m.LightAccumulation,
			&m.PumpDurationSec,
			&m.InferenceTimeMs,
			&m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan metric row: %w", err)
		}
		results = append(results, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return results, nil
}

// GetLatestMetric gets the most recent metric recorded for a node
func (r *Repository) GetLatestMetric(ctx context.Context, nodeID string) (*domain.DailyMetric, error) {
	query := `
		SELECT recorded_at, node_id, soil_moisture_pct, temperature_c, air_humidity_pct,
		       light_accumulation, pump_duration_sec, inference_time_ms, created_at
		FROM daily_metrics
		WHERE node_id = $1
		ORDER BY recorded_at DESC
		LIMIT 1;
	`
	var m domain.DailyMetric
	err := r.pool.QueryRow(ctx, query, nodeID).Scan(
		&m.RecordedAt,
		&m.NodeID,
		&m.SoilMoisturePct,
		&m.TemperatureC,
		&m.AirHumidityPct,
		&m.LightAccumulation,
		&m.PumpDurationSec,
		&m.InferenceTimeMs,
		&m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get latest metric: %w", err)
	}
	return &m, nil
}

// SaveReport stores or updates a daily LLM report
func (r *Repository) SaveReport(ctx context.Context, rep *domain.DailyReport) error {
	query := `
		INSERT INTO daily_reports (node_id, report_date, narrative, model, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (node_id, report_date) DO UPDATE SET
			narrative = EXCLUDED.narrative,
			model = EXCLUDED.model,
			created_at = EXCLUDED.created_at
		RETURNING id;
	`
	createdAt := rep.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	err := r.pool.QueryRow(ctx, query,
		rep.NodeID,
		rep.ReportDate,
		rep.Narrative,
		rep.Model,
		createdAt,
	).Scan(&rep.ID)
	if err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}
	return nil
}

// GetLatestReport returns the newest report for a node
func (r *Repository) GetLatestReport(ctx context.Context, nodeID string) (*domain.DailyReport, error) {
	query := `
		SELECT id, node_id, report_date, narrative, model, created_at
		FROM daily_reports
		WHERE node_id = $1
		ORDER BY report_date DESC
		LIMIT 1;
	`
	var rep domain.DailyReport
	err := r.pool.QueryRow(ctx, query, nodeID).Scan(
		&rep.ID,
		&rep.NodeID,
		&rep.ReportDate,
		&rep.Narrative,
		&rep.Model,
		&rep.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get latest report: %w", err)
	}
	return &rep, nil
}

// GetReportByDate returns report by specific date
func (r *Repository) GetReportByDate(ctx context.Context, nodeID string, date time.Time) (*domain.DailyReport, error) {
	query := `
		SELECT id, node_id, report_date, narrative, model, created_at
		FROM daily_reports
		WHERE node_id = $1 AND report_date = $2
		LIMIT 1;
	`
	var rep domain.DailyReport
	err := r.pool.QueryRow(ctx, query, nodeID, date).Scan(
		&rep.ID,
		&rep.NodeID,
		&rep.ReportDate,
		&rep.Narrative,
		&rep.Model,
		&rep.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get report by date: %w", err)
	}
	return &rep, nil
}

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- Table: daily_metrics
-- Note: TimescaleDB hypertable requires the partitioning time column to be part of the primary key / unique constraint.
CREATE TABLE IF NOT EXISTS daily_metrics (
    recorded_at TIMESTAMPTZ NOT NULL,
    node_id TEXT NOT NULL,
    soil_moisture_pct REAL NOT NULL,
    temperature_c REAL NOT NULL,
    air_humidity_pct REAL NOT NULL,
    light_accumulation INT NOT NULL,
    pump_duration_sec INT NOT NULL,
    inference_time_ms INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (recorded_at, node_id)
);

-- Convert daily_metrics to TimescaleDB hypertable partitioned by recorded_at
SELECT create_hypertable('daily_metrics', 'recorded_at', if_not_exists => TRUE);

-- Index for node-specific queries
CREATE INDEX IF NOT EXISTS idx_daily_metrics_node_time ON daily_metrics (node_id, recorded_at DESC);

-- Table: daily_reports (Gemini summary)
CREATE TABLE IF NOT EXISTS daily_reports (
    id BIGSERIAL PRIMARY KEY,
    node_id TEXT NOT NULL,
    report_date DATE NOT NULL,
    narrative TEXT NOT NULL,
    model TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (node_id, report_date)
);

CREATE INDEX IF NOT EXISTS idx_daily_reports_node_date ON daily_reports (node_id, report_date DESC);

package domain

import (
	"testing"
	"time"
)

func TestDailyMetric_Validate(t *testing.T) {
	now := time.Now()
	validMetric := &DailyMetric{
		RecordedAt:        now,
		NodeID:            "node1",
		SoilMoisturePct:   45.5,
		TemperatureC:      30.0,
		AirHumidityPct:    60.0,
		LightAccumulation: 8000,
		PumpDurationSec:   150,
		InferenceTimeMs:   12,
	}

	if err := validMetric.Validate(); err != nil {
		t.Fatalf("expected valid metric, got error: %v", err)
	}

	tests := []struct {
		name    string
		mutate  func(m *DailyMetric)
		wantErr bool
	}{
		{
			name: "missing node id",
			mutate: func(m *DailyMetric) {
				m.NodeID = ""
			},
			wantErr: true,
		},
		{
			name: "missing recorded_at",
			mutate: func(m *DailyMetric) {
				m.RecordedAt = time.Time{}
			},
			wantErr: true,
		},
		{
			name: "invalid soil moisture > 100",
			mutate: func(m *DailyMetric) {
				m.SoilMoisturePct = 105.0
			},
			wantErr: true,
		},
		{
			name: "invalid temperature < -10",
			mutate: func(m *DailyMetric) {
				m.TemperatureC = -15.0
			},
			wantErr: true,
		},
		{
			name: "invalid pump duration > 300 (exceeds PUMP_MAX_SEC)",
			mutate: func(m *DailyMetric) {
				m.PumpDurationSec = 350
			},
			wantErr: true,
		},
		{
			name: "negative light accumulation",
			mutate: func(m *DailyMetric) {
				m.LightAccumulation = -1
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := *validMetric
			tt.mutate(&m)
			err := m.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

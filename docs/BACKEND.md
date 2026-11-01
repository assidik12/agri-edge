# Backend (Golang)

Folder: `backend/` — **belum dibuat**. Konfigurasi via env (lihat [.env.example](../.env.example)).

## Struktur Rencana (Clean Architecture)
```
backend/
├── cmd/server/main.go
├── internal/
│   ├── domain/        # entity (DailyMetric, Report) + interface repository (tanpa dependensi eksternal)
│   ├── usecase/       # IngestMetric, GenerateDailyReport, ListReports
│   ├── adapter/
│   │   ├── mqtt/      # subscriber (paho.mqtt.golang)
│   │   ├── postgres/  # repository (pgx)
│   │   └── llm/       # klien Gemini
│   └── delivery/http/ # REST (chi/gin)
├── migrations/
└── Dockerfile
```

## Skema Database (rencana)
```sql
CREATE TABLE daily_metrics (
  id BIGSERIAL PRIMARY KEY,
  node_id TEXT NOT NULL,
  recorded_at TIMESTAMPTZ NOT NULL,
  soil_moisture_pct REAL, temperature_c REAL, air_humidity_pct REAL,
  light_accumulation INT,
  pump_duration_sec INT, inference_time_ms INT,
  UNIQUE (node_id, recorded_at)
);
CREATE TABLE daily_reports (
  id BIGSERIAL PRIMARY KEY,
  node_id TEXT NOT NULL, report_date DATE NOT NULL,
  narrative TEXT NOT NULL, model TEXT, created_at TIMESTAMPTZ DEFAULT now(),
  UNIQUE (node_id, report_date)
);
```

## REST API (rencana)
| Method | Path | Fungsi |
|---|---|---|
| GET | `/api/v1/nodes/{id}/metrics?from&to` | Riwayat metrik |
| GET | `/api/v1/nodes/{id}/reports/latest` | Laporan terbaru |
| GET | `/api/v1/nodes/{id}/reports?date=` | Laporan per tanggal |
| GET | `/healthz` | Health check |

## Pipeline Laporan
Cron (`REPORT_CRON`, default 19:00) → ambil metrik hari tsb → bangun prompt → Gemini → simpan. Idempotent (UNIQUE node+tanggal).

## Prompt LLM (template)
```
Anda asisten pertanian. Tulis ringkasan "Kondisi Tanaman Hari Ini" dalam Bahasa Indonesia, maks 120 kata.
HANYA gunakan angka pada data di bawah; jangan menambah angka/fakta lain. Jika data kurang, katakan tidak tersedia.
Data (JSON): {{metrics_json}}
```
Gunakan temperature rendah (≤0.3). Validasi: semua angka pada narasi harus ada pada input (cek otomatis) → mitigasi halusinasi (acceptance #4).

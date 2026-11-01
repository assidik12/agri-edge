# Roadmap & Tracking

## Fase
- [x] **0. Setup proyek & dokumentasi**
- [x] **1. Pipeline ML scaffold** (data sintetis, MLP, INT8, export .h) — *belum dijalankan/diverifikasi*
- [ ] 1b. Jalankan `train.py`, catat `metrics.json`, plot prediksi vs aktual
- [ ] 1c. Ganti/validasi dengan data lapangan atau acuan pakar
- [ ] **2. Firmware**: scaffold PlatformIO, sensor + kalibrasi, RTC alarm + deep sleep, TFLM inference, pompa, MQTT
- [ ] **3. Backend**: docker-compose (broker+TimescaleDB) [x], ingest [x], repository [x], REST API [x], cron+Gemini [ ]
- [ ] **4. Dashboard** web
- [ ] **5. Pengujian & soak 30 hari**, penulisan skripsi

## Acceptance Criteria (dari PRD)
| # | Kriteria | Status | Bukti |
|---|---|---|---|
| 1 | ESP32 ≥30 hari tanpa crash | ⬜ | log reset reason |
| 2 | MAE < 10% dari durasi maks | ⬜ (sintetis: belum dijalankan) | `ml/models/metrics.json` |
| 3 | `invoke()` < 100 ms di ESP32 | ⬜ | serial log |
| 4 | Laporan Gemini koheren tanpa halusinasi | ⬜ | cek angka otomatis |

## Risiko Terbuka
- Data sintetis ≠ data nyata (validitas ilmiah skripsi).
- Python 3.14 pada `.venv` tidak kompatibel TensorFlow.
- Drift sensor soil/LDR → perlu kalibrasi berkala.

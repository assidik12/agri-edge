# AgriEdge

Sistem irigasi cerdas berbasis **Edge AI (TinyML)** di ESP32 + backend Golang + laporan harian **Generative AI**.
Proyek skripsi S1. Sumber kebutuhan: [docs/PRD.md](docs/PRD.md).

```
ESP32 (sensor→MLP INT8→pompa) ──MQTT──► Broker ──► Go Worker ──► PostgreSQL
                                                      │ cron 19:00
                                                      ▼
                                              Gemini API → laporan narasi → REST API → Dashboard
```

## Status

| Komponen | Folder | Status |
|---|---|---|
| Pipeline ML (train, INT8, export .h) | [ml/](ml/) | ✅ Scaffold siap (data sintetis) |
| Firmware ESP32 | [firmware/](firmware/) | ⬜ Belum dimulai |
| Backend Golang | [backend/](backend/) | ⬜ Belum dimulai |
| Dashboard web | - | ⬜ Belum dimulai |

Detail progres: [docs/ROADMAP.md](docs/ROADMAP.md).

## Peta Dokumentasi (baca sesuai kebutuhan)

| Dokumen | Isi |
|---|---|
| [AGENTS.md](AGENTS.md) | **Mulai di sini bila Anda AI agent**: aturan & guardrail |
| [docs/PRD.md](docs/PRD.md) | Kebutuhan produk asli |
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Arsitektur, alur siklus, keputusan desain |
| [docs/AI_MODEL.md](docs/AI_MODEL.md) | Kontrak model: fitur, normalisasi, kuantisasi, metrik |
| [docs/MQTT_PAYLOAD.md](docs/MQTT_PAYLOAD.md) | Topic & skema JSON |
| [docs/HARDWARE.md](docs/HARDWARE.md) | Komponen, pin, catatan wiring |
| [docs/FIRMWARE.md](docs/FIRMWARE.md) | Aturan firmware & state machine deep sleep |
| [docs/BACKEND.md](docs/BACKEND.md) | Layer Clean Architecture, DB, LLM prompt, REST |
| [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) | Setup lingkungan, perintah, testing |
| [docs/ROADMAP.md](docs/ROADMAP.md) | Task & acceptance criteria tracking |
| [docs/GLOSSARY.md](docs/GLOSSARY.md) | Istilah |

## Quick Start (ML)

Butuh **Python 3.10–3.12** (TensorFlow belum mendukung Python 3.14; `.venv` yang ada memakai 3.14, buat ulang).

```powershell
py -3.11 -m venv .venv
.venv\Scripts\Activate.ps1
pip install -r ml/requirements.txt
cd ml
python src/train.py
```

Output di `ml/models/`: `agriedge_int8.tflite`, `agriedge_model_data.h`, `metrics.json`.

## Struktur Repo

```
agri_edge/
├── README.md  AGENTS.md  .env.example  .gitignore
├── docs/        # semua dokumentasi
├── ml/          # training & ekspor model (src/, data/, models/)
├── firmware/    # ESP32 PlatformIO (belum dibuat)
├── backend/     # Golang service (belum dibuat)
├── legacy/      # skrip lama bertema Smart Grid, REFERENSI SAJA
└── scratch/     # output sementara (plot, eksperimen), di-gitignore
```

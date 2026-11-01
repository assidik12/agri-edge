# Development Guide

## Prasyarat
- Python 3.10–3.12 (untuk TensorFlow), Git
- PlatformIO (firmware), Go ≥1.22 (backend), Docker (broker/DB)

## ML
```powershell
py -3.11 -m venv .venv          # ganti .venv bawaan (Python 3.14 tidak didukung TF)
.venv\Scripts\Activate.ps1
pip install -r ml/requirements.txt
cd ml; python src/train.py
```
Setelah training: salin `ml/models/agriedge_model_data.h` → `firmware/include/`.

## Firmware (rencana)
```powershell
cd firmware; pio run; pio run -t upload; pio device monitor
```

## Backend (rencana)
```powershell
docker compose up -d   # mosquitto + postgres (compose belum dibuat)
cd backend; go run ./cmd/server; go test ./...
```
Uji manual MQTT: `mosquitto_pub -t agriedge/node1/daily_metrics -f sample.json`.

## Konvensi
- Commit: Conventional Commits (`feat:`, `fix:`, `docs:`).
- Python: type hints, docstring, vectorized NumPy. Go: `gofmt`, `golangci-lint`. C++: tanpa `String`.
- Secret hanya di `.env` (gitignored).
- Output eksperimen/plot → `scratch/`.

## Checklist PR
- [ ] Tes lulus  - [ ] Docs sinkron (terutama kontrak model/payload)  - [ ] ROADMAP diperbarui

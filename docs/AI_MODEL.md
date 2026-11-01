# Kontrak Model AI

> Sumber kebenaran kode: [ml/src/config.py](../ml/src/config.py). Dokumen ini harus sinkron.

## Tugas
Regresi: memprediksi durasi penyiraman (detik) dari 4 fitur harian.

## Input (urutan tetap)
| # | Fitur | Satuan | Min | Max | Sumber |
|---|---|---|---|---|---|
| 0 | `soil_moisture_pct` | % | 0 | 100 | Capacitive v1.2 |
| 1 | `temperature_c` | °C | 15 | 45 | DHT22 |
| 2 | `air_humidity_pct` | % | 20 | 100 | DHT22 |
| 3 | `light_accumulation` | ADC sum | 0 | 32760 | LDR, akumulasi jam 08–16 |

Normalisasi: `x_norm = clip((x - min) / (max - min), 0, 1)`.

## Output
`y_norm ∈ [0,1]` (sigmoid) ⇒ `pump_duration_sec = y_norm × 300` (`PUMP_MAX_SEC`).

## Arsitektur
`Input(4) → Dense(16, ReLU) → Dense(8, ReLU) → Dense(1, sigmoid)`; loss MSE, optimizer Adam.

## Kuantisasi
Post-training INT8 penuh dengan representative dataset (200 sampel), input & output `int8`.
- Quantize input: `q = clip(round(x_norm / scale + zp), -128, 127)`
- Dequantize output: `y = (q - zp) × scale`
Scale/zero-point dibaca dari tensor di runtime (jangan hardcode).

## Artefak
| File | Penggunaan |
|---|---|
| `ml/models/agriedge_int8.tflite` | Evaluasi/benchmark Python |
| `ml/models/agriedge_model_data.h` | Di-include firmware (`g_agriedge_model`) |
| `ml/models/metrics.json` | Hasil evaluasi |

## Data
Saat ini **sintetis** ([ml/src/data.py](../ml/src/data.py)): label dari aturan agronomi + noise. Hasil MAE di data sintetis **belum membuktikan** akurasi lapangan. Untuk skripsi: kumpulkan data nyata atau label dari pakar/acuan kebutuhan air tanaman, simpan CSV di `ml/data/` dengan kolom = `FEATURES` + `pump_seconds`.

## Evaluasi (wajib dilaporkan)
MAE (detik & % dari 300 s), model size (KB), latensi `invoke()` **di ESP32**, tensor arena terpakai. Tambahan disarankan: RMSE, plot prediksi vs aktual (simpan ke `scratch/`).

## Acceptance
MAE < 10% dari PUMP_MAX_SEC; latensi invoke < 100 ms di ESP32.

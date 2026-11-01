# AGENTS.md — Panduan untuk AI Agent & Developer

Baca file ini dulu, lalu dokumen di [docs/](docs/) sesuai area kerja. Bahasa dokumen: Indonesia; kode & identifier: Inggris.

## Konteks Singkat
AgriEdge = ESP32 menjalankan MLP INT8 (TFLite Micro) untuk memutuskan durasi pompa; hasil dikirim via MQTT ke backend Go → PostgreSQL → laporan Gemini. Detail: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Sumber Kebenaran (jangan menduga, rujuk ini)
| Hal | File |
|---|---|
| Fitur, rentang normalisasi, PUMP_MAX_SEC | [ml/src/config.py](ml/src/config.py) |
| Skema MQTT | [docs/MQTT_PAYLOAD.md](docs/MQTT_PAYLOAD.md) |
| Kebutuhan & acceptance criteria | [docs/PRD.md](docs/PRD.md) |

## Aturan Wajib (Guardrails)
1. **Kontrak model**: urutan fitur `[soil_moisture_pct, temperature_c, air_humidity_pct, light_accumulation]`. Mengubah fitur/rentang/PUMP_MAX_SEC ⇒ update `ml/src/config.py`, `firmware/include/model_config.h`, `docs/AI_MODEL.md`, `docs/MQTT_PAYLOAD.md`, lalu latih ulang. Jangan ubah sebagian.
2. **Firmware**: dilarang `String` Arduino di loop/path berulang, hindari `malloc/new` dinamis setelah setup (heap fragmentation). Tensor arena statis. Lihat [docs/FIRMWARE.md](docs/FIRMWARE.md).
3. **Edge autonomy**: keputusan irigasi tidak boleh bergantung pada internet/NTP (waktu dari RTC DS3231). Kegagalan MQTT tidak boleh menghalangi penyiraman.
4. **Safety pompa**: durasi selalu di-clamp `[0, PUMP_MAX_SEC]`; ada timeout keras di firmware.
5. **LLM**: prompt hanya memakai angka dari DB; narasi tidak boleh mengarang data teknis (acceptance #4). Jangan commit API key; gunakan `.env` (lihat `.env.example`).
6. **Python**: gunakan `.venv` (Python 3.10–3.12), operasi vectorized NumPy. Laporkan metrik evaluasi (MAE, dll.) setiap mengubah model; simpan plot di `scratch/`.
7. **Backend**: Clean Architecture (domain → usecase → adapter → delivery); domain tidak boleh meng-import library eksternal.
8. `legacy/` hanya referensi (tema Smart Grid); jangan diimpor. File `ml/models/*.h|*.tflite` adalah artefak generate, jangan diedit manual.
9. Pertahankan komentar/docstring yang tidak terkait perubahan Anda.

## Alur Kerja Perubahan
1. Cek [docs/ROADMAP.md](docs/ROADMAP.md) untuk task aktif → 2. baca dokumen area terkait → 3. implementasi + tes → 4. update docs & centang ROADMAP.

## Perintah Cepat
```powershell
cd ml; python src/train.py          # latih + ekspor + metrics.json
```
(Perintah firmware/backend: lihat [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md).)

## Definition of Done
Kode jalan, tes lulus, acceptance criteria terkait dicek, dokumen ter-update, tidak ada secret ter-commit.

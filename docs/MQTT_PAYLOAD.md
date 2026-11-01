# MQTT Payload

- **Topic**: `agriedge/<node_id>/daily_metrics` (contoh `agriedge/node1/daily_metrics`)
- **QoS**: 1, **retain**: false
- **Frekuensi**: 1×/hari (±17:00 lokal)
- Backend subscribe: `agriedge/+/daily_metrics`

## Skema
```json
{
  "timestamp": "2026-10-24T17:00:00Z",
  "metrics": {
    "soil_moisture_pct": 42.5,
    "temperature_c": 33.2,
    "air_humidity_pct": 58.0,
    "light_accumulation": 8950
  },
  "ai_decision": {
    "pump_duration_sec": 150,
    "inference_time_ms": 12
  }
}
```

| Field | Tipe | Aturan validasi |
|---|---|---|
| `timestamp` | string ISO-8601 | wajib; dari RTC |
| `metrics.soil_moisture_pct` | float | 0–100 |
| `metrics.temperature_c` | float | -10–60 |
| `metrics.air_humidity_pct` | float | 0–100 |
| `metrics.light_accumulation` | int | ≥ 0 |
| `ai_decision.pump_duration_sec` | int | 0–300 |
| `ai_decision.inference_time_ms` | int | ≥ 0 |

Payload yang gagal validasi: backend menolak & mencatat log, tidak menyimpan.
Catatan: `timestamp` ditulis "Z" (UTC) di PRD; pastikan RTC diset UTC atau gunakan offset eksplisit (+07:00) secara konsisten. Firmware memformat dengan `snprintf` ke buffer statis (tanpa `String`).

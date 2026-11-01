"""Single source of truth untuk kontrak model AgriEdge.

PERHATIAN: Nilai di sini HARUS identik dengan firmware (firmware/include/model_config.h)
dan dokumen docs/AI_MODEL.md. Ubah ketiganya bersamaan.
"""
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DATA_DIR = ROOT / "data"
MODELS_DIR = ROOT / "models"

# Urutan fitur input (JANGAN diubah urutannya; firmware mengikuti urutan ini)
FEATURES = ["soil_moisture_pct", "temperature_c", "air_humidity_pct", "light_accumulation"]

# Rentang min/max untuk normalisasi ke [0.0, 1.0]:  x_norm = (x - min) / (max - min)
FEATURE_RANGES = {
    "soil_moisture_pct": (0.0, 100.0),
    "temperature_c": (15.0, 45.0),
    "air_humidity_pct": (20.0, 100.0),
    "light_accumulation": (0.0, 32760.0),  # 9 sampel/hari x ADC 12-bit (4095) ~ 36855; dibatasi
}

# Output: durasi pompa (detik). Target dinormalisasi: y_norm = y / PUMP_MAX_SEC
PUMP_MAX_SEC = 300.0

# Acceptance criteria (docs/PRD.md bagian 8)
MAE_TOLERANCE_PCT = 10.0   # MAE < 10% dari PUMP_MAX_SEC
MAX_INFERENCE_MS = 100.0

SEED = 42

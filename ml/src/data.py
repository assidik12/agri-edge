"""Dataset sintetis irigasi + normalisasi.

Dataset sintetis dipakai sampai data lapangan tersedia. Label dihasilkan dari aturan
agronomi sederhana + noise, sehingga model belajar fungsi yang dapat dijelaskan.
"""
import numpy as np

from config import FEATURES, FEATURE_RANGES, PUMP_MAX_SEC, SEED


def normalize(raw: np.ndarray) -> np.ndarray:
    """raw: (N, 4) nilai fisik -> (N, 4) float32 pada [0, 1]."""
    out = np.empty_like(raw, dtype=np.float32)
    for i, name in enumerate(FEATURES):
        lo, hi = FEATURE_RANGES[name]
        out[:, i] = np.clip((raw[:, i] - lo) / (hi - lo), 0.0, 1.0)
    return out


def _rule_pump_seconds(raw: np.ndarray, rng: np.random.Generator) -> np.ndarray:
    soil, temp, hum, light = raw.T
    # Kebutuhan air naik jika tanah kering, panas, udara kering, dan banyak cahaya.
    need = (
        0.55 * (1.0 - soil / 100.0)
        + 0.15 * np.clip((temp - 15.0) / 30.0, 0, 1)
        + 0.10 * (1.0 - np.clip((hum - 20.0) / 80.0, 0, 1))
        + 0.20 * np.clip(light / 32760.0, 0, 1)
    )
    need = np.where(soil > 80.0, 0.0, need)  # tanah sudah basah -> jangan siram
    sec = need * PUMP_MAX_SEC + rng.normal(0, 5.0, size=need.shape)
    return np.clip(sec, 0.0, PUMP_MAX_SEC)


def generate_synthetic(n: int = 5000, seed: int = SEED) -> tuple[np.ndarray, np.ndarray]:
    """Return (raw_features (N,4), pump_seconds (N,))."""
    rng = np.random.default_rng(seed)
    raw = np.column_stack([
        rng.uniform(10, 95, n),        # soil_moisture_pct
        rng.uniform(22, 40, n),        # temperature_c
        rng.uniform(35, 95, n),        # air_humidity_pct
        rng.uniform(1000, 30000, n),   # light_accumulation
    ])
    return raw, _rule_pump_seconds(raw, rng)


def make_splits(n: int = 5000, seed: int = SEED):
    """Return dict train/val/test berisi (X_norm, y_sec)."""
    raw, sec = generate_synthetic(n, seed)
    x = normalize(raw)
    idx = np.random.default_rng(seed).permutation(n)
    a, b = int(0.7 * n), int(0.85 * n)
    parts = {"train": idx[:a], "val": idx[a:b], "test": idx[b:]}
    return {k: (x[v], sec[v].astype(np.float32)) for k, v in parts.items()}

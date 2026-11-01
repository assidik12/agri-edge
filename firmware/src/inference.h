/**
 * @file inference.h
 * @brief Header modul inferensi TensorFlow Lite for Microcontrollers (TFLM) di ESP32.
 * 
 * Aturan:
 * - Tanpa dynamic allocation (malloc/new) setelah inisialisasi.
 * - Tensor arena statis.
 * - Zero String Arduino.
 */

#pragma once

#ifndef AGRIEDGE_INFERENCE_H
#define AGRIEDGE_INFERENCE_H

#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Struktur input data sensor mentah (sebelum dinormalisasi)
 */
typedef struct {
    float soil_moisture_pct;  // 0.0 - 100.0 %
    float temperature_c;      // 15.0 - 45.0 °C
    float air_humidity_pct;   // 20.0 - 100.0 %
    float light_accumulation; // 0.0 - 32760.0
} SensorFeatures;

/**
 * @brief Struktur hasil keputusan inferensi TinyML
 */
typedef struct {
    float pump_duration_sec;  // Durasi pompa rekomendasi (detik, clamped [0, 300])
    uint32_t latency_ms;      // Waktu eksekusi inferensi (milidetik)
    bool is_valid;            // Status keberhasilan inferensi
} InferenceResult;

/**
 * @brief Inisialisasi TFLM interpreter, tensor arena, dan registrasi model INT8.
 * Dipanggil sekali di setup().
 * 
 * @return true jika model berhasil dimuat dan dialokasikan, false jika gagal.
 */
bool inference_init(void);

/**
 * @brief Menjalankan inferensi model TinyML INT8 berdasarkan data sensor.
 * 
 * Melakukan:
 * 1. Normalisasi fitur mentah ke [0.0, 1.0]
 * 2. Kuantisasi float input ke int8 sesuai scale & zero_point tensor model
 * 3. Eksekusi invoke()
 * 4. De-kuantisasi output int8 kembali ke float
 * 5. Denormalisasi ke durasi detik pompa ([0, PUMP_MAX_SEC])
 * 
 * @param input Pointer ke data sensor
 * @param result Pointer ke struct hasil inferensi
 * @return true jika berhasil, false jika gagal
 */
bool inference_run(const SensorFeatures* input, InferenceResult* result);

#ifdef __cplusplus
}
#endif

#endif // AGRIEDGE_INFERENCE_H

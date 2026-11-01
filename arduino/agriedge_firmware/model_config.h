/**
 * @file model_config.h
 * @brief Konfigurasi rentang normalisasi fitur dan batas pompa (AgriEdge TinyML).
 * 
 * PERHATIAN: Nilai di file ini HARUS selalu sinkron dengan ml/src/config.py
 * dan dokumen docs/AI_MODEL.md.
 */

#pragma once

#ifndef MODEL_CONFIG_H
#define MODEL_CONFIG_H

// Jumlah fitur input ke model MLP
#define MODEL_NUM_INPUTS 4

// Durasi maksimum pompa (detik) untuk scaling output sigmoid [0.0, 1.0]
#define PUMP_MAX_SEC 300.0f

// Urutan dan batas normalisasi fitur: x_norm = (x - min) / (max - min)
// 0: Soil Moisture (%)
#define SOIL_MOISTURE_MIN 0.0f
#define SOIL_MOISTURE_MAX 100.0f

// 1: Temperature (°C)
#define TEMP_C_MIN 15.0f
#define TEMP_C_MAX 45.0f

// 2: Air Humidity (%)
#define AIR_HUMIDITY_MIN 20.0f
#define AIR_HUMIDITY_MAX 100.0f

// 3: Light Accumulation (ADC sum)
#define LIGHT_ACCUM_MIN 0.0f
#define LIGHT_ACCUM_MAX 32760.0f

#endif // MODEL_CONFIG_H

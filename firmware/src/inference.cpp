/**
 * @file inference.cpp
 * @brief Implementasi inferensi TinyML INT8 untuk ESP32 menggunakan TensorFlow Lite for Microcontrollers.
 */

#include "inference.h"
#include "model_config.h"
#include "agriedge_model_data.h"

#include <Arduino.h>
#include <math.h>

// Include TensorFlow Lite Micro headers
#include "tensorflow/lite/micro/all_ops_resolver.h"
#include "tensorflow/lite/micro/micro_interpreter.h"
#include "tensorflow/lite/schema/schema_generated.h"

namespace {
    // Alokasi Tensor Arena secara statis di BSS (menghindari heap fragmentation)
    // Model INT8 hanya ~3.2 KB, arena 4096 bytes sudah cukup aman.
    constexpr int kTensorArenaSize = 4096;
    alignas(16) uint8_t tensor_arena[kTensorArenaSize];

    const tflite::Model* model = nullptr;
    tflite::MicroInterpreter* interpreter = nullptr;
    TfLiteTensor* input_tensor = nullptr;
    TfLiteTensor* output_tensor = nullptr;

    // Helper clip float
    inline float clamp_float(float val, float min_val, float max_val) {
        if (val < min_val) return min_val;
        if (val > max_val) return max_val;
        return val;
    }

    // Helper normalisasi ke [0.0, 1.0]
    inline float normalize_val(float val, float min_val, float max_val) {
        float clipped = clamp_float(val, min_val, max_val);
        return (clipped - min_val) / (max_val - min_val);
    }
}

bool inference_init(void) {
    if (interpreter != nullptr) {
        return true; // Sudah diinisialisasi
    }

    // 1. Validasi model data
    model = tflite::GetModel(g_agriedge_model);
    if (model->version() != TFLITE_SCHEMA_VERSION) {
        return false;
    }

    // 2. Resolver operasi TFLite
    // Menggunakan AllOpsResolver untuk kepraktisan, atau MicroMutableOpResolver untuk menghemat flash
    static tflite::AllOpsResolver resolver;

    // 3. Bangun MicroInterpreter
    static tflite::MicroInterpreter static_interpreter(
        model, resolver, tensor_arena, kTensorArenaSize);
    interpreter = &static_interpreter;

    // 4. Alokasikan tensor ke tensor arena
    TfLiteStatus allocate_status = interpreter->AllocateTensors();
    if (allocate_status != kTfLiteOk) {
        return false;
    }

    // 5. Simpan pointer tensor input & output
    input_tensor = interpreter->input(0);
    output_tensor = interpreter->output(0);

    if (input_tensor == nullptr || output_tensor == nullptr) {
        return false;
    }

    // Verifikasi tipe input/output adalah INT8
    if (input_tensor->type != kTfLiteInt8 || output_tensor->type != kTfLiteInt8) {
        return false;
    }

    return true;
}

bool inference_run(const SensorFeatures* input, InferenceResult* result) {
    if (interpreter == nullptr || input == nullptr || result == nullptr) {
        if (result != nullptr) {
            result->is_valid = false;
            result->pump_duration_sec = 0.0f;
            result->latency_ms = 0;
        }
        return false;
    }

    // 1. Normalisasi data fisik ke rentang [0.0, 1.0] sesuai urutan kontrak model
    float norm_features[MODEL_NUM_INPUTS];
    norm_features[0] = normalize_val(input->soil_moisture_pct, SOIL_MOISTURE_MIN, SOIL_MOISTURE_MAX);
    norm_features[1] = normalize_val(input->temperature_c, TEMP_C_MIN, TEMP_C_MAX);
    norm_features[2] = normalize_val(input->air_humidity_pct, AIR_HUMIDITY_MIN, AIR_HUMIDITY_MAX);
    norm_features[3] = normalize_val(input->light_accumulation, LIGHT_ACCUM_MIN, LIGHT_ACCUM_MAX);

    // 2. Kuantisasi Float -> INT8 sesuai parameter tensor input
    // Rumus: q = round(float_val / scale) + zero_point, diclamp [-128, 127]
    const float in_scale = input_tensor->params.scale;
    const int32_t in_zero_point = input_tensor->params.zero_point;
    int8_t* in_data = input_tensor->data.int8;

    for (int i = 0; i < MODEL_NUM_INPUTS; ++i) {
        int32_t q = static_cast<int32_t>(roundf(norm_features[i] / in_scale)) + in_zero_point;
        if (q < -128) q = -128;
        if (q > 127) q = 127;
        in_data[i] = static_cast<int8_t>(q);
    }

    // 3. Jalankan inferensi & catat waktu eksekusi (latensi)
    uint32_t start_time = millis();
    TfLiteStatus invoke_status = interpreter->Invoke();
    uint32_t latency = millis() - start_time;

    if (invoke_status != kTfLiteOk) {
        result->is_valid = false;
        result->pump_duration_sec = 0.0f;
        result->latency_ms = latency;
        return false;
    }

    // 4. De-kuantisasi output model INT8 -> Float [0.0, 1.0]
    // Rumus: float_val = (q - zero_point) * scale
    const float out_scale = output_tensor->params.scale;
    const int32_t out_zero_point = output_tensor->params.zero_point;
    int8_t out_quant = output_tensor->data.int8[0];

    float y_norm = static_cast<float>(out_quant - out_zero_point) * out_scale;
    y_norm = clamp_float(y_norm, 0.0f, 1.0f);

    // 5. Konversi ke durasi pompa (detik) dan clamp batas keamanan [0, PUMP_MAX_SEC]
    float duration_sec = y_norm * PUMP_MAX_SEC;
    duration_sec = clamp_float(duration_sec, 0.0f, PUMP_MAX_SEC);

    result->pump_duration_sec = duration_sec;
    result->latency_ms = latency;
    result->is_valid = true;

    return true;
}

/**
 * @file agriedge_firmware.ino
 * @brief Sketch Arduino IDE untuk inferensi TinyML INT8 di ESP32 Dev Module
 */

#include <Arduino.h>
#include <TensorFlowLite_ESP32.h>
#include "tensorflow/lite/micro/all_ops_resolver.h"
#include "tensorflow/lite/micro/micro_error_reporter.h"
#include "tensorflow/lite/micro/micro_interpreter.h"
#include "tensorflow/lite/schema/schema_generated.h"

#include "model_config.h"
#include "agriedge_model_data.h"

// Konfigurasi Memori Model
#define NUMBER_OF_INPUTS 4
#define NUMBER_OF_OUTPUTS 1
constexpr int kTensorArenaSize = 8 * 1024;
uint8_t tensor_arena[kTensorArenaSize];

// Komponen TFLite Micro
tflite::ErrorReporter* error_reporter = nullptr;
const tflite::Model* model = nullptr;
tflite::MicroInterpreter* interpreter = nullptr;
TfLiteTensor* input = nullptr;
TfLiteTensor* output = nullptr;

// Helper clamp & normalisasi
inline float clamp_val(float val, float min_val, float max_val) {
    if (val < min_val) return min_val;
    if (val > max_val) return max_val;
    return val;
}

inline float normalize_val(float val, float min_val, float max_val) {
    float clipped = clamp_val(val, min_val, max_val);
    return (clipped - min_val) / (max_val - min_val);
}

// Fungsi inferensi model TinyML
float predict_pump_duration(float soil_moisture, float temp, float humidity, float light_accum, uint32_t &latency_ms) {
    // 1. Normalisasi data fisik
    float f_inputs[4];
    f_inputs[0] = normalize_val(soil_moisture, SOIL_MOISTURE_MIN, SOIL_MOISTURE_MAX);
    f_inputs[1] = normalize_val(temp, TEMP_C_MIN, TEMP_C_MAX);
    f_inputs[2] = normalize_val(humidity, AIR_HUMIDITY_MIN, AIR_HUMIDITY_MAX);
    f_inputs[3] = normalize_val(light_accum, LIGHT_ACCUM_MIN, LIGHT_ACCUM_MAX);

    // 2. Tulis data ke input tensor (sesuai tipe data model float atau quantized int8)
    if (input->type == kTfLiteInt8) {
        for (int i = 0; i < 4; i++) {
            input->data.int8[i] = (int8_t)(f_inputs[i] / input->params.scale + input->params.zero_point);
        }
    } else {
        for (int i = 0; i < 4; i++) {
            input->data.f[i] = f_inputs[i];
        }
    }

    // 3. Eksekusi inferensi & ukur latensi
    uint32_t start_time = millis();
    TfLiteStatus invoke_status = interpreter->Invoke();
    latency_ms = millis() - start_time;

    if (invoke_status != kTfLiteOk) {
        Serial.println("[ERROR] Gagal menjalankan Invoke()");
        return 0.0f;
    }

    // 4. Ambil output model
    float y_norm = 0.0f;
    if (output->type == kTfLiteInt8) {
        y_norm = (output->data.int8[0] - output->params.zero_point) * output->params.scale;
    } else {
        y_norm = output->data.f[0];
    }
    y_norm = clamp_val(y_norm, 0.0f, 1.0f);

    // 5. Denormalisasi ke durasi pompa
    float duration_sec = y_norm * PUMP_MAX_SEC;
    return clamp_val(duration_sec, 0.0f, PUMP_MAX_SEC);
}

void setup() {
    Serial.begin(115200);
    delay(2000);

    Serial.println("\n==============================================");
    Serial.println("  AgriEdge TinyML - Native TFLM ESP32         ");
    Serial.println("==============================================");

    static tflite::MicroErrorReporter micro_error_reporter;
    error_reporter = &micro_error_reporter;

    // Load model dari array agriedge_model_data.h
    model = tflite::GetModel(g_agriedge_model);
    if (model->version() != TFLITE_SCHEMA_VERSION) {
        Serial.println("[ERROR] Versi skema model tidak cocok!");
        while (1) { delay(1000); }
    }

    static tflite::AllOpsResolver resolver;
    static tflite::MicroInterpreter static_interpreter(
        model, resolver, tensor_arena, kTensorArenaSize, error_reporter);
    interpreter = &static_interpreter;

    if (interpreter->AllocateTensors() != kTfLiteOk) {
        Serial.println("[ERROR] AllocateTensors() gagal! Perbesar kTensorArenaSize.");
        while (1) { delay(1000); }
    }

    input = interpreter->input(0);
    output = interpreter->output(0);
    Serial.println("[OK] Model TinyML berhasil dimuat dan siap inferensi!");

    // Skenario Uji 1: Tanah Kering & Cuaca Terik
    Serial.println("\n--- Skenario 1: Tanah Kering & Panas ---");
    uint32_t lat1 = 0;
    float pump1 = predict_pump_duration(20.0f, 36.5f, 35.0f, 25000.0f, lat1);
    Serial.printf("Input: Soil=20%%, Temp=36.5C, Hum=35%%, Light=25000\n");
    Serial.printf("Keputusan Pompa : %.2f detik\n", pump1);
    Serial.printf("Latensi Eksekusi: %u ms\n", lat1);

    // Skenario Uji 2: Tanah Basah & Udara Lembap
    Serial.println("\n--- Skenario 2: Tanah Basah & Sejuk ---");
    uint32_t lat2 = 0;
    float pump2 = predict_pump_duration(85.0f, 23.0f, 80.0f, 5000.0f, lat2);
    Serial.printf("Input: Soil=85%%, Temp=23.0C, Hum=80%%, Light=5000\n");
    Serial.printf("Keputusan Pompa : %.2f detik\n", pump2);
    Serial.printf("Latensi Eksekusi: %u ms\n", lat2);

    Serial.println("\nUji coba selesai.");
}

void loop() {
    delay(10000);
}
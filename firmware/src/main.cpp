/**
 * @file main.cpp
 * @brief Contoh pengujian modul inferensi TinyML AgriEdge pada ESP32.
 */

#include <Arduino.h>
#include "inference.h"

void setup() {
    Serial.begin(115200);
    delay(2000);

    Serial.println("\n==========================================");
    Serial.println("   AgriEdge TinyML - ESP32 Inference Demo ");
    Serial.println("==========================================");

    // Inisialisasi TFLM Interpreter & model
    if (!inference_init()) {
        Serial.println("[ERROR] Gagal menginisialisasi TFLM Interpreter!");
        while (1) { delay(1000); }
    }
    Serial.println("[OK] TFLM Model INT8 berhasil dimuat.");

    // Skenario Uji 1: Kondisi Kering & Panas (Seharusnya durasi pompa tinggi)
    SensorFeatures test_hot_dry = {
        .soil_moisture_pct = 20.0f,
        .temperature_c = 36.5f,
        .air_humidity_pct = 35.0f,
        .light_accumulation = 25000.0f
    };

    // Skenario Uji 2: Kondisi Lembap & Sejuk (Seharusnya pompa 0 atau sangat kecil)
    SensorFeatures test_wet_cool = {
        .soil_moisture_pct = 85.0f,
        .temperature_c = 23.0f,
        .air_humidity_pct = 80.0f,
        .light_accumulation = 5000.0f
    };

    InferenceResult res;

    // Eksekusi Uji 1
    Serial.println("\n--- Skenario 1 (Tanah Kering & Terik) ---");
    if (inference_run(&test_hot_dry, &res)) {
        Serial.printf("Keputusan Pompa : %.2f detik\n", res.pump_duration_sec);
        Serial.printf("Latensi Inferensi: %u ms\n", res.latency_ms);
    } else {
        Serial.println("[ERROR] Inferensi gagal.");
    }

    // Eksekusi Uji 2
    Serial.println("\n--- Skenario 2 (Tanah Basah & Sejuk) ---");
    if (inference_run(&test_wet_cool, &res)) {
        Serial.printf("Keputusan Pompa : %.2f detik\n", res.pump_duration_sec);
        Serial.printf("Latensi Inferensi: %u ms\n", res.latency_ms);
    } else {
        Serial.println("[ERROR] Inferensi gagal.");
    }

    Serial.println("\nUji inferensi selesai.");
}

void loop() {
    // Demo hanya berjalan sekali di setup()
    delay(10000);
}

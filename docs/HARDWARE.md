# Hardware

| Komponen | Antarmuka | Pin ESP32 (usulan) | Catatan |
|---|---|---|---|
| ESP32 DevKit | - | - | Dual-core, TFLM |
| Capacitive Soil Moisture v1.2 | ADC1 | GPIO34 | Pakai **ADC1** (ADC2 bentrok dengan Wi-Fi). Lapisi anti-air. Kalibrasi nilai kering/basah |
| DHT22 | 1-wire | GPIO4 | Pull-up 10 kΩ; baca min. interval 2 s |
| LDR + resistor | ADC1 | GPIO35 | Voltage divider; sampel tiap jam, jumlahkan |
| RTC DS3231 | I²C | SDA 21 / SCL 22 | INT/SQW → GPIO33 (RTC GPIO) untuk wake `ext0` |
| Relay module | Digital | GPIO26 | Cek active-low/high; hold state saat sleep dengan `gpio_hold_en` |
| Pompa | via relay | - | Catu daya terpisah; dioda flyback jika beban induktif |

> Pin di atas **usulan**; sesuaikan dan perbarui tabel ini saat wiring final. Pin wake-up harus RTC-capable (mis. 32–39).

## Kalibrasi
- Soil: catat ADC di udara (kering) & air (basah) → petakan ke 0–100%.
- LDR: catat rentang ADC siang/malam; pastikan `light_accumulation` ≤ 32760 (rentang model).

## Keselamatan
Isolasi relay dari ESP32, pompa tidak boleh menyala saat boot/reset (pin default LOW/aman), batas durasi keras di firmware.

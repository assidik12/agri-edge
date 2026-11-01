# Glosarium

| Istilah | Arti |
|---|---|
| Edge AI / TinyML | Inferensi model ML langsung di mikrokontroler |
| TFLM | TensorFlow Lite for Microcontrollers |
| MLP | Multi-Layer Perceptron, jaringan saraf feed-forward |
| INT8 quantization | Bobot/aktivasi 8-bit agar kecil & cepat |
| Tensor arena | Buffer RAM statis untuk tensor TFLM |
| Deep Sleep | Mode hemat daya ESP32; RAM biasa hilang, RTC memory bertahan |
| `RTC_DATA_ATTR` | Variabel yang bertahan lintas deep sleep |
| Heap fragmentation | Memori terpecah akibat alokasi dinamis berulang → crash |
| MAE | Mean Absolute Error |
| light_accumulation | Jumlah pembacaan LDR 08:00–16:00 |
| Node | Satu unit ESP32 di lahan (`node_id`) |

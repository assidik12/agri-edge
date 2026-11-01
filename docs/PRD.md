# **PRODUCT REQUIREMENTS DOCUMENT (PRD)**

**Nama Proyek:** AgriEdge \- Sistem Irigasi Cerdas Berbasis Edge AI (TinyML) & Generative Reporting  
**Platform:** Embedded System (ESP32), Backend API (Golang), Web Dashboard  
**Status:** Draf Perencanaan Skripsi

## **1\. Ringkasan Eksekutif (Executive Summary)**

AgriEdge adalah ekosistem *Smart Agriculture* yang memindahkan beban komputasi pengambilan keputusan irigasi dari *cloud* ke mikrokontroler (ESP32) menggunakan teknologi Edge AI (TensorFlow Lite for Microcontrollers). Sistem ini menggunakan model *Multi-Layer Perceptron* (MLP) untuk memprediksi durasi penyiraman optimal berdasarkan akumulasi suhu, kelembapan udara, kelembapan tanah, dan intensitas cahaya harian. Selanjutnya, *backend* akan memproses metrik yang dikirim via MQTT untuk menghasilkan laporan analitik harian otomatis menggunakan *Generative AI* (LLM).

## **2\. Tujuan & Sasaran (Goals & Objectives)**

* **Edge Autonomy:** Sistem irigasi dapat mengambil keputusan secara lokal dan presisi tanpa bergantung pada stabilitas koneksi internet.  
* **Efisiensi Daya & Memori:** Menghindari *heap fragmentation* dan menghemat daya baterai mikrokontroler dengan mengimplementasikan siklus *Deep Sleep* terstruktur.  
* **AI-Driven Insight:** Menerjemahkan angka sensor mentah menjadi laporan harian berformat narasi (teks manusiawi) menggunakan *Generative AI* di sisi *backend*.  
* **Standar Industri Skripsi S1:** Membuktikan kompetensi arsitektur terdistribusi (IoT C++ \-\> MQTT \-\> Golang Backend \-\> Database \-\> GenAI).

## **3\. Arsitektur Sistem (System Architecture)**

Sistem dibagi menjadi dua domain komputasi utama:

* **Edge Node (Kebun/Sawah):** ESP32 bertindak sebagai *inference engine* (menjalankan .tflite), penjadwal harian (*Deep Sleep* via RTC), pembaca sensor, aktuator relai, dan MQTT *Publisher*.  
* **Backend & Cloud (Server):** *Message Broker* (Mosquitto/EMQX) menerima data. Servis Golang (mengadopsi prinsip *Clean Architecture* dari proyek Catalyst Anda) menampung data ke PostgreSQL, memicu API LLM (misalnya Gemini API) untuk membuat laporan harian, dan menyediakan *endpoint* REST API untuk *dashboard*.

## **4\. Kebutuhan Perangkat Keras (Hardware Requirements)**

| Komponen | Fungsi & Justifikasi |
| :---- | :---- |
| **ESP32 DevKit** | MCU utama dengan dual-core dan memori memadai untuk TFLite. |
| **Capacitive Soil Moisture v1.2** | Mengukur kelembapan tanah. Wajib diberi lapisan pelindung anti-korosi/air pada bagian komponen atas. |
| **DHT22** | Membaca suhu dan kelembapan udara lingkungan secara presisi. |
| **Sensor Cahaya LDR** | Menghitung akumulasi intensitas radiasi matahari harian. |
| **RTC DS3231** | *Trigger* waktu absolut agar jadwal komputasi dan irigasi harian tidak bergantung pada koneksi NTP/Internet. |
| **Relay Module & Pompa** | Aktuator fisik untuk menyalurkan air berdasarkan *output* inferensi AI. |

## **5\. Kebutuhan Perangkat Lunak & AI (Software & AI Requirements)**

* **Model Edge AI:** *Multi-Layer Perceptron* (MLP) dangkal, dilatih di Python (Keras), dikonversi ke TensorFlow Lite Int8, dan di-ekspor sebagai C++ byte array (.h).  
* **Firmware ESP32:** C++ dengan *framework* Arduino/PlatformIO. Wajib tanpa objek String pada fungsi *looping* demi keamanan memori.  
* **Backend Server:** Golang, MQTT Broker, Redis (opsional untuk *caching* status), dan PostgreSQL.  
* **LLM Integration:** Prompt *Engineering* yang dirancang untuk menerima payload JSON metrik pertanian dan memuntahkan narasi "Kondisi Tanaman Hari Ini".

## **6\. Alur Sistem (User & System Flow)**

**Siklus Siang (Pengumpulan Data):**

> 1. Setiap 1 jam (08:00 \- 16:00), RTC membangunkan ESP32 dari *Deep Sleep*.  
> 2. ESP32 membaca nilai LDR, menambahkannya ke variabel light\_accumulated di RTC Memory.  
> 3. ESP32 kembali *Deep Sleep*.

**Siklus Sore (Inferensi & Aktuasi):**

> 1. Pukul 17:00, RTC membangunkan ESP32.  
> 2. ESP32 membaca suhu (DHT22), kelembapan udara (DHT22), dan kelembapan tanah (Capacitive).  
> 3. MCU menormalisasi ke-4 data (termasuk total *light\_accumulated*) menjadi *array tensor* \[0.0 \- 1.0\].  
> 4. **Edge Inference:** TFLite memproses *tensor* dan memuntahkan *output* (misal: 150 detik).  
> 5. ESP32 mengaktifkan relai pompa selama 150 detik.  
> 6. ESP32 memformat data metrik dan durasi pompa ke JSON, lalu mengirimnya ke *Broker* MQTT.  
> 7. Variabel light\_accumulated direset, lalu ESP32 *Deep Sleep* hingga esok hari.

**Siklus Backend (Reporting):**

> 1. Golang *Worker* menerima JSON MQTT dan menyimpannya di PostgreSQL.  
> 2. Malam hari (misal pukul 19:00), *Cron Job* di Golang mengambil metrik harian.  
> 3. Sistem mengirim instruksi ke API Gemini: *"Buat ringkasan pertanian hari ini dari data berikut..."*  
> 4. Narasi laporan tersimpan dan siap dibaca pengguna via *Dashboard* web/mobile.

## **7\. Struktur Data (Payload MQTT)**

*Topic:* agriedge/node1/daily\_metrics

JSON  
{  
  "timestamp": "2026-10-24T17:00:00Z",  
  "metrics": {  
    "soil\_moisture\_pct": 42.5,  
    "temperature\_c": 33.2,  
    "air\_humidity\_pct": 58.0,  
    "light\_accumulation": 8950  
  },  
  "ai\_decision": {  
    "pump\_duration\_sec": 150,  
    "inference\_time\_ms": 12   
  }  
}

## **8\. Kriteria Penerimaan (Acceptance Criteria)**

> 1. **Edge:** ESP32 dapat beroperasi minimal 30 hari tanpa *crash/restart* berkat penggunaan *Deep Sleep* dan manajemen RAM tanpa String.  
> 2. **AI Accuracy:** *Mean Absolute Error* (MAE) durasi pompa pada dataset *testing* berada di bawah 10% toleransi durasi manual.  
> 3. **Latency:** Proses inferensi (invoke()) pada perangkat keras ESP32 selesai di bawah 100 milidetik.  
> 4. **Reporting:** *Backend* berhasil mengonsumsi MQTT *payload* dan men-generate laporan teks dari Gemini API yang koheren tanpa halusinasi data teknis.
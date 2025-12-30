# 🔄 Sync Mode Guide - Offline-First dengan Auto-Sync

## 📋 Apa itu Sync Mode?

Sync Mode adalah fitur yang memungkinkan aplikasi desktop Ritel-App untuk:
- ✅ **Bekerja OFFLINE** tanpa koneksi internet
- ✅ **Auto-sync** ke server ketika koneksi tersedia
- ✅ **Queue operations** selama offline
- ✅ **Automatic recovery** ketika kembali online

---

## 🎯 Cara Kerja

```
┌─────────────────────────────────────┐
│     DESKTOP APP (Ritel-App.exe)    │
│                                     │
│  ┌──────────────────────────────┐  │
│  │  SQLite Lokal (./ritel.db)   │  │ ← SELALU DIGUNAKAN (Cepat)
│  │  - Produk                     │  │
│  │  - Transaksi                  │  │
│  │  - Pelanggan                  │  │
│  └──────────────────────────────┘  │
│              ↓ ↑                    │
│  ┌──────────────────────────────┐  │
│  │     SYNC ENGINE              │  │
│  │  - Monitor koneksi           │  │
│  │  - Queue perubahan           │  │
│  │  - Auto-sync setiap 10 detik │  │
│  └──────────────────────────────┘  │
│              ↓ ↑                    │
└──────────────┼─┼────────────────────┘
               │ │
     Internet  │ │ (Optional)
               ↓ ↑
┌──────────────────────────────────────┐
│  PostgreSQL Server (Cloud/LAN)      │
│  - Data terpusat                    │
│  - Multi-device sync                │
│  - Backup otomatis                  │
└──────────────────────────────────────┘
```

---

## 🚀 Setup Sync Mode

### Langkah 1: Copy Template Configuration

```bash
# Di terminal/command prompt
cd C:\Users\Hp\Documents\Project\Ritel-App
copy .env.sync .env
```

Atau manual: Rename file `.env.sync` menjadi `.env`

### Langkah 2: Edit File .env

Buka file `.env` dan sesuaikan:

```env
# Aktifkan Sync Mode
SYNC_MODE=enabled

# Database Lokal (SQLite) - Tidak perlu diubah
SYNC_SQLITE_DSN=./ritel.db

# Database Server (PostgreSQL) - UBAH SESUAI SERVER ANDA
SYNC_POSTGRES_DSN=host=192.168.1.100 port=5432 user=ritel password=ritel123 dbname=ritel_db sslmode=disable
```

**Contoh konfigurasi untuk berbagai skenario:**

#### Skenario 1: Server di Komputer Lokal
```env
SYNC_POSTGRES_DSN=host=localhost port=5432 user=ritel password=ritel123 dbname=ritel_db sslmode=disable
```

#### Skenario 2: Server di Jaringan LAN (Toko)
```env
# Ganti 192.168.1.100 dengan IP komputer server Anda
SYNC_POSTGRES_DSN=host=192.168.1.100 port=5432 user=ritel password=ritel123 dbname=ritel_db sslmode=disable
```

#### Skenario 3: Server di Cloud (AWS, GCP, Azure)
```env
SYNC_POSTGRES_DSN=host=ritel-db.abc123.ap-southeast-1.rds.amazonaws.com port=5432 user=ritel password=yourpassword dbname=ritel_db sslmode=require
```

### Langkah 3: Jalankan Aplikasi

```bash
# Build aplikasi
wails build

# Jalankan
./build/bin/Ritel-App.exe
```

---

## 💡 Skenario Penggunaan

### Skenario A: Toko dengan Internet Stabil

**Setup:**
```env
SYNC_MODE=enabled
SYNC_POSTGRES_DSN=host=server.example.com port=5432...
```

**Apa yang Terjadi:**
1. User input transaksi → Langsung tersimpan di SQLite lokal (cepat!)
2. Sync Engine detect server online → Otomatis kirim ke PostgreSQL
3. Data tersinkron real-time
4. Jika server down sesaat, tidak masalah - data tetap aman di lokal

**Keuntungan:**
- ⚡ Performa tetap cepat (SQLite lokal)
- ☁️ Data otomatis backup ke cloud
- 🔄 Multi-device bisa sinkron

---

### Skenario B: Toko dengan Internet Tidak Stabil

**Setup:** (sama seperti Skenario A)

**Apa yang Terjadi:**

1. **Saat ONLINE:**
   - Input transaksi → SQLite lokal → Auto-sync ke server
   - Status: 🟢 ONLINE

2. **Internet PUTUS:**
   - Aplikasi tetap jalan normal
   - Semua input tersimpan di SQLite lokal
   - Perubahan di-queue untuk sync nanti
   - Status: 🔴 OFFLINE
   - **User tetap bisa bekerja seperti biasa!**

3. **Internet KEMBALI:**
   - Sync Engine otomatis detect
   - Push semua data yang pending ke server
   - Queue dikosongkan
   - Status: 🟢 ONLINE
   - **User tidak perlu lakukan apa-apa!**

**Keuntungan:**
- ✅ Tidak pernah kehilangan data
- ✅ Tidak mengganggu workflow kasir
- ✅ Auto-recovery tanpa intervensi user

---

### Skenario C: Toko Tanpa Internet (Full Offline)

**Setup:**
```env
SYNC_MODE=enabled
SYNC_POSTGRES_DSN=host=server-yang-tidak-ada.com port=5432...
```

**Apa yang Terjadi:**
1. Aplikasi start, gagal connect ke server
2. Status: 🔴 OFFLINE (permanent)
3. Semua operasi tetap jalan normal menggunakan SQLite
4. Data aman tersimpan lokal

**Keuntungan:**
- ✅ Tetap bisa digunakan sepenuhnya
- ✅ Data tersimpan di SQLite lokal
- ✅ Jika suatu saat ada internet, data akan auto-sync

---

## 📊 Monitoring Sync Status

### Via API (jika WEB_ENABLED=true)

#### 1. Cek Status Sync
```bash
curl http://localhost:8080/api/sync/status
```

Response:
```json
{
  "success": true,
  "data": {
    "online": true,
    "status": "🟢 ONLINE",
    "pending": 0,
    "synced": 145,
    "failed": 0
  }
}
```

#### 2. Lihat Antrian Sync
```bash
curl http://localhost:8080/api/sync/pending
```

Response:
```json
{
  "success": true,
  "data": [
    {
      "id": 123,
      "table_name": "transaksi",
      "operation": "INSERT",
      "record_id": 456,
      "created_at": "2025-12-20 10:30:00",
      "retry_count": 2,
      "last_error": "connection refused",
      "status": "pending"
    }
  ]
}
```

#### 3. Force Sync Manual
```bash
curl -X POST http://localhost:8080/api/sync/force
```

---

## 🛠️ Troubleshooting

### Problem 1: "Sync engine not initialized"

**Penyebab:** Sync mode tidak aktif di .env

**Solusi:**
```env
# Pastikan di file .env ada:
SYNC_MODE=enabled
```

---

### Problem 2: Selalu status OFFLINE padahal internet ada

**Penyebab:** PostgreSQL server tidak bisa diakses

**Solusi:**
1. Test koneksi PostgreSQL manual:
```bash
PGPASSWORD=ritel123 psql -U ritel -h server-ip -d ritel_db
```

2. Pastikan:
   - PostgreSQL server running
   - Firewall tidak block port 5432
   - IP address benar
   - Username/password benar

---

### Problem 3: Data tidak sync padahal sudah online

**Penyebab:** Mungkin ada error di queue

**Solusi:**
```bash
# Cek pending syncs untuk lihat error
curl http://localhost:8080/api/sync/pending

# Force sync manual
curl -X POST http://localhost:8080/api/sync/force
```

---

## 🔧 Advanced Configuration

### Mengubah Interval Sync

Edit file `internal/sync/sync_engine.go`:

```go
Engine = &SyncEngine{
    // ...
    syncInterval:    10 * time.Second,  // Ubah ini (default: 10 detik)
    healthInterval:  30 * time.Second,  // Ubah ini (default: 30 detik)
    // ...
}
```

### Retry Strategy

Jika sync gagal, sistem otomatis retry hingga 5 kali.
Setelah 5 kali gagal, status menjadi `failed` dan perlu manual intervention.

---

## 📈 Best Practices

### 1. Backup Berkala
Meskipun ada auto-sync, tetap backup SQLite lokal:
```bash
# Backup otomatis setiap hari
copy ritel.db ritel_backup_%date%.db
```

### 2. Monitor Sync Queue
Secara berkala, cek apakah ada failed syncs:
```bash
curl http://localhost:8080/api/sync/status
```

Jika ada `failed` > 0, investigate dan fix.

### 3. Clean Up Synced Data
Hapus data yang sudah sync (older than 7 days):
```bash
curl -X POST http://localhost:8080/api/sync/clear
```

### 4. Testing
Sebelum deploy:
1. Test skenario offline → online
2. Test putus koneksi di tengah transaksi
3. Verify data konsistensi antara SQLite dan PostgreSQL

---

## ❓ FAQ

**Q: Apakah data bisa konflik antara lokal dan server?**
A: Tidak, karena setiap device punya queue sendiri. Data di-push ke server secara sequential.

**Q: Bagaimana jika 2 kasir edit data yang sama saat offline?**
A: Last-write-wins. Yang terakhir sync akan overwrite. Untuk multi-kasir, gunakan server mode (bukan sync mode).

**Q: Apakah bisa untuk multi-cabang?**
A: Bisa, tapi setiap cabang sync ke server yang sama. Data akan terpusat di server PostgreSQL.

**Q: Data SQLite lokal apakah dihapus setelah sync?**
A: Tidak! SQLite tetap ada dan menjadi primary database. PostgreSQL hanya untuk backup dan sinkronisasi.

**Q: Apakah bisa mematikan sync mode?**
A: Ya, set `SYNC_MODE=disabled` di .env dan restart aplikasi. Tapi data di queue tidak akan di-sync.

---

## 🎯 Kesimpulan

Sync Mode memberikan **best of both worlds**:
- 🚀 Kecepatan offline (SQLite lokal)
- ☁️ Reliability online (PostgreSQL backup)
- 🔄 Auto-sync tanpa intervention
- 💪 Resilient terhadap gangguan koneksi

**Cocok untuk:**
- Toko dengan internet tidak stabil
- Multi-device deployment
- Butuh backup otomatis
- Ingin performa maksimal

**Tidak cocok untuk:**
- Multi-kasir real-time (gunakan web mode)
- Data harus SELALU konsisten antar device
- Tidak punya server PostgreSQL

---

**Butuh bantuan?** Buka issue di GitHub atau hubungi support.

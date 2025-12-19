# Ritel App - Indeks Dokumentasi

## 📖 Panduan Utama

### 🚀 [QUICK_START.md](QUICK_START.md) - **MULAI DI SINI!**
**Waktu baca: 5 menit**

Panduan lengkap untuk memulai aplikasi dengan cepat:
- ✅ Setup SQLite (paling mudah)
- ✅ Setup PostgreSQL (production)
- ✅ Setup Dual Mode (recommended production)
- ✅ FAQ & troubleshooting dasar
- ✅ Login default & tips keamanan

**Kapan membaca:** Saat pertama kali setup aplikasi

---

### 🔧 [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Solusi Masalah
**Waktu baca: 10 menit**

Solusi lengkap untuk error & masalah umum:
- ❌ PostgreSQL connection failed
- ❌ Database does not exist
- ❌ Password authentication failed
- ❌ File .env tidak ditemukan
- ❌ Migration failed
- ❌ Performance lambat
- ✅ Error yang aman diabaikan
- 🔍 Debug checklist
- 🔄 Reset aplikasi ke default

**Kapan membaca:** Saat mengalami error atau masalah

---

## 📊 Dokumentasi Database

### 💾 [DATABASE_SETUP.md](DATABASE_SETUP.md) - Setup PostgreSQL Lengkap
**Waktu baca: 15 menit**

Panduan detail setup PostgreSQL dari awal:
- 📥 Install PostgreSQL di Windows/Linux/Mac
- 🔧 Konfigurasi database & user
- 🗄️ Create schema & tables
- 🔐 Setup permissions
- 💾 Backup & restore strategies
- 📈 Optimasi performa

**Kapan membaca:** Saat akan menggunakan PostgreSQL untuk production

---

### ⚡ [QUICK_START_POSTGRESQL.md](QUICK_START_POSTGRESQL.md) - Quick Setup
**Waktu baca: 5 menit**

Quick start untuk setup PostgreSQL:
- 🚀 3 langkah setup
- 📋 Script siap pakai
- ✅ Verifikasi instalasi
- 🔧 Troubleshooting cepat

**Kapan membaca:** Sudah familiar PostgreSQL, butuh setup cepat

---

### 🔄 [DUAL_DATABASE.md](DUAL_DATABASE.md) - Dual Mode (PostgreSQL + SQLite)
**Waktu baca: 10 menit**

Panduan lengkap dual database mode:
- 🎯 Cara kerja dual mode
- ✅ Keuntungan & use cases
- ⚙️ Setup & konfigurasi
- 📊 Monitoring sinkronisasi
- 💡 Best practices
- 🔍 Troubleshooting dual mode

**Kapan membaca:** Ingin redundansi data & backup otomatis

---

### 📝 [LOGGING_DATABASE.md](LOGGING_DATABASE.md) - Database Logging
**Waktu baca: 8 menit**

Memahami logging system aplikasi:
- 📋 Format output logging
- ✓ Status file .env
- 📊 Konfigurasi database
- 🔗 Status koneksi
- 🔐 Password masking
- 🐛 Troubleshooting dengan logging
- 📸 Contoh output semua mode

**Kapan membaca:** Ingin memahami log aplikasi atau debugging

---

## 🛠️ Referensi Teknis

### 📄 [README.md](README.md) - Project Overview
**Waktu baca: 5 menit**

Overview project & quick reference:
- 🎯 Features aplikasi
- 🚀 Quick start semua mode
- 📚 Link ke semua dokumentasi
- 🏗️ Project structure
- 💻 Development setup

**Kapan membaca:** Overview aplikasi atau referensi cepat

---

## 🎯 Skenario Penggunaan

### Scenario 1: Baru Pertama Kali Setup

**Urutan baca:**
1. 🚀 [QUICK_START.md](QUICK_START.md) - Pilih mode database
2. 📄 [README.md](README.md) - Pahami features
3. 🔧 [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Jika ada error

**Estimasi waktu:** 15-20 menit

---

### Scenario 2: Setup Production dengan PostgreSQL

**Urutan baca:**
1. 💾 [DATABASE_SETUP.md](DATABASE_SETUP.md) - Install & setup PostgreSQL
2. 🚀 [QUICK_START.md](QUICK_START.md) - Setup aplikasi
3. 📝 [LOGGING_DATABASE.md](LOGGING_DATABASE.md) - Monitoring
4. 🔧 [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Referensi error

**Estimasi waktu:** 30-45 menit

---

### Scenario 3: Setup Production dengan Dual Mode (Recommended)

**Urutan baca:**
1. 🔄 [DUAL_DATABASE.md](DUAL_DATABASE.md) - Pahami dual mode
2. 💾 [DATABASE_SETUP.md](DATABASE_SETUP.md) - Setup PostgreSQL
3. 🚀 [QUICK_START.md](QUICK_START.md) - Setup aplikasi dual mode
4. 📝 [LOGGING_DATABASE.md](LOGGING_DATABASE.md) - Monitoring
5. 🔧 [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Referensi error

**Estimasi waktu:** 45-60 menit

---

### Scenario 4: Troubleshooting Error

**Urutan baca:**
1. 🔧 [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Cari error di index
2. 📝 [LOGGING_DATABASE.md](LOGGING_DATABASE.md) - Pahami log
3. 💾 [DATABASE_SETUP.md](DATABASE_SETUP.md) - Cek setup database (jika perlu)

**Estimasi waktu:** 10-15 menit

---

### Scenario 5: Migrasi dari SQLite ke PostgreSQL

**Urutan baca:**
1. 🔄 [DUAL_DATABASE.md](DUAL_DATABASE.md) - Gunakan dual mode untuk migrasi
2. 💾 [DATABASE_SETUP.md](DATABASE_SETUP.md) - Setup PostgreSQL
3. 🚀 [QUICK_START.md](QUICK_START.md) - Setup dual mode
4. 📝 [LOGGING_DATABASE.md](LOGGING_DATABASE.md) - Monitor migrasi

**Estimasi waktu:** 45 menit

---

## 📁 File Template Konfigurasi

### `.env.sqlite` - SQLite Mode
```env
DB_DRIVER=sqlite3
DB_DSN=./ritel.db
```
**Gunakan untuk:** Development, testing, toko kecil

---

### `.env.postgres` - PostgreSQL Mode
```env
DB_DRIVER=postgres
DB_DSN=host=localhost port=5432 user=postgres password=postgres dbname=ritel_db sslmode=disable
```
**Gunakan untuk:** Production, performa tinggi

---

### `.env.dual` - Dual Mode (Recommended Production)
```env
DB_DRIVER=dual
DB_POSTGRES_DSN=host=localhost port=5432 user=postgres password=postgres dbname=ritel_db sslmode=disable
DB_SQLITE_DSN=./ritel.db
```
**Gunakan untuk:** Production + backup otomatis

---

### `.env.example` - Template dengan Dokumentasi
File referensi lengkap dengan semua opsi konfigurasi dan penjelasan.

---

## 🔧 Setup Scripts

### `setup.bat` (Windows)
Script interaktif untuk setup database mode di Windows.

### `setup.sh` (Linux/Mac)
Script interaktif untuk setup database mode di Linux/Mac.

**Cara pakai:**
```bash
# Windows
setup.bat

# Linux/Mac
./setup.sh

# Pilih mode (1/2/3)
```

---

## 📊 Comparison Table - Mode Database

| Feature | SQLite | PostgreSQL | Dual Mode |
|---------|--------|------------|-----------|
| **Setup** | ⭐⭐⭐⭐⭐ Sangat Mudah | ⭐⭐⭐ Butuh Install | ⭐⭐⭐ Butuh Install |
| **Performance** | ⭐⭐⭐ Cukup | ⭐⭐⭐⭐⭐ Sangat Tinggi | ⭐⭐⭐⭐ Tinggi |
| **Concurrent Users** | ⭐⭐ 1-2 kasir | ⭐⭐⭐⭐⭐ Unlimited | ⭐⭐⭐⭐⭐ Unlimited |
| **Backup** | Manual | Manual | ⭐⭐⭐⭐⭐ Otomatis |
| **Data Safety** | ⭐⭐⭐ Cukup | ⭐⭐⭐⭐ Tinggi | ⭐⭐⭐⭐⭐ Maksimal |
| **Cost** | ✅ Free | ✅ Free | ✅ Free |
| **Recommended For** | Development | Production | Production + Safety |

---

## ❓ Quick FAQ Index

### Setup
- **Q: Mode mana yang harus saya pilih?** → [QUICK_START.md](QUICK_START.md#faq)
- **Q: Bagaimana cara install PostgreSQL?** → [DATABASE_SETUP.md](DATABASE_SETUP.md)
- **Q: File .env tidak ditemukan?** → [TROUBLESHOOTING.md](TROUBLESHOOTING.md#2-file-env-tidak-ditemukan)

### Database
- **Q: Bagaimana cara backup data?** → [DUAL_DATABASE.md](DUAL_DATABASE.md#backup-dan-restore)
- **Q: Bagaimana cara lihat data di database?** → [QUICK_START.md](QUICK_START.md#faq)
- **Q: Bisa pindah dari SQLite ke PostgreSQL?** → [QUICK_START.md](QUICK_START.md#faq)

### Error
- **Q: PostgreSQL connection failed?** → [TROUBLESHOOTING.md](TROUBLESHOOTING.md#3-postgresql-connection-failed)
- **Q: Password authentication failed?** → [TROUBLESHOOTING.md](TROUBLESHOOTING.md#5-password-authentication-failed)
- **Q: Database does not exist?** → [TROUBLESHOOTING.md](TROUBLESHOOTING.md#4-database-does-not-exist)

### Performance
- **Q: Aplikasi lambat di dual mode?** → [TROUBLESHOOTING.md](TROUBLESHOOTING.md#13-performance-lambat-dual-mode)
- **Q: Optimasi PostgreSQL?** → [DATABASE_SETUP.md](DATABASE_SETUP.md)

---

## 🎓 Learning Path

### Beginner (Baru Mulai)
1. ✅ Baca [QUICK_START.md](QUICK_START.md)
2. ✅ Pilih SQLite mode
3. ✅ Jalankan aplikasi
4. ✅ Explore features

**Estimasi:** 30 menit

---

### Intermediate (Siap Production)
1. ✅ Baca [DATABASE_SETUP.md](DATABASE_SETUP.md)
2. ✅ Install PostgreSQL
3. ✅ Setup PostgreSQL mode
4. ✅ Baca [LOGGING_DATABASE.md](LOGGING_DATABASE.md)
5. ✅ Monitor aplikasi

**Estimasi:** 1-2 jam

---

### Advanced (Production + Safety)
1. ✅ Baca [DUAL_DATABASE.md](DUAL_DATABASE.md)
2. ✅ Setup dual mode
3. ✅ Configure backup strategy
4. ✅ Setup monitoring
5. ✅ Performance tuning

**Estimasi:** 2-3 jam

---

## 📞 Butuh Bantuan?

1. 🔍 **Cari di dokumentasi** - Gunakan index ini
2. 🐛 **Cek Troubleshooting** - [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
3. 📝 **Cek Logging** - [LOGGING_DATABASE.md](LOGGING_DATABASE.md)
4. 📧 **Contact Support** - Dengan info error lengkap

---

## 🔄 Update Dokumentasi

**Versi:** 1.0.0
**Update terakhir:** 12 Desember 2024

**Changelog:**
- ✅ Dual database mode support
- ✅ Database logging system
- ✅ Comprehensive troubleshooting guide
- ✅ Quick start guide
- ✅ Setup scripts (Windows + Linux/Mac)

---

**Happy coding! 🚀**

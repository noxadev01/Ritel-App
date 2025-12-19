# Database Logging - Status Koneksi & Konfigurasi

## Deskripsi

Aplikasi Ritel-App sekarang menampilkan **logging detail** tentang:
- Status file .env (apakah berhasil dimuat atau tidak)
- Konfigurasi database yang digunakan
- Status koneksi database (berhasil/gagal)
- Mode database yang aktif (SQLite/PostgreSQL/Dual)

## Format Output Logging

### 1. Logging File .env

Saat aplikasi start, akan muncul info tentang file .env:

**Jika file .env ditemukan:**
```
✓ File .env berhasil dimuat dari: C:\Users\Hp\Documents\Project\Ritel-App\.env
```

**Jika file .env tidak ditemukan:**
```
⚠ File .env tidak ditemukan, menggunakan konfigurasi default
```

### 2. Logging Konfigurasi Database

#### Mode SQLite:

```
📋 KONFIGURASI DATABASE
========================================
✓ Sumber: File .env
  Lokasi: C:\Users\Hp\Documents\Project\Ritel-App\.env
🔧 Driver: sqlite3
🔗 DSN: ./ritel.db
========================================

========================================
💾 SQLITE MODE
========================================
💾 Menghubungkan ke SQLite... ✓ BERHASIL
📁 Lokasi database: C:\Users\Hp\ritel-app\ritel.db
⚙️  Menyiapkan schema database...
✓ Database siap digunakan!
========================================
```

#### Mode PostgreSQL:

```
📋 KONFIGURASI DATABASE
========================================
✓ Sumber: File .env
  Lokasi: C:\Users\Hp\Documents\Project\Ritel-App\.env
🔧 Driver: postgres
🔗 DSN: host=localhost port=5432 user=postgres password=**** dbname=ritel_db sslmode=disable
========================================

========================================
📊 POSTGRESQL MODE
========================================
📊 Menghubungkan ke PostgreSQL... ✓ BERHASIL
⚙️  Menyiapkan schema database...
✓ Database siap digunakan!
========================================
```

#### Mode Dual Database:

```
📋 KONFIGURASI DUAL DATABASE
========================================
✓ Sumber: File .env
  Lokasi: C:\Users\Hp\Documents\Project\Ritel-App\.env
----------------------------------------
📊 PostgreSQL (Primary):
   host=localhost port=5432 user=postgres password=**** dbname=ritel_db sslmode=disable
----------------------------------------
💾 SQLite (Backup):
   ./ritel.db
========================================

========================================
🔄 DUAL DATABASE MODE
========================================
📊 Menghubungkan ke PostgreSQL... ✓ BERHASIL
💾 Menghubungkan ke SQLite... ✓ BERHASIL
----------------------------------------
📍 PostgreSQL: Primary (Read/Write)
📍 SQLite: Backup (Write)
----------------------------------------
⚙️  Menyiapkan schema database...
✓ Dual database mode aktif!
========================================
```

### 3. Logging Error Koneksi

Jika koneksi database gagal:

**PostgreSQL gagal:**
```
📊 Menghubungkan ke PostgreSQL... ❌ GAGAL
Error: failed to initialize PostgreSQL: pq: password authentication failed for user "postgres"
```

**SQLite gagal:**
```
💾 Menghubungkan ke SQLite... ❌ GAGAL
Error: failed to initialize SQLite: unable to open database file
```

## Informasi yang Ditampilkan

### Status File .env:
- ✅ **Berhasil dimuat** - Menampilkan lokasi file .env
- ⚠️ **Tidak ditemukan** - Menggunakan konfigurasi default

### Konfigurasi Database:
- **Driver** - Mode database (sqlite3/postgres/dual)
- **DSN** - Connection string (password di-mask dengan ****)
- **Sumber** - Dari file .env atau default

### Status Koneksi:
- ✅ **BERHASIL** - Koneksi database sukses
- ❌ **GAGAL** - Koneksi database gagal (dengan detail error)

### Mode Database:
- 💾 **SQLite Mode** - Hanya SQLite
- 📊 **PostgreSQL Mode** - Hanya PostgreSQL
- 🔄 **Dual Database Mode** - PostgreSQL + SQLite

## Contoh Skenario

### Skenario 1: File .env Ada, Mode SQLite

File `.env`:
```env
DB_DRIVER=sqlite3
DB_DSN=./ritel.db
```

Output:
```
✓ File .env berhasil dimuat dari: C:\Users\Hp\Documents\Project\Ritel-App\.env

📋 KONFIGURASI DATABASE
========================================
✓ Sumber: File .env
  Lokasi: C:\Users\Hp\Documents\Project\Ritel-App\.env
🔧 Driver: sqlite3
🔗 DSN: ./ritel.db
========================================

========================================
💾 SQLITE MODE
========================================
💾 Menghubungkan ke SQLite... ✓ BERHASIL
📁 Lokasi database: C:\Users\Hp\ritel-app\ritel.db
⚙️  Menyiapkan schema database...
✓ Database siap digunakan!
========================================
```

### Skenario 2: File .env Ada, Mode Dual

File `.env`:
```env
DB_DRIVER=dual
DB_POSTGRES_DSN=host=localhost port=5432 user=postgres password=postgres dbname=ritel_db sslmode=disable
DB_SQLITE_DSN=./ritel.db
```

Output:
```
✓ File .env berhasil dimuat dari: C:\Users\Hp\Documents\Project\Ritel-App\.env

📋 KONFIGURASI DUAL DATABASE
========================================
✓ Sumber: File .env
  Lokasi: C:\Users\Hp\Documents\Project\Ritel-App\.env
----------------------------------------
📊 PostgreSQL (Primary):
   host=localhost port=5432 user=postgres password=**** dbname=ritel_db sslmode=disable
----------------------------------------
💾 SQLite (Backup):
   ./ritel.db
========================================

========================================
🔄 DUAL DATABASE MODE
========================================
📊 Menghubungkan ke PostgreSQL... ✓ BERHASIL
💾 Menghubungkan ke SQLite... ✓ BERHASIL
----------------------------------------
📍 PostgreSQL: Primary (Read/Write)
📍 SQLite: Backup (Write)
----------------------------------------
⚙️  Menyiapkan schema database...
✓ Dual database mode aktif!
========================================
```

### Skenario 3: File .env Tidak Ada (Default)

Tidak ada file `.env`

Output:
```
⚠ File .env tidak ditemukan, menggunakan konfigurasi default

📋 KONFIGURASI DATABASE
========================================
⚠ Sumber: Konfigurasi Default
  (File .env tidak ditemukan)
🔧 Driver: sqlite3
🔗 DSN: ./ritel.db
========================================

========================================
💾 SQLITE MODE
========================================
💾 Menghubungkan ke SQLite... ✓ BERHASIL
📁 Lokasi database: C:\Users\Hp\ritel-app\ritel.db
⚙️  Menyiapkan schema database...
✓ Database siap digunakan!
========================================
```

### Skenario 4: PostgreSQL Error (Password Salah)

File `.env`:
```env
DB_DRIVER=postgres
DB_DSN=host=localhost port=5432 user=postgres password=salah123 dbname=ritel_db sslmode=disable
```

Output:
```
✓ File .env berhasil dimuat dari: C:\Users\Hp\Documents\Project\Ritel-App\.env

📋 KONFIGURASI DATABASE
========================================
✓ Sumber: File .env
  Lokasi: C:\Users\Hp\Documents\Project\Ritel-App\.env
🔧 Driver: postgres
🔗 DSN: host=localhost port=5432 user=postgres password=**** dbname=ritel_db sslmode=disable
========================================

========================================
📊 POSTGRESQL MODE
========================================
📊 Menghubungkan ke PostgreSQL... ❌ GAGAL
Error: pq: password authentication failed for user "postgres"
```

## Keamanan

### Password Masking

Password dalam DSN **otomatis di-mask** dengan `****` saat ditampilkan di log:

**DSN Asli:**
```
host=localhost port=5432 user=postgres password=rahasia123 dbname=ritel_db
```

**DSN di Log:**
```
host=localhost port=5432 user=postgres password=**** dbname=ritel_db
```

Ini mencegah password terlihat di terminal atau log file.

## Troubleshooting dengan Logging

### Problem: Aplikasi tidak konek ke database

**Langkah 1: Cek file .env**
```
⚠ File .env tidak ditemukan, menggunakan konfigurasi default
```
→ **Solusi**: Buat file .env dari template

**Langkah 2: Cek konfigurasi**
```
🔧 Driver: postgres
🔗 DSN: host=localhost port=5432 user=postgres password=**** dbname=ritel_db
```
→ **Solusi**: Pastikan driver dan DSN sudah benar

**Langkah 3: Cek koneksi**
```
📊 Menghubungkan ke PostgreSQL... ❌ GAGAL
Error: dial tcp [::1]:5432: connectex: No connection could be made
```
→ **Solusi**: PostgreSQL belum running atau port salah

### Problem: Dual mode hanya connect ke salah satu

```
📊 Menghubungkan ke PostgreSQL... ✓ BERHASIL
💾 Menghubungkan ke SQLite... ❌ GAGAL
```
→ **Solusi**: Periksa path SQLite atau permission

## Summary

Dengan logging ini, Anda bisa dengan mudah:

✅ **Memverifikasi** file .env berhasil dimuat
✅ **Melihat** konfigurasi database yang digunakan
✅ **Mengidentifikasi** error koneksi database
✅ **Memantau** mode database yang aktif
✅ **Debugging** masalah konfigurasi dengan cepat

Semua informasi penting ditampilkan dengan jelas saat aplikasi start!

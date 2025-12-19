# Quick Start - Ritel App

## Langkah Cepat (5 Menit)

### Opsi 1: SQLite (Paling Mudah - Recommended untuk Testing)

```bash
# 1. Copy template SQLite
cp .env.sqlite .env

# 2. Jalankan aplikasi
./ritel-app.exe
```

**Output yang akan muncul:**
```
✓ File .env berhasil dimuat dari: .env

📋 KONFIGURASI DATABASE
========================================
✓ Sumber: File .env
🔧 Driver: sqlite3
🔗 DSN: ./ritel.db
========================================

💾 SQLITE MODE
💾 Menghubungkan ke SQLite... ✓ BERHASIL
📁 Lokasi database: C:\Users\Hp\ritel-app\ritel.db
⚙️  Menyiapkan schema database...
✓ Database siap digunakan!
========================================
```

✅ **Selesai!** Aplikasi langsung bisa digunakan.

---

### Opsi 2: PostgreSQL (Production)

```bash
# 1. Install PostgreSQL (jika belum)
# Download dari: https://www.postgresql.org/download/

# 2. Buat database
createdb ritel_db

# 3. Copy template PostgreSQL
cp .env.postgres .env

# 4. Jalankan aplikasi
./ritel-app.exe
```

**Output yang akan muncul:**
```
✓ File .env berhasil dimuat dari: .env

📋 KONFIGURASI DATABASE
========================================
✓ Sumber: File .env
🔧 Driver: postgres
🔗 DSN: host=localhost port=5432 user=postgres password=**** dbname=ritel_db
========================================

📊 POSTGRESQL MODE
📊 Menghubungkan ke PostgreSQL... ✓ BERHASIL
⚙️  Menyiapkan schema database...
✓ Database siap digunakan!
========================================
```

✅ **Selesai!** Data tersimpan di PostgreSQL.

---

### Opsi 3: Dual Mode (Recommended untuk Production)

**Keuntungan:** Data otomatis tersimpan ke PostgreSQL DAN SQLite sekaligus!

```bash
# 1. Install PostgreSQL (jika belum)

# 2. Buat database
createdb ritel_db

# 3. Copy template Dual
cp .env.dual .env

# 4. Jalankan aplikasi
./ritel-app.exe
```

**Output yang akan muncul:**
```
✓ File .env berhasil dimuat dari: .env

📋 KONFIGURASI DUAL DATABASE
========================================
✓ Sumber: File .env
----------------------------------------
📊 PostgreSQL (Primary):
   host=localhost port=5432 user=postgres password=****
----------------------------------------
💾 SQLite (Backup):
   ./ritel.db
========================================

🔄 DUAL DATABASE MODE
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

✅ **Selesai!** Data tersimpan di 2 database untuk keamanan maksimal!

---

## Troubleshooting

### ❌ Error: "File .env tidak ditemukan"

**Solusi:**
```bash
# Pilih salah satu template:
cp .env.sqlite .env     # Untuk SQLite
cp .env.postgres .env   # Untuk PostgreSQL
cp .env.dual .env       # Untuk Dual Mode
```

### ❌ Error: "Menghubungkan ke PostgreSQL... ❌ GAGAL"

**Penyebab:** PostgreSQL belum running atau database belum dibuat.

**Solusi:**
```bash
# Cek PostgreSQL status
# Windows: Buka Services, cari PostgreSQL
# Linux: sudo systemctl status postgresql
# Mac: brew services list

# Buat database
createdb ritel_db
```

### ❌ Error: "password authentication failed"

**Penyebab:** Password PostgreSQL salah.

**Solusi:**
Edit file `.env`, ubah password:
```env
DB_POSTGRES_DSN=host=localhost port=5432 user=postgres password=PASSWORD_ANDA dbname=ritel_db sslmode=disable
```

### ❌ Error: "database ritel_db does not exist"

**Solusi:**
```bash
createdb ritel_db
```

---

## FAQ

### Q: File .env mana yang harus saya pakai?

**A:** Tergantung kebutuhan:
- **Development/Testing** → `.env.sqlite` (paling mudah)
- **Production** → `.env.postgres` (performa tinggi)
- **Production + Backup** → `.env.dual` (recommended!)

### Q: Bagaimana cara melihat data di database?

**SQLite:**
```bash
sqlite3 ~/ritel-app/ritel.db
SELECT * FROM produk;
```

**PostgreSQL:**
```bash
psql -U postgres -d ritel_db
SELECT * FROM produk;
```

### Q: Apakah bisa pindah dari SQLite ke PostgreSQL?

**A:** Ya! Gunakan dual mode:
1. Aplikasi saat ini pakai SQLite dengan data
2. Ubah ke dual mode: `cp .env.dual .env`
3. Buat database PostgreSQL: `createdb ritel_db`
4. Jalankan aplikasi → Data akan otomatis di-sync ke PostgreSQL

### Q: Bagaimana cara backup data?

**Mode SQLite:**
```bash
# Copy file database
cp ~/ritel-app/ritel.db ~/backup/ritel_$(date +%Y%m%d).db
```

**Mode PostgreSQL:**
```bash
# Gunakan pg_dump
pg_dump -U postgres ritel_db > backup.sql
```

**Mode Dual:**
Tidak perlu backup manual! Data sudah otomatis ada di 2 database.

---

## Login Default

Username: `admin`
Password: `admin123`

⚠️ **PENTING:** Ganti password admin setelah login pertama kali!

---

## Template .env yang Tersedia

| File | Deskripsi | Kapan Digunakan |
|------|-----------|----------------|
| `.env.sqlite` | SQLite saja | Development, testing, toko kecil |
| `.env.postgres` | PostgreSQL saja | Production, toko besar |
| `.env.dual` | PostgreSQL + SQLite | Production + backup otomatis |
| `.env.example` | Template dengan dokumentasi | Referensi semua opsi |

---

## Mode Database - Ringkasan

### 💾 SQLite Mode
- ✅ Setup paling mudah (1 command)
- ✅ Tidak perlu install PostgreSQL
- ✅ Cocok untuk 1-2 kasir
- ⚠️ Performa terbatas untuk banyak user

### 📊 PostgreSQL Mode
- ✅ Performa tinggi
- ✅ Support banyak user concurrent
- ✅ Cocok untuk multiple kasir
- ⚠️ Perlu install & setup PostgreSQL

### 🔄 Dual Mode (RECOMMENDED)
- ✅ Semua keuntungan PostgreSQL
- ✅ Backup otomatis ke SQLite
- ✅ Data redundancy (aman!)
- ✅ Tidak perlu backup manual
- ⚠️ Sedikit lebih lambat (5-10%)

---

## Butuh Bantuan?

📚 Dokumentasi lengkap:
- [DATABASE_SETUP.md](DATABASE_SETUP.md) - Setup detail PostgreSQL
- [DUAL_DATABASE.md](DUAL_DATABASE.md) - Panduan dual mode
- [LOGGING_DATABASE.md](LOGGING_DATABASE.md) - Memahami logging

🎯 **Recommended Setup untuk Production:**
```bash
cp .env.dual .env
createdb ritel_db
./ritel-app.exe
```

**Selamat menggunakan Ritel App! 🎉**

# Dual Database Mode - PostgreSQL + SQLite

## Deskripsi

Aplikasi Ritel-App sekarang mendukung mode **Dual Database**, di mana semua data otomatis tersimpan ke **PostgreSQL** dan **SQLite** secara bersamaan.

### Keuntungan Dual Database Mode:

1. **Redundansi Data** - Data tersimpan di 2 database sekaligus
2. **Backup Otomatis** - SQLite menjadi backup lokal dari PostgreSQL
3. **Keamanan Tinggi** - Jika satu database error, data masih ada di database lain
4. **Migrasi Mudah** - Data sudah tersedia di kedua format
5. **PostgreSQL untuk Produksi** - Performa tinggi untuk banyak transaksi
6. **SQLite untuk Backup Lokal** - File database lokal yang mudah dicopy

## Cara Kerja

### Mode Operasi:

- **WRITE (INSERT/UPDATE/DELETE)**: Otomatis ditulis ke PostgreSQL **DAN** SQLite
- **READ (SELECT)**: Dibaca dari PostgreSQL (sebagai database utama)
- **TRANSACTION**: Transaksi dilakukan di kedua database secara atomik

### Arsitektur:

```
┌─────────────────────────────────────┐
│     Aplikasi Ritel-App              │
└────────────┬────────────────────────┘
             │
    ┌────────┴────────┐
    │                 │
┌───▼────────┐   ┌───▼────────┐
│ PostgreSQL │   │   SQLite   │
│  (Primary) │   │  (Backup)  │
└────────────┘   └────────────┘
```

## Setup Dual Database Mode

### 1. Persiapan PostgreSQL

```bash
# Pastikan PostgreSQL sudah terinstall dan berjalan

# Buat database baru
createdb ritel_db

# Atau via psql:
psql -U postgres
CREATE DATABASE ritel_db;
\q
```

### 2. Konfigurasi .env

**Pilihan A: Copy dari template**
```bash
cp .env.dual .env
```

**Pilihan B: Edit manual**

Buat file `.env` dengan isi:

```env
# Aktifkan Dual Mode
DB_DRIVER=dual

# PostgreSQL (Primary)
DB_POSTGRES_DSN=host=localhost port=5432 user=postgres password=postgres dbname=ritel_db sslmode=disable

# SQLite (Backup)
DB_SQLITE_DSN=./ritel.db
```

### 3. Jalankan Aplikasi

```bash
# Build
go build -o ritel-app.exe

# Run
./ritel-app.exe
```

Aplikasi akan otomatis:
1. Connect ke PostgreSQL
2. Connect ke SQLite
3. Membuat semua tabel di kedua database
4. Menjalankan migrations di kedua database

## Contoh Penggunaan Data

### Input Data (Otomatis Tersimpan ke Kedua Database)

Ketika Anda:
- Menambah produk baru
- Melakukan transaksi penjualan
- Update stok
- Menambah pelanggan

Data **otomatis tersimpan** ke:
- ✅ PostgreSQL database `ritel_db`
- ✅ SQLite file `ritel.db`

### Cek Data PostgreSQL

```bash
# Via psql
psql -U postgres -d ritel_db

# Lihat produk
SELECT * FROM produk;

# Lihat transaksi
SELECT * FROM transaksi ORDER BY tanggal DESC LIMIT 10;
```

### Cek Data SQLite

```bash
# Via sqlite3
sqlite3 ~/ritel-app/ritel.db

# Lihat produk
SELECT * FROM produk;

# Lihat transaksi
SELECT * FROM transaksi ORDER BY tanggal DESC LIMIT 10;
```

**Data di kedua database akan IDENTIK!**

## Mode Database Lainnya

### Mode SQLite Saja (Development)

```env
DB_DRIVER=sqlite3
DB_DSN=./ritel.db
```

- Hanya menggunakan SQLite
- Cocok untuk development/testing
- File database lokal

### Mode PostgreSQL Saja (Production)

```env
DB_DRIVER=postgres
DB_DSN=host=localhost port=5432 user=postgres password=postgres dbname=ritel_db sslmode=disable
```

- Hanya menggunakan PostgreSQL
- Cocok untuk production dengan traffic tinggi
- Tidak ada backup lokal otomatis

### Mode Dual (Recommended Production)

```env
DB_DRIVER=dual
DB_POSTGRES_DSN=host=localhost port=5432 user=postgres password=postgres dbname=ritel_db sslmode=disable
DB_SQLITE_DSN=./ritel.db
```

- Menggunakan PostgreSQL + SQLite
- Data redundancy
- Backup otomatis
- **RECOMMENDED untuk production**

## Transaksi Database

### Transaksi Sederhana

Transaksi otomatis berjalan di kedua database:

```go
// Semua operasi database.Exec() otomatis dual-write
database.Exec("INSERT INTO produk (nama, harga) VALUES (?, ?)", "Beras", 15000)
// ✅ Tersimpan di PostgreSQL
// ✅ Tersimpan di SQLite
```

### Transaksi Kompleks (dengan BeginDual)

Untuk transaksi yang lebih kompleks, gunakan `BeginDual()`:

```go
tx, err := database.BeginDual()
if err != nil {
    return err
}
defer tx.Rollback()

// Eksekusi query - otomatis ke kedua database
tx.Exec("INSERT INTO transaksi (...) VALUES (...)")
tx.Exec("UPDATE produk SET stok = stok - ? WHERE id = ?", qty, produkID)

// Commit ke kedua database sekaligus
if err := tx.Commit(); err != nil {
    return err
}
```

### Keamanan Transaksi

- Jika **PostgreSQL gagal**: Rollback kedua database
- Jika **SQLite gagal**: Rollback kedua database
- Jika **salah satu gagal**: Data tetap konsisten (tidak ada partial writes)

## Troubleshooting

### Error: "failed to initialize PostgreSQL"

**Penyebab**: PostgreSQL tidak berjalan atau tidak bisa diakses

**Solusi**:
```bash
# Cek status PostgreSQL
sudo systemctl status postgresql   # Linux
brew services list                  # Mac

# Start PostgreSQL
sudo systemctl start postgresql     # Linux
brew services start postgresql      # Mac
```

### Error: "database ritel_db does not exist"

**Penyebab**: Database belum dibuat

**Solusi**:
```bash
createdb ritel_db
```

### Error: "permission denied for schema public"

**Penyebab**: User tidak punya akses

**Solusi**:
```sql
psql -U postgres -d ritel_db

GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO postgres;
```

### Performa Lambat di Dual Mode

**Normal**: Dual write memang sedikit lebih lambat karena menulis ke 2 database

**Jika terlalu lambat**:
1. Pastikan PostgreSQL di server lokal (bukan remote)
2. Gunakan SSD untuk SQLite file
3. Pertimbangkan mode PostgreSQL saja untuk production dengan traffic sangat tinggi

## Backup dan Restore

### Backup PostgreSQL

```bash
# Backup database
pg_dump -U postgres ritel_db > backup_postgres.sql

# Restore
psql -U postgres ritel_db < backup_postgres.sql
```

### Backup SQLite

```bash
# SQLite file bisa langsung dicopy
cp ~/ritel-app/ritel.db ~/backup/ritel_backup_$(date +%Y%m%d).db

# Atau gunakan script backup
cd database
./backup.sh    # Linux/Mac
backup.bat     # Windows
```

### Restore dari Backup

**Jika PostgreSQL corrupt, restore dari SQLite:**

1. Switch ke mode SQLite saja:
   ```env
   DB_DRIVER=sqlite3
   DB_DSN=./ritel.db
   ```

2. Jalankan aplikasi dengan SQLite

3. Export data dari SQLite dan import ke PostgreSQL baru

**Jika SQLite corrupt, restore dari PostgreSQL:**

SQLite akan otomatis di-recreate saat aplikasi berjalan dalam dual mode.

## Monitoring

### Cek Sinkronisasi Data

```sql
-- PostgreSQL
psql -U postgres -d ritel_db -c "SELECT COUNT(*) FROM produk;"

-- SQLite
sqlite3 ~/ritel-app/ritel.db "SELECT COUNT(*) FROM produk;"
```

Jumlah rows harus sama!

### Cek Koneksi Database

Saat aplikasi start, tidak akan ada log (silent mode).
Jika ada error saat inisialisasi, akan muncul error message.

## FAQ

### Q: Apakah data benar-benar identik?
**A**: Ya, setiap operasi write (INSERT/UPDATE/DELETE) dilakukan ke kedua database. Jika salah satu gagal, keduanya di-rollback.

### Q: Database mana yang digunakan untuk membaca?
**A**: PostgreSQL digunakan sebagai primary untuk semua operasi read.

### Q: Apakah bisa pindah dari dual mode ke single mode?
**A**: Ya, cukup ubah `DB_DRIVER` di `.env`:
- `sqlite3` - hanya SQLite
- `postgres` - hanya PostgreSQL
- `dual` - keduanya

### Q: Apakah performa lebih lambat?
**A**: Sedikit lebih lambat pada operasi write (5-10%), tapi masih sangat cepat untuk aplikasi retail. Keuntungan data redundancy lebih besar dari cost performa.

### Q: Bagaimana jika salah satu database mati saat aplikasi berjalan?
**A**: Aplikasi akan error pada operasi write berikutnya. Anda perlu restart dengan mode single database.

### Q: Bisakah menggunakan PostgreSQL remote?
**A**: Ya, tapi performa akan bergantung pada latency jaringan. Recommended PostgreSQL di server lokal untuk dual mode.

## Kesimpulan

Dual database mode memberikan:

- ✅ **Keamanan data maksimal** - redundansi otomatis
- ✅ **Backup real-time** - tidak perlu cron job
- ✅ **Fleksibilitas** - bisa switch database kapan saja
- ✅ **Peace of mind** - data aman di 2 tempat

**Recommended** untuk production environment!

---

Untuk pertanyaan lebih lanjut, lihat dokumentasi di:
- `DATABASE_SETUP.md` - Setup PostgreSQL
- `QUICK_START_POSTGRESQL.md` - Quick start guide
- `.env.example` - Contoh konfigurasi lengkap

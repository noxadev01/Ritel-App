# Troubleshooting - Ritel App

## Error Umum & Solusi

### 1. Window Class Error (AMAN - Bisa Diabaikan)

**Error:**
```
[1212/193248.880:ERROR:ui\gfx\win\window_impl.cc:124] Failed to unregister class Chrome_WidgetWin_0. Error = 1411
```

**Penyebab:**
- Error internal dari Wails/Chromium saat aplikasi shutdown
- Window class masih terdaftar di Windows saat aplikasi close

**Status:** ⚠️ **AMAN - Bisa diabaikan**

**Dampak:**
- Tidak ada dampak pada aplikasi
- Data tetap aman
- Aplikasi tetap berfungsi normal

**Solusi:**
Tidak perlu action, error ini tidak berbahaya dan tidak mempengaruhi fungsi aplikasi.

---

### 2. File .env Tidak Ditemukan

**Error:**
```
⚠ File .env tidak ditemukan, menggunakan konfigurasi default
```

**Penyebab:**
File `.env` belum dibuat

**Solusi:**

**Opsi A - Pakai Setup Script (Recommended):**
```bash
# Windows
setup.bat

# Linux/Mac
./setup.sh
```

**Opsi B - Manual:**
```bash
# Pilih salah satu:
cp .env.sqlite .env      # Untuk SQLite
cp .env.postgres .env    # Untuk PostgreSQL
cp .env.dual .env        # Untuk Dual Mode
```

---

### 3. PostgreSQL Syntax Error (AUTOINCREMENT)

**Error:**
```
📊 Menghubungkan ke PostgreSQL... ✓ BERHASIL
⚙️  Menyiapkan schema database...
Error: pq: syntax error at or near "AUTOINCREMENT"
```

**Penyebab:**
Bug dalam kode - migrations table tidak ditranslate untuk PostgreSQL

**Status:** ✅ **SUDAH DIPERBAIKI** (versi terbaru)

**Solusi:**
Rebuild aplikasi dengan versi terbaru:
```bash
go build -o ritel-app.exe
./ritel-app.exe
```

---

### 4. PostgreSQL Connection Failed

**Error:**
```
📊 Menghubungkan ke PostgreSQL... ❌ GAGAL
Error: dial tcp [::1]:5432: connectex: No connection could be made
```

**Penyebab:**
PostgreSQL belum running atau tidak terinstall

**Solusi:**

**Windows:**
```bash
# Cek service PostgreSQL
# 1. Tekan Win+R
# 2. Ketik: services.msc
# 3. Cari "PostgreSQL"
# 4. Klik Start jika stopped

# Atau via command:
net start postgresql-x64-14
```

**Linux:**
```bash
# Cek status
sudo systemctl status postgresql

# Start jika stopped
sudo systemctl start postgresql
```

**Mac:**
```bash
# Cek status
brew services list

# Start jika stopped
brew services start postgresql
```

---

### 4. Database Does Not Exist

**Error:**
```
📊 Menghubungkan ke PostgreSQL... ❌ GAGAL
Error: pq: database "ritel_db" does not exist
```

**Penyebab:**
Database `ritel_db` belum dibuat

**Solusi:**
```bash
# Buat database
createdb ritel_db

# Atau via psql:
psql -U postgres
CREATE DATABASE ritel_db;
\q
```

---

### 5. Password Authentication Failed

**Error:**
```
📊 Menghubungkan ke PostgreSQL... ❌ GAGAL
Error: pq: password authentication failed for user "postgres"
```

**Penyebab:**
Password PostgreSQL salah di file `.env`

**Solusi:**

1. **Cari password PostgreSQL yang benar**
2. **Edit file `.env`:**
   ```env
   # Ubah password di baris ini:
   DB_POSTGRES_DSN=host=localhost port=5432 user=postgres password=PASSWORD_ANDA_YANG_BENAR dbname=ritel_db sslmode=disable
   ```
3. **Restart aplikasi**

**Lupa password PostgreSQL?**

**Windows:**
```bash
# Reset password via pgAdmin atau:
# 1. Buka pgAdmin
# 2. Klik kanan pada user 'postgres'
# 3. Properties > Definition > Password
```

**Linux/Mac:**
```bash
sudo -u postgres psql
ALTER USER postgres PASSWORD 'password_baru';
\q
```

---

### 6. Permission Denied (SQLite)

**Error:**
```
💾 Menghubungkan ke SQLite... ❌ GAGAL
Error: unable to open database file
```

**Penyebab:**
Tidak ada permission untuk create/write file database

**Solusi:**
```bash
# Pastikan folder ~/ritel-app ada dan writable
mkdir -p ~/ritel-app
chmod 755 ~/ritel-app
```

---

### 7. Port Already in Use

**Error:**
```
Error: listen tcp :34115: bind: address already in use
```

**Penyebab:**
Port 34115 sudah dipakai oleh instance aplikasi lain

**Solusi:**

**Windows:**
```bash
# Cari process yang pakai port 34115
netstat -ano | findstr :34115

# Kill process (ganti PID dengan hasil di atas)
taskkill /PID 1234 /F
```

**Linux/Mac:**
```bash
# Cari process yang pakai port 34115
lsof -i :34115

# Kill process
kill -9 PID
```

---

### 8. Dual Mode - PostgreSQL OK, SQLite GAGAL

**Error:**
```
📊 Menghubungkan ke PostgreSQL... ✓ BERHASIL
💾 Menghubungkan ke SQLite... ❌ GAGAL
```

**Penyebab:**
Path SQLite tidak bisa diakses atau permission denied

**Solusi:**
```bash
# Cek file .env, pastikan path SQLite valid:
DB_SQLITE_DSN=./ritel.db

# Pastikan folder ada dan writable
ls -la ~/ritel-app/
```

---

### 9. Dual Mode - SQLite OK, PostgreSQL GAGAL

**Error:**
```
📊 Menghubungkan ke PostgreSQL... ❌ GAGAL
💾 Menghubungkan ke SQLite... ✓ BERHASIL
```

**Penyebab:**
PostgreSQL tidak running atau koneksi salah

**Solusi:**
Lihat solusi #3 (PostgreSQL Connection Failed)

**Temporary Fix - Switch ke SQLite saja:**
```bash
# Edit .env
DB_DRIVER=sqlite3
DB_DSN=./ritel.db

# Restart aplikasi
```

---

### 10. Migration Failed

**Error:**
```
⚙️  Menyiapkan schema database...
Error: failed to run migrations: ...
```

**Penyebab:**
Schema database corrupt atau migration error

**Solusi:**

**Untuk SQLite:**
```bash
# Backup database lama
mv ~/ritel-app/ritel.db ~/ritel-app/ritel.db.backup

# Start fresh (database baru akan dibuat otomatis)
./ritel-app.exe
```

**Untuk PostgreSQL:**
```bash
# Drop dan recreate database
dropdb ritel_db
createdb ritel_db

# Restart aplikasi
./ritel-app.exe
```

---

### 11. Aplikasi Crash Saat Start

**Solusi:**

1. **Cek log terakhir:**
   ```bash
   # Lihat error message di terminal
   ./ritel-app.exe
   ```

2. **Reset konfigurasi:**
   ```bash
   # Backup .env lama
   mv .env .env.backup

   # Gunakan SQLite default
   cp .env.sqlite .env

   # Coba start
   ./ritel-app.exe
   ```

3. **Fresh install:**
   ```bash
   # Backup data
   cp ~/ritel-app/ritel.db ~/backup/

   # Hapus app directory
   rm -rf ~/ritel-app/

   # Start aplikasi (akan create fresh)
   ./ritel-app.exe
   ```

---

### 12. Data Tidak Tersimpan

**Gejala:**
Input data tapi saat cek database, data tidak ada

**Penyebab & Solusi:**

1. **Cek mode database:**
   ```bash
   cat .env | grep DB_DRIVER
   ```

2. **Cek database yang benar:**

   **SQLite Mode:**
   ```bash
   sqlite3 ~/ritel-app/ritel.db
   SELECT COUNT(*) FROM produk;
   ```

   **PostgreSQL Mode:**
   ```bash
   psql -U postgres -d ritel_db
   SELECT COUNT(*) FROM produk;
   ```

   **Dual Mode:**
   Cek kedua database di atas

3. **Cek error di terminal:**
   Lihat apakah ada error saat save data

---

### 13. Performance Lambat (Dual Mode)

**Gejala:**
Aplikasi lambat saat save data di dual mode

**Penyebab:**
Write ke 2 database memang sedikit lebih lambat

**Solusi:**

**Temporary - Switch ke PostgreSQL saja:**
```env
# Edit .env
DB_DRIVER=postgres
DB_DSN=host=localhost port=5432 user=postgres password=postgres dbname=ritel_db sslmode=disable
```

**Permanent - Optimasi:**
1. Pastikan PostgreSQL di server lokal (bukan remote)
2. Gunakan SSD untuk SQLite file
3. Untuk traffic sangat tinggi, pakai PostgreSQL mode saja

---

## Cek Status System

### Cek PostgreSQL Running

**Windows:**
```bash
# Via services
services.msc

# Via psql
psql -U postgres -c "SELECT version();"
```

**Linux:**
```bash
sudo systemctl status postgresql
```

**Mac:**
```bash
brew services list | grep postgresql
```

### Cek Database Exist

```bash
# List semua database
psql -U postgres -l

# Atau
psql -U postgres
\l
\q
```

### Cek File .env

```bash
# Cek file ada
ls -la .env

# Lihat isi
cat .env
```

### Cek Database File SQLite

```bash
# Cek file ada
ls -la ~/ritel-app/ritel.db

# Cek ukuran
du -h ~/ritel-app/ritel.db

# Cek isi
sqlite3 ~/ritel-app/ritel.db "SELECT COUNT(*) FROM produk;"
```

---

## Debug Mode

Untuk mendapatkan informasi lebih detail saat troubleshooting:

```bash
# Jalankan dengan verbose logging (coming soon)
./ritel-app.exe --debug
```

---

## Kontak Support

Jika masih mengalami masalah:

1. 📸 Screenshot error message
2. 📋 Copy output terminal
3. 📁 Share file `.env` (hapus password!)
4. 🔍 Jelaskan langkah yang dilakukan

---

## Error yang Aman Diabaikan

Error-error ini tidak berbahaya dan tidak perlu action:

✅ `Failed to unregister class Chrome_WidgetWin_0` - Chromium internal
✅ `DevTools listening on ws://...` - Debug info
✅ `[Violation] 'setTimeout' handler took...` - Performance info (dev mode)

---

## Quick Fix Checklist

Jika aplikasi tidak jalan, cek urut dari atas:

- [ ] File `.env` sudah ada?
- [ ] PostgreSQL running? (jika pakai postgres/dual)
- [ ] Database `ritel_db` sudah dibuat?
- [ ] Password PostgreSQL benar?
- [ ] Port 34115 tidak dipakai aplikasi lain?
- [ ] Folder `~/ritel-app/` ada dan writable?

---

## Reset Aplikasi ke Default

Jika semua tidak berhasil, reset total:

```bash
# 1. Backup data
cp ~/ritel-app/ritel.db ~/backup/ritel_backup_$(date +%Y%m%d).db

# 2. Hapus config
rm .env

# 3. Hapus app data
rm -rf ~/ritel-app/

# 4. Copy template default
cp .env.sqlite .env

# 5. Start fresh
./ritel-app.exe
```

Aplikasi akan start dengan database SQLite kosong yang fresh.

---

**Butuh bantuan lebih lanjut? Baca dokumentasi lengkap di:**
- [QUICK_START.md](QUICK_START.md)
- [DATABASE_SETUP.md](DATABASE_SETUP.md)
- [LOGGING_DATABASE.md](LOGGING_DATABASE.md)

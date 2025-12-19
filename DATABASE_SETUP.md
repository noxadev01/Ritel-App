# Database Setup - Ritel App

Ritel App sekarang mendukung dua jenis database:
- **SQLite** (default, untuk development)
- **PostgreSQL** (untuk production)

## SQLite (Default)

Secara default, aplikasi menggunakan SQLite tanpa konfigurasi tambahan.

```bash
# Tidak perlu environment variable
# Database akan dibuat di: ~/ritel-app/ritel.db
```

## PostgreSQL (Production)

### Cara 1: Menggunakan .env File (Recommended)

Cara paling mudah adalah menggunakan file `.env`:

1. **Copy template .env**
   ```bash
   cp .env.example .env
   ```

2. **Edit file .env**
   ```bash
   # Database Configuration
   DB_DRIVER=postgres
   DB_DSN=host=localhost port=5432 user=ritel password=your_password dbname=ritel_db sslmode=disable
   ```

3. **Run aplikasi**
   ```bash
   ./ritel-app
   # Aplikasi akan otomatis membaca konfigurasi dari .env
   ```

**Keuntungan menggunakan .env:**
- ✅ Tidak perlu set environment variable setiap kali
- ✅ Mudah switch antara development dan production
- ✅ Konfigurasi tersimpan di file (sudah ada di .gitignore)
- ✅ Portabel dan mudah di-deploy

### Cara 2: Menggunakan Environment Variables

Untuk menggunakan PostgreSQL dengan environment variables:

#### Windows (PowerShell)
```powershell
$env:DB_DRIVER="postgres"
$env:DB_DSN="host=localhost port=5432 user=ritel password=ritel_password dbname=ritel_db sslmode=disable"
```

#### Windows (Command Prompt)
```cmd
set DB_DRIVER=postgres
set DB_DSN=host=localhost port=5432 user=ritel password=ritel_password dbname=ritel_db sslmode=disable
```

#### Linux/Mac
```bash
export DB_DRIVER=postgres
export DB_DSN="host=localhost port=5432 user=ritel password=ritel_password dbname=ritel_db sslmode=disable"
```

## Setup PostgreSQL Database

### 1. Install PostgreSQL
Download dan install PostgreSQL dari https://www.postgresql.org/download/

### 2. Create Database dan User
```sql
-- Login sebagai postgres user
psql -U postgres

-- Buat database
CREATE DATABASE ritel_db;

-- Buat user (optional, bisa pakai user postgres)
CREATE USER ritel WITH PASSWORD 'ritel_password';

-- Berikan akses
GRANT ALL PRIVILEGES ON DATABASE ritel_db TO ritel;

-- Grant schema permissions (PostgreSQL 15+)
\c ritel_db
GRANT ALL ON SCHEMA public TO ritel;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ritel;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ritel;
```

### 3. Verifikasi Koneksi
```bash
# Test koneksi
psql -h localhost -p 5432 -U ritel -d ritel_db
```

### 4. Run Aplikasi

**Opsi A: Menggunakan .env file (Recommended)**
```bash
# 1. Buat file .env dari template
cp .env.example .env

# 2. Edit .env dengan konfigurasi Anda
# DB_DRIVER=postgres
# DB_DSN=host=localhost port=5432 user=ritel password=ritel_password dbname=ritel_db sslmode=disable

# 3. Run aplikasi
./ritel-app
```

**Opsi B: Menggunakan environment variables**
```bash
# Set environment variables
export DB_DRIVER=postgres
export DB_DSN="host=localhost port=5432 user=ritel password=ritel_password dbname=ritel_db sslmode=disable"

# Run aplikasi
./ritel-app
```

## Connection String Format

### SQLite
```
./ritel.db
/path/to/database.db
```

### PostgreSQL
```
host=localhost port=5432 user=username password=password dbname=database sslmode=disable
```

Parameter PostgreSQL:
- `host`: Server hostname (default: localhost)
- `port`: Port number (default: 5432)
- `user`: Database user
- `password`: User password
- `dbname`: Database name
- `sslmode`: SSL mode (disable, require, verify-ca, verify-full)

## Features yang Didukung

### ✅ Fully Supported (Kedua Database)
- Semua operasi CRUD
- Transaction management
- Foreign key constraints
- Migrations
- Promo system (curah/satuan)
- FIFO stock management
- Point of Sale
- Customer loyalty system

### ⚠️ Database-Specific Features
| Feature | SQLite | PostgreSQL |
|---------|--------|------------|
| File-based backup | ✅ | ❌ |
| Auto backup | ✅ | ❌ (gunakan pg_dump) |
| PRAGMA integrity_check | ✅ | ❌ (tidak diperlukan) |
| Database size query | PRAGMA | pg_database_size() |

## Migration Strategy

Aplikasi akan automatically run migrations saat start up. Semua 30+ migrations sudah kompatibel dengan kedua database.

### Migrasi Data dari SQLite ke PostgreSQL

Jika Anda ingin memindahkan data dari SQLite ke PostgreSQL:

1. Export data dari SQLite (gunakan tools seperti `sqlite3` CLI atau DB Browser)
2. Setup PostgreSQL database (lihat langkah di atas)
3. Run aplikasi dengan PostgreSQL untuk membuat schema
4. Import data menggunakan `COPY` atau `INSERT` statements

## Troubleshooting

### PostgreSQL: "permission denied for schema public"
```sql
-- Jalankan sebagai superuser (postgres)
GRANT ALL ON SCHEMA public TO ritel;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ritel;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ritel;
```

### PostgreSQL: "database does not exist"
```sql
-- Buat database dulu
CREATE DATABASE ritel_db;
```

### SQLite: "database is locked"
- Pastikan tidak ada proses lain yang membuka database
- Close semua connection sebelum backup

### Connection refused (PostgreSQL)
- Pastikan PostgreSQL service running
- Check `pg_hba.conf` untuk authentication settings
- Pastikan firewall tidak blocking port 5432

## Performance Tips

### SQLite (Development)
- Sudah optimal dengan WAL mode
- Backup otomatis sebelum migrations

### PostgreSQL (Production)
- Gunakan connection pooling (sudah dikonfigurasi: max 25 connections)
- Set `shared_buffers` di `postgresql.conf` (recommended: 25% of RAM)
- Enable `pg_stat_statements` untuk query monitoring
- Setup regular backups menggunakan `pg_dump`:
  ```bash
  pg_dump -U ritel ritel_db > backup_$(date +%Y%m%d).sql
  ```

## Environment Variables Reference

### Konfigurasi via .env File

Buat file `.env` di root directory aplikasi (atau directory executable):

```env
# Database Configuration
DB_DRIVER=postgres
DB_DSN=host=localhost port=5432 user=ritel password=your_password dbname=ritel_db sslmode=disable
```

**Lokasi file .env:**
1. Current directory (tempat Anda run aplikasi)
2. Directory executable (tempat file .exe berada)

### Daftar Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_DRIVER` | `sqlite3` | Database driver: `sqlite3` atau `postgres` |
| `DB_DSN` | `./ritel.db` | Connection string |

**Prioritas Konfigurasi:**
1. Environment variables (paling tinggi)
2. File .env
3. Default values (SQLite)

## Implementation Notes

Aplikasi menggunakan **Database Abstraction Layer** dengan pattern Adapter:
- Queries ditranslasi otomatis ke dialect yang sesuai
- `?` placeholders → `$1, $2, ...` untuk PostgreSQL
- `datetime('now')` → `NOW()` untuk PostgreSQL
- `INTEGER PRIMARY KEY AUTOINCREMENT` → `SERIAL PRIMARY KEY` untuk PostgreSQL
- PRAGMA statements diabaikan di PostgreSQL

Migrations dan queries sudah fully compatible dengan kedua database.

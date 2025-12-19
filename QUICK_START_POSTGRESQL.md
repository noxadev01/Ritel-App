# Quick Start: PostgreSQL Setup

Panduan cepat untuk setup PostgreSQL di Ritel App.

## ⚡ Quick Setup (5 Menit)

### 1. Setup PostgreSQL Database

**Opsi A: Manual Setup (SQL Commands)**

```sql
-- Login ke PostgreSQL
psql -U postgres

-- Buat database dan user
CREATE DATABASE ritel_db;
CREATE USER ritel WITH PASSWORD 'ritel';
GRANT ALL PRIVILEGES ON DATABASE ritel_db TO ritel;

-- Connect ke database
\c ritel_db

-- Grant permissions (PostgreSQL 15+)
GRANT ALL ON SCHEMA public TO ritel;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ritel;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ritel;

-- Exit
\q
```

**Opsi B: Menggunakan Schema File (Recommended)**

```bash
# 1. Buat database dan user dulu
psql -U postgres -c "CREATE DATABASE ritel_db;"
psql -U postgres -c "CREATE USER ritel WITH PASSWORD 'ritel';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE ritel_db TO ritel;"

# 2. Grant permissions
psql -U postgres -d ritel_db -c "GRANT ALL ON SCHEMA public TO ritel;"
psql -U postgres -d ritel_db -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ritel;"
psql -U postgres -d ritel_db -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ritel;"

# 3. Run schema file (membuat semua tables, indexes, triggers)
psql -U ritel -d ritel_db -f database/schema_postgres.sql

# 4. (Optional) Insert sample data
psql -U ritel -d ritel_db -f database/seed_data.sql
```

**Schema file akan otomatis membuat:**
- ✅ 18 tables (kategori, produk, transaksi, pelanggan, promo, dll)
- ✅ 30+ indexes untuk performance
- ✅ Foreign key constraints
- ✅ Auto-update triggers (updated_at columns)
- ✅ Table comments untuk dokumentasi

**Seed data akan insert:**
- Default admin user (username: admin, password: admin123)
- Sample product categories
- Sample products (curah & satuan)
- Sample customers
- Sample promos
- Default print settings
- Default loyalty points settings

### 2. Setup .env File

```bash
# Copy template
cp .env.example .env

# Atau untuk quick start, copy contoh PostgreSQL
cp .env.postgres .env
```

**Edit file .env:**
```env
DB_DRIVER=postgres
DB_DSN=host=localhost port=5432 user=ritel password=ritel dbname=ritel_db sslmode=disable
```

### 3. Run Aplikasi

```bash
./ritel-app
```

**Output yang diharapkan:**
```
Loaded configuration from .env file
Initializing database with driver: postgres
Connecting to PostgreSQL database
Database connection established
Database initialized successfully
```

## 🔄 Switch Between SQLite and PostgreSQL

### Switch ke PostgreSQL
```bash
# Edit .env
DB_DRIVER=postgres
DB_DSN=host=localhost port=5432 user=ritel password=ritel dbname=ritel_db sslmode=disable
```

### Switch ke SQLite
```bash
# Hapus atau rename .env file, atau edit:
DB_DRIVER=sqlite3
DB_DSN=./ritel.db

# Atau lebih mudah:
cp .env.sqlite .env
```

## 🐛 Troubleshooting

### Error: "password authentication failed"
```sql
-- Reset password
ALTER USER ritel WITH PASSWORD 'ritel';
```

### Error: "permission denied for schema public"
```sql
-- Login sebagai postgres
psql -U postgres -d ritel_db

-- Grant permissions
GRANT ALL ON SCHEMA public TO ritel;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ritel;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ritel;
```

### Error: "database does not exist"
```sql
-- Login sebagai postgres
psql -U postgres

-- Create database
CREATE DATABASE ritel_db;
```

### PostgreSQL service tidak running (Windows)
```powershell
# Check service status
Get-Service postgresql*

# Start service
Start-Service postgresql-x64-15  # Adjust version number
```

### PostgreSQL service tidak running (Linux)
```bash
# Check status
sudo systemctl status postgresql

# Start service
sudo systemctl start postgresql

# Enable auto-start
sudo systemctl enable postgresql
```

## 📊 Verifikasi Setup

### Test Connection
```bash
# Test koneksi PostgreSQL
psql -h localhost -p 5432 -U ritel -d ritel_db

# Jika berhasil, Anda akan masuk ke PostgreSQL prompt:
# ritel_db=>
```

### Check Tables
Setelah aplikasi pertama kali running, check apakah tables sudah dibuat:

```sql
-- List all tables
\dt

-- Expected output:
-- produk, kategori, transaksi, pelanggan, promo, stok_history, dll.
```

## 🚀 Production Deployment

### 1. Secure Password
```env
# Gunakan password yang kuat
DB_DRIVER=postgres
DB_DSN=host=localhost port=5432 user=ritel password=SecurePassword123! dbname=ritel_db sslmode=require
```

### 2. Remote Database
```env
# Untuk remote server
DB_DRIVER=postgres
DB_DSN=host=192.168.1.100 port=5432 user=ritel password=SecurePassword123! dbname=ritel_db sslmode=require
```

### 3. SSL/TLS
```env
# Enable SSL
DB_DRIVER=postgres
DB_DSN=host=db.example.com port=5432 user=ritel password=SecurePassword123! dbname=ritel_db sslmode=require
```

## 💾 Backup dan Restore

### Backup Database

**Linux/Mac:**
```bash
# Menggunakan script otomatis
chmod +x database/backup.sh
./database/backup.sh

# Manual backup
pg_dump -U ritel -d ritel_db -f backup_$(date +%Y%m%d_%H%M%S).sql
```

**Windows:**
```cmd
REM Menggunakan script otomatis
database\backup.bat

REM Manual backup
pg_dump -U ritel -d ritel_db -f backup_%date:~-4,4%%date:~-10,2%%date:~-7,2%.sql
```

**Backup akan disimpan di:** `database/backups/`

**Fitur script backup:**
- ✅ Timestamped backup files
- ✅ Auto cleanup (keep last 10 backups)
- ✅ Optional compression (gzip)
- ✅ Backup verification

### Restore Database

**⚠️ WARNING: Restore akan menghapus semua data!**

**Linux/Mac:**
```bash
# Menggunakan script otomatis
chmod +x database/restore.sh
./database/restore.sh database/backups/ritel_db_20250101_120000.sql

# Manual restore
psql -U postgres -c "DROP DATABASE IF EXISTS ritel_db;"
psql -U postgres -c "CREATE DATABASE ritel_db OWNER ritel;"
psql -U ritel -d ritel_db -f backup_file.sql
```

**Windows:**
```cmd
REM Menggunakan script otomatis
database\restore.bat database\backups\ritel_db_20250101_120000.sql

REM Manual restore
psql -U postgres -c "DROP DATABASE IF EXISTS ritel_db;"
psql -U postgres -c "CREATE DATABASE ritel_db OWNER ritel;"
psql -U ritel -d ritel_db -f backup_file.sql
```

**Fitur script restore:**
- ✅ Safety backup before restore
- ✅ Automatic database recreation
- ✅ Verification after restore
- ✅ Statistics display

## 📝 Tips

1. **Backup Schedule**
   ```bash
   # Setup cron job (Linux) untuk backup harian
   crontab -e
   # Tambahkan:
   0 2 * * * /path/to/database/backup.sh

   # Windows Task Scheduler
   # Buat scheduled task yang run database\backup.bat setiap hari
   ```

2. **Monitor Connections**
   ```sql
   -- Check active connections
   SELECT * FROM pg_stat_activity WHERE datname = 'ritel_db';
   ```

3. **Database Size**
   ```sql
   -- Check database size
   SELECT pg_size_pretty(pg_database_size('ritel_db'));
   ```

4. **Performance Tuning**
   Edit `postgresql.conf`:
   ```ini
   shared_buffers = 256MB        # 25% of RAM
   effective_cache_size = 1GB    # 50-75% of RAM
   work_mem = 16MB
   maintenance_work_mem = 128MB
   ```

## 🔗 Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Connection String Format](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING)
- [PostgreSQL on Windows](https://www.postgresql.org/download/windows/)
- [PostgreSQL on Linux](https://www.postgresql.org/download/linux/)

## 📞 Need Help?

Lihat dokumentasi lengkap di [DATABASE_SETUP.md](DATABASE_SETUP.md)

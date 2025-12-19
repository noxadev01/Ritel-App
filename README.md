# Ritel App

Aplikasi Point of Sale (POS) yang dibangun dengan Wails (Go + React) dengan support multi-database.

## Features

- 🛒 Point of Sale dengan promo system (curah/satuan)
- 📦 Inventory management dengan FIFO batch system
- 👥 Customer loyalty & point system
- 📊 Dashboard & reporting
- 🔄 Support SQLite (development) dan PostgreSQL (production)
- 💾 Auto backup (SQLite)
- 🎯 Printer thermal support

## Database Support

Aplikasi mendukung tiga mode database:

- **SQLite** (default) - Untuk development dan small business
- **PostgreSQL** - Untuk production dan scalability
- **Dual Mode** - PostgreSQL + SQLite bersamaan (recommended production)

### Quick Start - SQLite (Default)

```bash
# Tidak perlu konfigurasi, langsung run
./ritel-app
```

Output:
```
⚠ File .env tidak ditemukan, menggunakan konfigurasi default

📋 KONFIGURASI DATABASE
========================================
⚠ Sumber: Konfigurasi Default
🔧 Driver: sqlite3
🔗 DSN: ./ritel.db
========================================

💾 SQLITE MODE
💾 Menghubungkan ke SQLite... ✓ BERHASIL
✓ Database siap digunakan!
```

### Quick Start - PostgreSQL

1. **Setup database:**
   ```sql
   psql -U postgres
   CREATE DATABASE ritel_db;
   CREATE USER ritel WITH PASSWORD 'ritel';
   GRANT ALL PRIVILEGES ON DATABASE ritel_db TO ritel;
   ```

2. **Setup .env file:**
   ```bash
   cp .env.example .env
   # Edit .env:
   # DB_DRIVER=postgres
   # DB_DSN=host=localhost port=5432 user=ritel password=ritel dbname=ritel_db sslmode=disable
   ```

3. **Run aplikasi:**
   ```bash
   ./ritel-app
   ```

   Output:
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
   ```

### Quick Start - Dual Mode (Recommended)

**Dual mode** menyimpan data ke PostgreSQL dan SQLite secara bersamaan untuk redundansi.

1. **Copy template:**
   ```bash
   cp .env.dual .env
   ```

2. **Setup PostgreSQL:**
   ```bash
   createdb ritel_db
   ```

3. **Run aplikasi:**
   ```bash
   ./ritel-app
   ```

   Output:
   ```
   ✓ File .env berhasil dimuat dari: .env

   📋 KONFIGURASI DUAL DATABASE
   ========================================
   ✓ Sumber: File .env
   📊 PostgreSQL (Primary): host=localhost port=5432 user=postgres password=****
   💾 SQLite (Backup): ./ritel.db
   ========================================

   🔄 DUAL DATABASE MODE
   📊 Menghubungkan ke PostgreSQL... ✓ BERHASIL
   💾 Menghubungkan ke SQLite... ✓ BERHASIL
   📍 PostgreSQL: Primary (Read/Write)
   📍 SQLite: Backup (Write)
   ✓ Dual database mode aktif!
   ```

📚 **Dokumentasi:**
- 🚀 [QUICK_START.md](QUICK_START.md) - **START DI SINI!** Panduan 5 menit
- 🔧 [TROUBLESHOOTING.md](TROUBLESHOOTING.md) - Solusi error & masalah umum
- [DATABASE_SETUP.md](DATABASE_SETUP.md) - Setup lengkap PostgreSQL
- [QUICK_START_POSTGRESQL.md](QUICK_START_POSTGRESQL.md) - Quick start PostgreSQL
- [DUAL_DATABASE.md](DUAL_DATABASE.md) - Dual database mode (PostgreSQL + SQLite)
- [LOGGING_DATABASE.md](LOGGING_DATABASE.md) - Database logging & monitoring

## Development

### Prerequisites

- Go 1.23+
- Node.js 18+
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- PostgreSQL (optional, untuk testing dengan PostgreSQL)

### Live Development

```bash
# Run dengan hot reload
wails dev
```

Development server akan berjalan di:
- Frontend: Vite dev server
- Backend API: http://localhost:34115

### Building

```bash
# Build untuk production
wails build

# Build untuk platform tertentu
wails build -platform windows/amd64
wails build -platform linux/amd64
```

## Configuration

Aplikasi menggunakan .env file untuk konfigurasi:

```env
# Database Configuration
DB_DRIVER=sqlite3  # atau postgres
DB_DSN=./ritel.db  # atau PostgreSQL connection string
```

**Templates tersedia:**
- `.env.example` - Template dengan penjelasan lengkap
- `.env.sqlite` - Quick setup SQLite
- `.env.postgres` - Quick setup PostgreSQL
- `.env.dual` - Quick setup Dual mode (PostgreSQL + SQLite)

## Project Structure

```
.
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database layer & migrations
│   │   └── dialect/     # Database dialect abstraction
│   ├── models/          # Data models
│   ├── repository/      # Data access layer
│   └── service/         # Business logic
├── frontend/            # React frontend
│   └── src/
│       └── components/  # UI components
└── DATABASE_SETUP.md    # Database documentation
```

## Technologies

**Backend:**
- Go 1.23
- Wails v2 (desktop framework)
- SQLite / PostgreSQL
- Database abstraction layer

**Frontend:**
- React 18
- Tailwind CSS
- Context API for state management

## Contributing

This is a private project. For questions or issues, please contact the development team.

## License

Proprietary - All rights reserved.

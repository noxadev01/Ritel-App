# Panduan Testing - Dual Mode Application

## 📋 Persiapan Sebelum Testing

### 1. Pastikan Database Siap
```bash
# Cek apakah PostgreSQL sudah running (jika pakai PostgreSQL)
pg_isready

# Atau cek file database SQLite ada (jika masih pakai SQLite)
ls ritel.db
```

### 2. Pastikan Dependencies Terinstall
```bash
# Backend dependencies
go mod download

# Frontend dependencies
cd frontend
npm install
cd ..
```

---

## 🖥️ MODE 1: Testing Desktop Mode (Wails)

Desktop mode menggunakan Wails bindings (IPC) - **TIDAK** perlu HTTP server.

### Langkah Testing Desktop Mode:

#### A. Disable Web Server (Opsional)
Edit file `.env`:
```env
WEB_ENABLED=false
```

#### B. Jalankan Aplikasi Desktop
```bash
# Jalankan Wails dev mode
wails dev
```

#### C. Apa yang Harus Ditest:

**1. Login System**
- Buka aplikasi desktop
- Login dengan:
  - Username: `admin`
  - Password: `admin123`
- Pastikan berhasil login dan redirect ke dashboard

**2. POS/Transaksi (Sudah Updated)**
- Cari produk dengan search bar
- Scan barcode (jika ada scanner)
- Tambah produk ke keranjang
- Pilih pelanggan
- Terapkan promo
- Proses pembayaran
- Cetak struk

**3. Fitur Lain (Masih Pakai Wails Binding - Belum Updated)**
- Manajemen produk (CRUD)
- Manajemen pelanggan
- History transaksi
- Dashboard
- Laporan
- Settings

**Expected Result**:
✅ Semua fitur harus tetap bekerja normal seperti sebelumnya karena kita tidak mengubah Wails bindings, hanya menambahkan HTTP layer.

---

## 🌐 MODE 2: Testing Web Mode (Browser)

Web mode menggunakan HTTP REST API + JWT authentication.

### Langkah Testing Web Mode:

#### A. Enable Web Server
Edit file `.env`:
```env
WEB_ENABLED=true
WEB_PORT=8080
WEB_HOST=0.0.0.0
JWT_SECRET=ritel-app-secret-key-change-this-in-production-2024
JWT_EXPIRY_HOURS=24
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
CORS_ALLOW_CREDENTIALS=true

# Database config (penting!)
DB_DRIVER=postgres  # atau sqlite
DB_DSN=postgresql://user:password@localhost:5432/riteldb  # sesuaikan
```

#### B. Jalankan Backend Server

**Opsi 1: Tanpa Wails (Backend Only)**
```bash
# Build dan jalankan
go build -o ritel-app-server
./ritel-app-server

# Atau langsung run
go run .
```

**Opsi 2: Dengan Wails (Dual Mode)**
```bash
# Jalankan desktop + web server sekaligus
wails dev
```

Cek apakah server berjalan:
```bash
curl http://localhost:8080/api/auth/login
# Harus return response (bukan error connection refused)
```

#### C. Jalankan Frontend (Development Mode)

Buka terminal baru:
```bash
cd frontend
npm run dev
```

Frontend akan jalan di: `http://localhost:5173`

#### D. Test di Browser

**1. Buka Browser**
```
http://localhost:5173
```

**2. Test Login (PENTING - INI SUDAH UPDATED)**
- Masukkan username: `admin`
- Password: `admin123`
- Klik Login
- **Check Developer Console (F12):**
  - Cek Network tab - harus ada request ke `/api/auth/login`
  - Cek localStorage - harus ada `token` dan `user`
  - Cek Console - tidak boleh ada error

**Expected Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "namaLengkap": "Administrator",
      "role": "admin"
    }
  }
}
```

**3. Test POS/Transaksi (SUDAH UPDATED)**
- Pergi ke halaman Transaksi/Kasir
- **Check Developer Console:**
  - Harus ada request ke `/api/produk` (get products)
  - Harus ada request ke `/api/pelanggan` (get customers)
  - Harus ada request ke `/api/promo/active` (get promos)
  - Semua request harus include header `Authorization: Bearer <token>`

- Test workflow lengkap:
  - Search produk
  - Tambah ke keranjang
  - Pilih pelanggan
  - Apply promo
  - Bayar
  - **Check:** Request `/api/transaksi` (POST) harus berhasil

**4. Test Fitur Lain (BELUM UPDATED - AKAN ERROR)**
Fitur yang belum diupdate akan error karena masih menggunakan Wails imports:
- Manajemen produk page
- Manajemen pelanggan page
- Dashboard page
- Laporan page
- Settings page

**Expected Error:**
```
Cannot find module '../../wailsjs/go/main/App'
```

---

## 🔍 Cara Debug/Troubleshooting

### 1. Backend Server Tidak Jalan

**Check Port:**
```bash
# Windows
netstat -ano | findstr :8080

# Linux/Mac
lsof -i :8080
```

**Check Logs:**
Backend akan print log di console saat startup:
```
🚀 Starting Ritel App...
📊 Database driver: postgres
✅ Database connected successfully
🌐 HTTP Server enabled on port 8080
✨ Application started successfully!
```

### 2. Frontend Tidak Bisa Connect ke Backend

**Check CORS:**
- Pastikan `CORS_ALLOWED_ORIGINS` di `.env` include `http://localhost:5173`

**Check Browser Console:**
```
Access to fetch at 'http://localhost:8080/api/auth/login' from origin 'http://localhost:5173' has been blocked by CORS policy
```

**Fix:** Edit `.env`:
```env
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

### 3. Login Berhasil tapi API Call Gagal (401 Unauthorized)

**Check Token:**
- Buka Developer Console
- Cek localStorage: `localStorage.getItem('token')`
- Token harus ada dan format: `eyJhbGc...`

**Check Request Headers:**
- Buka Network tab
- Click request yang gagal
- Check Headers - harus ada:
  ```
  Authorization: Bearer eyJhbGc...
  ```

### 4. Database Error

**Error: "Connection refused"**
- PostgreSQL tidak running
- Fix: `pg_ctl start` atau start PostgreSQL service

**Error: "Database does not exist"**
```bash
createdb riteldb
```

**Error: "relation does not exist"**
- Migration belum jalan
- Database akan auto-create tables saat pertama kali start

### 5. Component Error di Browser

**Error: "Cannot find module wailsjs"**
- Component belum diupdate ke API modules
- Solusi: Update component sesuai panduan di `MIGRATION_STATUS.md`

**Error: "Network Error"**
- Backend tidak jalan
- Check: `curl http://localhost:8080/api/auth/login`

---

## 📊 Testing Checklist

### Desktop Mode (Wails)
```
✅ Login works
✅ POS/Transaksi works
⏳ Product management (not yet updated)
⏳ Customer management (not yet updated)
⏳ Transaction history (not yet updated)
⏳ Dashboard (not yet updated)
⏳ Reports (not yet updated)
⏳ Settings (not yet updated)
```

### Web Mode (Browser)
```
✅ Backend server starts (port 8080)
✅ Frontend connects to backend
✅ Login works (JWT token saved)
✅ POS/Transaksi works
⏳ Product management (not yet updated)
⏳ Customer management (not yet updated)
⏳ Transaction history (not yet updated)
⏳ Dashboard (not yet updated)
⏳ Reports (not yet updated)
⏳ Settings (not yet updated)
```

---

## 🎯 Test Scenarios

### Scenario 1: Basic Login & POS (Web Mode)

1. Start backend: `go run .`
2. Start frontend: `cd frontend && npm run dev`
3. Open browser: `http://localhost:5173`
4. Login as admin
5. Go to Transaksi/POS page
6. Add product to cart
7. Complete transaction
8. **Expected**: Transaction saved, receipt printed/shown

### Scenario 2: Concurrent Access (Web Mode)

1. Login di Chrome dengan user `admin`
2. Login di Firefox dengan user `staff` (jika ada)
3. Buat transaksi di Chrome
4. Buat transaksi di Firefox
5. **Expected**: Both transactions saved correctly, no conflicts

### Scenario 3: Desktop + Web Simultaneous (Dual Mode)

1. Start: `wails dev` (dengan `WEB_ENABLED=true`)
2. Desktop app opens automatically
3. Open browser: `http://localhost:5173`
4. Login di desktop app
5. Login di browser
6. Create transaction di desktop
7. Create transaction di browser
8. **Expected**: Both transactions visible in database

---

## 🛠️ Quick Test Commands

### Test Backend API Directly (Terminal)

```bash
# Test login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# Save token dari response
TOKEN="eyJhbGc..."

# Test get products (dengan auth)
curl -X GET http://localhost:8080/api/produk \
  -H "Authorization: Bearer $TOKEN"

# Test get transactions
curl -X GET http://localhost:8080/api/transaksi \
  -H "Authorization: Bearer $TOKEN"
```

### Test Database Connection

```bash
# PostgreSQL
PGPASSWORD=postgres psql -h localhost -U postgres -d riteldb -c "SELECT COUNT(*) FROM produk;"

# SQLite
sqlite3 ritel.db "SELECT COUNT(*) FROM produk;"
```

---

## 📝 Next Steps After Testing

1. **Jika Desktop Mode OK**: Desktop app masih bekerja normal ✅
2. **Jika Web Mode Login OK**: Backend + Auth sudah bekerja ✅
3. **Jika POS OK di Web**: Core functionality sudah porting ✅
4. **Update Component Lain**: Gunakan pattern yang sama seperti di `Transaksi.jsx`

---

## ❓ FAQ

**Q: Apakah harus update semua component sekarang?**
A: Tidak. Yang sudah updated (Login + POS) sudah cukup untuk demo web mode. Component lain bisa diupdate bertahap.

**Q: Bagaimana cara switch antara desktop dan web mode?**
A:
- Desktop: Jalankan `wails dev`
- Web: Jalankan `go run .` + `npm run dev` di folder frontend

**Q: Apakah bisa jalan desktop dan web bersamaan?**
A: Ya! Set `WEB_ENABLED=true` lalu jalankan `wails dev`. Desktop app + web server akan jalan bersamaan.

**Q: Database mana yang dipakai?**
A: Kedua mode pakai database yang sama (sesuai `DB_DSN` di `.env`). Data di-share antara desktop dan web.

**Q: Bagaimana cara deploy web mode ke production?**
A:
```bash
# Build backend
go build -o ritel-app-server

# Build frontend
cd frontend && npm run build

# Copy build hasil ke folder yang bisa diakses web server
# Jalankan backend
WEB_ENABLED=true ./ritel-app-server
```

Frontend build hasil akan di-serve oleh Gin di `/` route.

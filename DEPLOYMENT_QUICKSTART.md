# 🚀 Deployment Quick Start

Panduan cepat untuk deploy Ritel-App ke production dalam 10 menit.

## 📋 Pilih Skenario Anda

### 🖥️ Skenario 1: Toko Tunggal (Desktop Mode)
**Setup:** 5 menit | **Requirements:** Windows PC/Laptop
```bash
# Windows
1. Download ritel-app.exe
2. Run: deploy-windows.bat
3. Pilih SQLite
4. Done!
```
✅ Perfect untuk: Toko retail kecil, 1-2 kasir

---

### 🌐 Skenario 2: Multi-User Web Server
**Setup:** 15 menit | **Requirements:** Ubuntu Server, Domain name
```bash
# Ubuntu Server
sudo ./deploy-production.sh
# Follow prompts, enter domain & password
# Done!
```
✅ Perfect untuk: Retail chain, banyak cabang, remote access

---

### 🔄 Skenario 3: Hybrid (Desktop + Web)
**Setup:** 5 menit | **Requirements:** Windows + PostgreSQL
```bash
# Windows
1. Install PostgreSQL
2. Run: deploy-windows.bat
3. Pilih PostgreSQL
4. Enable WEB_ENABLED=true in .env
5. Done!
```
✅ Perfect untuk: Toko dengan POS utama + tablet/mobile

---

## 🎯 Quick Deploy Scripts

### Windows (Single Store)
```batch
REM Download aplikasi ke C:\ritel-app\
cd C:\ritel-app
deploy-windows.bat

REM Atau manual:
1. Extract ritel-app.exe ke folder
2. Create .env (copy dari .env.example)
3. Edit DB_DRIVER=sqlite3
4. Run ritel-app.exe
```

### Linux/Ubuntu (Web Server)
```bash
# 1. Build aplikasi
go build -o ritel-app .

# 2. Run deploy script
chmod +x deploy-production.sh
sudo ./deploy-production.sh

# Script akan otomatis:
# - Install PostgreSQL, Nginx, SSL
# - Setup database & user
# - Configure systemd service
# - Setup firewall
# - Schedule automated backups
```

### Docker (Coming Soon)
```bash
# Build image
docker build -t ritel-app .

# Run container
docker-compose up -d

# Access
http://localhost:8080
```

---

## ⚙️ Konfigurasi Cepat

### SQLite (Standalone)
```env
DB_DRIVER=sqlite3
DB_DSN=./ritel.db
WEB_ENABLED=false
```

### PostgreSQL (Production)
```env
DB_DRIVER=postgres
DB_DSN=host=localhost port=5432 user=ritel password=SecurePass123! dbname=ritel_db sslmode=disable
WEB_ENABLED=true
WEB_PORT=8080
JWT_SECRET=ChangeThisToRandomString
```

### Dual Mode (Backup)
```env
DB_DRIVER=dual
DB_POSTGRES_DSN=host=localhost port=5432 user=ritel password=SecurePass123! dbname=ritel_db sslmode=disable
DB_SQLITE_DSN=./ritel_backup.db
WEB_ENABLED=true
```

---

## 🔒 Security Checklist (2 menit)

✅ **After First Deploy:**
```bash
# 1. Change admin password
Login → Settings → Change Password

# 2. Generate JWT secret
openssl rand -base64 32
# Add to .env: JWT_SECRET=generated_value

# 3. Change database password (PostgreSQL)
ALTER USER ritel WITH PASSWORD 'NewSecurePassword123!';
# Update .env with new password

# 4. Enable SSL (Web server only)
sudo certbot --nginx -d yourdomain.com

# 5. Setup firewall
sudo ufw enable
sudo ufw allow 22,80,443/tcp
```

---

## 💾 Backup Setup (1 menit)

### Windows
```batch
REM Manual backup
backup.bat

REM Scheduled (already set by deploy script)
REM Daily at 2:00 AM
```

### Linux
```bash
# Manual backup
./backup.sh

# Scheduled (already set by deploy script)
# Daily at 2:00 AM via cron
```

### Backup Locations
```
Windows: D:\backups\ritel-app\
Linux:   /backups/ritel-app/
```

---

## 🧪 Testing Deployment

### Test 1: Health Check
```bash
# Web mode
curl http://localhost:8080/health
# Expected: {"status":"ok"}

# Desktop mode
# Just open the app, should show login screen
```

### Test 2: Database Connection
```bash
# PostgreSQL
psql -h localhost -U ritel -d ritel_db -c "SELECT COUNT(*) FROM users;"
# Expected: 1 (admin user)

# SQLite
# App will auto-create on first run
```

### Test 3: Login
```
URL: http://localhost:8080 (web) or open desktop app
Username: admin
Password: admin123
Expected: Successfully logged in
```

### Test 4: API
```bash
# Get token
TOKEN=$(curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.data.token')

# Test API
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/produk
# Expected: {"success":true, "data":[]}
```

---

## 🆘 Troubleshooting

### App Won't Start
```bash
# Check logs
Windows: Check Event Viewer or app console
Linux: journalctl -u ritel-app -f

# Common issues:
1. Database not running → Start PostgreSQL
2. Port in use → Change WEB_PORT in .env
3. Missing .env → Copy from .env.example
```

### Database Connection Failed
```bash
# Test connection
Windows: psql -h localhost -U ritel -d ritel_db
Linux: sudo -u postgres psql -d ritel_db

# Fix permissions
sudo -u postgres psql -d ritel_db
GRANT ALL ON SCHEMA public TO ritel;
```

### Can't Access Web Interface
```bash
# Check if service running
Linux: systemctl status ritel-app
Windows: nssm status ritel-app

# Check firewall
Linux: sudo ufw status
Windows: netsh advfirewall firewall show rule name="Ritel-App Web Server"

# Test port
netstat -an | grep 8080
# Should show LISTENING
```

---

## 📊 Monitoring

### Check Status
```bash
# Service status
systemctl status ritel-app  # Linux
nssm status ritel-app       # Windows

# Logs
journalctl -u ritel-app -f  # Linux
# Windows: Check app console or Event Viewer

# Resource usage
top -p $(pgrep ritel-app)   # Linux
tasklist /FI "IMAGENAME eq ritel-app.exe" /FO TABLE /NH  # Windows
```

### Performance Metrics
```sql
-- Database size
SELECT pg_size_pretty(pg_database_size('ritel_db'));

-- Table statistics
SELECT schemaname, tablename,
       pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

---

## 🔄 Updates & Maintenance

### Update Application
```bash
# 1. Stop application
systemctl stop ritel-app  # Linux
nssm stop ritel-app       # Windows

# 2. Backup database
./backup.sh  # or backup.bat

# 3. Replace binary
cp new-ritel-app /opt/ritel-app/ritel-app

# 4. Start application
systemctl start ritel-app
```

### Database Maintenance
```sql
-- Optimize database (run monthly)
VACUUM ANALYZE;
REINDEX DATABASE ritel_db;

-- Check for issues
SELECT * FROM pg_stat_activity WHERE state = 'idle in transaction';
```

---

## 📞 Need Help?

### Quick Links
- **Full Documentation:** [PRODUCTION_DEPLOYMENT.md](PRODUCTION_DEPLOYMENT.md)
- **PostgreSQL Setup:** [QUICK_START_POSTGRESQL.md](QUICK_START_POSTGRESQL.md)
- **Database Guide:** [DATABASE_SETUP.md](DATABASE_SETUP.md)

### Support
- GitHub Issues: [Create Issue](https://github.com/yourusername/ritel-app/issues)
- Email: support@your-domain.com

---

## 📝 Summary

| Scenario | Setup Time | Best For | Complexity |
|----------|-----------|----------|------------|
| Desktop (SQLite) | 5 min | Single store | ⭐ Easy |
| Web Server (PostgreSQL) | 15 min | Multi-store | ⭐⭐ Medium |
| Hybrid (Dual Mode) | 10 min | POS + Backup | ⭐⭐ Medium |

**Recommendation:**
- Small store (1-2 kasir): **Desktop Mode**
- Growing business: **Hybrid Mode**
- Multiple locations: **Web Server Mode**

---

**Ready to deploy?** Choose your scenario and run the deployment script! 🚀

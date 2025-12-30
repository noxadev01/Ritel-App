# 🚀 Production Deployment Guide - Ritel-App

Panduan lengkap untuk deploy Ritel-App ke production environment.

## 📋 Table of Contents

1. [Deployment Scenarios](#deployment-scenarios)
2. [Persiapan Production](#persiapan-production)
3. [Option A: Single POS Desktop](#option-a-single-pos-desktop)
4. [Option B: Web Server Multi-User](#option-b-web-server-multi-user)
5. [Option C: Hybrid Desktop + Web](#option-c-hybrid-desktop--web)
6. [Security Checklist](#security-checklist)
7. [Backup & Recovery](#backup--recovery)
8. [Monitoring & Maintenance](#monitoring--maintenance)

---

## 🎯 Deployment Scenarios

### Scenario 1: Single Store POS (Desktop)
**Use Case:** Toko retail tunggal dengan 1-2 kasir
- ✅ Desktop app standalone
- ✅ SQLite database lokal
- ✅ Tidak perlu server/internet
- ✅ Setup paling mudah

### Scenario 2: Multi-Store / Multi-User (Web Server)
**Use Case:** Retail chain dengan banyak cabang
- ✅ Centralized PostgreSQL database
- ✅ Web-based access dari mana saja
- ✅ Real-time inventory sync
- ✅ Centralized reporting

### Scenario 3: Hybrid (Desktop + Web)
**Use Case:** Toko dengan POS utama + backup device
- ✅ Desktop app sebagai POS utama
- ✅ Web access untuk mobile/tablet
- ✅ Dual database (PostgreSQL + SQLite backup)
- ✅ Best of both worlds

---

## 🔧 Persiapan Production

### 1. Build Production Binary

**Windows:**
```bash
# Build executable
go build -ldflags="-s -w" -o ritel-app.exe .

# Atau build dengan Wails (recommended)
wails build -clean -upx

# Output: build/bin/ritel-app.exe (~40MB)
```

**Linux:**
```bash
# Build executable
go build -ldflags="-s -w" -o ritel-app .

# Dengan Wails
wails build -clean -platform linux/amd64

# Output: build/bin/ritel-app
```

**macOS:**
```bash
# Build executable
go build -ldflags="-s -w" -o ritel-app .

# Dengan Wails
wails build -clean -platform darwin/universal

# Output: build/bin/ritel-app.app
```

### 2. Setup PostgreSQL Database (Production)

**Install PostgreSQL:**
```bash
# Windows: Download installer dari postgresql.org
# Linux (Ubuntu/Debian):
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib

# Linux (CentOS/RHEL):
sudo yum install postgresql-server postgresql-contrib
sudo postgresql-setup initdb
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

**Create Database:**
```bash
# Login sebagai postgres user
sudo -u postgres psql

# Atau di Windows:
psql -U postgres

# Jalankan SQL commands:
CREATE DATABASE ritel_db;
CREATE USER ritel WITH PASSWORD 'YourSecurePassword123!';
GRANT ALL PRIVILEGES ON DATABASE ritel_db TO ritel;

# Connect ke database
\c ritel_db

# Grant schema permissions (PostgreSQL 15+)
GRANT ALL ON SCHEMA public TO ritel;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ritel;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ritel;

# Exit
\q
```

**Run Schema File:**
```bash
psql -U ritel -d ritel_db -f database/schema_postgres.sql
psql -U ritel -d ritel_db -f database/seed_data.sql
```

### 3. Setup .env File

```bash
# Copy template
cp .env.example .env

# Edit dengan nano/vim/notepad
nano .env
```

---

## 📦 Option A: Single POS Desktop

### Setup (Windows)

**1. Directory Structure:**
```
C:\ritel-app-production\
├── ritel-app.exe          # Main executable
├── .env                   # Configuration file
└── data\
    └── ritel.db          # SQLite database (auto-created)
```

**2. Create .env File:**
```env
# Database Configuration
DB_DRIVER=sqlite3
DB_DSN=./data/ritel.db

# Web Server (Disabled)
WEB_ENABLED=false

# JWT Secret (change this!)
JWT_SECRET=YourRandomSecretKey123!@#
JWT_EXPIRY_HOURS=24
```

**3. Run Application:**
```cmd
# Create data directory
mkdir data

# Run application
ritel-app.exe
```

**4. Auto-start on Windows Boot:**
```cmd
# Method 1: Startup Folder
# Press Win+R, type: shell:startup
# Copy shortcut to ritel-app.exe here

# Method 2: Task Scheduler
# Create task yang run ritel-app.exe on login
```

**5. Backup Strategy:**
```cmd
# Manual backup
xcopy C:\ritel-app-production\data\ritel.db D:\backups\ /Y

# Automated backup script (backup.bat):
@echo off
set SOURCE=C:\ritel-app-production\data\ritel.db
set DEST=D:\backups\ritel_backup_%date:~-4,4%%date:~-10,2%%date:~-7,2%_%time:~0,2%%time:~3,2%.db
xcopy %SOURCE% %DEST% /Y
echo Backup completed: %DEST%

# Schedule dengan Task Scheduler (daily at 11 PM)
```

---

## 🌐 Option B: Web Server Multi-User

### Setup (Ubuntu/Debian Server)

**1. Server Requirements:**
- Ubuntu 20.04 LTS or newer
- 2GB RAM minimum
- 10GB disk space
- PostgreSQL 12+
- Nginx (reverse proxy)

**2. Install Dependencies:**
```bash
# Update system
sudo apt-get update && sudo apt-get upgrade -y

# Install PostgreSQL
sudo apt-get install postgresql postgresql-contrib -y

# Install Nginx
sudo apt-get install nginx -y

# Install certbot (for SSL)
sudo apt-get install certbot python3-certbot-nginx -y
```

**3. Setup Application Directory:**
```bash
# Create app directory
sudo mkdir -p /opt/ritel-app
cd /opt/ritel-app

# Copy files
sudo cp ~/ritel-app /opt/ritel-app/
sudo cp ~/database/schema_postgres.sql /opt/ritel-app/
sudo cp ~/database/seed_data.sql /opt/ritel-app/

# Make executable
sudo chmod +x /opt/ritel-app/ritel-app
```

**4. Create .env File:**
```bash
sudo nano /opt/ritel-app/.env
```

```env
# Database Configuration (PostgreSQL)
DB_DRIVER=postgres
DB_DSN=host=localhost port=5432 user=ritel password=YourSecurePassword123! dbname=ritel_db sslmode=disable

# Web Server Configuration
WEB_ENABLED=true
WEB_PORT=8080
WEB_HOST=127.0.0.1  # Only accessible via nginx

# JWT Configuration
JWT_SECRET=YourRandomSecretKey123!@#$%^&*
JWT_EXPIRY_HOURS=8

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://your-domain.com,https://www.your-domain.com
CORS_ALLOW_CREDENTIALS=true
```

**5. Create Systemd Service:**
```bash
sudo nano /etc/systemd/system/ritel-app.service
```

```ini
[Unit]
Description=Ritel-App POS System
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/ritel-app
ExecStart=/opt/ritel-app/ritel-app
Restart=always
RestartSec=10

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/ritel-app

# Environment
Environment="HOME=/opt/ritel-app"

# Logging
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

**6. Setup Nginx Reverse Proxy:**
```bash
sudo nano /etc/nginx/sites-available/ritel-app
```

```nginx
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;

    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com www.your-domain.com;

    # SSL certificates (will be configured by certbot)
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # SSL settings
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # API Proxy
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check
    location /health {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
    }

    # Frontend static files
    location / {
        root /opt/ritel-app/frontend/dist;
        try_files $uri $uri/ /index.html;

        # Cache static assets
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # Logging
    access_log /var/log/nginx/ritel-app-access.log;
    error_log /var/log/nginx/ritel-app-error.log;
}
```

**7. Build Frontend:**
```bash
cd frontend
npm install
npm run build

# Copy dist to server
sudo cp -r dist /opt/ritel-app/frontend/
```

**8. Enable and Start Services:**
```bash
# Enable nginx site
sudo ln -s /etc/nginx/sites-available/ritel-app /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

# Enable and start ritel-app
sudo systemctl daemon-reload
sudo systemctl enable ritel-app
sudo systemctl start ritel-app

# Check status
sudo systemctl status ritel-app
sudo journalctl -u ritel-app -f  # View logs
```

**9. Setup SSL Certificate:**
```bash
# Get SSL certificate from Let's Encrypt
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# Auto-renewal is configured by default
# Test renewal:
sudo certbot renew --dry-run
```

**10. Firewall Configuration:**
```bash
# Allow HTTP, HTTPS, SSH
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable

# PostgreSQL (if allowing remote access)
# sudo ufw allow from trusted_ip to any port 5432
```

---

## 🔄 Option C: Hybrid Desktop + Web

### Use Case:
- Desktop app di kasir utama
- Web access untuk tablet/mobile di toko
- Database sync antara PostgreSQL (server) dan SQLite (lokal)

**1. Setup .env (Dual Mode):**
```env
# Dual Database Mode
DB_DRIVER=dual

# PostgreSQL (Primary - untuk reads)
DB_POSTGRES_DSN=host=localhost port=5432 user=ritel password=SecurePass123! dbname=ritel_db sslmode=disable

# SQLite (Backup - for local redundancy)
DB_SQLITE_DSN=./data/ritel_backup.db

# Enable Web Server
WEB_ENABLED=true
WEB_PORT=8080
WEB_HOST=0.0.0.0  # Allow local network access

# JWT Configuration
JWT_SECRET=YourRandomSecretKey123!
JWT_EXPIRY_HOURS=24

# CORS (allow local network)
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://192.168.1.*
CORS_ALLOW_CREDENTIALS=true
```

**2. Run Application:**
```bash
# Desktop app will start + Web server on port 8080
./ritel-app.exe
```

**3. Access from Other Devices:**
```
# Desktop app: Main POS
# Tablet: http://192.168.1.100:8080 (kasir IP)
# Mobile: http://192.168.1.100:8080
```

**Benefits:**
- ✅ Data redundancy (PostgreSQL + SQLite)
- ✅ Desktop performance untuk kasir utama
- ✅ Web access untuk device tambahan
- ✅ Jika PostgreSQL down, SQLite tetap jalan

---

## 🔒 Security Checklist

### Database Security

**PostgreSQL:**
```sql
-- Strong password
ALTER USER ritel WITH PASSWORD 'YourVerySecurePassword123!@#';

-- Limit connections
ALTER USER ritel CONNECTION LIMIT 10;

-- Revoke public access
REVOKE ALL ON DATABASE ritel_db FROM PUBLIC;
```

**pg_hba.conf:**
```conf
# Local connections only
local   ritel_db    ritel                           md5
host    ritel_db    ritel    127.0.0.1/32          md5

# Remote connections (if needed)
host    ritel_db    ritel    192.168.1.0/24        md5
```

### Application Security

**1. Change JWT Secret:**
```bash
# Generate strong secret
openssl rand -base64 32

# Add to .env
JWT_SECRET=GeneratedSecretHere
```

**2. Change Default Admin Password:**
```sql
-- After first login, change admin password
UPDATE users
SET password = '$2a$10$NewHashedPasswordHere'
WHERE username = 'admin';
```

**3. Disable Debug Mode:**
```env
# No debug/verbose logging in production
# Check main.go or config for debug flags
```

**4. File Permissions:**
```bash
# Linux/Mac
chmod 600 .env              # Only owner can read/write
chmod 755 ritel-app         # Executable
chown www-data:www-data *   # Correct ownership

# Windows
# Use file properties > Security > Advanced
# Remove inheritance, grant access only to app user
```

### Network Security

**1. Firewall Rules:**
```bash
# Allow only necessary ports
# SSH: 22
# HTTP: 80 (redirect to HTTPS)
# HTTPS: 443
# PostgreSQL: 5432 (only from trusted IPs)
```

**2. HTTPS Only:**
- Use SSL certificate (Let's Encrypt)
- Redirect HTTP to HTTPS
- Enable HSTS header

**3. Rate Limiting (Nginx):**
```nginx
# Add to nginx config
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;

location /api/ {
    limit_req zone=api burst=20 nodelay;
    # ... rest of config
}
```

---

## 💾 Backup & Recovery

### Automated Backup Script

**Linux (backup.sh):**
```bash
#!/bin/bash
# Ritel-App Automated Backup Script

BACKUP_DIR="/backups/ritel-app"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="ritel_db"
DB_USER="ritel"
RETENTION_DAYS=30

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup PostgreSQL
echo "Backing up PostgreSQL..."
pg_dump -U $DB_USER -d $DB_NAME | gzip > $BACKUP_DIR/postgres_$DATE.sql.gz

# Backup SQLite (if exists)
if [ -f "/opt/ritel-app/data/ritel.db" ]; then
    echo "Backing up SQLite..."
    cp /opt/ritel-app/data/ritel.db $BACKUP_DIR/sqlite_$DATE.db
    gzip $BACKUP_DIR/sqlite_$DATE.db
fi

# Backup .env file
echo "Backing up configuration..."
cp /opt/ritel-app/.env $BACKUP_DIR/env_$DATE.txt

# Remove old backups (older than RETENTION_DAYS)
echo "Cleaning old backups..."
find $BACKUP_DIR -name "*.gz" -mtime +$RETENTION_DAYS -delete
find $BACKUP_DIR -name "*.txt" -mtime +$RETENTION_DAYS -delete

echo "Backup completed: $BACKUP_DIR"
```

**Windows (backup.bat):**
```batch
@echo off
REM Ritel-App Automated Backup Script

set BACKUP_DIR=D:\backups\ritel-app
set DATE=%date:~-4,4%%date:~-10,2%%date:~-7,2%_%time:~0,2%%time:~3,2%
set DATE=%DATE: =0%
set DB_NAME=ritel_db
set DB_USER=ritel

REM Create backup directory
if not exist "%BACKUP_DIR%" mkdir "%BACKUP_DIR%"

REM Backup PostgreSQL
echo Backing up PostgreSQL...
pg_dump -U %DB_USER% -d %DB_NAME% > "%BACKUP_DIR%\postgres_%DATE%.sql"

REM Backup SQLite
echo Backing up SQLite...
xcopy "C:\ritel-app-production\data\ritel.db" "%BACKUP_DIR%\sqlite_%DATE%.db*" /Y

REM Backup .env
echo Backing up configuration...
copy "C:\ritel-app-production\.env" "%BACKUP_DIR%\env_%DATE%.txt"

echo Backup completed: %BACKUP_DIR%
```

### Schedule Automated Backups

**Linux (Cron):**
```bash
# Edit crontab
crontab -e

# Add backup schedule (daily at 2 AM)
0 2 * * * /opt/ritel-app/backup.sh >> /var/log/ritel-backup.log 2>&1

# Or with absolute path
0 2 * * * /bin/bash /opt/ritel-app/backup.sh
```

**Windows (Task Scheduler):**
```powershell
# Create scheduled task
schtasks /create /tn "Ritel-App Backup" /tr "D:\ritel-app\backup.bat" /sc daily /st 02:00
```

### Restore Procedure

**From PostgreSQL Backup:**
```bash
# Stop application
sudo systemctl stop ritel-app

# Restore database
psql -U postgres -c "DROP DATABASE IF EXISTS ritel_db;"
psql -U postgres -c "CREATE DATABASE ritel_db OWNER ritel;"
gunzip -c postgres_20250101_020000.sql.gz | psql -U ritel -d ritel_db

# Start application
sudo systemctl start ritel-app
```

**From SQLite Backup:**
```bash
# Stop application
sudo systemctl stop ritel-app

# Restore database
gunzip -c sqlite_20250101_020000.db.gz > /opt/ritel-app/data/ritel.db

# Start application
sudo systemctl start ritel-app
```

---

## 📊 Monitoring & Maintenance

### Health Checks

**1. Application Health:**
```bash
# Check if application is running
curl http://localhost:8080/health

# Expected response:
# {"status":"ok"}
```

**2. Database Health:**
```sql
-- PostgreSQL connection count
SELECT count(*) FROM pg_stat_activity WHERE datname = 'ritel_db';

-- Database size
SELECT pg_size_pretty(pg_database_size('ritel_db'));

-- Table sizes
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

**3. System Resources:**
```bash
# CPU and Memory usage
top -p $(pgrep ritel-app)

# Disk space
df -h

# Check logs
sudo journalctl -u ritel-app -n 100 --no-pager
```

### Log Monitoring

**Application Logs:**
```bash
# View real-time logs
sudo journalctl -u ritel-app -f

# View last 100 lines
sudo journalctl -u ritel-app -n 100

# Filter by time
sudo journalctl -u ritel-app --since "1 hour ago"

# Filter by priority (errors only)
sudo journalctl -u ritel-app -p err
```

**Nginx Logs:**
```bash
# Access logs
sudo tail -f /var/log/nginx/ritel-app-access.log

# Error logs
sudo tail -f /var/log/nginx/ritel-app-error.log
```

### Performance Monitoring

**Setup Prometheus (Optional):**
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'ritel-app'
    static_configs:
      - targets: ['localhost:8080']
```

### Maintenance Tasks

**Daily:**
- ✅ Check application health
- ✅ Monitor disk space
- ✅ Review error logs

**Weekly:**
- ✅ Verify backups
- ✅ Check database size
- ✅ Review slow queries

**Monthly:**
- ✅ Update dependencies
- ✅ Security patches
- ✅ Performance optimization
- ✅ Test backup restore

---

## 🆘 Troubleshooting

### Application Won't Start

**Check logs:**
```bash
sudo journalctl -u ritel-app -n 50
```

**Common issues:**
- Database connection failed → Check .env DSN
- Port already in use → Check `sudo lsof -i :8080`
- Permission denied → Check file ownership
- Missing .env file → Create from .env.example

### Database Connection Issues

**Test PostgreSQL:**
```bash
psql -h localhost -p 5432 -U ritel -d ritel_db
```

**Check PostgreSQL service:**
```bash
sudo systemctl status postgresql
sudo systemctl restart postgresql
```

### Performance Issues

**Check connections:**
```sql
SELECT count(*) FROM pg_stat_activity;
```

**Optimize PostgreSQL:**
```sql
-- Analyze tables
ANALYZE;

-- Reindex
REINDEX DATABASE ritel_db;

-- Vacuum
VACUUM ANALYZE;
```

---

## 📞 Support

For issues or questions:
- Check logs first
- Review this documentation
- Contact: your-support-email@domain.com

---

## 📝 Changelog

### Version 1.0.0 (2025-12-19)
- Initial production deployment guide
- Multi-scenario deployment options
- Security hardening guidelines
- Automated backup procedures
- Monitoring and maintenance procedures

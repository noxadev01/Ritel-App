#!/bin/bash
# Ritel-App Production Deployment Script
# For Ubuntu/Debian servers

set -e  # Exit on error

echo "=========================================="
echo "🚀 Ritel-App Production Deployment"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="ritel-app"
APP_DIR="/opt/ritel-app"
DB_NAME="ritel_db"
DB_USER="ritel"
DOMAIN=""

# Functions
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    print_error "Please run as root (sudo ./deploy-production.sh)"
    exit 1
fi

# Get domain name
read -p "Enter your domain name (e.g., pos.example.com): " DOMAIN
if [ -z "$DOMAIN" ]; then
    print_error "Domain name is required"
    exit 1
fi

# Get database password
read -sp "Enter PostgreSQL password for 'ritel' user: " DB_PASSWORD
echo
if [ -z "$DB_PASSWORD" ]; then
    print_error "Database password is required"
    exit 1
fi

# Get JWT secret
print_info "Generating JWT secret..."
JWT_SECRET=$(openssl rand -base64 32)

echo ""
echo "=========================================="
echo "📦 Installing Dependencies"
echo "=========================================="

# Update system
print_info "Updating system packages..."
apt-get update -qq
apt-get upgrade -y -qq

# Install PostgreSQL
if ! command -v psql &> /dev/null; then
    print_info "Installing PostgreSQL..."
    apt-get install -y postgresql postgresql-contrib
    print_success "PostgreSQL installed"
else
    print_success "PostgreSQL already installed"
fi

# Install Nginx
if ! command -v nginx &> /dev/null; then
    print_info "Installing Nginx..."
    apt-get install -y nginx
    print_success "Nginx installed"
else
    print_success "Nginx already installed"
fi

# Install Certbot
if ! command -v certbot &> /dev/null; then
    print_info "Installing Certbot..."
    apt-get install -y certbot python3-certbot-nginx
    print_success "Certbot installed"
else
    print_success "Certbot already installed"
fi

echo ""
echo "=========================================="
echo "🗄️  Setting up PostgreSQL Database"
echo "=========================================="

# Start PostgreSQL
systemctl start postgresql
systemctl enable postgresql
print_success "PostgreSQL service started"

# Create database and user
print_info "Creating database and user..."
sudo -u postgres psql << EOF
-- Drop existing if any
DROP DATABASE IF EXISTS $DB_NAME;
DROP USER IF EXISTS $DB_USER;

-- Create new
CREATE DATABASE $DB_NAME;
CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;

-- Connect to database
\c $DB_NAME

-- Grant schema permissions
GRANT ALL ON SCHEMA public TO $DB_USER;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO $DB_USER;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO $DB_USER;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $DB_USER;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO $DB_USER;
EOF

print_success "Database created: $DB_NAME"

# Run schema
if [ -f "database/schema_postgres.sql" ]; then
    print_info "Creating database schema..."
    sudo -u postgres psql -d $DB_NAME -f database/schema_postgres.sql
    print_success "Schema created"
fi

# Run seed data
if [ -f "database/seed_data.sql" ]; then
    print_info "Inserting seed data..."
    sudo -u postgres psql -d $DB_NAME -f database/seed_data.sql
    print_success "Seed data inserted"
fi

echo ""
echo "=========================================="
echo "📁 Setting up Application Directory"
echo "=========================================="

# Create app directory
mkdir -p $APP_DIR
print_success "Created directory: $APP_DIR"

# Copy application files
if [ -f "build/bin/ritel-app" ]; then
    cp build/bin/ritel-app $APP_DIR/
    print_success "Copied application binary"
elif [ -f "ritel-app" ]; then
    cp ritel-app $APP_DIR/
    print_success "Copied application binary"
else
    print_error "Application binary not found! Please build first."
    exit 1
fi

# Make executable
chmod +x $APP_DIR/ritel-app
print_success "Made binary executable"

# Create .env file
cat > $APP_DIR/.env << EOF
# Database Configuration
DB_DRIVER=postgres
DB_DSN=host=localhost port=5432 user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable

# Web Server Configuration
WEB_ENABLED=true
WEB_PORT=8080
WEB_HOST=127.0.0.1

# JWT Configuration
JWT_SECRET=$JWT_SECRET
JWT_EXPIRY_HOURS=8

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://$DOMAIN,https://www.$DOMAIN
CORS_ALLOW_CREDENTIALS=true
EOF

chmod 600 $APP_DIR/.env
print_success "Created .env configuration"

echo ""
echo "=========================================="
echo "🔧 Setting up Systemd Service"
echo "=========================================="

# Create systemd service
cat > /etc/systemd/system/$APP_NAME.service << EOF
[Unit]
Description=Ritel-App POS System
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=$APP_DIR
ExecStart=$APP_DIR/ritel-app
Restart=always
RestartSec=10

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$APP_DIR

# Environment
Environment="HOME=$APP_DIR"

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=$APP_NAME

[Install]
WantedBy=multi-user.target
EOF

print_success "Created systemd service"

# Set ownership
chown -R www-data:www-data $APP_DIR

# Reload systemd
systemctl daemon-reload

# Enable and start service
systemctl enable $APP_NAME
systemctl start $APP_NAME
print_success "Service enabled and started"

# Wait for app to start
sleep 3

# Check service status
if systemctl is-active --quiet $APP_NAME; then
    print_success "Application is running"
else
    print_error "Application failed to start. Check logs: journalctl -u $APP_NAME"
    exit 1
fi

echo ""
echo "=========================================="
echo "🌐 Setting up Nginx Reverse Proxy"
echo "=========================================="

# Build frontend if exists
if [ -d "frontend" ]; then
    print_info "Building frontend..."
    cd frontend
    if command -v npm &> /dev/null; then
        npm install --silent
        npm run build --silent
        mkdir -p $APP_DIR/frontend
        cp -r dist $APP_DIR/frontend/
        print_success "Frontend built and copied"
    else
        print_info "npm not found, skipping frontend build"
    fi
    cd ..
fi

# Create nginx config
cat > /etc/nginx/sites-available/$APP_NAME << 'NGINXEOF'
server {
    listen 80;
    server_name DOMAIN_PLACEHOLDER www.DOMAIN_PLACEHOLDER;

    # Redirect to HTTPS (will be enabled after SSL setup)
    # return 301 https://$server_name$request_uri;

    # Temporary: Allow HTTP for initial setup
    location / {
        root /opt/ritel-app/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

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
    }

    location /health {
        proxy_pass http://127.0.0.1:8080;
    }
}
NGINXEOF

# Replace domain placeholder
sed -i "s/DOMAIN_PLACEHOLDER/$DOMAIN/g" /etc/nginx/sites-available/$APP_NAME

# Enable site
ln -sf /etc/nginx/sites-available/$APP_NAME /etc/nginx/sites-enabled/

# Test nginx config
if nginx -t; then
    print_success "Nginx configuration valid"
    systemctl restart nginx
    print_success "Nginx restarted"
else
    print_error "Nginx configuration error"
    exit 1
fi

echo ""
echo "=========================================="
echo "🔒 Setting up SSL Certificate"
echo "=========================================="

print_info "Obtaining SSL certificate from Let's Encrypt..."
print_info "Make sure your domain $DOMAIN points to this server's IP"
read -p "Press Enter when ready to continue..."

# Get SSL certificate
if certbot --nginx -d $DOMAIN -d www.$DOMAIN --non-interactive --agree-tos --register-unsafely-without-email; then
    print_success "SSL certificate obtained"
else
    print_info "SSL setup failed. You can run manually later: certbot --nginx -d $DOMAIN -d www.$DOMAIN"
fi

echo ""
echo "=========================================="
echo "🔥 Setting up Firewall"
echo "=========================================="

# Setup UFW firewall
if command -v ufw &> /dev/null; then
    print_info "Configuring firewall..."
    ufw --force enable
    ufw allow 22/tcp    # SSH
    ufw allow 80/tcp    # HTTP
    ufw allow 443/tcp   # HTTPS
    print_success "Firewall configured"
else
    print_info "UFW not installed, skipping firewall setup"
fi

echo ""
echo "=========================================="
echo "💾 Setting up Automated Backups"
echo "=========================================="

# Create backup directory
BACKUP_DIR="/backups/ritel-app"
mkdir -p $BACKUP_DIR

# Create backup script
cat > $APP_DIR/backup.sh << 'BACKUPEOF'
#!/bin/bash
BACKUP_DIR="/backups/ritel-app"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="ritel_db"
DB_USER="ritel"

mkdir -p $BACKUP_DIR
pg_dump -U $DB_USER -d $DB_NAME | gzip > $BACKUP_DIR/backup_$DATE.sql.gz
find $BACKUP_DIR -name "*.gz" -mtime +30 -delete
echo "Backup completed: $BACKUP_DIR/backup_$DATE.sql.gz"
BACKUPEOF

chmod +x $APP_DIR/backup.sh
print_success "Backup script created"

# Add to crontab (daily at 2 AM)
(crontab -l 2>/dev/null; echo "0 2 * * * $APP_DIR/backup.sh >> /var/log/ritel-backup.log 2>&1") | crontab -
print_success "Backup scheduled daily at 2 AM"

echo ""
echo "=========================================="
echo "✅ Deployment Complete!"
echo "=========================================="
echo ""
print_success "Application URL: https://$DOMAIN"
print_success "API Endpoint: https://$DOMAIN/api"
print_success "Default Login: admin / admin123"
echo ""
print_info "Next Steps:"
echo "  1. Visit https://$DOMAIN and login"
echo "  2. Change admin password immediately"
echo "  3. Configure printer settings"
echo "  4. Add products and start selling!"
echo ""
print_info "Useful Commands:"
echo "  Status:  systemctl status $APP_NAME"
echo "  Logs:    journalctl -u $APP_NAME -f"
echo "  Restart: systemctl restart $APP_NAME"
echo "  Backup:  $APP_DIR/backup.sh"
echo ""
print_info "Configuration file: $APP_DIR/.env"
print_info "Database: $DB_NAME (PostgreSQL)"
echo "=========================================="

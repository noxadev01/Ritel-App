@echo off
REM Ritel-App Windows Production Setup Script
REM For Windows Server or Desktop deployment

echo ==========================================
echo 🚀 Ritel-App Production Setup (Windows)
echo ==========================================
echo.

REM Check for admin privileges
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo ✗ Error: Please run as Administrator
    pause
    exit /b 1
)

REM Configuration
set APP_NAME=ritel-app
set INSTALL_DIR=C:\ritel-app-production
set DATA_DIR=%INSTALL_DIR%\data
set BACKUP_DIR=D:\backups\ritel-app

echo 📋 Setup Configuration:
echo    Installation: %INSTALL_DIR%
echo    Data Dir:     %DATA_DIR%
echo    Backup Dir:   %BACKUP_DIR%
echo.

REM Ask for database choice
echo Choose database type:
echo   1. SQLite (Standalone, no setup needed)
echo   2. PostgreSQL (Multi-user, requires PostgreSQL installed)
choice /c 12 /n /m "Enter choice (1 or 2): "
set DB_CHOICE=%errorlevel%

echo.
echo ==========================================
echo 📁 Creating Directories
echo ==========================================

REM Create directories
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
if not exist "%DATA_DIR%" mkdir "%DATA_DIR%"
if not exist "%BACKUP_DIR%" mkdir "%BACKUP_DIR%"
echo ✓ Directories created

echo.
echo ==========================================
echo 📦 Copying Application Files
echo ==========================================

REM Copy executable
if exist "build\bin\ritel-app.exe" (
    copy /Y "build\bin\ritel-app.exe" "%INSTALL_DIR%\"
    echo ✓ Copied from build\bin\ritel-app.exe
) else if exist "ritel-app.exe" (
    copy /Y "ritel-app.exe" "%INSTALL_DIR%\"
    echo ✓ Copied ritel-app.exe
) else (
    echo ✗ Error: ritel-app.exe not found!
    echo    Please build the application first.
    pause
    exit /b 1
)

echo.
echo ==========================================
echo ⚙️  Creating Configuration File
echo ==========================================

if %DB_CHOICE%==1 (
    REM SQLite configuration
    (
        echo # Database Configuration
        echo DB_DRIVER=sqlite3
        echo DB_DSN=./data/ritel.db
        echo.
        echo # Web Server Configuration
        echo WEB_ENABLED=false
        echo WEB_PORT=8080
        echo WEB_HOST=0.0.0.0
        echo.
        echo # JWT Configuration
        echo JWT_SECRET=ChangeThisToSecureRandomString123!
        echo JWT_EXPIRY_HOURS=24
        echo.
        echo # CORS Configuration
        echo CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
        echo CORS_ALLOW_CREDENTIALS=true
    ) > "%INSTALL_DIR%\.env"
    echo ✓ Created .env for SQLite
) else (
    REM PostgreSQL configuration
    set /p PG_HOST="Enter PostgreSQL host (default: localhost): "
    if "%PG_HOST%"=="" set PG_HOST=localhost

    set /p PG_PORT="Enter PostgreSQL port (default: 5432): "
    if "%PG_PORT%"=="" set PG_PORT=5432

    set /p PG_USER="Enter PostgreSQL user (default: ritel): "
    if "%PG_USER%"=="" set PG_USER=ritel

    set /p PG_PASSWORD="Enter PostgreSQL password: "
    if "%PG_PASSWORD%"=="" (
        echo ✗ Password cannot be empty
        pause
        exit /b 1
    )

    set /p PG_DATABASE="Enter database name (default: ritel_db): "
    if "%PG_DATABASE%"=="" set PG_DATABASE=ritel_db

    (
        echo # Database Configuration
        echo DB_DRIVER=postgres
        echo DB_DSN=host=%PG_HOST% port=%PG_PORT% user=%PG_USER% password=%PG_PASSWORD% dbname=%PG_DATABASE% sslmode=disable
        echo.
        echo # Web Server Configuration
        echo WEB_ENABLED=true
        echo WEB_PORT=8080
        echo WEB_HOST=0.0.0.0
        echo.
        echo # JWT Configuration
        echo JWT_SECRET=ChangeThisToSecureRandomString123!
        echo JWT_EXPIRY_HOURS=8
        echo.
        echo # CORS Configuration
        echo CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
        echo CORS_ALLOW_CREDENTIALS=true
    ) > "%INSTALL_DIR%\.env"
    echo ✓ Created .env for PostgreSQL
)

echo.
echo ==========================================
echo 🔧 Setting up Windows Service
echo ==========================================

REM Create service wrapper using NSSM (if available)
where nssm >nul 2>&1
if %errorLevel% equ 0 (
    echo ℹ Installing Windows Service...
    nssm install %APP_NAME% "%INSTALL_DIR%\ritel-app.exe"
    nssm set %APP_NAME% AppDirectory "%INSTALL_DIR%"
    nssm set %APP_NAME% DisplayName "Ritel-App POS System"
    nssm set %APP_NAME% Description "Point of Sale and Inventory Management System"
    nssm set %APP_NAME% Start SERVICE_AUTO_START
    nssm set %APP_NAME% AppExit Default Restart
    echo ✓ Service installed: %APP_NAME%
    echo ℹ Starting service...
    nssm start %APP_NAME%
    echo ✓ Service started
) else (
    echo ℹ NSSM not found. Service not installed.
    echo    You can install NSSM from: https://nssm.cc/download
    echo    Or run the app manually with startup shortcut.
)

echo.
echo ==========================================
echo 🔄 Creating Startup Shortcut
echo ==========================================

REM Create startup shortcut
set STARTUP_FOLDER=%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup
powershell -Command "$WS = New-Object -ComObject WScript.Shell; $Shortcut = $WS.CreateShortcut('%STARTUP_FOLDER%\Ritel-App.lnk'); $Shortcut.TargetPath = '%INSTALL_DIR%\ritel-app.exe'; $Shortcut.WorkingDirectory = '%INSTALL_DIR%'; $Shortcut.Save()"
echo ✓ Startup shortcut created

echo.
echo ==========================================
echo 💾 Creating Backup Script
echo ==========================================

REM Create backup script
(
    echo @echo off
    echo REM Automated Backup Script
    echo set BACKUP_DIR=%BACKUP_DIR%
    echo set DATE=%%date:~-4,4%%%%date:~-10,2%%%%date:~-7,2%%_%%time:~0,2%%%%time:~3,2%%
    echo set DATE=%%DATE: =0%%
    echo.
    echo echo Backing up Ritel-App data...
    echo.
    echo REM Backup SQLite database
    echo if exist "%DATA_DIR%\ritel.db" (
    echo     copy /Y "%DATA_DIR%\ritel.db" "%%BACKUP_DIR%%\sqlite_%%DATE%%.db"
    echo     echo ✓ SQLite backup: %%BACKUP_DIR%%\sqlite_%%DATE%%.db
    echo ^)
    echo.
    echo REM Backup .env file
    echo copy /Y "%INSTALL_DIR%\.env" "%%BACKUP_DIR%%\env_%%DATE%%.txt"
    echo echo ✓ Config backup: %%BACKUP_DIR%%\env_%%DATE%%.txt
    echo.
    echo REM Cleanup old backups (keep last 30 days^)
    echo forfiles /p "%%BACKUP_DIR%%" /m *.db /d -30 /c "cmd /c del @path" 2^>nul
    echo forfiles /p "%%BACKUP_DIR%%" /m *.txt /d -30 /c "cmd /c del @path" 2^>nul
    echo.
    echo echo Backup completed!
) > "%INSTALL_DIR%\backup.bat"

echo ✓ Backup script created: %INSTALL_DIR%\backup.bat

echo.
echo ℹ Setting up scheduled backup...
schtasks /create /tn "Ritel-App Backup" /tr "%INSTALL_DIR%\backup.bat" /sc daily /st 02:00 /f >nul 2>&1
if %errorLevel% equ 0 (
    echo ✓ Backup scheduled daily at 2:00 AM
) else (
    echo ⚠ Could not schedule backup. Run manually: %INSTALL_DIR%\backup.bat
)

echo.
echo ==========================================
echo 🔥 Windows Firewall Configuration
echo ==========================================

if %DB_CHOICE%==2 (
    echo ℹ Adding firewall rule for web access...
    netsh advfirewall firewall add rule name="Ritel-App Web Server" dir=in action=allow protocol=TCP localport=8080 >nul 2>&1
    echo ✓ Firewall rule added for port 8080
)

echo.
echo ==========================================
echo 📝 Creating README
echo ==========================================

(
    echo Ritel-App Production Installation
    echo ==================================
    echo.
    echo Installation Directory: %INSTALL_DIR%
    echo Data Directory: %DATA_DIR%
    echo Backup Directory: %BACKUP_DIR%
    echo.
    echo How to Use:
    echo -----------
    echo 1. Double-click ritel-app.exe to start
    echo 2. Login with: admin / admin123
    echo 3. IMPORTANT: Change admin password after first login!
    echo.
    echo Configuration:
    echo -------------
    echo - Config file: %INSTALL_DIR%\.env
    echo - To change settings, edit .env and restart app
    echo.
    echo Backup:
    echo -------
    echo - Manual backup: Run %INSTALL_DIR%\backup.bat
    echo - Automated: Daily at 2:00 AM
    echo - Backup location: %BACKUP_DIR%
    echo.
    echo Service Management:
    echo ------------------
    echo - Start:   nssm start %APP_NAME%
    echo - Stop:    nssm stop %APP_NAME%
    echo - Restart: nssm restart %APP_NAME%
    echo - Status:  nssm status %APP_NAME%
    echo.
    echo Troubleshooting:
    echo ---------------
    echo - Check logs in: %INSTALL_DIR%\logs\
    echo - Test database: Try opening the app
    echo - Service issues: Check Windows Event Viewer
    echo.
    echo Web Access (if enabled^):
    echo ------------------------
    echo - URL: http://localhost:8080
    echo - Or: http://YOUR_IP:8080 (from other devices^)
    echo.
) > "%INSTALL_DIR%\README.txt"

echo ✓ README created: %INSTALL_DIR%\README.txt

echo.
echo ==========================================
echo ✅ Installation Complete!
echo ==========================================
echo.
echo 🎉 Ritel-App has been installed successfully!
echo.
echo 📍 Installation: %INSTALL_DIR%
if %DB_CHOICE%==1 (
    echo 💾 Database: SQLite (Local)
) else (
    echo 💾 Database: PostgreSQL (%PG_HOST%:%PG_PORT%)
    echo 🌐 Web Access: http://localhost:8080
)
echo.
echo 🔐 Default Login:
echo    Username: admin
echo    Password: admin123
echo.
echo ⚠️  IMPORTANT:
echo    1. Change admin password immediately after first login
echo    2. Backup regularly: %INSTALL_DIR%\backup.bat
echo    3. Keep .env file secure (contains passwords)
echo.
echo 📚 Next Steps:
echo    1. Run the application: %INSTALL_DIR%\ritel-app.exe
echo    2. Login and change password
echo    3. Configure printer settings
echo    4. Add your products
echo    5. Start selling!
echo.
echo 📖 Documentation: %INSTALL_DIR%\README.txt
echo.
echo ==========================================

REM Open README
start notepad "%INSTALL_DIR%\README.txt"

echo.
echo Press any key to start the application...
pause >nul

REM Start the application
start "" "%INSTALL_DIR%\ritel-app.exe"

echo.
echo Application started! Check the system tray.
echo.
pause

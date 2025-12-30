@echo off
echo ========================================
echo Starting Ritel-App Web Server
echo ========================================
echo.

REM Check if PostgreSQL is running
echo [1/4] Checking PostgreSQL connection...
pg_isready -h localhost -p 5432 >nul 2>&1
if errorlevel 1 (
    echo [WARNING] PostgreSQL is not running or not accessible
    echo [INFO] Make sure PostgreSQL is installed and running on localhost:5432
    echo.
    echo Continue anyway? (Press Ctrl+C to cancel, or
    pause
)

REM Backup current .env
echo [2/4] Preparing configuration...
if exist .env (
    copy /Y .env .env.backup >nul
    echo [INFO] Current .env backed up to .env.backup
)

REM Copy web configuration
copy /Y .env.web .env >nul
echo [INFO] Using web-only configuration (.env.web)

REM Build and run
echo.
echo [3/4] Building web server...
go build -o web-server.exe web-server.go
if errorlevel 1 (
    echo [ERROR] Build failed!
    pause
    exit /b 1
)

echo [4/4] Starting web server...
echo.
echo ========================================
echo Web Server Starting...
echo Access at: http://localhost:8080
echo Press Ctrl+C to stop
echo ========================================
echo.

web-server.exe

REM Restore original .env on exit
echo.
echo ========================================
echo Restoring original configuration...
if exist .env.backup (
    copy /Y .env.backup .env >nul
    del .env.backup >nul
    echo [INFO] Original .env restored
)
echo ========================================
pause

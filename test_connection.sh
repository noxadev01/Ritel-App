#!/bin/bash
# Script untuk test koneksi database

echo "========================================="
echo "   TEST DATABASE CONNECTION"
echo "========================================="
echo ""

# Check .env file
if [ ! -f ".env" ]; then
    echo "❌ File .env tidak ditemukan!"
    echo ""
    echo "Solusi:"
    echo "  cp .env.dual .env      # Untuk dual mode"
    echo "  cp .env.postgres .env  # Untuk PostgreSQL"
    echo "  cp .env.sqlite .env    # Untuk SQLite"
    exit 1
fi

echo "✓ File .env ditemukan"
echo ""

# Show current config
echo "Konfigurasi saat ini:"
grep "^DB_DRIVER" .env
echo ""

# Try to run app briefly
echo "Mencoba menjalankan aplikasi..."
echo "========================================="
timeout 5 ./ritel-app.exe 2>&1 &
PID=$!

# Wait a bit
sleep 3

# Check if still running
if ps -p $PID > /dev/null; then
    echo ""
    echo "✓ Aplikasi berhasil start!"
    echo "========================================="
    kill $PID 2>/dev/null
else
    echo ""
    echo "⚠️  Aplikasi mungkin ada error"
    echo "Cek output di atas untuk detail error"
    echo "========================================="
fi

echo ""
echo "Untuk test lengkap, jalankan: ./ritel-app.exe"

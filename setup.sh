#!/bin/bash

echo "========================================"
echo "   SETUP RITEL APP"
echo "========================================"
echo ""

echo "Pilih mode database:"
echo "1. SQLite (Paling Mudah - Recommended untuk Testing)"
echo "2. PostgreSQL (Production)"
echo "3. Dual Mode (PostgreSQL + SQLite - Recommended untuk Production)"
echo ""

read -p "Pilihan Anda (1/2/3): " choice

case $choice in
    1)
        echo ""
        echo "[1/2] Membuat file .env untuk SQLite..."
        cp .env.sqlite .env
        echo "✓ File .env berhasil dibuat"
        echo ""
        echo "[2/2] Setup selesai!"
        echo ""
        echo "========================================"
        echo "   SETUP SELESAI - SQLite Mode"
        echo "========================================"
        echo ""
        echo "Mode: SQLite"
        echo "Database: ~/ritel-app/ritel.db"
        echo ""
        echo "Jalankan aplikasi dengan: ./ritel-app"
        echo "========================================"
        ;;
    2)
        echo ""
        echo "[1/3] Membuat file .env untuk PostgreSQL..."
        cp .env.postgres .env
        echo "✓ File .env berhasil dibuat"
        echo ""
        echo "[2/3] Cek PostgreSQL..."
        echo "⚠ PASTIKAN PostgreSQL sudah terinstall dan running"
        echo ""
        echo "[3/3] Buat database..."
        echo "Jalankan command: createdb ritel_db"
        echo ""
        echo "========================================"
        echo "   SETUP SELESAI - PostgreSQL Mode"
        echo "========================================"
        echo ""
        echo "Mode: PostgreSQL"
        echo "Database: ritel_db"
        echo ""
        echo "NEXT STEPS:"
        echo "1. Pastikan PostgreSQL running"
        echo "2. Buat database: createdb ritel_db"
        echo "3. Jalankan aplikasi: ./ritel-app"
        echo "========================================"
        ;;
    3)
        echo ""
        echo "[1/3] Membuat file .env untuk Dual Mode..."
        cp .env.dual .env
        echo "✓ File .env berhasil dibuat"
        echo ""
        echo "[2/3] Cek PostgreSQL..."
        echo "⚠ PASTIKAN PostgreSQL sudah terinstall dan running"
        echo ""
        echo "[3/3] Buat database..."
        echo "Jalankan command: createdb ritel_db"
        echo ""
        echo "========================================"
        echo "   SETUP SELESAI - Dual Mode"
        echo "========================================"
        echo ""
        echo "Mode: Dual Database (PostgreSQL + SQLite)"
        echo "Primary: PostgreSQL (ritel_db)"
        echo "Backup: SQLite (~/ritel-app/ritel.db)"
        echo ""
        echo "KEUNTUNGAN DUAL MODE:"
        echo "✓ Data otomatis tersimpan di 2 database"
        echo "✓ Backup real-time ke SQLite"
        echo "✓ Keamanan data maksimal"
        echo ""
        echo "NEXT STEPS:"
        echo "1. Pastikan PostgreSQL running"
        echo "2. Buat database: createdb ritel_db"
        echo "3. Jalankan aplikasi: ./ritel-app"
        echo "========================================"
        ;;
    *)
        echo ""
        echo "❌ Pilihan tidak valid!"
        echo ""
        ;;
esac

echo ""

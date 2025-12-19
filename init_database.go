// init_database.go
// Initialize SQLite database with complete schema
package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Check if database exists
	if _, err := os.Stat("./ritel.db"); err == nil {
		log.Println("⚠️  Database ritel.db already exists")
		log.Print("Do you want to recreate it? This will DELETE all existing data! (y/N): ")
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil || (response != "y" && response != "Y") {
			log.Println("Cancelled. Database not modified.")
			return
		}
		os.Remove("./ritel.db")
		log.Println("Deleted existing database")
	}

	// Open database (will create if not exists)
	db, err := sql.Open("sqlite3", "./ritel.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Connected to database successfully")

	// Schema SQLite (converted from PostgreSQL schema)
	schema := `
-- Migrations tracking table
CREATE TABLE IF NOT EXISTS migrations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Kategori table
CREATE TABLE IF NOT EXISTS kategori (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT UNIQUE NOT NULL,
    deskripsi TEXT,
    icon TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Produk table
CREATE TABLE IF NOT EXISTS produk (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sku TEXT UNIQUE NOT NULL,
    barcode TEXT UNIQUE,
    nama TEXT NOT NULL,
    kategori TEXT,
    berat REAL DEFAULT 0,
    harga_beli INTEGER DEFAULT 0,
    harga_jual INTEGER NOT NULL,
    stok REAL DEFAULT 0,
    satuan TEXT DEFAULT 'kg',
    jenis_produk TEXT DEFAULT 'curah',
    kadaluarsa TEXT,
    tanggal_masuk TEXT,
    deskripsi TEXT,
    gambar TEXT,
    masa_simpan_hari INTEGER DEFAULT 0,
    hari_pemberitahuan_kadaluarsa INTEGER DEFAULT 30,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Keranjang table
CREATE TABLE IF NOT EXISTS keranjang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    produk_id INTEGER NOT NULL,
    jumlah INTEGER DEFAULT 1,
    harga_beli INTEGER DEFAULT 0,
    subtotal INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (produk_id) REFERENCES produk(id) ON DELETE RESTRICT
);

-- Pelanggan table
CREATE TABLE IF NOT EXISTS pelanggan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL,
    telepon TEXT UNIQUE NOT NULL,
    email TEXT,
    tipe TEXT DEFAULT 'reguler',
    level INTEGER DEFAULT 1,
    poin INTEGER DEFAULT 0,
    diskon_persen INTEGER DEFAULT 0,
    total_transaksi INTEGER DEFAULT 0,
    total_belanja INTEGER DEFAULT 0,
    alamat TEXT,
    last_transaction_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    nama_lengkap TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'staff',
    status TEXT DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Transaksi table
CREATE TABLE IF NOT EXISTS transaksi (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nomor_transaksi TEXT UNIQUE NOT NULL,
    tanggal TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    pelanggan_id INTEGER,
    pelanggan_nama TEXT,
    pelanggan_telp TEXT,
    staff_id INTEGER,
    staff_nama TEXT,
    subtotal INTEGER DEFAULT 0,
    diskon INTEGER DEFAULT 0,
    diskon_promo INTEGER DEFAULT 0,
    diskon_pelanggan INTEGER DEFAULT 0,
    diskon_poin INTEGER DEFAULT 0,
    poin_ditukar INTEGER DEFAULT 0,
    total INTEGER NOT NULL,
    total_bayar INTEGER DEFAULT 0,
    kembalian INTEGER DEFAULT 0,
    status TEXT DEFAULT 'selesai',
    catatan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (staff_id) REFERENCES users(id) ON DELETE RESTRICT
);

-- Transaksi Item table
CREATE TABLE IF NOT EXISTS transaksi_item (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaksi_id INTEGER NOT NULL,
    produk_id INTEGER,
    produk_sku TEXT,
    produk_nama TEXT,
    produk_kategori TEXT,
    harga_satuan INTEGER NOT NULL,
    jumlah INTEGER DEFAULT 1,
    beratgram REAL DEFAULT 0,
    subtotal INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaksi_id) REFERENCES transaksi(id) ON DELETE RESTRICT
);

-- Promo table
CREATE TABLE IF NOT EXISTS promo (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL,
    kode TEXT UNIQUE,
    tipe TEXT NOT NULL,
    tipe_promo TEXT DEFAULT 'diskon_produk',
    tipe_produk_berlaku TEXT DEFAULT 'semua',
    nilai INTEGER NOT NULL,
    min_quantity INTEGER DEFAULT 0,
    max_diskon INTEGER DEFAULT 0,
    buy_quantity INTEGER DEFAULT 0,
    get_quantity INTEGER DEFAULT 0,
    tipe_buy_get TEXT DEFAULT 'sama',
    harga_bundling INTEGER DEFAULT 0,
    tipe_bundling TEXT DEFAULT 'harga_tetap',
    diskon_bundling INTEGER DEFAULT 0,
    produk_x INTEGER,
    produk_y INTEGER,
    tanggal_mulai TIMESTAMP,
    tanggal_selesai TIMESTAMP,
    status TEXT DEFAULT 'aktif',
    deskripsi TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (produk_x) REFERENCES produk(id) ON DELETE SET NULL,
    FOREIGN KEY (produk_y) REFERENCES produk(id) ON DELETE SET NULL
);

-- Promo Produk table
CREATE TABLE IF NOT EXISTS promo_produk (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    promo_id INTEGER NOT NULL,
    produk_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (promo_id) REFERENCES promo(id) ON DELETE CASCADE,
    FOREIGN KEY (produk_id) REFERENCES produk(id) ON DELETE CASCADE
);

-- Batch table
CREATE TABLE IF NOT EXISTS batch (
    id TEXT PRIMARY KEY,
    produk_id INTEGER NOT NULL,
    qty REAL NOT NULL DEFAULT 0,
    qty_tersisa REAL NOT NULL DEFAULT 0,
    tanggal_restok TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    masa_simpan_hari INTEGER DEFAULT 0,
    tanggal_kadaluarsa TIMESTAMP,
    status TEXT DEFAULT 'fresh',
    supplier TEXT,
    keterangan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (produk_id) REFERENCES produk(id) ON DELETE CASCADE
);

-- Stok History table
CREATE TABLE IF NOT EXISTS stok_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    produk_id INTEGER NOT NULL,
    stok_sebelum REAL NOT NULL,
    stok_sesudah REAL NOT NULL,
    perubahan REAL NOT NULL,
    jenis_perubahan TEXT NOT NULL,
    tipe_kerugian TEXT,
    nilai_kerugian INTEGER DEFAULT 0,
    keterangan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (produk_id) REFERENCES produk(id) ON DELETE CASCADE
);

-- Pembayaran table
CREATE TABLE IF NOT EXISTS pembayaran (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaksi_id INTEGER NOT NULL,
    metode TEXT NOT NULL,
    jumlah INTEGER NOT NULL,
    referensi TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaksi_id) REFERENCES transaksi(id) ON DELETE RESTRICT
);

-- Returns table
CREATE TABLE IF NOT EXISTS returns (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaksi_id INTEGER NOT NULL,
    no_transaksi TEXT NOT NULL,
    return_date TIMESTAMP NOT NULL,
    reason TEXT NOT NULL,
    type TEXT NOT NULL,
    replacement_product_id INTEGER,
    refund_amount INTEGER DEFAULT 0,
    refund_method TEXT,
    refund_status TEXT DEFAULT 'pending',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Return Items table
CREATE TABLE IF NOT EXISTS return_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    return_id INTEGER NOT NULL,
    product_id INTEGER,
    product_name TEXT,
    quantity INTEGER NOT NULL,
    price INTEGER NOT NULL,
    subtotal INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (return_id) REFERENCES returns(id) ON DELETE CASCADE
);

-- Print Settings table
CREATE TABLE IF NOT EXISTS print_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    printer_name TEXT,
    paper_width INTEGER DEFAULT 80,
    company_name TEXT,
    company_address TEXT,
    company_phone TEXT,
    show_logo INTEGER DEFAULT 0,
    logo_path TEXT,
    footer_text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Poin Settings table
CREATE TABLE IF NOT EXISTS poin_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    min_transaction INTEGER DEFAULT 25000,
    point_value INTEGER DEFAULT 500,
    min_redeem_points INTEGER DEFAULT 100,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_produk_deleted_at ON produk(deleted_at);
CREATE INDEX IF NOT EXISTS idx_pelanggan_deleted_at ON pelanggan(deleted_at);
CREATE INDEX IF NOT EXISTS idx_promo_deleted_at ON promo(deleted_at);
CREATE INDEX IF NOT EXISTS idx_kategori_deleted_at ON kategori(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_transaksi_tanggal ON transaksi(tanggal);
CREATE INDEX IF NOT EXISTS idx_transaksi_pelanggan ON transaksi(pelanggan_id);
CREATE INDEX IF NOT EXISTS idx_transaksi_staff ON transaksi(staff_id);
`

	// Execute schema
	log.Println("📝 Creating tables and indexes...")
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	log.Println("✅ Schema created successfully")

	// Create default admin user
	log.Println("👤 Creating default admin user...")
	// Password: admin123 (bcrypt hash)
	hashedPassword := "$2a$10$ZK5p3h9vJ8rZ5VmKXJ5fVeYqY3Z9K5p3h9vJ8rZ5VmKXJ5fVe.6K2"
	_, err = db.Exec(`
		INSERT INTO users (username, password, nama_lengkap, role, status)
		VALUES (?, ?, ?, ?, ?)
	`, "admin", hashedPassword, "Administrator", "admin", "active")
	if err != nil {
		log.Printf("⚠️  Could not create admin user (may already exist): %v", err)
	} else {
		log.Println("✅ Default admin user created (username: admin, password: admin123)")
	}

	// Verify tables
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		log.Fatalf("Failed to query tables: %v", err)
	}
	defer rows.Close()

	log.Println("\n📊 Created tables:")
	tableCount := 0
	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		log.Printf("  ✅ %s", tableName)
		tableCount++
	}

	log.Printf("\n🎉 Database initialized successfully!")
	log.Printf("Total tables created: %d", tableCount)
	log.Println("\n▶️  You can now start your application!")
}

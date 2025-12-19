package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"time"

// 	_ "github.com/mattn/go-sqlite3"
// )

// func main() {
// 	// Get user home directory
// 	homeDir, err := os.UserHomeDir()
// 	if err != nil {
// 		log.Fatal("Failed to get home directory:", err)
// 	}

// 	// Database path: ~/ritel-app/ritel.db
// 	dbPath := filepath.Join(homeDir, "ritel-app", "ritel.db")
// 	log.Printf("Database path: %s\n", dbPath)

// 	// Connect to database
// 	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}
// 	defer db.Close()

// 	// Test connection
// 	if err := db.Ping(); err != nil {
// 		log.Fatal("Failed to ping database:", err)
// 	}

// 	log.Println("✅ Connected to database successfully")

// 	// Start transaction
// 	tx, err := db.Begin()
// 	if err != nil {
// 		log.Fatal("Failed to begin transaction:", err)
// 	}
// 	defer tx.Rollback()

// 	// 1. Insert Categories
// 	log.Println("\n📦 Inserting Categories...")
// 	categories := []struct {
// 		nama      string
// 		deskripsi string
// 		icon      string
// 	}{
// 		{"Sayuran", "Berbagai jenis sayuran segar", "🥬"},
// 		{"Buah", "Buah-buahan segar dan berkualitas", "🍎"},
// 		{"Bumbu", "Bumbu dapur dan rempah-rempah", "🌶️"},
// 		{"Umbi", "Umbi-umbian segar", "🥔"},
// 		{"Jamur", "Aneka jamur segar", "🍄"},
// 		{"Kacang", "Kacang-kacangan berkualitas", "🥜"},
// 		{"Biji-bijian", "Biji-bijian dan sereal", "🌾"},
// 		{"Herbal", "Tanaman herbal dan obat", "🌿"},
// 	}

// 	categoryIDs := make(map[string]int)
// 	for _, cat := range categories {
// 		result, err := tx.Exec(`
// 			INSERT INTO kategori (nama, deskripsi, icon, created_at, updated_at)
// 			VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
// 		`, cat.nama, cat.deskripsi, cat.icon)

// 		if err != nil {
// 			log.Printf("⚠️  Category '%s' might already exist, skipping...", cat.nama)
// 			// Get existing ID
// 			var id int
// 			tx.QueryRow("SELECT id FROM kategori WHERE nama = ?", cat.nama).Scan(&id)
// 			categoryIDs[cat.nama] = id
// 		} else {
// 			id, _ := result.LastInsertId()
// 			categoryIDs[cat.nama] = int(id)
// 			log.Printf("   ✓ Added category: %s (ID: %d)", cat.nama, id)
// 		}
// 	}

// 	// 2. Insert Products
// 	log.Println("\n🛒 Inserting Products...")
// 	products := []struct {
// 		sku                         string
// 		barcode                     string
// 		nama                        string
// 		kategori                    string
// 		hargaBeli                   int
// 		hargaJual                   int
// 		stok                        float64
// 		satuan                      string
// 		jenisProduk                 string
// 		masaSimpanHari              int
// 		hariPemberitahuanKadaluarsa int
// 		deskripsi                   string
// 	}{
// 		// Sayuran
// 		{"SKU-SAY-001", "8991002100015", "Bayam Segar", "Sayuran", 3000, 5000, 50, "ikat", "satuan", 3, 1, "Bayam segar pilihan"},
// 		{"SKU-SAY-002", "8991002100022", "Kangkung Segar", "Sayuran", 2500, 4000, 60, "ikat", "satuan", 3, 1, "Kangkung segar"},
// 		{"SKU-SAY-003", "8991002100039", "Sawi Hijau", "Sayuran", 3500, 5500, 40, "ikat", "satuan", 4, 1, "Sawi hijau berkualitas"},
// 		{"SKU-SAY-004", "8991002100046", "Wortel", "Sayuran", 8000, 12000, 0, "kg", "curah", 14, 3, "Wortel segar per kg"},
// 		{"SKU-SAY-005", "8991002100053", "Tomat", "Sayuran", 7000, 10000, 0, "kg", "curah", 7, 2, "Tomat segar per kg"},

// 		// Buah
// 		{"SKU-BUA-001", "8991003100014", "Apel Fuji", "Buah", 15000, 22000, 0, "kg", "curah", 30, 7, "Apel fuji import"},
// 		{"SKU-BUA-002", "8991003100021", "Pisang Cavendish", "Buah", 10000, 15000, 0, "kg", "curah", 7, 2, "Pisang cavendish premium"},
// 		{"SKU-BUA-003", "8991003100038", "Jeruk Medan", "Buah", 12000, 18000, 0, "kg", "curah", 14, 3, "Jeruk medan manis"},
// 		{"SKU-BUA-004", "8991003100045", "Mangga Harum Manis", "Buah", 18000, 25000, 0, "kg", "curah", 10, 2, "Mangga harum manis"},

// 		// Bumbu
// 		{"SKU-BUM-001", "8991004100013", "Cabai Merah Keriting", "Bumbu", 40000, 55000, 0, "kg", "curah", 7, 2, "Cabai merah keriting segar"},
// 		{"SKU-BUM-002", "8991004100020", "Bawang Merah", "Bumbu", 30000, 45000, 0, "kg", "curah", 30, 7, "Bawang merah berkualitas"},
// 		{"SKU-BUM-003", "8991004100037", "Bawang Putih", "Bumbu", 35000, 50000, 0, "kg", "curah", 60, 14, "Bawang putih kating"},
// 		{"SKU-BUM-004", "8991004100044", "Jahe Merah", "Bumbu", 20000, 30000, 0, "kg", "curah", 30, 7, "Jahe merah segar"},

// 		// Umbi
// 		{"SKU-UMB-001", "8991005100012", "Kentang", "Umbi", 8000, 13000, 0, "kg", "curah", 60, 14, "Kentang granola"},
// 		{"SKU-UMB-002", "8991005100029", "Ubi Cilembu", "Umbi", 15000, 22000, 0, "kg", "curah", 30, 7, "Ubi cilembu asli"},

// 		// Jamur
// 		{"SKU-JAM-001", "8991006100011", "Jamur Tiram", "Jamur", 12000, 18000, 25, "pack", "satuan", 5, 1, "Jamur tiram segar 250g"},
// 		{"SKU-JAM-002", "8991006100028", "Jamur Shitake", "Jamur", 25000, 35000, 15, "pack", "satuan", 7, 2, "Jamur shitake premium"},

// 		// Kacang
// 		{"SKU-KAC-001", "8991007100010", "Kacang Panjang", "Kacang", 4000, 7000, 0, "kg", "curah", 5, 1, "Kacang panjang segar"},
// 		{"SKU-KAC-002", "8991007100027", "Kacang Merah", "Kacang", 18000, 25000, 0, "kg", "curah", 365, 30, "Kacang merah kering"},

// 		// Herbal
// 		{"SKU-HER-001", "8991008100019", "Daun Mint", "Herbal", 8000, 12000, 20, "pack", "satuan", 7, 2, "Daun mint segar"},
// 		{"SKU-HER-002", "8991008100026", "Serai", "Herbal", 3000, 5000, 30, "ikat", "satuan", 14, 3, "Serai segar"},
// 	}

// 	productIDs := make(map[string]int)
// 	for _, p := range products {
// 		result, err := tx.Exec(`
// 			INSERT INTO produk (
// 				sku, barcode, nama, kategori, berat, harga_beli, harga_jual,
// 				stok, satuan, jenis_produk, masa_simpan_hari, hari_pemberitahuan_kadaluarsa,
// 				deskripsi, tanggal_masuk, created_at, updated_at
// 			) VALUES (?, ?, ?, ?, 0, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
// 		`, p.sku, p.barcode, p.nama, p.kategori, p.hargaBeli, p.hargaJual,
// 			p.stok, p.satuan, p.jenisProduk, p.masaSimpanHari, p.hariPemberitahuanKadaluarsa,
// 			p.deskripsi, time.Now().Format("2006-01-02"))

// 		if err != nil {
// 			log.Printf("⚠️  Product '%s' might already exist, skipping...", p.nama)
// 			// Get existing ID
// 			var id int
// 			tx.QueryRow("SELECT id FROM produk WHERE barcode = ?", p.barcode).Scan(&id)
// 			productIDs[p.barcode] = id
// 		} else {
// 			id, _ := result.LastInsertId()
// 			productIDs[p.barcode] = int(id)
// 			log.Printf("   ✓ Added product: %s - %s (ID: %d)", p.sku, p.nama, id)
// 		}
// 	}

// 	// 3. Insert Batches for products with curah type
// 	log.Println("\n📊 Inserting Batches...")

// 	// Get product IDs for batch creation
// 	rows, err := tx.Query(`
// 		SELECT id, barcode, nama, masa_simpan_hari, jenis_produk
// 		FROM produk
// 		WHERE jenis_produk = 'curah'
// 	`)
// 	if err != nil {
// 		log.Fatal("Failed to query products:", err)
// 	}
// 	defer rows.Close()

// 	batchCount := 0
// 	for rows.Next() {
// 		var produkID int
// 		var barcode, nama string
// 		var masaSimpanHari int
// 		var jenisProduk string

// 		rows.Scan(&produkID, &barcode, &nama, &masaSimpanHari, &jenisProduk)

// 		// Create 2 batches per product (different dates)
// 		batches := []struct {
// 			qty            float64
// 			tanggalRestok  time.Time
// 			supplier       string
// 		}{
// 			{50.0, time.Now().AddDate(0, 0, -5), "Supplier A"},
// 			{30.0, time.Now().AddDate(0, 0, -2), "Supplier B"},
// 		}

// 		for _, batch := range batches {
// 			tanggalKadaluarsa := batch.tanggalRestok.AddDate(0, 0, masaSimpanHari)

// 			// Calculate status
// 			daysUntilExpiry := int(tanggalKadaluarsa.Sub(time.Now()).Hours() / 24)
// 			var status string
// 			if daysUntilExpiry <= 0 {
// 				status = "expired"
// 			} else if daysUntilExpiry <= 7 {
// 				status = "hampir_expired"
// 			} else {
// 				status = "fresh"
// 			}

// 			_, err := tx.Exec(`
// 				INSERT INTO batch (
// 					id, produk_id, qty, qty_tersisa, tanggal_restok, masa_simpan_hari,
// 					tanggal_kadaluarsa, status, supplier, created_at, updated_at
// 				) VALUES (
// 					lower(hex(randomblob(16))), ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
// 				)
// 			`, produkID, batch.qty, batch.qty, batch.tanggalRestok.Format("2006-01-02"),
// 				masaSimpanHari, tanggalKadaluarsa.Format("2006-01-02"), status, batch.supplier)

// 			if err != nil {
// 				log.Printf("⚠️  Failed to create batch for %s: %v", nama, err)
// 			} else {
// 				batchCount++
// 			}

// 			// Update product stock
// 			tx.Exec("UPDATE produk SET stok = stok + ? WHERE id = ?", batch.qty, produkID)
// 		}
// 	}
// 	log.Printf("   ✓ Created %d batches", batchCount)

// 	// 4. Insert Sample Transactions
// 	log.Println("\n💰 Inserting Sample Transactions...")

// 	// First, ensure we have a default staff user
// 	var staffID int
// 	err = tx.QueryRow("SELECT id FROM users WHERE role = 'staff' LIMIT 1").Scan(&staffID)
// 	if err != nil {
// 		// Create default staff
// 		result, _ := tx.Exec(`
// 			INSERT INTO users (username, password, nama_lengkap, role, status, created_at, updated_at)
// 			VALUES ('staff', '$2a$10$8K1p/a0dL3LKzOWQ5tWZS.qHnC3gJ7wHJxZv9vZm5BqH9VQnPCZES', 'Staff Demo', 'staff', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
// 		`)
// 		id, _ := result.LastInsertId()
// 		staffID = int(id)
// 		log.Printf("   ✓ Created demo staff user (ID: %d)", staffID)
// 	}

// 	// Create 5 sample transactions
// 	transactionCount := 0
// 	for i := 1; i <= 5; i++ {
// 		transaksiDate := time.Now().AddDate(0, 0, -i)
// 		nomorTransaksi := fmt.Sprintf("TRX-%s-%04d", transaksiDate.Format("20060102"), 1000+i)

// 		// Insert transaction header
// 		result, err := tx.Exec(`
// 			INSERT INTO transaksi (
// 				nomor_transaksi, tanggal, staff_id, subtotal, total, total_bayar,
// 				kasir, created_at
// 			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
// 		`, nomorTransaksi, transaksiDate.Format("2006-01-02 15:04:05"), staffID,
// 			0, 0, 0, "Staff Demo", transaksiDate.Format("2006-01-02 15:04:05"))

// 		if err != nil {
// 			log.Printf("⚠️  Failed to create transaction: %v", err)
// 			continue
// 		}

// 		transaksiID, _ := result.LastInsertId()

// 		// Add 2-4 random items to each transaction
// 		itemCount := 2 + i%3
// 		var subtotal int

// 		// Get random products
// 		prodRows, _ := tx.Query("SELECT id, nama, harga_jual, jenis_produk, satuan FROM produk ORDER BY RANDOM() LIMIT ?", itemCount)
// 		for prodRows.Next() {
// 			var prodID int
// 			var prodNama, jenisProduk, satuan string
// 			var hargaJual int

// 			prodRows.Scan(&prodID, &prodNama, &hargaJual, &jenisProduk, &satuan)

// 			var jumlah float64
// 			var beratGram sql.NullInt64
// 			var itemSubtotal int

// 			if jenisProduk == "curah" {
// 				// Random weight between 0.5kg - 3kg
// 				jumlah = 1.0
// 				berat := 500 + (i*200)%2500
// 				beratGram = sql.NullInt64{Int64: int64(berat), Valid: true}
// 				itemSubtotal = (hargaJual * berat) / 1000
// 			} else {
// 				// Random quantity 1-5
// 				jumlah = float64(1 + i%5)
// 				beratGram = sql.NullInt64{Valid: false}
// 				itemSubtotal = hargaJual * int(jumlah)
// 			}

// 			tx.Exec(`
// 				INSERT INTO transaksi_item (
// 					transaksi_id, produk_id, produk_nama, produk_sku, produk_kategori,
// 					harga_satuan, jumlah, beratgram, subtotal
// 				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
// 			`, transaksiID, prodID, prodNama, "", "", hargaJual, jumlah, beratGram, itemSubtotal)

// 			subtotal += itemSubtotal
// 		}
// 		prodRows.Close()

// 		// Update transaction totals
// 		tx.Exec(`
// 			UPDATE transaksi
// 			SET subtotal = ?, total = ?, total_bayar = ?
// 			WHERE id = ?
// 		`, subtotal, subtotal, subtotal, transaksiID)

// 		// Insert payment
// 		tx.Exec(`
// 			INSERT INTO pembayaran (transaksi_id, metode, jumlah)
// 			VALUES (?, 'tunai', ?)
// 		`, transaksiID, subtotal)

// 		transactionCount++
// 		log.Printf("   ✓ Created transaction: %s (Rp %s)", nomorTransaksi, formatRupiah(subtotal))
// 	}

// 	log.Printf("   ✓ Created %d transactions", transactionCount)

// 	// Commit transaction
// 	if err := tx.Commit(); err != nil {
// 		log.Fatal("Failed to commit transaction:", err)
// 	}

// 	log.Println("\n✅ All sample data inserted successfully!")
// 	log.Println("\n📊 Summary:")
// 	log.Printf("   - Categories: %d", len(categories))
// 	log.Printf("   - Products: %d", len(products))
// 	log.Printf("   - Batches: %d", batchCount)
// 	log.Printf("   - Transactions: %d", transactionCount)
// }

// func formatRupiah(amount int) string {
// 	if amount < 1000 {
// 		return fmt.Sprintf("%d", amount)
// 	}
// 	return fmt.Sprintf("%d.%03d", amount/1000, amount%1000)
// }

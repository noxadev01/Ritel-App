package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"

// 	_ "github.com/mattn/go-sqlite3"
// )

// func main() {
// 	// Get user home directory
// 	homeDir, err := os.UserHomeDir()
// 	if err != nil {
// 		log.Fatal("Failed to get home directory:", err)
// 	}

// 	// Database path
// 	dbPath := filepath.Join(homeDir, "ritel-app", "ritel.db")

// 	// Connect to database
// 	db, err := sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}
// 	defer db.Close()

// 	fmt.Println("📊 DATABASE VERIFICATION REPORT")
// 	fmt.Println("═══════════════════════════════════════════════════════")

// 	// 1. Categories
// 	var categoryCount int
// 	db.QueryRow("SELECT COUNT(*) FROM kategori").Scan(&categoryCount)
// 	fmt.Printf("\n✅ KATEGORI: %d categories\n", categoryCount)

// 	rows, _ := db.Query("SELECT id, nama, icon FROM kategori ORDER BY id")
// 	for rows.Next() {
// 		var id int
// 		var nama, icon string
// 		rows.Scan(&id, &nama, &icon)
// 		fmt.Printf("   %d. %s %s\n", id, icon, nama)
// 	}
// 	rows.Close()

// 	// 2. Products
// 	var productCount int
// 	db.QueryRow("SELECT COUNT(*) FROM produk").Scan(&productCount)
// 	fmt.Printf("\n✅ PRODUK: %d products\n", productCount)

// 	rows, _ = db.Query(`
// 		SELECT p.nama, p.kategori, p.stok, p.satuan, p.jenis_produk, p.harga_jual
// 		FROM produk p
// 		ORDER BY p.kategori, p.nama
// 		LIMIT 10
// 	`)
// 	fmt.Println("   Top 10 products:")
// 	for rows.Next() {
// 		var nama, kategori, satuan, jenisProduk string
// 		var stok float64
// 		var harga int
// 		rows.Scan(&nama, &kategori, &stok, &satuan, &jenisProduk, &harga)
// 		if jenisProduk == "curah" {
// 			fmt.Printf("   - %s (%s) - Stok: %.1f %s - Rp %s/%s\n",
// 				nama, kategori, stok, satuan, formatRupiah(harga), satuan)
// 		} else {
// 			fmt.Printf("   - %s (%s) - Stok: %.0f %s - Rp %s\n",
// 				nama, kategori, stok, satuan, formatRupiah(harga))
// 		}
// 	}
// 	rows.Close()

// 	// 3. Batches
// 	var batchCount int
// 	db.QueryRow("SELECT COUNT(*) FROM batch").Scan(&batchCount)
// 	fmt.Printf("\n✅ BATCH: %d batches\n", batchCount)

// 	// Batch status breakdown
// 	var freshCount, hampirCount, expiredCount int
// 	db.QueryRow("SELECT COUNT(*) FROM batch WHERE status = 'fresh'").Scan(&freshCount)
// 	db.QueryRow("SELECT COUNT(*) FROM batch WHERE status = 'hampir_expired'").Scan(&hampirCount)
// 	db.QueryRow("SELECT COUNT(*) FROM batch WHERE status = 'expired'").Scan(&expiredCount)
// 	fmt.Printf("   - Fresh: %d\n", freshCount)
// 	fmt.Printf("   - Hampir Expired: %d\n", hampirCount)
// 	fmt.Printf("   - Expired: %d\n", expiredCount)

// 	// Sample batches
// 	fmt.Println("\n   Sample batches (showing 5):")
// 	rows, _ = db.Query(`
// 		SELECT b.id, p.nama, b.qty_tersisa, b.tanggal_kadaluarsa, b.status
// 		FROM batch b
// 		JOIN produk p ON b.produk_id = p.id
// 		ORDER BY b.tanggal_kadaluarsa
// 		LIMIT 5
// 	`)
// 	for rows.Next() {
// 		var id, nama, tanggal, status string
// 		var qty float64
// 		rows.Scan(&id, &nama, &qty, &tanggal, &status)
// 		fmt.Printf("   - %s: %.1fkg - Exp: %s [%s]\n", nama, qty, tanggal, status)
// 	}
// 	rows.Close()

// 	// 4. Transactions
// 	var transactionCount int
// 	var totalRevenue int
// 	db.QueryRow("SELECT COUNT(*) FROM transaksi").Scan(&transactionCount)
// 	db.QueryRow("SELECT COALESCE(SUM(total), 0) FROM transaksi").Scan(&totalRevenue)
// 	fmt.Printf("\n✅ TRANSAKSI: %d transactions\n", transactionCount)
// 	fmt.Printf("   Total Revenue: Rp %s\n", formatRupiah(totalRevenue))

// 	if transactionCount > 0 {
// 		fmt.Println("\n   Recent transactions:")
// 		rows, _ = db.Query(`
// 			SELECT nomor_transaksi, tanggal, total, kasir
// 			FROM transaksi
// 			ORDER BY tanggal DESC
// 			LIMIT 5
// 		`)
// 		for rows.Next() {
// 			var nomor, tanggal, kasir string
// 			var total int
// 			rows.Scan(&nomor, &tanggal, &total, &kasir)
// 			fmt.Printf("   - %s | %s | Rp %s | by %s\n",
// 				nomor, tanggal[:10], formatRupiah(total), kasir)
// 		}
// 		rows.Close()
// 	}

// 	// 5. Transaction Items
// 	var itemCount int
// 	db.QueryRow("SELECT COUNT(*) FROM transaksi_item").Scan(&itemCount)
// 	fmt.Printf("\n✅ TRANSAKSI ITEMS: %d items sold\n", itemCount)

// 	// 6. Payments
// 	var paymentCount int
// 	db.QueryRow("SELECT COUNT(*) FROM pembayaran").Scan(&paymentCount)
// 	fmt.Printf("\n✅ PEMBAYARAN: %d payment records\n", paymentCount)

// 	// Payment method breakdown
// 	rows, _ = db.Query(`
// 		SELECT metode, COUNT(*), COALESCE(SUM(jumlah), 0)
// 		FROM pembayaran
// 		GROUP BY metode
// 	`)
// 	fmt.Println("   Payment methods:")
// 	for rows.Next() {
// 		var metode string
// 		var count, total int
// 		rows.Scan(&metode, &count, &total)
// 		fmt.Printf("   - %s: %d payments (Rp %s)\n", metode, count, formatRupiah(total))
// 	}
// 	rows.Close()

// 	fmt.Println("\n═══════════════════════════════════════════════════════")
// 	fmt.Println("✅ Verification complete!")
// 	fmt.Println("\n💡 Tip: You can now run the application and see the data!")
// }

// func formatRupiah(amount int) string {
// 	if amount < 1000 {
// 		return fmt.Sprintf("%d", amount)
// 	}
// 	if amount < 1000000 {
// 		return fmt.Sprintf("%d.%03d", amount/1000, amount%1000)
// 	}
// 	return fmt.Sprintf("%d.%03d.%03d", amount/1000000, (amount/1000)%1000, amount%1000)
// }

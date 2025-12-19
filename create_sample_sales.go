package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"time"

// 	_ "github.com/mattn/go-sqlite3"
// )

// func main() {
// 	// Open database
// 	db, err := sql.Open("sqlite3", "./ritel.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// Get all products
// 	rows, err := db.Query("SELECT id, sku, nama, kategori, harga_jual FROM produk LIMIT 15")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	type Product struct {
// 		ID         int
// 		SKU        string
// 		Nama       string
// 		Kategori   string
// 		HargaJual  int
// 	}

// 	var products []Product
// 	for rows.Next() {
// 		var p Product
// 		err := rows.Scan(&p.ID, &p.SKU, &p.Nama, &p.Kategori, &p.HargaJual)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		products = append(products, p)
// 	}
// 	rows.Close()

// 	if len(products) == 0 {
// 		log.Fatal("No products found in database. Please add products first.")
// 	}

// 	fmt.Printf("Found %d products\n", len(products))

// 	// Create transactions for the current month (November 2025)
// 	startDate := time.Date(2025, 11, 1, 9, 0, 0, 0, time.Local)
// 	endDate := time.Date(2025, 11, 28, 20, 0, 0, 0, time.Local)

// 	// Product weights for random selection (some products more popular than others)
// 	weights := []int{30, 25, 20, 18, 15, 12, 10, 8, 7, 6, 5, 4, 3, 2, 1}

// 	rand.Seed(time.Now().UnixNano())

// 	transactionCount := 0

// 	// Create 100 transactions spread across the month
// 	for i := 0; i < 100; i++ {
// 		// Random date/time between start and end
// 		randomSeconds := rand.Int63n(int64(endDate.Sub(startDate).Seconds()))
// 		transaksiDate := startDate.Add(time.Duration(randomSeconds) * time.Second)

// 		// Random number of items per transaction (1-5)
// 		numItems := rand.Intn(5) + 1

// 		// Create transaction
// 		nomorTransaksi := fmt.Sprintf("TRX-%s-%04d", transaksiDate.Format("20060102"), i+1)

// 		// Insert transaction header
// 		subtotal := 0
// 		var selectedProducts []struct {
// 			Product Product
// 			Qty     int
// 		}

// 		// Select random products with weights
// 		for j := 0; j < numItems; j++ {
// 			// Weighted random selection
// 			totalWeight := 0
// 			for k := 0; k < len(products) && k < len(weights); k++ {
// 				totalWeight += weights[k]
// 			}

// 			randomWeight := rand.Intn(totalWeight)
// 			selectedIdx := 0
// 			currentWeight := 0

// 			for k := 0; k < len(products) && k < len(weights); k++ {
// 				currentWeight += weights[k]
// 				if randomWeight < currentWeight {
// 					selectedIdx = k
// 					break
// 				}
// 			}

// 			qty := rand.Intn(3) + 1 // 1-3 items
// 			selectedProducts = append(selectedProducts, struct {
// 				Product Product
// 				Qty     int
// 			}{products[selectedIdx], qty})

// 			subtotal += products[selectedIdx].HargaJual * qty
// 		}

// 		diskon := 0
// 		if rand.Float32() < 0.3 { // 30% chance of discount
// 			diskon = subtotal * rand.Intn(10) / 100 // 0-10% discount
// 		}

// 		total := subtotal - diskon

// 		result, err := db.Exec(`
// 			INSERT INTO transaksi (
// 				nomor_transaksi, tanggal, pelanggan_id, pelanggan_nama,
// 				subtotal, diskon, diskon_promo, diskon_pelanggan, total,
// 				total_bayar, kembalian, status, kasir, staff_id, staff_nama
// 			) VALUES (?, ?, 0, 'Umum', ?, ?, 0, 0, ?, ?, 0, 'selesai', 'Admin', 1, 'Admin')
// 		`, nomorTransaksi, transaksiDate, subtotal, diskon, total, total)

// 		if err != nil {
// 			log.Printf("Error creating transaction: %v", err)
// 			continue
// 		}

// 		transaksiID, _ := result.LastInsertId()

// 		// Insert transaction items
// 		for _, item := range selectedProducts {
// 			itemSubtotal := item.Product.HargaJual * item.Qty
// 			_, err = db.Exec(`
// 				INSERT INTO transaksi_item (
// 					transaksi_id, produk_id, produk_sku, produk_nama,
// 					produk_kategori, harga_satuan, jumlah, subtotal
// 				) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
// 			`, transaksiID, item.Product.ID, item.Product.SKU, item.Product.Nama,
// 				item.Product.Kategori, item.Product.HargaJual, item.Qty, itemSubtotal)

// 			if err != nil {
// 				log.Printf("Error creating transaction item: %v", err)
// 			}
// 		}

// 		// Insert payment
// 		paymentMethods := []string{"tunai", "qris", "debit", "kredit"}
// 		paymentMethod := paymentMethods[rand.Intn(len(paymentMethods))]

// 		_, err = db.Exec(`
// 			INSERT INTO pembayaran (transaksi_id, metode, jumlah)
// 			VALUES (?, ?, ?)
// 		`, transaksiID, paymentMethod, total)

// 		if err != nil {
// 			log.Printf("Error creating payment: %v", err)
// 		}

// 		transactionCount++
// 		if (i+1) % 20 == 0 {
// 			fmt.Printf("Created %d transactions...\n", i+1)
// 		}
// 	}

// 	fmt.Printf("\n✅ Successfully created %d transactions!\n", transactionCount)
// 	fmt.Println("\nTop selling products preview:")

// 	// Query top products
// 	topRows, err := db.Query(`
// 		SELECT
// 			ti.produk_nama,
// 			ti.produk_kategori,
// 			SUM(ti.jumlah) as total_terjual,
// 			SUM(ti.subtotal) as total_omset
// 		FROM transaksi_item ti
// 		INNER JOIN transaksi t ON t.id = ti.transaksi_id
// 		WHERE DATE(t.tanggal) BETWEEN '2025-11-01' AND '2025-11-28'
// 		GROUP BY ti.produk_id
// 		ORDER BY total_terjual DESC
// 		LIMIT 10
// 	`)

// 	if err != nil {
// 		log.Printf("Error querying top products: %v", err)
// 		return
// 	}
// 	defer topRows.Close()

// 	fmt.Println("\nRank | Produk | Kategori | Qty Terjual | Omset")
// 	fmt.Println("-----|--------|----------|-------------|-------")
// 	rank := 1
// 	for topRows.Next() {
// 		var nama, kategori string
// 		var qty, omset int
// 		topRows.Scan(&nama, &kategori, &qty, &omset)
// 		fmt.Printf("%2d   | %s | %s | %d | Rp %d\n", rank, nama, kategori, qty, omset)
// 		rank++
// 	}
// }

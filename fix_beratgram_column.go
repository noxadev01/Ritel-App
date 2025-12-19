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
// 	log.Printf("Database path: %s\n", dbPath)

// 	// Connect to database
// 	db, err := sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}
// 	defer db.Close()

// 	// Test connection
// 	if err := db.Ping(); err != nil {
// 		log.Fatal("Failed to ping database:", err)
// 	}

// 	log.Println("✅ Connected to database successfully")

// 	// Check if beratgram column exists
// 	var hasBeratgram bool
// 	rows, err := db.Query("PRAGMA table_info(transaksi_item)")
// 	if err != nil {
// 		log.Fatal("Failed to get table info:", err)
// 	}

// 	for rows.Next() {
// 		var cid int
// 		var name, colType string
// 		var notNull, pk int
// 		var dfltValue sql.NullString

// 		rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk)
// 		if name == "beratgram" {
// 			hasBeratgram = true
// 			log.Printf("✅ Column 'beratgram' exists (type: %s)", colType)
// 		}
// 	}
// 	rows.Close()

// 	if !hasBeratgram {
// 		log.Println("⚠️  Column 'beratgram' not found, adding it...")

// 		// Add beratgram column
// 		_, err = db.Exec("ALTER TABLE transaksi_item ADD COLUMN beratgram REAL DEFAULT 0")
// 		if err != nil {
// 			log.Printf("Error adding column: %v", err)
// 			log.Println("This might be okay if the column already exists with a different name")
// 		} else {
// 			log.Println("✅ Column 'beratgram' added successfully")
// 		}
// 	}

// 	// Check if old berat_gram column exists and migrate data if needed
// 	var hasBerat_gram bool
// 	rows, err = db.Query("PRAGMA table_info(transaksi_item)")
// 	if err != nil {
// 		log.Fatal("Failed to get table info:", err)
// 	}

// 	for rows.Next() {
// 		var cid int
// 		var name, colType string
// 		var notNull, pk int
// 		var dfltValue sql.NullString

// 		rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk)
// 		if name == "berat_gram" {
// 			hasBerat_gram = true
// 			log.Printf("⚠️  Old column 'berat_gram' found (type: %s)", colType)
// 		}
// 	}
// 	rows.Close()

// 	// If both columns exist, migrate data from old to new
// 	if hasBeratgram && hasBerat_gram {
// 		log.Println("📊 Migrating data from 'berat_gram' to 'beratgram'...")

// 		result, err := db.Exec("UPDATE transaksi_item SET beratgram = berat_gram WHERE berat_gram IS NOT NULL AND berat_gram > 0")
// 		if err != nil {
// 			log.Printf("⚠️  Error migrating data: %v", err)
// 		} else {
// 			rowsAffected, _ := result.RowsAffected()
// 			log.Printf("✅ Migrated %d rows", rowsAffected)
// 		}

// 		// Note: We don't drop the old column automatically for safety
// 		log.Println("ℹ️  Old column 'berat_gram' is kept for safety. You can manually drop it later if needed.")
// 	}

// 	// Verify final schema
// 	log.Println("\n📋 Final transaksi_item schema:")
// 	rows, err = db.Query("PRAGMA table_info(transaksi_item)")
// 	if err != nil {
// 		log.Fatal("Failed to get table info:", err)
// 	}

// 	fmt.Println("   Columns:")
// 	for rows.Next() {
// 		var cid int
// 		var name, colType string
// 		var notNull, pk int
// 		var dfltValue sql.NullString

// 		rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk)
// 		dflt := "NULL"
// 		if dfltValue.Valid {
// 			dflt = dfltValue.String
// 		}
// 		fmt.Printf("   - %s (%s, default: %s)\n", name, colType, dflt)
// 	}
// 	rows.Close()

// 	log.Println("\n✅ Database schema fix completed!")
// }

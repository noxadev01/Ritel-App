// run_migration.go
// Quick script to run soft delete migration
package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open database
	db, err := sql.Open("sqlite3", "./ritel.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Check current tables
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Fatalf("Failed to query tables: %v", err)
	}
	defer rows.Close()

	log.Println("Existing tables:")
	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		log.Printf("  - %s", tableName)
	}

	// Run migration
	log.Println("\n=== Running Soft Delete Migration ===")

	migrations := []struct {
		name  string
		query string
	}{
		{
			name:  "Add deleted_at to pelanggan",
			query: "ALTER TABLE pelanggan ADD COLUMN deleted_at TIMESTAMP NULL",
		},
		{
			name:  "Add deleted_at to produk",
			query: "ALTER TABLE produk ADD COLUMN deleted_at TIMESTAMP NULL",
		},
		{
			name:  "Add deleted_at to promo",
			query: "ALTER TABLE promo ADD COLUMN deleted_at TIMESTAMP NULL",
		},
		{
			name:  "Add deleted_at to kategori",
			query: "ALTER TABLE kategori ADD COLUMN deleted_at TIMESTAMP NULL",
		},
		{
			name:  "Create index on pelanggan.deleted_at",
			query: "CREATE INDEX IF NOT EXISTS idx_pelanggan_deleted_at ON pelanggan(deleted_at)",
		},
		{
			name:  "Create index on produk.deleted_at",
			query: "CREATE INDEX IF NOT EXISTS idx_produk_deleted_at ON produk(deleted_at)",
		},
		{
			name:  "Create index on promo.deleted_at",
			query: "CREATE INDEX IF NOT EXISTS idx_promo_deleted_at ON promo(deleted_at)",
		},
		{
			name:  "Create index on kategori.deleted_at",
			query: "CREATE INDEX IF NOT EXISTS idx_kategori_deleted_at ON kategori(deleted_at)",
		},
	}

	successCount := 0
	for _, mig := range migrations {
		log.Printf("Running: %s", mig.name)
		_, err := db.Exec(mig.query)
		if err != nil {
			// Ignore "duplicate column" error (column already exists)
			if contains(err.Error(), "duplicate column") {
				log.Printf("  ⚠️  Column already exists, skipping")
				continue
			}
			log.Printf("  ❌ Error: %v", err)
		} else {
			log.Printf("  ✅ Success")
			successCount++
		}
	}

	log.Printf("\n=== Migration Complete ===")
	log.Printf("Successfully applied %d migrations", successCount)

	// Verify columns exist
	log.Println("\n=== Verifying pelanggan table structure ===")
	rows2, err := db.Query("PRAGMA table_info(pelanggan)")
	if err != nil {
		log.Printf("Failed to get table info: %v", err)
	} else {
		defer rows2.Close()
		log.Println("Columns in pelanggan:")
		for rows2.Next() {
			var cid int
			var name, colType string
			var notNull, pk int
			var dfltValue interface{}
			rows2.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk)
			log.Printf("  - %s (%s)", name, colType)
		}
	}

	log.Println("\n✅ Migration script completed successfully!")
	log.Println("You can now restart your application.")
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

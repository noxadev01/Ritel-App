package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"
// 	_ "github.com/lib/pq"
// 	_ "github.com/mattn/go-sqlite3"
// )

// func main() {
// 	// Load environment variables
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("No .env file found, using default settings")
// 	}

// 	dbMode := os.Getenv("DB_MODE")
// 	if dbMode == "" {
// 		dbMode = "sqlite"
// 	}

// 	fmt.Printf("Running alignment migration for DB_MODE: %s\n", dbMode)

// 	var db *sql.DB
// 	var err error

// 	if dbMode == "postgresql" || dbMode == "dual" {
// 		// PostgreSQL connection
// 		dbHost := os.Getenv("DB_HOST")
// 		dbPort := os.Getenv("DB_PORT")
// 		dbUser := os.Getenv("DB_USER")
// 		dbPassword := os.Getenv("DB_PASSWORD")
// 		dbName := os.Getenv("DB_NAME")

// 		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 			dbHost, dbPort, dbUser, dbPassword, dbName)

// 		db, err = sql.Open("postgres", connStr)
// 		if err != nil {
// 			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
// 		}
// 		defer db.Close()

// 		fmt.Println("Connected to PostgreSQL")
// 		err = migratePostgreSQL(db)
// 	} else {
// 		// SQLite connection
// 		homeDir, _ := os.UserHomeDir()
// 		dbPath := homeDir + "/ritel-app/ritel.db"

// 		db, err = sql.Open("sqlite3", dbPath)
// 		if err != nil {
// 			log.Fatalf("Failed to connect to SQLite: %v", err)
// 		}
// 		defer db.Close()

// 		fmt.Printf("Connected to SQLite at: %s\n", dbPath)
// 		err = migrateSQLite(db)
// 	}

// 	if err != nil {
// 		log.Fatalf("Migration failed: %v", err)
// 	}

// 	fmt.Println("✅ Migration completed successfully!")
// 	fmt.Println("\nNew columns added to print_settings table:")
// 	fmt.Println("  - header_alignment (TEXT, default: 'center')")
// 	fmt.Println("  - title_alignment (TEXT, default: 'center')")
// 	fmt.Println("  - footer_alignment (TEXT, default: 'center')")
// 	fmt.Println("\nSekarang Anda bisa mengatur alignment untuk Header, Title, dan Footer!")
// }

// func migrateSQLite(db *sql.DB) error {
// 	// Check if columns already exist
// 	rows, err := db.Query("PRAGMA table_info(print_settings)")
// 	if err != nil {
// 		return fmt.Errorf("failed to get table info: %v", err)
// 	}
// 	defer rows.Close()

// 	existingColumns := make(map[string]bool)
// 	for rows.Next() {
// 		var cid int
// 		var name, typ string
// 		var notnull, pk int
// 		var dfltValue sql.NullString

// 		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
// 			return fmt.Errorf("failed to scan column info: %v", err)
// 		}
// 		existingColumns[name] = true
// 	}

// 	// Add header_alignment if not exists
// 	if !existingColumns["header_alignment"] {
// 		fmt.Println("Adding column: header_alignment")
// 		_, err = db.Exec("ALTER TABLE print_settings ADD COLUMN header_alignment TEXT DEFAULT 'center'")
// 		if err != nil {
// 			return fmt.Errorf("failed to add header_alignment: %v", err)
// 		}
// 	} else {
// 		fmt.Println("Column header_alignment already exists, skipping")
// 	}

// 	// Add title_alignment if not exists
// 	if !existingColumns["title_alignment"] {
// 		fmt.Println("Adding column: title_alignment")
// 		_, err = db.Exec("ALTER TABLE print_settings ADD COLUMN title_alignment TEXT DEFAULT 'center'")
// 		if err != nil {
// 			return fmt.Errorf("failed to add title_alignment: %v", err)
// 		}
// 	} else {
// 		fmt.Println("Column title_alignment already exists, skipping")
// 	}

// 	// Add footer_alignment if not exists
// 	if !existingColumns["footer_alignment"] {
// 		fmt.Println("Adding column: footer_alignment")
// 		_, err = db.Exec("ALTER TABLE print_settings ADD COLUMN footer_alignment TEXT DEFAULT 'center'")
// 		if err != nil {
// 			return fmt.Errorf("failed to add footer_alignment: %v", err)
// 		}
// 	} else {
// 		fmt.Println("Column footer_alignment already exists, skipping")
// 	}

// 	return nil
// }

// func migratePostgreSQL(db *sql.DB) error {
// 	// PostgreSQL supports IF NOT EXISTS in ALTER TABLE
// 	migrations := []string{
// 		"ALTER TABLE print_settings ADD COLUMN IF NOT EXISTS header_alignment VARCHAR(10) DEFAULT 'center'",
// 		"ALTER TABLE print_settings ADD COLUMN IF NOT EXISTS title_alignment VARCHAR(10) DEFAULT 'center'",
// 		"ALTER TABLE print_settings ADD COLUMN IF NOT EXISTS footer_alignment VARCHAR(10) DEFAULT 'center'",
// 	}

// 	for _, migration := range migrations {
// 		fmt.Println("Running migration:", migration)
// 		_, err := db.Exec(migration)
// 		if err != nil {
// 			return fmt.Errorf("failed to run migration: %v", err)
// 		}
// 	}

// 	return nil
// }

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

// 	fmt.Printf("Running migration for DB_MODE: %s\n", dbMode)

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
// 	fmt.Println("  - paper_width (INTEGER, default: 48)")
// 	fmt.Println("  - dash_line_char (TEXT, default: '-')")
// 	fmt.Println("  - double_line_char (TEXT, default: '=')")
// 	fmt.Println("\nYou can now use the new features in Format Struk settings!")
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

// 	// Add paper_width if not exists
// 	if !existingColumns["paper_width"] {
// 		fmt.Println("Adding column: paper_width")
// 		_, err = db.Exec("ALTER TABLE print_settings ADD COLUMN paper_width INTEGER DEFAULT 48")
// 		if err != nil {
// 			return fmt.Errorf("failed to add paper_width: %v", err)
// 		}
// 	} else {
// 		fmt.Println("Column paper_width already exists, skipping")
// 	}

// 	// Add dash_line_char if not exists
// 	if !existingColumns["dash_line_char"] {
// 		fmt.Println("Adding column: dash_line_char")
// 		_, err = db.Exec("ALTER TABLE print_settings ADD COLUMN dash_line_char TEXT DEFAULT '-'")
// 		if err != nil {
// 			return fmt.Errorf("failed to add dash_line_char: %v", err)
// 		}
// 	} else {
// 		fmt.Println("Column dash_line_char already exists, skipping")
// 	}

// 	// Add double_line_char if not exists
// 	if !existingColumns["double_line_char"] {
// 		fmt.Println("Adding column: double_line_char")
// 		_, err = db.Exec("ALTER TABLE print_settings ADD COLUMN double_line_char TEXT DEFAULT '='")
// 		if err != nil {
// 			return fmt.Errorf("failed to add double_line_char: %v", err)
// 		}
// 	} else {
// 		fmt.Println("Column double_line_char already exists, skipping")
// 	}

// 	return nil
// }

// func migratePostgreSQL(db *sql.DB) error {
// 	// PostgreSQL supports IF NOT EXISTS in ALTER TABLE
// 	migrations := []string{
// 		"ALTER TABLE print_settings ADD COLUMN IF NOT EXISTS paper_width INTEGER DEFAULT 48",
// 		"ALTER TABLE print_settings ADD COLUMN IF NOT EXISTS dash_line_char VARCHAR(10) DEFAULT '-'",
// 		"ALTER TABLE print_settings ADD COLUMN IF NOT EXISTS double_line_char VARCHAR(10) DEFAULT '='",
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

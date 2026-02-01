package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"gorm.io/gorm"
)

// SchemaMigration tracks which migrations have been executed (like TypeORM)
type SchemaMigration struct {
	ID        uint      `gorm:"primaryKey"`
	Version   int       `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"size:255;not null"`
	AppliedAt time.Time `gorm:"not null;default:now()"`
}

func (SchemaMigration) TableName() string {
	return "schema_migrations"
}

// RunMigrations executes pending migrations (TypeORM-style)
// 1. Creates schema_migrations table if not exists
// 2. Reads all .up.sql files from migrations folder
// 3. Checks which ones have already run
// 4. Executes pending migrations one by one in order
func RunMigrations(db *gorm.DB, migrationsDir string) error {
	log.Println("🔄 Running migrations (TypeORM-style)...")

	// Create schema_migrations table to track executed migrations
	if err := db.AutoMigrate(&SchemaMigration{}); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Get list of already executed migrations
	var executedMigrations []SchemaMigration
	if err := db.Order("version ASC").Find(&executedMigrations).Error; err != nil {
		return fmt.Errorf("failed to fetch executed migrations: %w", err)
	}

	executedVersions := make(map[int]bool)
	for _, m := range executedMigrations {
		executedVersions[m.Version] = true
	}

	// Get all .up.sql files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("📁 No migrations directory found, skipping...")
			return nil
		}
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort migration files
	var migrationFiles []string
	for _, f := range files {
		// Support Atlas generated .sql files, ignore atlas.sum
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") && f.Name() != "atlas.sum" {
			migrationFiles = append(migrationFiles, f.Name())
		}
	}
	sort.Strings(migrationFiles)

	if len(migrationFiles) == 0 {
		log.Println("📁 No migration files found")
		return nil
	}

	// Execute pending migrations
	pendingCount := 0
	for _, fileName := range migrationFiles {
		// Parse version from filename (e.g., "20231010123000_initial_schema.sql" or "000001_initial.up.sql")
		var version int
		// Try parsing integer from the start of the string
		_, err := fmt.Sscanf(fileName, "%d", &version)
		if err != nil {
			log.Printf("⚠️  Skipping file %s: cannot parse version", fileName)
			continue
		}

		// Skip if already executed
		if executedVersions[version] {
			log.Printf("⏭️  Skipping migration %d (already applied)", version)
			continue
		}

		// Read SQL file
		filePath := filepath.Join(migrationsDir, fileName)
		sqlContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		// Execute migration in a transaction
		log.Printf("▶️  Running migration: %s", fileName)
		err = db.Transaction(func(tx *gorm.DB) error {
			// Execute the SQL
			if err := tx.ExecfileName // Use full filename{
				return fmt.Errorf("migration failed: %w", err)
			}

			// Record the migration
			migrationName := strings.TrimSuffix(fileName, ".up.sql")
			migration := SchemaMigration{
				Version:   version,
				Name:      migrationName,
				AppliedAt: time.Now(),
			}
			if err := tx.Create(&migration).Error; err != nil {
				return fmt.Errorf("failed to record migration: %w", err)
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("migration %s failed: %w", fileName, err)
		}

		log.Printf("✅ Migration %s applied successfully", fileName)
		pendingCount++
	}

	if pendingCount == 0 {
		log.Println("✅ Database is up to date, no pending migrations")
	} else {
		log.Printf("✅ Applied %d migration(s) successfully", pendingCount)
	}

	return nil
}

// GetMigrationStatus returns the current migration status
func GetMigrationStatus(db *gorm.DB) ([]SchemaMigration, error) {
	var migrations []SchemaMigration
	err := db.Order("version ASC").Find(&migrations).Error
	return migrations, err
}

// AutoMigrate runs GORM auto-migration for all models at runtime
// This will create tables, add missing columns, and create indexes
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running auto-migration...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.ProductImages{},
		&models.Cart{},
		&models.CartItems{},
		&models.Order{},
		&models.OrderItem{},
	)

	if err != nil {
		log.Printf("Auto-migration failed: %v", err)
		return err
	}

	log.Println("Auto-migration completed successfully")
	return nil
}

// GenerateMigrationFiles creates SQL migration files based on models (TypeORM-style)
// Usage: Call this function to generate .up.sql and .down.sql files
func GenerateMigrationFiles(name string) error {
	migrationsDir := "migrations"

	if name == "" {
		name = "initial_schema"
	}

	// Create migrations directory
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Get next version number
	version := getNextVersion(migrationsDir)
	timestamp := time.Now().Format(time.RFC3339)

	// Create file names
	upFile := filepath.Join(migrationsDir, fmt.Sprintf("%06d_%s.up.sql", version, name))
	downFile := filepath.Join(migrationsDir, fmt.Sprintf("%06d_%s.down.sql", version, name))

	// Generate UP migration SQL
	upSQL := generateUpSQL(timestamp)
	if err := os.WriteFile(upFile, []byte(upSQL), 0644); err != nil {
		return fmt.Errorf("failed to write up migration: %w", err)
	}
	log.Printf("✓ Created: %s\n", upFile)

	// Generate DOWN migration SQL
	downSQL := generateDownSQL(timestamp)
	if err := os.WriteFile(downFile, []byte(downSQL), 0644); err != nil {
		return fmt.Errorf("failed to write down migration: %w", err)
	}
	log.Printf("✓ Created: %s\n", downFile)

	log.Printf("Migration '%s' (version %d) generated successfully!\n", name, version)
	return nil
}

func getNextVersion(dir string) int {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 1
	}

	maxVersion := 0
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		var v int
		fmt.Sscanf(f.Name(), "%06d", &v)
		if v > maxVersion {
			maxVersion = v
		}
	}
	return maxVersion + 1
}

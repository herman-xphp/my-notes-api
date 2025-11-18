package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type Migration struct {
	Version string
	Name    string
	UpSQL   string
	DownSQL string
}

// RunMigrations executes all pending migrations
func RunMigrations(db *gorm.DB, migrationsPath string) error {
	log.Println("🔄 Running database migrations...")

	// Get all migration files
	migrations, err := loadMigrations(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	if len(migrations) == 0 {
		log.Println("ℹ️  No migrations found")
		return nil
	}

	// Execute migrations in order
	for _, migration := range migrations {
		log.Printf("⬆️  Applying migration: %s - %s", migration.Version, migration.Name)

		if err := db.Exec(migration.UpSQL).Error; err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		log.Printf("✅ Migration applied: %s", migration.Version)
	}

	log.Println("✅ All migrations completed successfully")
	return nil
}

// RollbackMigrations rolls back all migrations
func RollbackMigrations(db *gorm.DB, migrationsPath string) error {
	log.Println("🔄 Rolling back database migrations...")

	// Get all migration files
	migrations, err := loadMigrations(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	if len(migrations) == 0 {
		log.Println("ℹ️  No migrations found")
		return nil
	}

	// Execute rollbacks in reverse order
	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]
		log.Printf("⬇️  Rolling back migration: %s - %s", migration.Version, migration.Name)

		if err := db.Exec(migration.DownSQL).Error; err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", migration.Version, err)
		}

		log.Printf("✅ Migration rolled back: %s", migration.Version)
	}

	log.Println("✅ All migrations rolled back successfully")
	return nil
}

// loadMigrations loads all migration files from the migrations directory
func loadMigrations(migrationsPath string) ([]Migration, error) {
	var migrations []Migration

	// Read migrations directory
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Group migrations by version
	migrationMap := make(map[string]*Migration)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()

		// Parse migration filename: 001_create_users_table.up.sql
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			continue
		}

		version := parts[0]

		// Get migration name (remove version and .up.sql or .down.sql)
		nameWithExt := strings.Join(parts[1:], "_")
		name := strings.TrimSuffix(strings.TrimSuffix(nameWithExt, ".up.sql"), ".down.sql")

		// Read file content
		content, err := os.ReadFile(filepath.Join(migrationsPath, filename))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Initialize migration if not exists
		if migrationMap[version] == nil {
			migrationMap[version] = &Migration{
				Version: version,
				Name:    name,
			}
		}

		// Set up or down SQL
		if strings.HasSuffix(filename, ".up.sql") {
			migrationMap[version].UpSQL = string(content)
		} else if strings.HasSuffix(filename, ".down.sql") {
			migrationMap[version].DownSQL = string(content)
		}
	}

	// Convert map to slice and sort by version
	for _, migration := range migrationMap {
		migrations = append(migrations, *migration)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

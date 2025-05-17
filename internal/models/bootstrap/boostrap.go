package bootstrap

import (
	"bufio"
	// "database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/linuxswords/hat/internal/models"
	"gorm.io/gorm"
)

func initBowclasses(db *gorm.DB) error {
	var count int64
	db.Model(&models.BowClass{}).Count(&count)
	if count == 0 {
		if err := execSQLFile(db, "internal/models/bootstrap/initial_bowclasses.sql"); err != nil {
			return fmt.Errorf("failed to load bow classes: %w", err)
		}
		fmt.Println("✅ Bow classes loaded")
	} else {
		fmt.Println("ℹ️ Bow classes already exist — skipping")
	}
	return nil
}

func initHandyCapSet(db *gorm.DB) error {
	var existingSet models.HandycapSet
	if err := db.Where("name = ?", "3D 2021-2025 3-Arrow Round").First(&existingSet).Error; err == nil {
		fmt.Println("ℹ️ Handicap set already exists — skipping")
		return nil
	}

	// Create HandycapSet & entries
	if err := execSQLFile(db, "internal/models/bootstrap/initial_handycapset.sql"); err != nil {
		return fmt.Errorf("failed to load inital handycap set: %w", err)
	}
	fmt.Println("✅ Handicap set and entries loaded")
	return nil
}

// BootstrapData loads initial SQL and CSV data into the database
func BootstrapData(db *gorm.DB) error {
	// Migrate the schema
	if err := db.AutoMigrate(&models.BowClass{}, &models.Archer{}, &models.HandycapSet{}, &models.HandycapEntry{}, &models.Tournament{}); err != nil {
		return fmt.Errorf("auto-migrate failed: %w", err)
	}

	// Check if bow classes exist
	if err := initBowclasses(db); err != nil {
		return fmt.Errorf("Error initializing bow classes: %w", err)
	}

	// Check if the handicap set already exists
	if err := initHandyCapSet(db); err != nil {
		return fmt.Errorf("Error initializing handycap classes: %w", err)
	}
	return nil
}

// execSQLFile executes raw SQL statements from a file
func execSQLFile(db *gorm.DB, filepath string) error {
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	sqlStmts := string(sqlBytes)
	return db.Exec(sqlStmts).Error
}

// loadCSVEntries parses the CSV and attaches entries to the handycap set
func loadCSVEntries(db *gorm.DB, set *models.HandycapSet, csvPath string) error {
	f, err := os.Open(filepath.Clean(csvPath))
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, rec := range records {
		if len(rec) != 2 {
			continue // malformed row
		}
		code := strings.TrimSpace(rec[0])
		valueStr := strings.TrimSpace(rec[1])

		var bc models.BowClass
		if err := db.Where("code = ?", code).First(&bc).Error; err != nil {
			fmt.Printf("⚠️ Skipping unknown BowClass code: %s\n", code)
			continue
		}

		var value float64
		fmt.Sscanf(valueStr, "%f", &value)

		set.HandycapEntries = append(set.HandycapEntries, models.HandycapEntry{
			BowClassID: bc.ID,
			Value:      value,
		})
	}

	return nil
}

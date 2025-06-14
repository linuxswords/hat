package bootstrap

import (
	"fmt"
	"os"

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

func initHandiCapSet(db *gorm.DB) error {
	var existingSet models.HandicapSet
	if err := db.Where("name = ?", "3D 2021-2025 3-Arrow Round (KI generated)").First(&existingSet).Error; err == nil {
		fmt.Println("ℹ️ Handicap set already exists — skipping")
		return nil
	}

	// Create HandicapSet & entries
	if err := execSQLFile(db, "internal/models/bootstrap/initial_handicapset.sql"); err != nil {
		return fmt.Errorf("failed to load inital handicap set: %w", err)
	}
	fmt.Println("✅ Handicap set and entries loaded")
	return nil
}

func initArchers(db *gorm.DB) error {
	var count int64
	db.Model(&models.Archer{}).Count(&count)
	if count == 0 {
		if err := execSQLFile(db, "internal/models/bootstrap/initial_archers.sql"); err != nil {
			return fmt.Errorf("failed to load archers: %w", err)
		}
		fmt.Println("✅ Archers loaded")
	} else {
		fmt.Println("ℹ️ Archers already exist — skipping")
	}
	return nil
}

// BootstrapData loads initial SQL and CSV data into the database
func BootstrapData(db *gorm.DB) error {
	// Migrate the schema
	if err := db.AutoMigrate(&models.BowClass{}, &models.Archer{}, &models.HandicapSet{}, &models.HandicapEntry{}, &models.Tournament{}, &models.Score{}, &models.TournamentArcher{}); err != nil {
		return fmt.Errorf("auto-migrate failed: %w", err)
	}

	// Check if bow classes exist and create if not
	if err := initBowclasses(db); err != nil {
		return fmt.Errorf("Error initializing bow classes: %w", err)
	}

	// Check if the handicap set already exists and create if not
	if err := initHandiCapSet(db); err != nil {
		return fmt.Errorf("Error initializing handicap classes: %w", err)
	}

	// Check if the archers already exists and create if not
	if err := initArchers(db); err != nil {
		return fmt.Errorf("Error initializing archers classes: %w", err)
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

package database

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"gorm.io/gorm"
	"github.com/linuxswords/hat/internal/models"
)

// TestDB holds a test database instance
type TestDB struct {
	DB     *gorm.DB
	DBName string
}

// SetupTestDatabase creates a fresh test database for each test
func SetupTestDatabase(t *testing.T) *TestDB {
	// Generate unique database name for this test
	dbName := fmt.Sprintf("hat_test_%d_%s", time.Now().Unix(), t.Name())
	
	// Connect to postgres to create the test database
	config := GetTestConfig()
	config.DBName = "postgres" // Connect to default database first
	
	db, err := Connect(config)
	if err != nil {
		t.Fatalf("Failed to connect to postgres: %v", err)
	}
	
	// Create test database
	createSQL := fmt.Sprintf("CREATE DATABASE %s", dbName)
	if err := db.Exec(createSQL).Error; err != nil {
		t.Fatalf("Failed to create test database %s: %v", dbName, err)
	}
	
	// Close connection to postgres
	sqlDB, _ := db.DB()
	sqlDB.Close()
	
	// Connect to the new test database
	config.DBName = dbName
	testDB, err := Connect(config)
	if err != nil {
		t.Fatalf("Failed to connect to test database %s: %v", dbName, err)
	}
	
	// Run migrations
	if err := Migrate(testDB); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	
	return &TestDB{
		DB:     testDB,
		DBName: dbName,
	}
}

// SetupTestDatabaseWithData creates a test database and seeds it with test data
func SetupTestDatabaseWithData(t *testing.T) *TestDB {
	testDB := SetupTestDatabase(t)
	
	// Seed with test data
	if err := SeedTestData(testDB.DB); err != nil {
		t.Fatalf("Failed to seed test database: %v", err)
	}
	
	return testDB
}

// Cleanup drops the test database
func (tdb *TestDB) Cleanup(t *testing.T) {
	// Close connection to test database
	sqlDB, _ := tdb.DB.DB()
	sqlDB.Close()
	
	// Connect to postgres to drop the test database
	config := GetTestConfig()
	config.DBName = "postgres"
	
	db, err := Connect(config)
	if err != nil {
		t.Logf("Failed to connect to postgres for cleanup: %v", err)
		return
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()
	
	// Terminate any active connections to the test database
	terminateSQL := fmt.Sprintf(`
		SELECT pg_terminate_backend(pid) 
		FROM pg_stat_activity 
		WHERE datname = '%s' AND pid <> pg_backend_pid()
	`, tdb.DBName)
	db.Exec(terminateSQL)
	
	// Drop test database
	dropSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s", tdb.DBName)
	if err := db.Exec(dropSQL).Error; err != nil {
		t.Logf("Failed to drop test database %s: %v", tdb.DBName, err)
	}
}

// SeedTestData seeds the database with minimal test data
func SeedTestData(db *gorm.DB) error {
	// Create handicap set
	handicapSet := models.HandicapSet{
		Name:        "Test Handicaps",
		Description: "Test handicap set for automated tests",
		IsActive:    true,
	}
	if err := db.Create(&handicapSet).Error; err != nil {
		return fmt.Errorf("failed to create test handicap set: %w", err)
	}
	
	// Create a few test handicaps
	handicaps := []models.Handicap{
		{BowClassID: "AMLB", Factor: 1.0, SetID: handicapSet.ID},
		{BowClassID: "AFLB", Factor: 0.95, SetID: handicapSet.ID},
		{BowClassID: "AMTR", Factor: 0.90, SetID: handicapSet.ID},
	}
	
	for _, handicap := range handicaps {
		if err := db.Create(&handicap).Error; err != nil {
			return fmt.Errorf("failed to create test handicap: %w", err)
		}
	}
	
	// Create test archers
	archers := []models.Archer{
		{Name: "Test Archer 1", Gender: "Male", BowClass: "AMLB", Email: "test1@example.com"},
		{Name: "Test Archer 2", Gender: "Female", BowClass: "AFLB", Email: "test2@example.com"},
		{Name: "Test Archer 3", Gender: "Male", BowClass: "AMTR", Email: "test3@example.com"},
	}
	
	for _, archer := range archers {
		if err := db.Create(&archer).Error; err != nil {
			return fmt.Errorf("failed to create test archer: %w", err)
		}
	}
	
	// Create test tournament
	tournament := models.Tournament{
		Name:          "Test Tournament",
		Location:      "Test Location",
		StartDate:     time.Now().AddDate(0, 0, 1), // Tomorrow
		EndDate:       time.Now().AddDate(0, 0, 2), // Day after tomorrow
		HandicapSetID: handicapSet.ID,
	}
	if err := db.Create(&tournament).Error; err != nil {
		return fmt.Errorf("failed to create test tournament: %w", err)
	}
	
	return nil
}

// InitializeTestEnvironment sets up the test environment
func InitializeTestEnvironment() {
	// Set environment variables for testing
	os.Setenv("ENV", "test")
	os.Setenv("TEST_DB_HOST", getEnv("TEST_DB_HOST", "localhost"))
	os.Setenv("TEST_DB_PORT", getEnv("TEST_DB_PORT", "5432"))
	os.Setenv("TEST_DB_USER", getEnv("TEST_DB_USER", "hat_user"))
	os.Setenv("TEST_DB_PASSWORD", getEnv("TEST_DB_PASSWORD", "hat_password"))
	
	// Disable GORM logging in tests
	log.SetOutput(os.Stderr)
}
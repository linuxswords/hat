package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/linuxswords/hat/internal/models"
)

// DB is the global database instance
var DB *gorm.DB

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetDefaultConfig returns default database configuration
func GetDefaultConfig() Config {
	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "hat_user"),
		Password: getEnv("DB_PASSWORD", "hat_password"),
		DBName:   getEnv("DB_NAME", "hat_development"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// GetTestConfig returns test database configuration
func GetTestConfig() Config {
	return Config{
		Host:     getEnv("TEST_DB_HOST", "localhost"),
		Port:     getEnv("TEST_DB_PORT", "5432"),
		User:     getEnv("TEST_DB_USER", "hat_user"),
		Password: getEnv("TEST_DB_PASSWORD", "hat_password"),
		DBName:   getEnv("TEST_DB_NAME", "hat_test"),
		SSLMode:  getEnv("TEST_DB_SSLMODE", "disable"),
	}
}

// Connect establishes a database connection
func Connect(config Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// Configure GORM logger based on environment
	var gormLogger logger.Interface
	if getEnv("ENV", "development") == "test" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// Initialize sets up the global database connection
func Initialize() error {
	config := GetDefaultConfig()
	db, err := Connect(config)
	if err != nil {
		return err
	}

	DB = db
	log.Printf("Connected to database: %s", config.DBName)
	return nil
}

// InitializeTest sets up a test database connection
func InitializeTest() (*gorm.DB, error) {
	config := GetTestConfig()
	return Connect(config)
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	
	err := db.AutoMigrate(
		&models.Archer{},
		&models.Tournament{},
		&models.TournamentParticipation{},
		&models.HandicapSet{},
		&models.Handicap{},
		&models.Score{},
	)
	
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	
	log.Println("Database migrations completed successfully")
	return nil
}

// Seed inserts initial test data
func Seed(db *gorm.DB) error {
	log.Println("Seeding database with initial data...")
	
	// Check if data already exists
	var count int64
	db.Model(&models.Archer{}).Count(&count)
	if count > 0 {
		log.Println("Database already contains data, skipping seed")
		return nil
	}
	
	// Create handicap set
	handicapSet := models.HandicapSet{
		Name:        "HAT Standard Handicaps 2024",
		Description: "Standard handicap factors for HAT tournaments",
		IsActive:    true,
	}
	if err := db.Create(&handicapSet).Error; err != nil {
		return fmt.Errorf("failed to create handicap set: %w", err)
	}
	
	// Create handicaps
	handicaps := []models.Handicap{
		{BowClassID: "AMLB", Factor: 1.0, SetID: handicapSet.ID},     // Adult Male Longbow (baseline)
		{BowClassID: "AFLB", Factor: 0.95, SetID: handicapSet.ID},    // Adult Female Longbow
		{BowClassID: "AMTR", Factor: 0.90, SetID: handicapSet.ID},    // Adult Male Traditional Recurve
		{BowClassID: "AFTR", Factor: 0.85, SetID: handicapSet.ID},    // Adult Female Traditional Recurve
		{BowClassID: "AMBHR", Factor: 0.88, SetID: handicapSet.ID},   // Adult Male Bowhunter Recurve
		{BowClassID: "AFBHR", Factor: 0.83, SetID: handicapSet.ID},   // Adult Female Bowhunter Recurve
		{BowClassID: "AMFU", Factor: 0.75, SetID: handicapSet.ID},    // Adult Male Freestyle Unlimited
		{BowClassID: "AFFU", Factor: 0.70, SetID: handicapSet.ID},    // Adult Female Freestyle Unlimited
		{BowClassID: "VMTR", Factor: 0.92, SetID: handicapSet.ID},    // Veteran Male Traditional Recurve
		{BowClassID: "VFTR", Factor: 0.87, SetID: handicapSet.ID},    // Veteran Female Traditional Recurve
	}
	
	for _, handicap := range handicaps {
		if err := db.Create(&handicap).Error; err != nil {
			return fmt.Errorf("failed to create handicap %s: %w", handicap.BowClassID, err)
		}
	}
	
	// Create sample archers
	archers := []models.Archer{
		{Name: "Alice Johnson", Gender: "Female", BowClass: "AFLB", Email: "alice.johnson@example.com"},
		{Name: "Bob Smith", Gender: "Male", BowClass: "AMLB", Email: "bob.smith@example.com"},
		{Name: "Carol Williams", Gender: "Female", BowClass: "AFTR", Email: "carol.williams@example.com"},
		{Name: "David Brown", Gender: "Male", BowClass: "AMTR", Email: "david.brown@example.com"},
		{Name: "Eva Davis", Gender: "Female", BowClass: "AFBHR", Email: "eva.davis@example.com"},
	}
	
	for _, archer := range archers {
		if err := db.Create(&archer).Error; err != nil {
			return fmt.Errorf("failed to create archer %s: %w", archer.Name, err)
		}
	}
	
	// Create sample tournament
	tournament := models.Tournament{
		Name:         "Spring Championship 2024",
		Location:     "Central Archery Range",
		StartDate:    time.Now().AddDate(0, 1, 0), // Next month
		EndDate:      time.Now().AddDate(0, 1, 2), // Next month + 2 days
		HandicapSetID: handicapSet.ID,
	}
	if err := db.Create(&tournament).Error; err != nil {
		return fmt.Errorf("failed to create tournament: %w", err)
	}
	
	log.Println("Database seeding completed successfully")
	return nil
}

// Bootstrap runs migrations and seeding
func Bootstrap(db *gorm.DB) error {
	if err := Migrate(db); err != nil {
		return err
	}
	
	if err := Seed(db); err != nil {
		return err
	}
	
	return nil
}

// CleanDatabase removes all data from the database (for testing)
func CleanDatabase(db *gorm.DB) error {
	// Delete in reverse order of dependencies
	tables := []interface{}{
		&models.Score{},
		&models.TournamentParticipation{},
		&models.Handicap{},
		&models.HandicapSet{},
		&models.Tournament{},
		&models.Archer{},
	}
	
	for _, table := range tables {
		if err := db.Unscoped().Delete(table, "1 = 1").Error; err != nil {
			return fmt.Errorf("failed to clean table %T: %w", table, err)
		}
	}
	
	return nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	
	return sqlDB.Close()
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
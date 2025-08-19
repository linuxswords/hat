# Database Setup for HAT Application

This document describes how to set up and work with the PostgreSQL database for the HAT (Handicap Archery Tournament) application.

## Quick Start

### 1. Start the Database

```bash
# Start PostgreSQL containers
make db-setup
```

This will:
- Start PostgreSQL containers for development and testing
- Create databases: `hat_development` and `hat_test`
- Set up user credentials

### 2. Run the Application

```bash
# Build and run the application
make build
make run
```

The application will automatically:
- Connect to the database
- Run migrations to create tables
- Seed initial data if the database is empty

## Database Configuration

### Development Database
- **Host**: localhost
- **Port**: 5432
- **Database**: hat_development
- **User**: hat_user
- **Password**: hat_password

### Test Database
- **Host**: localhost
- **Port**: 5433
- **Database**: hat_test
- **User**: hat_user
- **Password**: hat_password

## Environment Variables

You can customize database settings using environment variables:

```bash
# Development database
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=hat_user
export DB_PASSWORD=hat_password
export DB_NAME=hat_development

# Test database
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5433
export TEST_DB_USER=hat_user
export TEST_DB_PASSWORD=hat_password
export TEST_DB_NAME=hat_test
```

## Database Schema

The application uses GORM for ORM and the following main tables:

### Core Tables
- **archers** - Archer profiles and information
- **tournaments** - Tournament details and scheduling
- **handicap_sets** - Handicap configuration sets
- **handicaps** - Individual handicap factors by bow class

### Relationship Tables
- **tournament_participations** - Links archers to tournaments
- **scores** - Individual scores for tournament/archer combinations

### Key Features
- **Soft deletes** - Records are marked as deleted, not physically removed
- **Timestamps** - Automatic created_at, updated_at tracking
- **Foreign keys** - Proper relationships between entities
- **Indexes** - Optimized queries for common lookups

## Migrations

Migrations run automatically when the application starts. The system uses GORM's AutoMigrate feature which:

- Creates tables if they don't exist
- Adds new columns
- Adds missing indexes
- **Does NOT** delete or modify existing columns (safe for production)

## Seeding Data

Initial data is automatically seeded when the application starts if the database is empty:

- **Handicap Sets** - Default handicap configuration
- **Handicaps** - Standard factors for common bow classes
- **Sample Archers** - Test archer profiles
- **Sample Tournament** - Example tournament for testing

## Testing

### Test Database Isolation

Tests automatically create isolated databases:

```go
func TestSomething(t *testing.T) {
    // Each test gets a fresh database
    testDB := database.SetupTestDatabase(t)
    defer testDB.Cleanup(t)
    
    // Use testDB.DB for your test operations
}
```

### Running Tests

```bash
# Run backend tests (with database)
make test

# Run frontend tests (will use test database)
make test-e2e

# Run all tests
make test-all
```

## Development Workflow

### 1. Database Changes

When adding new models or fields:

1. Update the model structs in `internal/models/`
2. Add GORM tags for validation and relationships
3. Restart the application - migrations run automatically
4. Update repository interfaces and implementations if needed

### 2. Test Data Management

For consistent testing:

1. Use `database.SetupTestDatabaseWithData()` for tests that need data
2. Each test gets a completely isolated database
3. Test databases are automatically cleaned up
4. No data pollution between tests

### 3. Repository Pattern

The application uses repository interfaces for clean separation:

```go
// Interface (in handlers/handlers.go)
type ArcherRepository interface {
    GetAll() []models.Archer
    GetByID(id uint) (*models.Archer, error)
    // ... other methods
}

// Implementation (in repositories/db_archer_repository.go)
type DBArcherRepository struct {
    db *gorm.DB
}
```

## Production Considerations

### Security
- Use environment variables for database credentials
- Enable SSL/TLS for database connections
- Use connection pooling (already configured)
- Regular security updates

### Performance
- Connection pooling is pre-configured
- Indexes are automatically created for foreign keys
- Use `Preload()` for eager loading relationships
- Monitor slow queries

### Backup
- Set up regular PostgreSQL backups
- Test restore procedures
- Consider point-in-time recovery

## Troubleshooting

### Connection Issues

```bash
# Check if PostgreSQL is running
make db-setup

# Check container status
docker-compose ps

# View logs
docker-compose logs postgres
```

### Migration Issues

```bash
# Check application logs
make run

# Manually inspect database
docker exec -it hat_postgres psql -U hat_user -d hat_development
```

### Test Issues

```bash
# Ensure test database is available
docker-compose ps postgres_test

# Run tests with verbose output
make test-verbose
```

## Advanced Usage

### Custom Migrations

For complex schema changes, you can create custom migration functions:

```go
func CustomMigration(db *gorm.DB) error {
    // Custom migration logic
    return db.Exec("ALTER TABLE ...").Error
}
```

### Performance Monitoring

Enable GORM logging for development:

```go
db.Logger = logger.Default.LogMode(logger.Info)
```

### Multiple Environments

Use different environment variables for staging, production, etc.:

```bash
# Staging
export DB_NAME=hat_staging
export DB_HOST=staging-db.example.com

# Production  
export DB_NAME=hat_production
export DB_HOST=prod-db.example.com
```
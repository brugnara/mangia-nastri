package src

import (
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	config "mangia_nastri/config"
)

// `getSqliteDialector` returns a SQLite dialector for the GORM database connection.
// It ensures that the database file is created if it does not exist using the options
// provided via flags or their default values. (Ref: config/options.go)
//
// Returns
//   A SQLite dialector for the GORM database connection.

func getSqliteDialector() gorm.Dialector {
	if !strings.Contains(config.DatabaseName, ".") {
		config.DatabaseName += ".sqlite"
	}

	e := os.MkdirAll(config.DatabasePath, os.ModePerm)
	handleError("Failed to connect database", e)

	databasePath := config.DatabasePath + "/" + config.DatabaseName
	return sqlite.Open(databasePath)
}

// `CreateDatabase` initializes and returns a GORM database instance.
// It uses the SQLite dialector to connect to the database and sets the connection
// pool settings.
//
// Returns
//   A GORM database instance.

func CreateDatabase() (*gorm.DB, error) {
	database, e := gorm.Open(getSqliteDialector(), &gorm.Config{})
	handleError("Failed to connect database", e)

	db, e := database.DB()
	handleError("Failed to get database instance", e)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(0)

	return database, nil
}

package migrate

import (
	"log"

	utility "github.com/TripressoCTS/cts-server-utility"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migrationsV2 = []*gormigrate.Migration{}

var migrateOptionsV2 = &gormigrate.Options{
	TableName:                 "go_migrations",
	IDColumnName:              "id",
	IDColumnSize:              255,
	UseTransaction:            false,
	ValidateUnknownMigrations: false,
}

// MigrateDB performs database migrations for the dbmodels
func MigrateDB(db *gorm.DB) error {
	log.Println("Starting database migration...")

	if err := gormigrate.New(db, migrateOptionsV2, migrationsV2).Migrate(); err != nil {
		utility.SugarLogger.Errorln("Error during migration", "error", err)
		return err
	}

	utility.SugarLogger.Infoln("Database migration completed successfully")
	return nil
}

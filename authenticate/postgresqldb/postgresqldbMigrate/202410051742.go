package postgresqldbMigrate

import (
	"authenticate/model/dbmodel"

	utility "github.com/TripressoCTS/cts-server-utility"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func v202410051742Migration(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&dbmodel.User{}); err != nil {
			utility.SugarLogger.Error("AutoMigrate User failed", err)
			return err
		}
		if err := tx.AutoMigrate(&dbmodel.Session{}); err != nil {
			utility.SugarLogger.Error("AutoMigrate Session failed", err)
			return err
		}
		if err := tx.AutoMigrate(&dbmodel.Place{}); err != nil {
			utility.SugarLogger.Error("AutoMigrate Place failed", err)
			return err
		}
		if err := tx.AutoMigrate(&dbmodel.Favorite{}); err != nil {
			utility.SugarLogger.Error("AutoMigrate Favorite failed", err)
			return err
		}
		return nil
	})
}

var v202410051742 = &gormigrate.Migration{
	ID:      "v202410051742",
	Migrate: v202410051742Migration,
}

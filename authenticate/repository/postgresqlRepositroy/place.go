package postgresqlRepository

import (
	"authenticate/model/dbmodel"
	"authenticate/model/entitymodel"
	"authenticate/repository"
	"authenticate/utility"
	"context"

	"gorm.io/gorm"
)

func NewPlaceHandle(db *gorm.DB) repository.PlaceHandler {
	return &PlaceHandle{db: db}
}

type PlaceHandle struct {
	db *gorm.DB
}

func (handle *PlaceHandle) CreatePlace(ctx context.Context, place entitymodel.Place) error {
	dbModel := dbmodel.Place{
		ID:                       place.ID,
		GoogleID:                 place.GoogleID,
		TWDisplayName:            place.TWDisplayName,
		JPDisplayName:            place.JPDisplayName,
		TWFormattedAddress:       place.TWFormattedAddress,
		TWWeekdayDescriptions:    place.TWWeekdayDescriptions,
		AdministrativeAreaLevel1: place.AdministrativeAreaLevel1,
		Country:                  place.Country,
		GoogleMapURI:             place.GoogleMapURI,
		InternationalPhoneNumber: place.InternationalPhoneNumber,
		Lat:                      place.Lat,
		Lng:                      place.Lng,
		PrimaryType:              place.PrimaryType,
		Rating:                   place.Rating,
		Types:                    place.Types,
		UserRatingCount:          place.UserRatingCount,
		WebsiteURI:               place.WebsiteURI,
		PlaceVersion:             place.PlaceVersion,
		CreatedAt:                place.CreatedAt,
		UpdatedAt:                place.UpdatedAt,
	}

	if err := handle.db.WithContext(ctx).
		Create(&dbModel).Error; err != nil {
		utility.SugarLogger.Errorln("create place error", err)
		return err
	}

	return nil
}

func (handle *PlaceHandle) GetPlaceByGoogleID(ctx context.Context, googleID string) (entitymodel.Place, error) {
	var place dbmodel.Place
	if err := handle.db.WithContext(ctx).
		Where("google_id = ?", googleID).
		Find(&place).Error; err != nil {
		utility.SugarLogger.Errorln("get place by google id error", err)
		return entitymodel.Place{}, err
	}
	return place.ToEntityModel(), nil
}

func (handle *PlaceHandle) UpdatePlace(ctx context.Context, placeID string, updates map[string]interface{}) error {
	if err := handle.db.WithContext(ctx).
		Model(&dbmodel.Place{}).
		Where("id = ?", placeID).
		Updates(updates).Error; err != nil {
		utility.SugarLogger.Errorln("update place error", err)
		return err
	}
	return nil
}

func (handle *PlaceHandle) DeletePlace(ctx context.Context, placeID string) error {
	if err := handle.db.WithContext(ctx).
		Where("id = ?", placeID).
		Delete(&dbmodel.Place{}).Error; err != nil {
		utility.SugarLogger.Errorln("delete place error", err)
		return err
	}
	return nil
}

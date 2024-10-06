package postgresqlRepository

import (
	"authenticate/model/dbmodel"
	"authenticate/model/entitymodel"
	"authenticate/model/enum"
	"authenticate/repository"
	"authenticate/utility"
	"context"

	"gorm.io/gorm"
)

func NewFavoriteHandler(db *gorm.DB) repository.FavoriteHandler {
	return FavoriteHandle{db: db}
}

type FavoriteHandle struct {
	db *gorm.DB
}

func (f FavoriteHandle) CreateFavorite(ctx context.Context, favorite entitymodel.Favorite) error {
	dbModel := dbmodel.Favorite{
		ID:         favorite.ID,
		IsFavorite: favorite.IsFavorite,
		UserID:     favorite.UserID,
		PlaceID:    favorite.PlaceID,
		CreatedAt:  favorite.CreatedAt,
		UpdatedAt:  favorite.UpdatedAt,
	}
	if err := f.db.WithContext(ctx).
		Create(&dbModel).Error; err != nil {
		utility.SugarLogger.Errorln("error creating favorite:", err)
		return err
	}
	return nil
}

func (f FavoriteHandle) GetFavoriteByUserIDAndPlaceID(ctx context.Context, userID, placeID string) (entitymodel.Favorite, error) {
	var dbModel dbmodel.Favorite
	if err := f.db.WithContext(ctx).
		Where("user_id = ? AND place_id = ?", userID, placeID).
		Find(&dbModel).Error; err != nil {
		utility.SugarLogger.Errorln("error getting favorite by user id and place id:", err)
		return entitymodel.Favorite{}, err
	}
	return dbModel.ToEntityModel(), nil
}

func (f FavoriteHandle) GetUserFavoritePlaces(ctx context.Context, userID string, country *string, administrativeAreaLevel1 *string, orderBy *enum.FavoriteOrderBy) (entitymodel.UserFavoritePlaces, error) {
	var dbModel dbmodel.PreloadFavoritePlaceUserSlice
	sql := f.db.WithContext(ctx).Model(&dbmodel.PreloadFavoritePlaceUser{})
	if country != nil {
		if administrativeAreaLevel1 != nil {
			sql = sql.Preload("Place", "country = ? AND administrative_area_level_1 = ?", country, administrativeAreaLevel1)
		} else {
			sql = sql.Preload("Place", "country = ?", country)
		}
	} else {
		if administrativeAreaLevel1 != nil {
			sql = sql.Preload("Place", "administrative_area_level_1 = ?", administrativeAreaLevel1)
		} else {
			sql = sql.Preload("Place")
		}
	}
	sql = sql.Preload("User").Where("user_id = ?", userID)
	if orderBy != nil {
		switch *orderBy {
		case enum.OrderByCreatedAtDESC:
			sql = sql.Order("created_at DESC")
		case enum.OrderByCreatedAtASC:
			sql = sql.Order("created_at ASC")
		}
	}
	if err := sql.Find(&dbModel).Error; err != nil {
		utility.SugarLogger.Errorln("error getting user favorite places:", err)
		return entitymodel.UserFavoritePlaces{}, err
	}
	entity := dbModel.ToUserFavoritePlacesSliceEntityModel()
	if len(entity) == 0 {
		return entitymodel.UserFavoritePlaces{}, nil
	}
	return entity[0], nil
}

func (f FavoriteHandle) ToggleFavorite(ctx context.Context, userID, placeID string) error {
	if err := f.db.WithContext(ctx).
		Model(&dbmodel.Favorite{}).
		Where("user_id = ? AND place_id = ?", userID, placeID).
		Update("is_favorite", gorm.Expr("NOT is_favorite")).Error; err != nil {
		utility.SugarLogger.Errorln("error toggling favorite by user id and place id:", err)
		return err
	}
	return nil
}

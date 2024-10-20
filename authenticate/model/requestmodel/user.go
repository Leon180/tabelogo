package requestmodel

import (
	"authenticate/model/entitymodel"
	"authenticate/model/enum"
	"authenticate/utility"
	"time"
)

type RegistUserRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (req *RegistUserRequest) ToEntity() (entitymodel.User, error) {
	hashedPassword, err := utility.HashedPassword(req.Password)
	if err != nil {
		return entitymodel.User{}, err
	}
	return entitymodel.User{
		ID:             utility.GenDefaultUUID(),
		Nickname:       req.Nickname,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		Active:         true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (req *LoginUserRequest) ToEntity() (entitymodel.User, error) {
	hashedPassword, err := utility.HashedPassword(req.Password)
	if err != nil {
		return entitymodel.User{}, err
	}
	return entitymodel.User{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}, nil
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type SaveFavoriteRequest struct {
	IsFavorite    bool   `json:"is_favorite" binding:"required"`
	UserID        string `json:"user_id" binding:"required"`
	PlaceGoogleID string `json:"place_google_id" binding:"required"`
}

func (req *SaveFavoriteRequest) ToEntity() entitymodel.Favorite {
	return entitymodel.Favorite{
		ID:            utility.GenDefaultUUID(),
		UserID:        req.UserID,
		PlaceGoogleID: req.PlaceGoogleID,
		IsFavorite:    req.IsFavorite,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

type GetUserFavoritesRequest struct {
	UserID                   string  `json:"user_id" binding:"required"`
	Country                  *string `json:"country"`
	AdministrativeAreaLevel1 *string `json:"administrative_area_level_1"`
	OrderBy                  *int    `json:"order_by"`
}

func (req *GetUserFavoritesRequest) ToEntity() entitymodel.GetUserFavoritesRequest {
	entity := entitymodel.GetUserFavoritesRequest{
		UserID:                   req.UserID,
		Country:                  req.Country,
		AdministrativeAreaLevel1: req.AdministrativeAreaLevel1,
	}
	if req.OrderBy != nil {
		enumOrderBy := enum.FavoriteOrderBy(*req.OrderBy)
		entity.OrderBy = &enumOrderBy
	}
	return entity
}

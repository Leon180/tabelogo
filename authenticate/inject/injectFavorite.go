package inject

import (
	"authenticate/repository"
	"authenticate/repository/postgresqlRepository"

	"gorm.io/gorm"
)

func provideFavoriteRepository(
	db *gorm.DB,
) repository.FavoriteHandler {
	return postgresqlRepository.NewFavoriteHandler(db)
}

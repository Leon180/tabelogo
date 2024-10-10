package postgresqlRepository

import (
	"authenticate/model/dbmodel"
	"authenticate/model/entitymodel"
	"authenticate/repository"
	"authenticate/utility"
	"context"

	"gorm.io/gorm"
)

func NewUserHandler(db *gorm.DB) repository.UserHandler {
	return &UserHandle{db: db}
}

type UserHandle struct {
	db *gorm.DB
}

func (handle *UserHandle) CreateUser(ctx context.Context, user entitymodel.User) error {
	dbModel := dbmodel.User{
		ID:             user.ID,
		Nickname:       user.Nickname,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		Active:         user.Active,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	if err := handle.db.WithContext(ctx).
		Create(&dbModel).Error; err != nil {
		utility.SugarLogger.Errorln("error creating user:", err)
		return err
	}

	return nil
}

func (handle *UserHandle) GetUserByEmail(ctx context.Context, email string) (entitymodel.User, error) {
	var user dbmodel.User

	if err := handle.db.WithContext(ctx).
		Where("email = ?", email).
		Find(&user).Error; err != nil {
		utility.SugarLogger.Errorln("error getting user by email:", err)
		return entitymodel.User{}, err
	}

	return user.ToEntityModel(), nil
}

func (handle *UserHandle) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {
	if err := handle.db.WithContext(ctx).
		Model(&dbmodel.User{}).
		Where("id = ?", userID).
		Updates(updates).Error; err != nil {
		utility.SugarLogger.Errorln("error updating user:", err)
		return err
	}
	return nil
}

func (handle *UserHandle) DeleteUser(ctx context.Context, userID string) error {
	if err := handle.db.WithContext(ctx).
		Where("id = ?", userID).
		Delete(&dbmodel.User{}).Error; err != nil {
		utility.SugarLogger.Errorln("error deleting user:", err)
		return err
	}
	return nil
}

func NewUserWithTransactionHandler(db *gorm.DB) repository.UserWithTransactionHandler {
	return &UserWithTransactionHandle{
		db:         db,
		UserHandle: UserHandle{db: db},
	}
}

type UserWithTransactionHandle struct {
	db *gorm.DB
	UserHandle
}

func (handle *UserWithTransactionHandle) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return handle.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func NewUserAndSessionWithTransactionHandler(db *gorm.DB) repository.UserAndSessionWithTransactionHandler {
	return &UserAndSessionWithTransactionHandle{
		db:            db,
		UserHandle:    UserHandle{db: db},
		SessionHandle: SessionHandle{db: db},
	}
}

type UserAndSessionWithTransactionHandle struct {
	db *gorm.DB
	UserHandle
	SessionHandle
}

func (handle *UserAndSessionWithTransactionHandle) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return handle.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

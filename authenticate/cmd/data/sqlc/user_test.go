package db

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the CreateUser method
	expectedParams := CreateUserParams{
		Email:          "test@example.com",
		HashedPassword: "hashedpassword123",
	}

	createdAt := time.Now()
	updatedAt := createdAt
	// Rows to be returned by the mock database
	rows := sqlmock.NewRows([]string{
		"user_id", "email", "hashed_password", "active", "created_at", "updated_at",
	}).AddRow(
		1, expectedParams.Email, expectedParams.HashedPassword, true,
		createdAt, updatedAt,
	)

	mock.ExpectQuery(regexp.QuoteMeta(createUser)).
		WithArgs(expectedParams.Email, expectedParams.HashedPassword).
		WillReturnRows(rows)

	// Call the CreateUser method
	createdUser, err := q.CreateUser(context.Background(), expectedParams)

	// Check for errors
	assert.NoError(t, err, "Error executing CreateUser")

	// Expected User result
	expectedUser := User{
		UserID:         1,
		Email:          expectedParams.Email,
		HashedPassword: expectedParams.HashedPassword,
		Active:         true,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	// Compare the actual and expected results using testify/assert
	assert.Equal(t, expectedUser, createdUser, "Unexpected User returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestGetUser(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	// Create a new instance of your Queries struct with the mock database
	q := New(db)

	// Mock the expected database query and result
	createdAt := time.Now()
	updatedAt := createdAt
	rows := sqlmock.NewRows([]string{"user_id", "email", "hashed_password", "active", "created_at", "updated_at"}).
		AddRow(1, "test@example.com", "hashed123", true, createdAt, updatedAt)

	// Expect the query to be executed with the given parameters
	mock.ExpectQuery(regexp.QuoteMeta(getUser)).
		WithArgs("test@example.com").
		WillReturnRows(rows)

	// Call the GetUser function with the mock database
	user, err := q.GetUser(context.Background(), "test@example.com")

	// Check if there were any errors during the execution
	assert.NoError(t, err, "Error executing GetUser")

	// Check if the returned user matches the expected values
	expectedUser := User{
		UserID:         1,
		Email:          "test@example.com",
		HashedPassword: "hashed123",
		Active:         true,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	// Compare the actual and expected results using testify/assert
	assert.Equal(t, expectedUser, user, "Unexpected User returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

package db

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the CreateSession method
	expectedParams := CreateSessionParams{
		SessionID:    uuid.New(),
		Email:        "test@example.com",
		RefreshToken: "refresh123",
		UserAgent:    "Mozilla/5.0",
		ClientIp:     "127.0.0.1",
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	createdAt := time.Now()
	// Rows to be returned by the mock database
	rows := sqlmock.NewRows([]string{
		"session_id", "email", "refresh_token", "user_agent", "client_ip", "is_blocked", "expires_at", "created_at",
	}).AddRow(
		expectedParams.SessionID, expectedParams.Email, expectedParams.RefreshToken,
		expectedParams.UserAgent, expectedParams.ClientIp, expectedParams.IsBlocked,
		expectedParams.ExpiresAt, createdAt,
	)

	// Expectation for the CreateSession query
	mock.ExpectQuery(regexp.QuoteMeta(createSession)).
		WithArgs(
			expectedParams.SessionID,
			expectedParams.Email,
			expectedParams.RefreshToken,
			expectedParams.UserAgent,
			expectedParams.ClientIp,
			expectedParams.IsBlocked,
			expectedParams.ExpiresAt,
		).
		WillReturnRows(rows)

	// Call the CreateSession method
	createdSession, err := q.CreateSession(context.Background(), expectedParams)

	// Check for errors
	assert.NoError(t, err, "Error executing CreateSession")

	// Expected Session result
	expectedSession := Session{
		SessionID:    expectedParams.SessionID,
		Email:        expectedParams.Email,
		RefreshToken: expectedParams.RefreshToken,
		UserAgent:    expectedParams.UserAgent,
		ClientIp:     expectedParams.ClientIp,
		IsBlocked:    expectedParams.IsBlocked,
		ExpiresAt:    expectedParams.ExpiresAt,
		CreatedAt:    createdAt,
	}

	// Compare the actual and expected results using testify/assert
	assert.Equal(t, expectedSession, createdSession, "Unexpected Session returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDeleteSession(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the DeleteSession method
	expectedSessionID := uuid.New()

	// Expectation for the DeleteSession query
	mock.ExpectExec(regexp.QuoteMeta(deleteSession)).
		WithArgs(expectedSessionID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// Call the DeleteSession method
	err = q.DeleteSession(context.Background(), expectedSessionID)

	// Check for errors
	assert.NoError(t, err, "Error executing DeleteSession")

	// Compare the actual and expected results using testify/assert
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestGetSession(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the GetSession method
	expectedSessionID := uuid.New()

	// Rows to be returned by the mock database
	expiredAt := time.Now().Add(time.Hour)
	createdAt := time.Now()
	rows := sqlmock.NewRows([]string{
		"session_id", "email", "refresh_token", "user_agent", "client_ip", "is_blocked", "expires_at", "created_at",
	}).AddRow(
		expectedSessionID, "test@example.com", "refresh_token_value", "user_agent_value", "127.0.0.1", false,
		expiredAt, createdAt,
	)

	// Expectation for the GetSession query
	mock.ExpectQuery(regexp.QuoteMeta(getSession)).
		WithArgs(expectedSessionID).
		WillReturnRows(rows)

	// Call the GetSession method
	resultSession, err := q.GetSession(context.Background(), expectedSessionID)

	// Check for errors
	assert.NoError(t, err, "Error executing GetSession")

	// Expected Session result
	expectedSession := Session{
		SessionID:    expectedSessionID,
		Email:        "test@example.com",
		RefreshToken: "refresh_token_value",
		UserAgent:    "user_agent_value",
		ClientIp:     "127.0.0.1",
		IsBlocked:    false,
		ExpiresAt:    expiredAt,
		CreatedAt:    createdAt,
	}

	// Compare the actual and expected results using testify/assert
	assert.Equal(t, expectedSession, resultSession, "Unexpected Session returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

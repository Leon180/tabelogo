package db

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePlace(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	// Create a new instance of your Queries struct with the mock database
	q := New(db)

	// Define the expected parameters for CreatePlace
	expectedParams := CreatePlaceParams{
		GoogleID:                 "example-google-id",
		TwDisplayName:            "Example Place",
		TwFormattedAddress:       "123 Example St, City",
		TwWeekdayDescriptions:    []string{"Monday", "Tuesday"},
		AdministrativeAreaLevel1: "Example Admin Area",
		Country:                  "Example Country",
		GoogleMapUri:             "https://maps.google.com/example",
		InternationalPhoneNumber: "123-456-789",
		Lat:                      "12.345",
		Lng:                      "67.890",
		PrimaryType:              "restaurant",
		Rating:                   "4.5",
		Types:                    []string{"food", "bar"},
		UserRatingCount:          100,
		WebsiteUri:               "https://example.com",
	}

	// Mock the expected database query and result
	createdAt := time.Now()
	updatedAt := createdAt
	rows := sqlmock.NewRows([]string{
		"google_id", "tw_display_name", "jp_display_name", "tw_formatted_address", "tw_weekday_descriptions",
		"administrative_area_level_1", "country", "google_map_uri", "international_phone_number", "lat", "lng",
		"primary_type", "rating", "types", "user_rating_count", "website_uri", "place_version", "created_at", "updated_at",
	}).AddRow(
		expectedParams.GoogleID, expectedParams.TwDisplayName, "", expectedParams.TwFormattedAddress, pq.Array(expectedParams.TwWeekdayDescriptions),
		expectedParams.AdministrativeAreaLevel1, expectedParams.Country, expectedParams.GoogleMapUri,
		expectedParams.InternationalPhoneNumber, expectedParams.Lat, expectedParams.Lng, expectedParams.PrimaryType,
		expectedParams.Rating, pq.Array(expectedParams.Types), expectedParams.UserRatingCount, expectedParams.WebsiteUri,
		1, createdAt, updatedAt,
	)

	// Expect the query to be executed with the given parameters
	mock.ExpectQuery(regexp.QuoteMeta(createPlace)).
		WithArgs(
			expectedParams.GoogleID, expectedParams.TwDisplayName, expectedParams.TwFormattedAddress, pq.Array(expectedParams.TwWeekdayDescriptions),
			expectedParams.AdministrativeAreaLevel1, expectedParams.Country, expectedParams.GoogleMapUri,
			expectedParams.InternationalPhoneNumber, expectedParams.Lat, expectedParams.Lng, expectedParams.PrimaryType,
			expectedParams.Rating, pq.Array(expectedParams.Types), expectedParams.UserRatingCount, expectedParams.WebsiteUri,
		).
		WillReturnRows(rows)

	// Call the CreatePlace function with the mock database
	place, err := q.CreatePlace(context.Background(), expectedParams)

	// Check if there were any errors during the execution
	assert.NoError(t, err, "Error executing CreatePlace")

	// Check if the returned place matches the expected values
	expectedPlace := Place{
		GoogleID:                 expectedParams.GoogleID,
		TwDisplayName:            expectedParams.TwDisplayName,
		TwFormattedAddress:       expectedParams.TwFormattedAddress,
		TwWeekdayDescriptions:    expectedParams.TwWeekdayDescriptions,
		AdministrativeAreaLevel1: expectedParams.AdministrativeAreaLevel1,
		Country:                  expectedParams.Country,
		GoogleMapUri:             expectedParams.GoogleMapUri,
		InternationalPhoneNumber: expectedParams.InternationalPhoneNumber,
		Lat:                      expectedParams.Lat,
		Lng:                      expectedParams.Lng,
		PrimaryType:              expectedParams.PrimaryType,
		Rating:                   expectedParams.Rating,
		Types:                    expectedParams.Types,
		UserRatingCount:          expectedParams.UserRatingCount,
		WebsiteUri:               expectedParams.WebsiteUri,
		PlaceVersion:             1,
		CreatedAt:                createdAt,
		UpdatedAt:                updatedAt,
	}

	assert.Equal(t, place, expectedPlace, "Unexpected place returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDeletePlace(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	// Create a new instance of your Queries struct with the mock database
	q := New(db)

	// Define the expected parameters for DeletePlace
	expectedParams := DeletePlaceParams{
		GoogleID:     "example-google-id",
		PlaceVersion: 1,
	}

	// Mock the expected database query and result
	mock.ExpectExec(regexp.QuoteMeta(deletePlace)).
		WithArgs(expectedParams.GoogleID, expectedParams.PlaceVersion).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the DeletePlace function with the mock database
	err = q.DeletePlace(context.Background(), expectedParams)

	assert.NoError(t, err, "Error executing DeletePlace")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestGetPlaceByGoogleId(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	// Create a new instance of your Queries struct with the mock database
	q := New(db)

	// Define the expected parameters for GetPlaceByGoogleId
	expectedGoogleID := "example-google-id"

	// Mock the expected database query and result
	createdAt := time.Now()
	updatedAt := createdAt
	rows := sqlmock.NewRows([]string{
		"google_id", "tw_display_name", "jp_display_name", "tw_formatted_address", "tw_weekday_descriptions",
		"administrative_area_level_1", "country", "google_map_uri", "international_phone_number", "lat", "lng",
		"primary_type", "rating", "types", "user_rating_count", "website_uri", "place_version", "created_at", "updated_at",
	}).AddRow(
		expectedGoogleID, "Example Place", "", "123 Example St, City", pq.Array([]string{"Monday", "Tuesday"}),
		"Example Admin Area", "Example Country", "https://maps.google.com/example", "123-456-789", "12.345", "67.890",
		"restaurant", "4.5", pq.Array([]string{"food", "bar"}), 100, "https://example.com", 1, createdAt, updatedAt,
	)

	// Expect the query to be executed with the given parameters
	mock.ExpectQuery(regexp.QuoteMeta(getPlaceByGoogleId)).
		WithArgs(expectedGoogleID).
		WillReturnRows(rows)

	// Call the GetPlaceByGoogleId function with the mock database
	place, err := q.GetPlaceByGoogleId(context.Background(), expectedGoogleID)

	// Check if there were any errors during the execution
	assert.NoError(t, err, "Error executing GetPlaceByGoogleId")

	// Check if the returned place matches the expected values
	expectedPlace := Place{
		GoogleID:                 expectedGoogleID,
		TwDisplayName:            "Example Place",
		TwFormattedAddress:       "123 Example St, City",
		TwWeekdayDescriptions:    []string{"Monday", "Tuesday"},
		AdministrativeAreaLevel1: "Example Admin Area",
		Country:                  "Example Country",
		GoogleMapUri:             "https://maps.google.com/example",
		InternationalPhoneNumber: "123-456-789",
		Lat:                      "12.345",
		Lng:                      "67.890",
		PrimaryType:              "restaurant",
		Rating:                   "4.5",
		Types:                    []string{"food", "bar"},
		UserRatingCount:          100,
		WebsiteUri:               "https://example.com",
		PlaceVersion:             1,
		CreatedAt:                createdAt,
		UpdatedAt:                updatedAt,
	}

	assert.Equal(t, place, expectedPlace, "Unexpected place returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestUpdatePlace(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	q := New(db)

	expectedParams := UpdatePlaceParams{
		GoogleID:              sql.NullString{String: "example-google-id", Valid: true},
		TwDisplayName:         sql.NullString{String: "Updated Place Name", Valid: true},
		TwWeekdayDescriptions: []string{},
		Types:                 []string{},
		UserRatingCount:       sql.NullInt32{Int32: 42, Valid: true},
		PlaceVersion:          1,
	}

	createdAt := time.Now()
	updatedAt := createdAt
	rows := sqlmock.NewRows([]string{
		"google_id", "tw_display_name", "jp_display_name", "tw_formatted_address", "tw_weekday_descriptions",
		"administrative_area_level_1", "country", "google_map_uri", "international_phone_number", "lat", "lng",
		"primary_type", "rating", "types", "user_rating_count", "website_uri", "place_version", "created_at", "updated_at",
	}).AddRow(
		expectedParams.GoogleID.String, expectedParams.TwDisplayName.String, "", "", pq.Array([]string{}),
		"", "", "", "", "", "", "", "", pq.Array([]string{}), int32(expectedParams.UserRatingCount.Int32), "",
		expectedParams.PlaceVersion, createdAt, updatedAt,
	)

	mock.ExpectQuery(regexp.QuoteMeta(updatePlace)).
		WithArgs(
			expectedParams.GoogleID,
			expectedParams.TwDisplayName,
			sql.NullString{},
			sql.NullString{},
			pq.Array([]string{}), // Fix: Ensure pq.Array is used correctly for empty array
			sql.NullString{},
			sql.NullString{},
			sql.NullString{},
			sql.NullString{},
			sql.NullString{},
			sql.NullString{},
			sql.NullString{},
			sql.NullString{},
			pq.Array([]string{}), // Fix: Ensure pq.Array is used correctly for empty array
			expectedParams.UserRatingCount,
			sql.NullString{},
			expectedParams.PlaceVersion,
		).
		WillReturnRows(rows)

	updatedPlace, err := q.UpdatePlace(context.Background(), expectedParams)
	assert.NoError(t, err, "Error executing UpdatePlace")

	expectedPlace := Place{
		GoogleID:        expectedParams.GoogleID.String,
		TwDisplayName:   expectedParams.TwDisplayName.String,
		UserRatingCount: int32(expectedParams.UserRatingCount.Int32),
		PlaceVersion:    expectedParams.PlaceVersion,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}

	if updatedPlace.GoogleID != expectedPlace.GoogleID ||
		updatedPlace.TwDisplayName != expectedPlace.TwDisplayName ||
		updatedPlace.UserRatingCount != expectedPlace.UserRatingCount ||
		updatedPlace.PlaceVersion != expectedPlace.PlaceVersion ||
		!updatedPlace.CreatedAt.Equal(expectedPlace.CreatedAt) ||
		!updatedPlace.UpdatedAt.Equal(expectedPlace.UpdatedAt) {
		t.Errorf("Unexpected place returned. Expected: %v, Got: %v", expectedPlace, updatedPlace)
	}

	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

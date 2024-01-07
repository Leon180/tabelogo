package db

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFavorite(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the CreateFavorite method
	expectedParams := CreateFavoriteParams{
		UserEmail: "test@example.com",
		GoogleID:  "example-google-id",
	}

	createdAt := time.Now()
	updatedAt := createdAt
	// Rows to be returned by the mock database
	rows := sqlmock.NewRows([]string{
		"is_favorite", "user_email", "google_id", "created_at", "updated_at",
	}).AddRow(
		true, expectedParams.UserEmail, expectedParams.GoogleID, createdAt, updatedAt,
	)

	// Expectation for the CreateFavorite query
	mock.ExpectQuery(regexp.QuoteMeta(createFavorite)).
		WithArgs(expectedParams.UserEmail, expectedParams.GoogleID).
		WillReturnRows(rows)

	// Call the CreateFavorite method
	createdFavorite, err := q.CreateFavorite(context.Background(), expectedParams)

	// Check for errors
	assert.NoError(t, err, "Error executing CreateFavorite")

	// Expected Favorite result
	expectedFavorite := Favorite{
		IsFavorite: true,
		UserEmail:  expectedParams.UserEmail,
		GoogleID:   expectedParams.GoogleID,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	// Compare the actual and expected results using testify/assert
	assert.Equal(t, expectedFavorite, createdFavorite, "Unexpected Favorite returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestGetCountryList(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the GetCountryList method
	expectedUserEmail := "test@example.com"

	// Rows to be returned by the mock database
	rows := sqlmock.NewRows([]string{"country"}).
		AddRow("Country1").
		AddRow("Country2").
		AddRow("Country3")

	// Expectation for the GetCountryList query
	mock.ExpectQuery(regexp.QuoteMeta(getCountryList)).
		WithArgs(expectedUserEmail).
		WillReturnRows(rows)

	// Call the GetCountryList method
	countryList, err := q.GetCountryList(context.Background(), expectedUserEmail)

	// Check for errors
	assert.NoError(t, err, "Error executing GetCountryList")

	// Expected country list result
	expectedCountryList := []string{"Country1", "Country2", "Country3"}

	// Compare the actual and expected results using testify/assert
	assert.Equal(t, expectedCountryList, countryList, "Unexpected country list returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestGetFavorite(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the GetFavorite method
	expectedUserEmail := "test@example.com"
	expectedGoogleID := "example-google-id"

	// Rows to be returned by the mock database
	createdAt := time.Now()
	updatedAt := createdAt
	rows := sqlmock.NewRows([]string{"is_favorite", "user_email", "google_id", "created_at", "updated_at"}).
		AddRow(true, expectedUserEmail, expectedGoogleID, createdAt, updatedAt)

	// Expectation for the GetFavorite query
	mock.ExpectQuery(regexp.QuoteMeta(getFavorite)).
		WithArgs(expectedUserEmail, expectedGoogleID).
		WillReturnRows(rows)

	// Call the GetFavorite method
	favorite, err := q.GetFavorite(context.Background(), GetFavoriteParams{
		UserEmail: expectedUserEmail,
		GoogleID:  expectedGoogleID,
	})

	// Check for errors
	assert.NoError(t, err, "Error executing GetFavorite")

	// Expected favorite result
	expectedFavorite := Favorite{
		IsFavorite: true,
		UserEmail:  expectedUserEmail,
		GoogleID:   expectedGoogleID,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	// Compare the actual and expected results using testify/assert
	assert.Equal(t, expectedFavorite, favorite, "Unexpected favorite returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestGetRegionList(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the GetRegionList method
	expectedUserEmail := "test@example.com"
	expectedCountry := "example-country"

	// Rows to be returned by the mock database
	rows := sqlmock.NewRows([]string{"administrative_area_level_1"}).
		AddRow("Region1").
		AddRow("Region2").
		AddRow("Region3")

	// Expectation for the GetRegionList query
	mock.ExpectQuery(regexp.QuoteMeta(getRegionList)).
		WithArgs(expectedUserEmail, expectedCountry).
		WillReturnRows(rows)

	// Call the GetRegionList method
	regionList, err := q.GetRegionList(context.Background(), GetRegionListParams{
		UserEmail: expectedUserEmail,
		Country:   expectedCountry,
	})

	// Check for errors
	assert.NoError(t, err, "Error executing GetRegionList")

	// Expected region list result
	expectedRegionList := []string{"Region1", "Region2", "Region3"}

	// Compare the actual and expected results using testify/assert
	assert.Equal(t, expectedRegionList, regionList, "Unexpected region list returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestListFavoritesByCountrAndRegion(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the ListFavoritesByCountrAndRegion method
	expectedParams := ListFavoritesByCountrAndRegionParams{
		UserEmail:                "test@example.com",
		Country:                  "US",
		AdministrativeAreaLevel1: "CA",
		Limit:                    10,
		Offset:                   0,
	}

	// Rows to be returned by the mock database
	rows := sqlmock.NewRows([]string{
		"google_id", "tw_display_name", "tw_formatted_address", "tw_weekday_descriptions",
		"administrative_area_level_1", "country", "google_map_uri", "international_phone_number",
		"lat", "lng", "primary_type", "rating", "types", "user_rating_count", "website_uri",
	}).AddRow(
		"example-google-id", "Place Name", "Formatted Address", pq.Array([]string{"Monday,Tuesday,Wednesday"}), // Update weekday descriptions
		"CA", "US", "http://maps.google.com", "123456789", "37.7749", "-122.4194", "restaurant",
		"4.5", pq.Array([]string{"type1", "type2"}), 100, "http://example.com",
	)

	// Expectation for the ListFavoritesByCountrAndRegion query
	mock.ExpectQuery(regexp.QuoteMeta(listFavoritesByCountrAndRegion)).
		WithArgs(expectedParams.UserEmail, expectedParams.Country, expectedParams.AdministrativeAreaLevel1, expectedParams.Limit, expectedParams.Offset).
		WillReturnRows(rows)

	// Call the ListFavoritesByCountrAndRegion method
	list, err := q.ListFavoritesByCountrAndRegion(context.Background(), expectedParams)

	// Check for errors
	assert.NoError(t, err, "Error executing ListFavoritesByCountrAndRegion")

	// Expected ListFavoritesByCountrAndRegionRow result
	expectedList := []ListFavoritesByCountrAndRegionRow{
		{
			GoogleID:                 "example-google-id",
			TwDisplayName:            "Place Name",
			TwFormattedAddress:       "Formatted Address",
			TwWeekdayDescriptions:    []string{"Monday,Tuesday,Wednesday"}, // Update weekday descriptions
			AdministrativeAreaLevel1: "CA",
			Country:                  "US",
			GoogleMapUri:             "http://maps.google.com",
			InternationalPhoneNumber: "123456789",
			Lat:                      "37.7749",
			Lng:                      "-122.4194",
			PrimaryType:              "restaurant",
			Rating:                   "4.5",
			Types:                    []string{"type1", "type2"},
			UserRatingCount:          100,
			WebsiteUri:               "http://example.com",
		},
	}

	// Compare the actual and expected results using testify/assert
	assert.Equal(t, expectedList, list, "Unexpected list returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestListFavoritesByCountry(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the ListFavoritesByCountry method
	expectedUserEmail := "test@example.com"
	expectedCountry := "Country1"
	expectedLimit := int32(10)
	expectedOffset := int32(0)

	// Rows to be returned by the mock database
	rows := sqlmock.NewRows([]string{
		"google_id", "tw_display_name", "tw_formatted_address", "tw_weekday_descriptions",
		"administrative_area_level_1", "country", "google_map_uri", "international_phone_number",
		"lat", "lng", "primary_type", "rating", "types", "user_rating_count", "website_uri",
	}).
		AddRow(
			"GoogleID1", "Place1", "Address1", pq.Array([]string{"Mon", "Tue"}), "AdminArea1",
			"Country1", "MapUri1", "123456", "12.345", "67.890", "Type1", "4.5",
			pq.Array([]string{"TypeA", "TypeB"}), 100, "http://example.com/place1",
		).
		AddRow(
			"GoogleID2", "Place2", "Address2", pq.Array([]string{"Wed", "Thu"}), "AdminArea2",
			"Country1", "MapUri2", "789012", "23.456", "78.901", "Type2", "3.7",
			pq.Array([]string{"TypeC", "TypeD"}), 50, "http://example.com/place2",
		)

	// Expectation for the ListFavoritesByCountry query
	mock.ExpectQuery(regexp.QuoteMeta(listFavoritesByCountry)).
		WithArgs(expectedUserEmail, expectedCountry, expectedLimit, expectedOffset).
		WillReturnRows(rows)

	// Call the ListFavoritesByCountry method
	favoritesList, err := q.ListFavoritesByCountry(context.Background(), ListFavoritesByCountryParams{
		UserEmail: expectedUserEmail,
		Country:   expectedCountry,
		Limit:     expectedLimit,
		Offset:    expectedOffset,
	})

	// Check for errors
	assert.NoError(t, err, "Error executing ListFavoritesByCountry")

	// Expected favorites list result
	expectedFavoritesList := []ListFavoritesByCountryRow{
		{
			GoogleID:                 "GoogleID1",
			TwDisplayName:            "Place1",
			TwFormattedAddress:       "Address1",
			TwWeekdayDescriptions:    []string{"Mon", "Tue"},
			AdministrativeAreaLevel1: "AdminArea1",
			Country:                  "Country1",
			GoogleMapUri:             "MapUri1",
			InternationalPhoneNumber: "123456",
			Lat:                      "12.345",
			Lng:                      "67.890",
			PrimaryType:              "Type1",
			Rating:                   "4.5",
			Types:                    []string{"TypeA", "TypeB"},
			UserRatingCount:          100,
			WebsiteUri:               "http://example.com/place1",
		},
		{
			GoogleID:                 "GoogleID2",
			TwDisplayName:            "Place2",
			TwFormattedAddress:       "Address2",
			TwWeekdayDescriptions:    []string{"Wed", "Thu"},
			AdministrativeAreaLevel1: "AdminArea2",
			Country:                  "Country1",
			GoogleMapUri:             "MapUri2",
			InternationalPhoneNumber: "789012",
			Lat:                      "23.456",
			Lng:                      "78.901",
			PrimaryType:              "Type2",
			Rating:                   "3.7",
			Types:                    []string{"TypeC", "TypeD"},
			UserRatingCount:          50,
			WebsiteUri:               "http://example.com/place2",
		},
	}

	// Compare the actual and expected results using testify/assert
	assert.Equal(t, expectedFavoritesList, favoritesList, "Unexpected favorites list returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestListFavoritesByCreateTime(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the ListFavoritesByCreateTime method
	expectedUserEmail := "test@example.com"
	expectedLimit := int32(10)
	expectedOffset := int32(0)

	// Rows to be returned by the mock database
	createdAt1 := time.Now()
	updatedAt1 := createdAt1
	createdAt2 := time.Now()
	updatedAt2 := createdAt2
	rows := sqlmock.NewRows([]string{
		"google_id", "tw_display_name", "tw_formatted_address", "tw_weekday_descriptions",
		"administrative_area_level_1", "country", "google_map_uri", "international_phone_number",
		"lat", "lng", "primary_type", "rating", "types", "user_rating_count", "website_uri",
		"created_at", "updated_at",
	}).
		AddRow(
			"GoogleID1", "Place1", "Address1", pq.Array([]string{"Mon", "Tue"}), "AdminArea1",
			"Country1", "MapUri1", "123456", "12.345", "67.890", "Type1", "4.5",
			pq.Array([]string{"TypeA", "TypeB"}), 100, "http://example.com/place1",
			createdAt1, updatedAt1,
		).
		AddRow(
			"GoogleID2", "Place2", "Address2", pq.Array([]string{"Wed", "Thu"}), "AdminArea2",
			"Country1", "MapUri2", "789012", "23.456", "78.901", "Type2", "3.7",
			pq.Array([]string{"TypeC", "TypeD"}), 50, "http://example.com/place2",
			createdAt2, updatedAt2,
		)

	// Expectation for the ListFavoritesByCreateTime query
	mock.ExpectQuery(regexp.QuoteMeta(listFavoritesByCreateTime)).
		WithArgs(expectedUserEmail, expectedLimit, expectedOffset).
		WillReturnRows(rows)

	// Call the ListFavoritesByCreateTime method
	favoritesList, err := q.ListFavoritesByCreateTime(context.Background(), ListFavoritesByCreateTimeParams{
		UserEmail: expectedUserEmail,
		Limit:     expectedLimit,
		Offset:    expectedOffset,
	})

	// Check for errors
	assert.NoError(t, err, "Error executing ListFavoritesByCreateTime")

	// Expected favorites list result
	expectedFavoritesList := []ListFavoritesByCreateTimeRow{
		{
			GoogleID:                 "GoogleID1",
			TwDisplayName:            "Place1",
			TwFormattedAddress:       "Address1",
			TwWeekdayDescriptions:    []string{"Mon", "Tue"},
			AdministrativeAreaLevel1: "AdminArea1",
			Country:                  "Country1",
			GoogleMapUri:             "MapUri1",
			InternationalPhoneNumber: "123456",
			Lat:                      "12.345",
			Lng:                      "67.890",
			PrimaryType:              "Type1",
			Rating:                   "4.5",
			Types:                    []string{"TypeA", "TypeB"},
			UserRatingCount:          100,
			WebsiteUri:               "http://example.com/place1",
			CreatedAt:                createdAt1,
			UpdatedAt:                updatedAt1,
		},
		{
			GoogleID:                 "GoogleID2",
			TwDisplayName:            "Place2",
			TwFormattedAddress:       "Address2",
			TwWeekdayDescriptions:    []string{"Wed", "Thu"},
			AdministrativeAreaLevel1: "AdminArea2",
			Country:                  "Country1",
			GoogleMapUri:             "MapUri2",
			InternationalPhoneNumber: "789012",
			Lat:                      "23.456",
			Lng:                      "78.901",
			PrimaryType:              "Type2",
			Rating:                   "3.7",
			Types:                    []string{"TypeC", "TypeD"},
			UserRatingCount:          50,
			WebsiteUri:               "http://example.com/place2",
			CreatedAt:                createdAt2,
			UpdatedAt:                updatedAt2,
		},
	}

	// Compare the actual and expected results using testify/assert
	assert.Equal(t, expectedFavoritesList, favoritesList, "Unexpected favorites list returned")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestRemoveFavorite(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the RemoveFavorite method
	expectedUserEmail := "test@example.com"
	expectedGoogleID := "google-id-123"

	// Expectation for the RemoveFavorite query
	mock.ExpectExec(regexp.QuoteMeta(removeFavorite)).
		WithArgs(expectedUserEmail, expectedGoogleID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// Call the RemoveFavorite method
	err = q.RemoveFavorite(context.Background(), RemoveFavoriteParams{
		UserEmail: expectedUserEmail,
		GoogleID:  expectedGoogleID,
	})

	// Check for errors
	assert.NoError(t, err, "Error executing RemoveFavorite")

	// Check if the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestToggleFavorite(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	// Create a new instance of your queries with the mock database
	q := New(db)

	// Expected parameters for the ToggleFavorite method
	expectedUserEmail := "test@example.com"
	expectedGoogleID := "google-id-123"

	createdAt := time.Now()
	updatedAt := time.Now()

	// Rows to be returned by the mock database
	rows := sqlmock.NewRows([]string{
		"is_favorite", "user_email", "google_id", "created_at", "updated_at",
	}).AddRow(
		true, expectedUserEmail, expectedGoogleID, createdAt, updatedAt,
	)

	// Expectation for the ToggleFavorite query
	mock.ExpectQuery(regexp.QuoteMeta(toggleFavorite)).
		WithArgs(expectedUserEmail, expectedGoogleID).
		WillReturnRows(rows)

	// Call the ToggleFavorite method
	toggledFavorite, err := q.ToggleFavorite(context.Background(), ToggleFavoriteParams{
		UserEmail: expectedUserEmail,
		GoogleID:  expectedGoogleID,
	})

	// Check for errors
	assert.NoError(t, err, "Error executing ToggleFavorite")

	// Expected Favorite result
	expectedFavorite := Favorite{
		IsFavorite: true,
		UserEmail:  expectedUserEmail,
		GoogleID:   expectedGoogleID,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	// Compare the actual and expected results using testify/assert
	assert.Equal(t, expectedFavorite, toggledFavorite, "Unexpected Favorite returned")

	// Check if the expectations were met
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

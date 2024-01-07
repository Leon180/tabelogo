package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	// Create a new mock database connection
	db, _, err := sqlmock.New()
	require.NoError(t, err, "Error creating mock database")
	defer db.Close()

	// Call the NewStore function with the mock database
	store := NewStore(db)

	// Check if the returned store is not nil
	assert.NotNil(t, store, "NewStore returned a nil store")

	// Check if the Queries field of the store is not nil
	assert.NotNil(t, store.(*SQLStore).Queries, "NewStore returned a store with a nil Queries field")

	// Check if the DB field of the store is the expected database connection
	assert.Equal(t, store.(*SQLStore).db, db, "NewStore returned a store with an unexpected DB field")
}

package main

import (
	db "authenticate/cmd/data/sqlc"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store) *Server {

	// Add your routes to the router
	key, err := generateRandomString(32)
	require.NoError(t, err, "Error generating random string")
	config := Config{TokenSymmetricKey: key}
	server, err := NewServer(config, store, nil, CacheInstance{})
	require.NoError(t, err, "Error creating server")

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

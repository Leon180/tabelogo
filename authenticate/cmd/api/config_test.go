package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary configuration file for testing
	viper.Reset()
	viper.SetConfigType("env")
	viper.Set("DSN_TEST", "test_dsn")
	viper.Set("DSN_DEPLOYMENT", "deployment_dsn")
	viper.Set("TOKEN_SYMMETRIC_KEY", "symmetric_key")
	viper.Set("ACCESS_TOKEN_DURATION", "1h")
	viper.Set("REFRESH_TOKEN_DURATION", "24h")
	viper.Set("REFRESH_DURATION", "30m")
	viper.Set("RABBITMQ_CONNECT", "rabbitmq_connection")
	viper.Set("REDIS_CONNECT_SESSION", "redis_session_connection")
	viper.Set("REDIS_CONNECT_PLACE", "redis_place_connection")

	// Create a temporary directory for the configuration file
	tempDir := t.TempDir()

	// Set the name of the configuration file
	viper.SetConfigName("app")

	// Write the configuration file to the temporary directory
	err := viper.SafeWriteConfigAs(filepath.Join(tempDir, "app.env"))
	if err != nil {
		t.Fatalf("Error writing config file: %v", err)
	}
	defer func() {
		// Remove the temporary configuration file after the test
		err := os.Remove(filepath.Join(tempDir, "app.env"))
		if err != nil {
			t.Fatalf("Error removing temp config file: %v", err)
		}
	}()

	// Load the configuration
	config, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	// Assert the values
	assert.Equal(t, "test_dsn", config.DSNTest)
	assert.Equal(t, "deployment_dsn", config.DSNDeployment)
	assert.Equal(t, "symmetric_key", config.TokenSymmetricKey)
	assert.Equal(t, 1*time.Hour, config.AccessTokenDuration)
	assert.Equal(t, 24*time.Hour, config.RefreshTokenDuration)
	assert.Equal(t, 30*time.Minute, config.RefreshDuration)
	assert.Equal(t, "rabbitmq_connection", config.RabbitMQConnect)
	assert.Equal(t, "redis_session_connection", config.RedisConnectSession)
	assert.Equal(t, "redis_place_connection", config.RedisConnectPlace)
}

package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterSecrets_NoFilter(t *testing.T) {
	secrets := map[string]map[string]interface{}{
		"secret/app/db": {"password": "s3cr3t"},
		"secret/app/api": {"key": "abc123"},
	}

	result := FilterSecrets(secrets, FilterOptions{})
	assert.Equal(t, secrets, result)
}

func TestFilterSecrets_Prefix(t *testing.T) {
	secrets := map[string]map[string]interface{}{
		"secret/app/db":    {"password": "s3cr3t"},
		"secret/other/key": {"token": "xyz"},
	}

	result := FilterSecrets(secrets, FilterOptions{Prefix: "secret/app"})
	assert.Len(t, result, 1)
	_, ok := result["secret/app/db"]
	assert.True(t, ok)
}

func TestFilterSecrets_ExcludeKeys(t *testing.T) {
	secrets := map[string]map[string]interface{}{
		"secret/app/db": {"password": "s3cr3t", "username": "admin"},
	}

	result := FilterSecrets(secrets, FilterOptions{ExcludeKeys: []string{"password"}})
	assert.Equal(t, map[string]interface{}{"username": "admin"}, result["secret/app/db"])
}

func TestFilterSecrets_ExcludeAllKeysDropsPath(t *testing.T) {
	secrets := map[string]map[string]interface{}{
		"secret/app/db": {"password": "s3cr3t"},
	}

	result := FilterSecrets(secrets, FilterOptions{ExcludeKeys: []string{"password"}})
	assert.Empty(t, result)
}

func TestFilterSecrets_PrefixAndExclude(t *testing.T) {
	secrets := map[string]map[string]interface{}{
		"secret/app/db":    {"password": "s3cr3t", "host": "localhost"},
		"secret/other/key": {"token": "xyz"},
	}

	result := FilterSecrets(secrets, FilterOptions{
		Prefix:      "secret/app",
		ExcludeKeys: []string{"password"},
	})

	assert.Len(t, result, 1)
	assert.Equal(t, map[string]interface{}{"host": "localhost"}, result["secret/app/db"])
}

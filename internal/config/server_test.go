package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPort(t *testing.T) {
	t.Run("default_port", func(t *testing.T) {
		os.Unsetenv("PORT")
		assert.Equal(t, "3000", GetPort())
	})

	t.Run("custom_port", func(t *testing.T) {
		os.Setenv("PORT", "8080")
		defer os.Unsetenv("PORT")
		assert.Equal(t, "8080", GetPort())
	})
}

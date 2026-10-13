package config

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	t.Run("default_level", func(t *testing.T) {
		_ = os.Unsetenv("LOG_LEVEL")
		InitLogger()
		assert.Equal(t, logrus.InfoLevel, logrus.GetLevel())
	})

	t.Run("custom_valid_level", func(t *testing.T) {
		_ = os.Setenv("LOG_LEVEL", "warn")
		defer os.Unsetenv("LOG_LEVEL")
		InitLogger()
		assert.Equal(t, logrus.WarnLevel, logrus.GetLevel())
	})

	t.Run("invalid_level_fallback_debug", func(t *testing.T) {
		_ = os.Setenv("LOG_LEVEL", "invalid_level_xyz")
		defer os.Unsetenv("LOG_LEVEL")
		InitLogger()
		assert.Equal(t, logrus.DebugLevel, logrus.GetLevel())
	})
}

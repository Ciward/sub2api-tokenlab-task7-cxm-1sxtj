package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidOpenMode(t *testing.T) {
	t.Run("valid_open_mode_values", func(t *testing.T) {
		// Empty string is valid (means default/iframe for backward compatibility)
		assert.True(t, IsValidOpenMode(""))
		// Explicit iframe mode
		assert.True(t, IsValidOpenMode("iframe"))
		// Explicit redirect mode
		assert.True(t, IsValidOpenMode("redirect"))
	})

	t.Run("invalid_open_mode_values", func(t *testing.T) {
		assert.False(t, IsValidOpenMode("invalid"))
		assert.False(t, IsValidOpenMode("popup"))
		assert.False(t, IsValidOpenMode("new_tab"))
		assert.False(t, IsValidOpenMode("external"))
		assert.False(t, IsValidOpenMode("IFRAME")) // case-sensitive
	})
}

func TestCustomMenuItemOpenMode(t *testing.T) {
	t.Run("default_open_mode_is_empty_for_backward_compat", func(t *testing.T) {
		// Zero value CustomMenuItem has empty OpenMode
		item := CustomMenuItem{
			ID:    "test",
			Label: "Test",
			URL:   "https://example.com",
		}
		assert.Equal(t, CustomMenuOpenMode(""), item.OpenMode)
		// Empty is considered valid (defaults to iframe behavior)
		assert.True(t, IsValidOpenMode(string(item.OpenMode)))
	})

	t.Run("explicit_iframe_mode", func(t *testing.T) {
		item := CustomMenuItem{
			ID:       "test",
			Label:    "Test",
			URL:      "https://example.com",
			OpenMode: CustomMenuOpenModeIframe,
		}
		assert.Equal(t, CustomMenuOpenModeIframe, item.OpenMode)
	})

	t.Run("explicit_redirect_mode", func(t *testing.T) {
		item := CustomMenuItem{
			ID:       "test",
			Label:    "Test",
			URL:      "https://example.com",
			OpenMode: CustomMenuOpenModeRedirect,
		}
		assert.Equal(t, CustomMenuOpenModeRedirect, item.OpenMode)
	})
}

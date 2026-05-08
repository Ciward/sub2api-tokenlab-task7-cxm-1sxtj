package service

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCustomMenuItemURLs(t *testing.T) {
	t.Run("returns_nil_for_empty_input", func(t *testing.T) {
		assert.Nil(t, parseCustomMenuItemURLs(""))
		assert.Nil(t, parseCustomMenuItemURLs("  "))
		assert.Nil(t, parseCustomMenuItemURLs("[]"))
	})

	t.Run("returns_nil_for_invalid_json", func(t *testing.T) {
		assert.Nil(t, parseCustomMenuItemURLs("not json"))
		assert.Nil(t, parseCustomMenuItemURLs("{invalid}"))
	})

	t.Run("extracts_urls_from_iframe_items", func(t *testing.T) {
		// Default (no open_mode) means iframe
		input := `[
			{"id": "1", "url": "https://app1.example.com/page"},
			{"id": "2", "url": "https://app2.example.com/dashboard", "open_mode": "iframe"},
			{"id": "3", "url": ""}
		]`
		result := parseCustomMenuItemURLs(input)
		sort.Strings(result)
		assert.Equal(t, []string{
			"https://app1.example.com/page",
			"https://app2.example.com/dashboard",
		}, result)
	})

	t.Run("excludes_redirect_mode_items", func(t *testing.T) {
		// redirect mode items should NOT be included in CSP frame-src
		input := `[
			{"id": "1", "url": "https://iframe.example.com", "open_mode": "iframe"},
			{"id": "2", "url": "https://redirect.example.com", "open_mode": "redirect"},
			{"id": "3", "url": "https://default.example.com"}
		]`
		result := parseCustomMenuItemURLs(input)
		sort.Strings(result)
		assert.Equal(t, []string{
			"https://default.example.com",
			"https://iframe.example.com",
		}, result)
		// Important: redirect URL should NOT be in the result
		assert.NotContains(t, result, "https://redirect.example.com")
	})

	t.Run("handles_all_redirect_items", func(t *testing.T) {
		input := `[
			{"id": "1", "url": "https://ext1.example.com", "open_mode": "redirect"},
			{"id": "2", "url": "https://ext2.example.com", "open_mode": "redirect"}
		]`
		result := parseCustomMenuItemURLs(input)
		// All redirect items should be excluded from CSP
		assert.Empty(t, result)
	})

	t.Run("backward_compatible_with_legacy_format", func(t *testing.T) {
		// Legacy format: no open_mode field at all
		// Should be treated as iframe mode for backward compatibility
		input := `[
			{"id": "1", "url": "https://legacy1.example.com/page", "label": "Page 1"},
			{"id": "2", "url": "https://legacy2.example.com/dash", "label": "Dashboard"}
		]`
		result := parseCustomMenuItemURLs(input)
		sort.Strings(result)
		assert.Equal(t, []string{
			"https://legacy1.example.com/page",
			"https://legacy2.example.com/dash",
		}, result)
	})
}

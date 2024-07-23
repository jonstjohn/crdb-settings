package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSettingsCompareRegex(t *testing.T) {
	url := "/settings/compare/v23.1.5..23.1.6"
	assert.True(t, SettingsCompareReWithReleases.Match([]byte(url)))

	matches := SettingsCompareReWithReleases.FindStringSubmatch(url)
	assert.Len(t, matches, 3)
}

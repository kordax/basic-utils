package umap_test

import (
	"testing"

	"git.casinomodule.org/casino27/basic-utils/v4/umap"
	"github.com/stretchr/testify/assert"
)

func TestMapTransformHelpers(t *testing.T) {
	values := map[string]int{"a": 1, "bb": 2, "cc": 2}

	assert.Equal(t, map[string]int{"bb": 2, "cc": 2}, umap.Filter(values, func(_ string, v int) bool {
		return v == 2
	}))
	assert.Equal(t, map[string]string{"a": "a:1", "bb": "bb:2", "cc": "cc:2"}, umap.MapValues(values, func(k string, v int) string {
		return k + ":" + string(rune('0'+v))
	}))
	assert.Equal(t, map[int]int{1: 1, 2: 2}, umap.MapKeys(values, func(k string, v int) int {
		return len(k)
	}))
	assert.ElementsMatch(t, []string{"bb", "cc"}, umap.Invert(values)[2])
}

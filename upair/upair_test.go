/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package upair

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPair(t *testing.T) {
	pair := NewPair(1, "test")

	require.NotNil(t, pair, "NewPair should not return nil")
	assert.Equal(t, 1, pair.GetLeft(), "GetLeft should return the correct left value")
	assert.Equal(t, "test", pair.GetRight(), "GetRight should return the correct right value")
}

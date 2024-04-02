/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uos_test

import (
	"testing"

	"github.com/kordax/basic-utils/uos"
	"github.com/stretchr/testify/assert"
)

func TestGetCPUs_Stub(t *testing.T) {
	assert.NotZero(t, uos.GetCPUs())
}

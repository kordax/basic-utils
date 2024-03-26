/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package os_utils_test

import (
	"testing"

	osutils "github.com/kordax/basic-utils/os-utils"
	"github.com/stretchr/testify/assert"
)

func TestGetCPUs_Stub(t *testing.T) {
	assert.NotZero(t, osutils.GetCPUs())
}

/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset_test

import (
	"testing"

	"github.com/kordax/basic-utils/v2/uset"
)

func Test_Implements(t *testing.T) {
	var _ uset.Set[int] = (*uset.HashSet[int])(nil)
	var _ uset.Set[int] = (*uset.ConcurrentHashSet[int])(nil)
	var _ uset.Set[testElement] = (*uset.ComparableHashSet[testElement, int])(nil)
	var _ uset.OrderedSet[testElement, int] = (*uset.OrderedHashSet[testElement, int])(nil)
}

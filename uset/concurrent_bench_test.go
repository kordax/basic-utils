/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package uset_test

import (
	"testing"

	"github.com/kordax/basic-utils/uset"
)

func BenchmarkConcurrentHashSet_Add(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		set := uset.NewConcurrentHashSet[int]()
		b.StartTimer()
		for i := 0; i < benchSetSize; i++ {
			set.Add(i)
		}
	}
}

func BenchmarkConcurrentHashSet_Contains(b *testing.B) {
	set := uset.NewConcurrentHashSet[int]()
	for i := 0; i < benchSetSize; i++ {
		set.Add(i)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchSetSize; i++ {
			set.Contains(i)
		}
	}
}

func BenchmarkConcurrentHashSet_Remove(b *testing.B) {
	set := uset.NewConcurrentHashSet[int]()
	for i := 0; i < benchSetSize; i++ {
		set.Add(i)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchSetSize; i++ {
			set.Remove(i)
		}
	}
}

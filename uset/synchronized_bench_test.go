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

func BenchmarkSynchronizedHashSet_Add(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		set := uset.NewSynchronizedHashSet[int]()
		b.StartTimer()
		for i := 0; i < benchSetSize; i++ {
			set.Add(i)
		}
	}
}

func BenchmarkSynchronizedHashSet_Contains(b *testing.B) {
	set := uset.NewSynchronizedHashSet[int]()
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

func BenchmarkSynchronizedHashSet_Remove(b *testing.B) {
	set := uset.NewSynchronizedHashSet[int]()
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

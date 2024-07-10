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

const benchSetSize = 100

func BenchmarkHashSet_Add(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		set := uset.NewHashSet[int]()
		b.StartTimer()
		for i := 0; i < benchSetSize; i++ {
			set.Add(i)
		}
	}
}

func BenchmarkHashSet_Contains(b *testing.B) {
	set := uset.NewHashSet[int]()
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

func BenchmarkHashSet_Remove(b *testing.B) {
	set := uset.NewHashSet[int]()
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

func BenchmarkMap_Add(b *testing.B) {
	for n := 0; n < b.N; n++ {
		m := make(map[int]struct{})
		for i := 0; i < benchSetSize; i++ {
			m[i] = struct{}{}
		}
	}
}

func BenchmarkMap_Contains(b *testing.B) {
	m := make(map[int]struct{})
	for i := 0; i < benchSetSize; i++ {
		m[i] = struct{}{}
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchSetSize; i++ {
			_ = m[i]
		}
	}
}

func BenchmarkMap_Remove(b *testing.B) {
	m := make(map[int]struct{})
	for i := 0; i < benchSetSize; i++ {
		m[i] = struct{}{}
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < benchSetSize; i++ {
			delete(m, i)
		}
	}
}

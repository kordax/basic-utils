/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ucache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kordax/basic-utils/ucache"
	"github.com/kordax/basic-utils/uopt"
)

func BenchmarkFarmHash64Entity(b *testing.B) {
	type testEntity struct {
		MyInteger int
		MyString  string
		MyFloat   float64
		MyMap     map[string]int
	}

	numItems := 10000
	cache := ucache.NewInMemoryHashMapCache[*ucache.FarmHash64Entity, int](uopt.Null[time.Duration]())
	keys := make([]*ucache.FarmHash64Entity, numItems)
	for i := 0; i < numItems; i++ {
		str := fmt.Sprint(i)
		keys[i] = ucache.Hashed(testEntity{
			MyInteger: i,
			MyString:  str,
			MyFloat:   float64(i),
			MyMap: map[string]int{
				str: i,
			},
		})
		cache.Set(keys[i], i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Get(keys[i%numItems])
	}
}

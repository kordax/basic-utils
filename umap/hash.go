/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package umap

import (
	"encoding/binary"
	"fmt"
	"hash"
	"reflect"
	"sort"
)

func computeHash[V any](h hash.Hash, v V) int64 {
	h.Reset()
	val := reflect.ValueOf(v)
	writeHashData(h, val)
	hashBytes := h.Sum(nil)
	return int64(binary.LittleEndian.Uint64(hashBytes[:8]))
}

func writeHashData(h hash.Hash, val reflect.Value) {
	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			writeHashData(h, val.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			writeHashData(h, val.Index(i))
		}
	case reflect.Map:
		keys := val.MapKeys()
		sortKeys := make([]reflectMapKey, len(keys))
		for i, key := range keys {
			sortKeys[i] = reflectMapKey{
				value: key,
				sort:  fmt.Sprint(key),
			}
		}
		sort.Slice(sortKeys, func(i, j int) bool {
			return sortKeys[i].sort < sortKeys[j].sort
		})
		for _, key := range sortKeys {
			writeHashData(h, key.value)
			writeHashData(h, val.MapIndex(key.value))
		}
	case reflect.Pointer, reflect.Interface:
		if !val.IsNil() {
			writeHashData(h, val.Elem())
		}
	default:
		_, _ = fmt.Fprintf(h, "%v", val.Interface())
	}
}

type reflectMapKey struct {
	value reflect.Value
	sort  string
}

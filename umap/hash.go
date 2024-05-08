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
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprintf("%v", keys[i]) < fmt.Sprintf("%v", keys[j])
		})
		for _, key := range keys {
			writeHashData(h, key)
			writeHashData(h, val.MapIndex(key))
		}
	case reflect.Ptr, reflect.Interface:
		if !val.IsNil() {
			writeHashData(h, val.Elem())
		}
	default:
		h.Write([]byte(fmt.Sprintf("%v", val.Interface())))
	}
}

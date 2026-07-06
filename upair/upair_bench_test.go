/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package upair

import "testing"

var pairSink *Pair[int, string]

func BenchmarkOf(b *testing.B) {
	var result Pair[int, string]

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result = Of(i, "value")
	}

	if result.Right == "" {
		b.Fatal("unexpected empty result")
	}
}

func BenchmarkNewPair(b *testing.B) {
	var result *Pair[int, string]

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result = NewPair(i, "value")
	}

	pairSink = result
}

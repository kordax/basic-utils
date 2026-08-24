package ustream_test

import (
	"slices"
	"testing"

	"git.casinomodule.org/casino27/basic-utils/v4/ustream"
)

const comparisonStreamSize = 10_000

var (
	benchmarkComparisonBool   bool
	benchmarkComparisonInt    int
	benchmarkComparisonValues []int
)

func BenchmarkStreamComparison_Map(b *testing.B) {
	cheapValues := comparisonValues(comparisonStreamSize)
	cheapStream := ustream.Of(cheapValues)
	cheapMapper := func(value int) int { return value * 2 }

	b.Run("cheap/for-loop", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result := make([]int, len(cheapValues))
			for index, value := range cheapValues {
				result[index] = cheapMapper(value)
			}
			benchmarkComparisonValues = result
		}
	})
	b.Run("cheap/stream", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkComparisonValues = cheapStream.Map(cheapMapper).Collect()
		}
	})
	b.Run("cheap/parallel-stream-4", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkComparisonValues = cheapStream.ParallelMap(cheapMapper, 4).Collect()
		}
	})

	cpuValues := comparisonValues(1_024)
	cpuStream := ustream.Of(cpuValues)
	cpuMapper := func(value int) int {
		for range 2_048 {
			value = value*1_664_525 + 1_013_904_223
		}
		return value
	}

	b.Run("cpu/for-loop", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result := make([]int, len(cpuValues))
			for index, value := range cpuValues {
				result[index] = cpuMapper(value)
			}
			benchmarkComparisonValues = result
		}
	})
	b.Run("cpu/stream", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkComparisonValues = cpuStream.Map(cpuMapper).Collect()
		}
	})
	b.Run("cpu/parallel-stream-4", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkComparisonValues = cpuStream.ParallelMap(cpuMapper, 4).Collect()
		}
	})
}

func BenchmarkStreamComparison_FusedPipeline(b *testing.B) {
	values := comparisonValues(comparisonStreamSize)
	expected := comparisonForPipeline(values)
	pipeline := ustream.Of(values).
		Filter(func(value int) bool { return value%2 == 0 }).
		Map(func(value int) int { return value * 2 }).
		Skip(10).
		Limit(1_000)
	actual := pipeline.Reduce(0, func(result, value int) int { return result + value })
	if actual != expected {
		b.Fatalf("pipeline mismatch: stream=%d for-loop=%d", actual, expected)
	}

	b.Run("for-loop", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkComparisonInt = comparisonForPipeline(values)
		}
	})
	b.Run("staged-slices", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkComparisonInt = comparisonStagedPipeline(values)
		}
	})
	b.Run("lazy-stream-build-and-run", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkComparisonInt = ustream.Of(values).
				Filter(func(value int) bool { return value%2 == 0 }).
				Map(func(value int) int { return value * 2 }).
				Skip(10).
				Limit(1_000).
				Reduce(0, func(result, value int) int { return result + value })
		}
	})
	b.Run("lazy-stream-reused", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkComparisonInt = pipeline.Reduce(0, func(result, value int) int {
				return result + value
			})
		}
	})
}

func BenchmarkStreamComparison_ShortCircuit(b *testing.B) {
	values := comparisonValues(comparisonStreamSize)
	stream := ustream.Of(values)
	positions := []struct {
		name   string
		target int
	}{
		{name: "early", target: 10},
		{name: "middle", target: comparisonStreamSize / 2},
		{name: "late", target: comparisonStreamSize - 1},
		{name: "missing", target: -1},
	}

	for _, position := range positions {
		b.Run(position.name, func(b *testing.B) {
			b.Run("for-loop", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					found := false
					for _, value := range values {
						if value == position.target {
							found = true
							break
						}
					}
					benchmarkComparisonBool = found
				}
			})
			b.Run("stream-any-match", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					benchmarkComparisonBool = stream.AnyMatch(func(value int) bool {
						return value == position.target
					})
				}
			})
		})
	}
}

func BenchmarkStreamComparison_BufferingBarriers(b *testing.B) {
	values := comparisonValues(comparisonStreamSize)

	b.Run("reverse", func(b *testing.B) {
		b.Run("for-loop", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				result := make([]int, len(values))
				for index, value := range values {
					result[len(values)-index-1] = value
				}
				benchmarkComparisonValues = result
			}
		})
		b.Run("stream", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkComparisonValues = ustream.Of(values).Reverse().Collect()
			}
		})
	})

	unsorted := make([]int, comparisonStreamSize)
	for index := range unsorted {
		unsorted[index] = index * 7_919 % 10_007
	}
	b.Run("sort", func(b *testing.B) {
		b.Run("slices-sort", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				result := slices.Clone(unsorted)
				slices.Sort(result)
				benchmarkComparisonValues = result
			}
		})
		b.Run("stream", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchmarkComparisonValues = ustream.Of(unsorted).
					Sort(func(left, right int) bool { return left < right }).
					Collect()
			}
		})
	})
}

func BenchmarkStreamComparison_Distinct(b *testing.B) {
	cases := []struct {
		name   string
		values []int
	}{
		{name: "all-unique", values: comparisonValues(comparisonStreamSize)},
		{name: "repeated-512", values: comparisonRepeatedValues(comparisonStreamSize, 512)},
	}

	for _, benchmarkCase := range cases {
		b.Run(benchmarkCase.name, func(b *testing.B) {
			b.Run("for-loop", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					seen := make(map[int]struct{}, len(benchmarkCase.values))
					result := make([]int, 0, len(benchmarkCase.values))
					for _, value := range benchmarkCase.values {
						if _, exists := seen[value]; exists {
							continue
						}
						seen[value] = struct{}{}
						result = append(result, value)
					}
					benchmarkComparisonValues = result
				}
			})
			b.Run("stream", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					benchmarkComparisonValues = ustream.Of(benchmarkCase.values).
						DistinctBy(func(value int) int { return value }).
						Collect()
				}
			})
		})
	}
}

func BenchmarkStreamComparison_FlatMap(b *testing.B) {
	values := comparisonValues(comparisonStreamSize)

	b.Run("for-loop", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result := make([]int, 0, len(values)*2)
			for _, value := range values {
				result = append(result, value, -value)
			}
			benchmarkComparisonValues = result
		}
	})
	b.Run("stream", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkComparisonValues = ustream.Of(values).
				FlatMap(func(value int) []int { return []int{value, -value} }).
				Collect()
		}
	})
}

func comparisonValues(size int) []int {
	values := make([]int, size)
	for index := range values {
		values[index] = index
	}
	return values
}

func comparisonRepeatedValues(size, cardinality int) []int {
	values := make([]int, size)
	for index := range values {
		values[index] = index % cardinality
	}
	return values
}

func comparisonForPipeline(values []int) int {
	result := 0
	skipped := 0
	accepted := 0
	for _, value := range values {
		if value%2 != 0 {
			continue
		}
		value *= 2
		if skipped < 10 {
			skipped++
			continue
		}
		if accepted == 1_000 {
			break
		}
		result += value
		accepted++
	}
	return result
}

func comparisonStagedPipeline(values []int) int {
	filtered := make([]int, 0, len(values)/2)
	for _, value := range values {
		if value%2 == 0 {
			filtered = append(filtered, value)
		}
	}
	mapped := make([]int, len(filtered))
	for index, value := range filtered {
		mapped[index] = value * 2
	}

	start := min(10, len(mapped))
	end := min(start+1_000, len(mapped))
	result := 0
	for _, value := range mapped[start:end] {
		result += value
	}
	return result
}
